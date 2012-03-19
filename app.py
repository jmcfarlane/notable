# Python imports
import logging
import os

# Third party imports
from bottle import route, run, static_file, view

# Project imports
import db

# Constants
root = os.path.join(os.path.dirname(__file__), 'static')

@route('/')
@view('index')
def homepage():
    return dict()

@route('/static/<filename>')
def static(filename):
    return static_file(filename, root=root)

@route('/api/list/<tags:re:.*>')
def api_list(tags=None):
    return db.fetch_by_tags(tags)

if __name__ == '__main__':
    logging.basicConfig()
    db.prepare()
    run(host='localhost', port=8082)

