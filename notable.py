#! /usr/bin/env python

# Python imports
import logging
import os
import threading
import time
import webbrowser

# Third party imports
from bottle import get, post, request, route, run, static_file, view

# Project imports
import db, editor

# Constants, and help with template lookup
root = os.path.abspath(os.path.dirname(__file__))
static = os.path.join(root, 'static')
os.chdir(root)

@route('/')
@view('index')
def homepage():
    return dict()

@route('/static/<filename>')
def htdocs(filename):
    return static_file(filename, root=static)

@post('/api/decrypt')
def decrypt():
    password = request.forms.get('password')
    uid = request.forms.get('uid')
    return db.get_content(uid, password)

@post('/api/launch_editor')
def launch_editor():
    uid = request.forms.get('uid')
    content = request.forms.get('content')
    return editor.launch(uid, content)

@get('/api/from_disk/<uid>')
def from_disk(uid=None):
    path = os.path.join('/tmp', uid)
    if os.path.exists(path):
        return open(path).read()
    else:
        return 'missing'

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

def launcher():
    time.sleep(1)
    webbrowser.open_new_tab('http://localhost:8082')

if __name__ == '__main__':
    logging.basicConfig()
    db.prepare()
    threading.Thread(target=launcher).start()
    run(host='localhost', port=8082)

