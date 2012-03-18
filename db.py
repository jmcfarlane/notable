# Python imports
import logging
import os
import sqlite3

# Constants
path = os.path.expanduser('~/.notes/notes.sqlite3')
log = logging.getLogger(__name__)

def conn():
    d = os.path.dirname(path)
    _ = os.makedirs(d) if not os.path.exists(d) else None
    return sqlite3.connect(path)

def create_schema(c):
    sql = """
    CREATE TABLE notes (
        created text,
        updated text,
        tags text,
        subject text,
        password text
    );
    """
    try:
        c.execute(sql)
    except sqlite3.OperationalError, ex:
        log.warning(ex)

def prepare():
    c = conn()
    create_schema(c)
