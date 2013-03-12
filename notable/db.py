# Python imports
import datetime
import logging
import os
import re
import sqlite3
import sys
import uuid

# Project imports
from . import crypt

# Constants
log = logging.getLogger(__name__)
mod = sys.modules.get(__name__)
path = os.path.expanduser('~/.notable/notes.sqlite3')

def note(exclude=None, actual=False):
    exclude = exclude or []
    model = [('uid', 'string'),
             ('created', 'string'),
             ('updated', 'string'),
             ('subject', 'string'),
             ('tags', 'string'),
             ('content', 'string'),
             ('encrypted', 'integer'),
            ]
    n = dict((k, v) for k, v in model if not k in exclude)
    if actual:
        uid = uuid.uuid4().hex
        n.update(uid=uid, created=now(), updated=now())
    return n

def create_note(n, password=None):
    c = conn()
    n = encrypt(n, password)
    sql = 'INSERT INTO notes (%s) VALUES (%s);'
    cols = ','.join(k for k, t in list(n.items()) if not t is None)
    values = ','.join('?' for k, t in list(n.items()) if not t is None)
    sql = sql % (cols, values)
    values = [v for v in list(n.values()) if not v is None]
    log.debug('%s, %s' % (sql, values))
    c.execute(sql, values)
    c.commit()
    n.pop('content')
    return n

def delete_note(uid):
    c = conn()
    sql = 'DELETE FROM notes WHERE uid = ?'
    log.debug('%s, %s' % (sql, uid))
    c.execute(sql, (uid,))
    c.commit()
    return True

def encrypt(n, password):
    encrypted = False
    content = n['content']
    if password:
        content = crypt.encrypt(content, password)
        encrypted = True
    n.update(dict(content=content, encrypted=encrypted))
    return n

def now():
    return datetime.datetime.now().strftime('%Y-%m-%d %H:%M')

def update_note(n, password=None):
    c = conn()
    n = encrypt(n, password)
    n['updated'] = now()
    sql = """
      UPDATE notes SET
        tags = ?,
        subject = ?,
        content = ?,
        encrypted = ?,
        updated = ?
      WHERE uid = ?
      """
    values = [n.get('tags'),
              n.get('subject'),
              n.get('content'),
              n.get('encrypted'),
              n.get('updated'),
              n.get('uid')]
    log.debug('%s, %s' % (sql, values))
    c.execute(sql, values)
    c.commit()
    n.pop('content')
    return n

def conn():
    d = os.path.dirname(path)
    _ = os.makedirs(d) if not os.path.exists(d) else None
    return sqlite3.connect(path)

def dict_factory(c, row):
    return dict((col[0], row[i]) for i, col in enumerate(c.description))

def create_schema(c):
    pairs = [' '.join(pair) for pair in list(note().items()) if pair[1]]
    sql ='CREATE TABLE notes (%s);' % ','.join(pairs)
    try:
        c.execute(sql)
    except sqlite3.OperationalError as ex:
        log.debug(ex)

def get_content(uid, password):
    c = conn()
    c.row_factory = dict_factory
    cursor = c.cursor()
    cursor.execute('SELECT content FROM notes WHERE uid = ?', [uid])
    # Assume python3 and fallback to python2
    content = getattr(cursor, 'fetchone', 'next')().get('content')
    return crypt.decrypt(content, password) if password else content

def migrate_data(c):
    # Migrate 1 encrypted contents to column + dedicated content block
    sql ="""
        UPDATE notes SET
            subject = ?,
            content = ?,
            encrypted = ?
        WHERE uid = ?
    """
    c.rollback()
    for n in search(''):
        if n['subject']:
            continue
        encrypted = False
        content = n['content'].split('\n')
        if len(content) > 1 and re.search(r'([^ ]{32,})', content[1]):
            encrypted = True
        values = (content.pop(0), '\n'.join(content).strip(), encrypted, n['uid'])
        c.execute(sql, values)
        log.debug('Migrated: %s', values[0])
    c.commit()

def migrate_schema(c):
    sqls =['ALTER TABLE notes ADD COLUMN encrypted INTEGER DEFAULT 0;',
           ' ALTER TABLE notes ADD COLUMN subject TEXT;']
    try:
        for sql in sqls:
            c.execute(sql)
    except sqlite3.OperationalError as ex:
        log.debug('schema migration: %s', ex)

def search(s, exclude=None):
    terms = s.split() if s else []
    n = note(exclude=exclude)
    naive = "(content LIKE '%{0}%' OR tags LIKE '%{0}%')"
    where = ['1=1'] + [naive.format(t) for t in terms]
    sql = 'SELECT %s FROM notes WHERE %s ORDER BY updated DESC;'
    sql = sql % (','.join(k for k, v in list(n.items()) if v), ' AND '.join(where))
    log.debug(sql)
    c = conn()
    c.row_factory = dict_factory
    return c.cursor().execute(sql)

def prepare():
    create_schema(conn())
    migrate_schema(conn())
    migrate_data(conn())
