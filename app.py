# Python imports
import logging
import os

# Third party imports
from bottle import post, request, route, run, static_file, view

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

@post('/api/persist')
def persist(tags=None):
    n = db.note(actual=True)
    form = dict((k, v) for k, v in request.forms.items() if k in n)
    n.update(form)
    return dict(success=db.create_note(n))

if __name__ == '__main__':
    logging.basicConfig()
    db.prepare()
    run(host='localhost', port=8082)

