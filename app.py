#! /usr/bin/env python

# Python imports
import logging
import os

# Third party imports
from bottle import get, post, request, route, run, static_file, view

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

@post('/api/decrypt')
def decrypt():
    password = request.forms.get('password')
    uid = request.forms.get('uid')
    return db.get_content(uid, password)

@get('/api/list')
def api_list():
    return db.search(request.query.get('s'))

@post('/api/persist')
def persist():
    n = db.note(actual=True)
    password = request.forms.get('password')
    form = dict((k, v) for k, v in request.forms.items() if k in n)
    if form.get('uid') == '':
        fcn = db.create_note
        form.pop('uid')
    else:
        fcn = db.update_note

    n.update(form)
    return dict(success=fcn(n, password=password))

if __name__ == '__main__':
    logging.basicConfig()
    db.prepare()
    run(host='localhost', port=8082, reloader=True)

