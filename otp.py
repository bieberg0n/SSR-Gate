import queue
from dataclasses import dataclass
from typing import Any
from threading import Thread


def log(*args):
    print(*args)


def spawn(target, args=(), daemon: bool = True):
    t = Thread(target=target, args=args)
    t.setDaemon(daemon)
    t.start()
    return t


class Service:
    services = {}
    name = ''

    def __init__(self):
        self.queue = queue.Queue()
        self.handle_map = {}
        self.services = {}
        self.states = {}

    def bind(self, state: str):
        self.handle_map[state] = lambda q: q.put(self.states.get(state))
        def set(args):
            self.states[state] = args
        self.handle_map['set-' + state] = set

    @classmethod
    def get(cls, method: str):
        q = queue.Queue()
        cls.services.get(cls.name).put((method, q))
        return q.get()

    @classmethod
    def emit(cls, method: str, *args):
        cls.services.get(cls.name).put((method, *args))

    def run(self):
        while True:
            (method, *args) = self.queue.get()
            handle = self.handle_map.get(method)
            if handle:
                handle(*args)
            else:
                log(f'[{self.name}] unknown method: {method}')

    def start(self):
        spawn(target=self.run, daemon=False)


class Supervisor:
    def __init__(self):
        self.services = {}

    def start(self, services: [Service]):
        for service in services:
            log(f'start service: [{service.name}]')
            service.services = self.services
            s: Service = service()
            self.services[s.name] = s.queue
            s.start()
