# Python imports
from collections import OrderedDict
import logging
import os
import sqlite3

# Constants
path = os.path.expanduser('~/.notes/notes.sqlite3')
log = logging.getLogger(__name__)

def note(exclude=None):
    exclude = exclude or []
    fields = ['created', 'updated', 'tags', 'subject', 'content', 'password']
    return OrderedDict((k,'string') for k in fields if not k in exclude)

def add_sample_notes(c):
    d = ('today', 'now', 'food water', 'subject', None, 'Content of note')
    sql = 'INSERT INTO notes VALUES (%s);' % ','.join('?' for _ in d)
    c.execute(sql, d)
    c.commit()

def conn():
    d = os.path.dirname(path)
    _ = os.makedirs(d) if not os.path.exists(d) else None
    return sqlite3.connect(path)

def columns(note):
    for k, v in note.items():
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

def fetch_by_tags(tags):
    n = note(exclude=['password'])
    sql = 'SELECT %s FROM notes;' % ','.join(n.keys())
    return dict(cols=list(columns(n)),
                rows=list(rows(conn().cursor().execute(sql))))

def prepare():
    c = conn()
    create_schema(c)
    add_sample_notes(c)
