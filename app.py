# Python imports
import logging
import os

# Third party imports
from bottle import route, run, static_file

# Project imports
import crypt
import db

# Constants
root = os.path.join(os.path.dirname(__file__), 'static')

@route('/hello/:name')
def index(name='World'):
    return '<b>Hello %s!</b>' % name

@route('/')
def homepage():
    return static_file('index.html', root=root)

@route('/static/<filename>')
def static(filename):
    return static_file(filename, root=root)

if __name__ == '__main__':
    logging.basicConfig()
    db.prepare()
    run(host='localhost', port=8082)

