# Python imports
from collections import OrderedDict
import datetime
import logging
import os
import re
import sqlite3
import sys
import uuid

# Project imports
import crypt

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
    n = OrderedDict((k, v) for k, v in model if not k in exclude)
    if actual:
        now = datetime.datetime.now().strftime('%Y-%m-%d %H:%M')
        uid = uuid.uuid4().hex
        n.update(uid=uid, created=now, updated=now)
    return n

def create_note(n, password=None):
    c = conn()
    n = encrypt(n, password)
    sql = 'INSERT INTO notes (%s) VALUES (%s);'
    cols = ','.join(k for k, t in n.items() if not t is None)
    values = ','.join('?' for k, t in n.items() if not t is None)
    sql = sql % (cols, values)
    values = [v for v in n.values() if not v is None]
    log.debug('%s, %s' % (sql, values))
    c.execute(sql, values)
    return c.commit()

def delete_note(uid, password=None):
    c = conn()
    sql = 'DELETE FROM notes WHERE uid = ?'
    log.debug('%s, %s' % (sql, uid))
    c.execute(sql, (uid,))
    return c.commit()

def encrypt(n, password):
    if password:
        subject = n['content'].split('\n')[0]
        n['content'] = '\n'.join((subject,
                                  crypt.encrypt(n['content'], password)))
    return n

def update_note(n, password=None):
    c = conn()
    n = encrypt(n, password)
    sql = 'UPDATE notes SET tags = ?, content = ?, updated = ? WHERE uid = ?'
    values = [n.get('tags'), n.get('content'), n.get('updated'), n.get('uid')]
    log.debug('%s, %s' % (sql, values))
    c.execute(sql, values)
    return c.commit()

def conn():
    d = os.path.dirname(path)
    _ = os.makedirs(d) if not os.path.exists(d) else None
    return sqlite3.connect(path)

def columns(n, row):
    # TODO: This can be deleted
    for k, _ in n.items():
        yield row.get(k)

def dict_factory(c, row):
    return dict((col[0], row[i]) for i, col in enumerate(c.description))

def fields(n):
    for k, v in n.items():
        yield dict(id=k, label=k.capitalize(), type=v)

def rows(n, rs):
    for row in rs:
        yield list(columns(n, row))

def create_schema(c):
    pairs = [' '.join(pair) for pair in note().items() if pair[1]]
    sql ='CREATE TABLE notes (%s);' % ','.join(pairs)
    try:
        c.execute(sql)
    except sqlite3.OperationalError as ex:
        log.debug(ex)

def get_content(uid, password):
    c = conn()
    c.row_factory = dict_factory
    sql = 'SELECT content FROM notes WHERE uid = ?'
    content = c.cursor().execute(sql, [uid]).next().get('content')
    return crypt.decrypt(content.split('\n', 1)[1], password)

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
        encrypted = 0
        content = n['content'].split('\n')
        if len(content) > 1 and re.search(r'([^ ]{32,})', content[1]):
            encrypted = 1
        c.execute(sql, (content.pop(0), '\n'.join(content), encrypted, n['uid']))
    c.commit()

def migrate_schema(c):
    sqls =['ALTER TABLE notes ADD COLUMN encrypted INTEGER DEFAULT 0;',
           ' ALTER TABLE notes ADD COLUMN subject TEXT;']
    try:
        for sql in sqls:
            c.execute(sql)
    except sqlite3.OperationalError as ex:
        log.debug('schema migration: %s', ex)

def search(s):
    terms = s.split() if s else []
    n = note()
    naive = "(content LIKE '%{0}%' OR tags LIKE '%{0}%')"
    where = ['1=1'] + [naive.format(t) for t in terms]
    sql = 'SELECT %s FROM notes WHERE %s ORDER BY updated DESC;'
    sql = sql % (','.join(k for k, v in n.items() if v), ' AND '.join(where))
    log.debug(sql)
    c = conn()
    c.row_factory = dict_factory
    cols = list(fields(n))
    for row in rows(n, c.cursor().execute(sql)):
        n = {}
        for i, col in enumerate(cols):
            n[cols[i]['id']] = row[i]
        yield n

def prepare():
    create_schema(conn())
    migrate_schema(conn())
    migrate_data(conn())
