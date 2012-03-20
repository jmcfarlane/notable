# Python imports
from collections import OrderedDict
import datetime
import logging
import os
import sqlite3

# Constants
path = os.path.expanduser('~/.notes/notes.sqlite3')
log = logging.getLogger(__name__)

def note(exclude=None, actual=False):
    exclude = exclude or []
    fields = ['created', 'updated', 'tags', 'subject', 'content', 'password']
    n = OrderedDict((k,'string') for k in fields if not k in exclude)
    if actual:
        now = datetime.datetime.now()
        n.update(created=now, updated=now)
    return n

def create_note(n):
    c = conn()
    sql = 'INSERT INTO notes VALUES (%s);' % ','.join('?' for _ in n.keys())
    c.execute(sql, n.values())
    c.commit()
    return True

def conn():
    d = os.path.dirname(path)
    _ = os.makedirs(d) if not os.path.exists(d) else None
    return sqlite3.connect(path)

def columns(n):
    for k, v in n.items():
        yield dict(id=k, label=k.capitalize(), type=v)

def rows(rs):
    for row in rs:
        yield dict(c=[ dict(v=d) for d in row ])

def create_schema(c):
    pairs = [' '.join(pair) for pair in note().items()]
    sql ='CREATE TABLE notes (%s);' % ','.join(pairs)
    try:
        c.execute(sql)
    except sqlite3.OperationalError, ex:
        log.warning(ex)

def search(s):
    terms = s.split() if s else []
    n = note(exclude=['created', 'password'])
    where = ['1=1'] + ["content LIKE '%{0}%'".format(t) for t in terms]
    sql = 'SELECT %s FROM notes WHERE %s;'
    sql = sql % (','.join(n.keys()), ' AND '.join(where))
    log.warn(sql)
    return dict(cols=list(columns(n)),
                rows=list(rows(conn().cursor().execute(sql))))

def prepare():
    c = conn()
    create_schema(c)
    #create_note(note(actual=True))
