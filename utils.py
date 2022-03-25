import base64


def b64decode(s):
    s += len(s) % 4 * '='
    s = s.replace('-', '+').replace('_', '/')
    return base64.b64decode(s).decode()
