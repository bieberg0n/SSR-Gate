import base64
from threading import Thread


def b64decode(s):
    s += len(s) % 4 * '='
    s = s.replace('-', '+').replace('_', '/')
    return base64.b64decode(s).decode()


def log(*args):
    print(*args)

    
def spawn(target, args=(), daemon: bool = True):
    t = Thread(target=target, args=args)
    t.setDaemon(daemon)
    t.start()
