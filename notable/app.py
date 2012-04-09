"""Notable - a simple not taking application"""

# Python imports
import httplib
import logging
import optparse
import os
import signal
import sys
import threading
import time
import webbrowser
try:
    import urllib2
except ImportError:
    from urllib import request as urllib2

root = os.path.abspath(os.path.dirname(__file__))
sys.path = [os.path.join(root, '..')] + sys.path

# Project imports
from notable import bottle, db, editor

# Constants, and help with template lookup
host = 'localhost'
version = '0.0.6'
static = os.path.join(root, 'static')
bottle.TEMPLATE_PATH.insert(0, os.path.join(root, 'views'))
log = logging.getLogger(__name__)

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
        fcn, _ = db.create_note, form.pop('uid')
    else:
        fcn = db.update_note

    n.update(form)
    return dict(success=fcn(n, password=password))

@bottle.get('/pid')
def getpid():
    return str(os.getpid())

def browser(opts):
    time.sleep(1)
    webbrowser.open_new_tab('http://localhost:%s' % opts.port)

def fork_and_exit():
    pid = os.fork()
    if pid > 0:
        sys.exit(0)
    return pid

def fork():
    fork_and_exit()
    os.setsid()
    os.umask(0)
    fork_and_exit()

def getopts():
    parser = optparse.OptionParser(__doc__.strip())
    parser.add_option('-b', '--browser',
                      action='store_true',
                      dest='browser',
                      help='Launch a browser')
    parser.add_option('-d', '--debug',
                      action='store_true',
                      dest='debug',
                      help='Debug using a debug db')
    parser.add_option('-f', '--fork',
                      action='store_true',
                      dest='fork',
                      help='Start the server in the background (fork)')
    parser.add_option('-p', '--port',
                      default=8082,
                      dest='port',
                      help='TCP port to start the server on')
    parser.add_option('-r', '--restart',
                      action='store_true',
                      dest='restart',
                      help='Restart if already running')
    return parser.parse_args(), parser

def running(opts):
    url = 'http://%s:%s/pid' % (host, opts.port)
    try:
        return int(urllib2.urlopen(url).read())
    except (httplib.BadStatusLine, urllib2.URLError, ValueError):
        return False

def run(opts):
    pid = running(opts)
    if pid and opts.restart:
        os.kill(pid, signal.SIGTERM)
    elif pid:
        return

    db.path = db.path + '.debug' if opts.debug else db.path
    db.prepare()
    bottle.run(host=host, port=opts.port)

def main():
    logging.basicConfig(level=logging.DEBUG)
    (opts, _), _ = getopts()
    opts.port = int(opts.port) + 1 if opts.debug else int(opts.port)
    if opts.browser:
        threading.Thread(target=browser, args=[opts]).start()
    if opts.fork:
        fork()
    run(opts)

if __name__ == '__main__':
    main()
