# Python imports
import os
import subprocess
import threading
import time
import uuid

def launcher(path):
    subprocess.call(['gvim', '--nofork', path])
    time.sleep(2.5)
    os.remove(path)

def launch(uid, content):
    uid = uid or uuid.uuid4().hex
    path = os.path.join('/tmp', uid)
    with open(path, 'w') as fh:
        fh.write(content)
    threading.Thread(args=(path,), target=launcher).start()
    return uid
