# Python imports
import logging
import os
import sys
import threading
import time
import webbrowser

root = os.path.abspath(os.path.dirname(__file__))
sys.path = [os.path.join(root, '..')] + sys.path

# Project imports
from notable import bottle, db, editor

# Constants, and help with template lookup
version = '0.0.2a'
static = os.path.join(root, 'static')
bottle.TEMPLATE_PATH.insert(0, os.path.join(root, 'views'))

@bottle.route('/')
@bottle.view('index')
def homepage():
    return dict()

@bottle.route('/static/<filename>')
def htdocs(filename):
    return bottle.static_file(filename, root=static)

@bottle.post('/api/decrypt')
def decrypt():
    password = bottle.request.forms.get('password')
    uid = bottle.request.forms.get('uid')
    return db.get_content(uid, password)

@bottle.post('/api/delete')
def delete():
    password = bottle.request.forms.get('password')
    uid = bottle.request.forms.get('uid')
    return dict(success=db.delete_note(uid, password=password))

@bottle.post('/api/launch_editor')
def launch_editor():
    uid = bottle.request.forms.get('uid')
    content = bottle.request.forms.get('content')
    return editor.launch(uid, content)

@bottle.get('/api/from_disk/<uid>')
def from_disk(uid=None):
    path = os.path.join('/tmp', uid)
    if os.path.exists(path):
        return open(path).read()
    else:
        return 'missing'

@bottle.get('/api/list')
def api_list():
    return db.search(bottle.request.query.get('s'))

@bottle.post('/api/persist')
def persist():
    n = db.note(actual=True)
    password = bottle.request.forms.get('password')
    form = dict((k, v) for k, v in bottle.request.forms.items() if k in n)
    if form.get('uid') == '':
        fcn = db.create_note
        form.pop('uid')
    else:
        fcn = db.update_note

    n.update(form)
    return dict(success=fcn(n, password=password))

def main():
    logging.basicConfig()
    db.prepare()
    #threading.Thread(target=launcher).start()
    bottle.run(host='localhost', port=8082)

def launcher():
    time.sleep(1)
    webbrowser.open_new_tab('http://localhost:8082')

if __name__ == '__main__':
    main()
