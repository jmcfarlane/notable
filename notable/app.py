"""Notable - a simple not taking application"""

# Python imports
from time import gmtime, sleep, strftime, time
import json
import logging
import optparse
import os
import re
import signal
import sys
import threading
import webbrowser

# Assume python3 and fall back to python2
try:
    import http.client as http
    from urllib.request import urlopen
    from urllib.error import URLError
except ImportError:
    import httplib as http
    from urllib2 import urlopen, URLError

root = os.path.abspath(os.path.dirname(__file__))
sys.path = [os.path.join(root, '..')] + sys.path

# Project imports
from notable import bottle, db

# Constants, and help with template lookup
host = 'localhost'
version = '0.3.0'
static = os.path.join(root, 'static')
bottle.TEMPLATE_PATH.insert(0, os.path.join(root, 'static/templates'))
log = logging.getLogger()

@bottle.route('/')
@bottle.view('index')
def homepage():
    return dict()

@bottle.route('/static/<filename:re:.*>')
def htdocs(filename):
    response = bottle.static_file(filename, root=static)
    expires = strftime('%a, %d %b %Y %H:%M:%S GMT', gmtime(time()))
    response.set_header('Expires', expires)
    return response

@bottle.post('/api/note/content/<uid>')
def content(uid):
    password = bottle.request.forms.get('password')
    content = db.get_content(uid, password)
    if smells_encrypted(content):
        bottle.response.status = 403
        return 'Nope, try again'
    return content

@bottle.post('/api/decrypt')
def decrypt():
    password = bottle.request.forms.get('password')
    uid = bottle.request.forms.get('uid')
    return db.get_content(uid, password)

@bottle.delete('/api/note/<uid>')
def delete(uid):
    return dict(success=db.delete_note(uid))

@bottle.get('/api/notes/list')
def listing():
    exclude = ['content']
    notes = db.search(bottle.request.query.get('s'), exclude=exclude)
    return json.dumps(list(notes), indent=2)

@bottle.put('/api/note/<uid>')
def update_note(uid):
    return _persist(db.update_note)

@bottle.post('/api/note/create')
def create_note():
    return _persist(db.create_note)

@bottle.get('/pid')
def getpid():
    return str(os.getpid())

def _persist(func):
    n = db.note(actual=True)
    note = json.loads(bottle.request.body.getvalue().decode())
    password = note.pop('password') if 'password' in note else ''
    n.update(dict((k, v) for k, v in list(note.items()) if k in n))
    return func(n, password=password)

def browser(opts):
    sleep(1)
    webbrowser.open_new_tab('http://localhost:%s' % opts.port)

def fork_and_exit():
    pid = os.fork() # Children return 0 here, parents > 0
    if pid > 0:
        log.debug('Forked, new pid is: %s', pid)
        sys.exit(0)
    return pid

def fork(opts):
    pid = running(opts)
    if pid:
        if opts.restart:
            stop(opts, pid)
        else:
            log.debug('Already running, noop (pid:%s)', pid)
            return 0
    fork_and_exit()
    os.setsid()
    os.umask(0)
    fork_and_exit()

def getopts():
    parser = optparse.OptionParser(__doc__.strip())
    parser.add_option('-d', '--debug',
                      action='store_true',
                      help='Debug using a debug db')
    parser.add_option('-f', '--foreground',
                      action='store_true',
                      help='Start the server in the foreground')
    parser.add_option('-n', '--no-browser',
                      action='store_true',
                      help='Do not launch a browser')
    parser.add_option('-l', '--log',
                      default=os.path.expanduser('~/.notable/notable.log'),
                      help='Log file to be used')
    parser.add_option('-p', '--port',
                      default=8082,
                      help='TCP port to start the server on')
    parser.add_option('-r', '--restart',
                      action='store_true',
                      help='Restart if already running')
    parser.add_option('-s', '--stop',
                      action='store_true',
                      help='Stop if already running, and exit')
    parser.add_option('-v', '--version',
                      action='store_true',
                      help='Print version number')
    return parser.parse_args(), parser

def running(opts):
    url = 'http://%s:%s/pid' % (host, opts.port)
    try:
        return int(urlopen(url).read())
    except (http.BadStatusLine, URLError, ValueError):
        return False

def run(opts):
    pid = running(opts)
    if pid and opts.restart:
        os.kill(pid, signal.SIGTERM)
    elif pid:
        return

    reloader = True if opts.debug and opts.foreground else False
    db.path = db.path + '.debug' if opts.debug else db.path
    db.prepare()
    bottle.run(host=host, port=opts.port, reloader=reloader)
    return 0

def setup_logging(opts):
    if not os.path.exists(os.path.dirname(opts.log)):
        os.makedirs(os.path.dirname(opts.log))
    handler = logging.FileHandler(opts.log)
    log.addHandler(handler)
    # Forcibly take over SimpleHTTPServer's writing to sys.std{out,err}
    for i, fd in enumerate(('stdout', 'stderr'), start=1):
        _log = '%s.%s' % (opts.log, fd)
        open(_log, 'w').close()
        os.dup2(os.open(_log, os.O_RDWR), i)

def smells_encrypted(content):
    if sys.version_info >= (3, 0):
        normal = float(len(re.findall(r'[\x00-\x7F]+', content)))
        special = float(len(re.findall(r'[^\x00-\x7F]+', content)))
    else:
        normal = float(len([s for s in content if ord(s) < 128]))
        special = float(len([s for s in content if ord(s) >= 128]))
    return False if not normal else special / normal > 0.40

def stop(opts, pid=None):
    pid = pid or running(opts)
    if pid:
        log.debug('Stopping instance (pid:%s)', pid)
        os.kill(pid, signal.SIGTERM)
    return 0

def main():
    logging.basicConfig(level=logging.DEBUG)
    (opts, _), _ = getopts()
    opts.port = int(opts.port) + 1 if opts.debug else int(opts.port)
    setup_logging(opts)
    if opts.version:
        print('Notable version %s' % version)
        return 0
    if opts.stop:
        return stop(opts)
    if not opts.no_browser:
        threading.Thread(target=browser, args=[opts]).start()
    if not opts.foreground:
        fork(opts)
    return run(opts)

if __name__ == '__main__':
    sys.exit(main())
