from typing import Tuple
import time
import socket
from queue import Queue
import requests
from ssr import SSR, SSRParam
from utils import log, spawn


class Checker:
    def __init__(self, ssr_supervisor):
        self.ssr_supervisor = ssr_supervisor
        self.ssr_params_standby = Queue()

    def reload(self):
        self.ssr_params_standby = Queue()
        self.push_ssr_param()
    
    def ping_once(self, addr: Tuple[str, int]):
        old_time = time.time() * 1000
        s = socket.socket()
        s.settimeout(3)
        try:
            s.connect(addr)
        except Exception as e:
            log(e)
            return -1

        ttl = int(time.time() * 1000 - old_time)
        return ttl

    def ping(self, ssr_param: SSRParam):
        keyword = self.ssr_supervisor.keyword
        if keyword not in ssr_param.remarks:
            log(ssr_param.remarks, 'BAN')
            ssr_param.ttl = -1
            
        else:
            addr = (ssr_param.host, ssr_param.port)
            ttls = [self.ping_once(addr) for _ in range(3)]
            min_ttl = min(ttls)
            ssr_param.ttl = min_ttl
            log(ssr_param.remarks, 'TCP TTL:', ssr_param.ttl)

    def ping_all(self):
        for ssr_param in self.ssr_supervisor.ssr_params:
            self.ping(ssr_param)
    
    def http_ping_once(self):
        old_time = time.time() * 1000
        listen_port = self.ssr_supervisor.listen_port
        try:
            r = requests.head('http://google.com',
                              proxies=dict(http=f'socks5h://127.0.0.1:{listen_port}',
                                           https=f'socks5h://127.0.0.1:{listen_port}'),
                              timeout=3)
        except Exception as e:
            log(e)
            return -1

        ttl = int(time.time() * 1000 - old_time)
        return ttl
    
    def http_ping(self):
        for _ in range(3):
            ttl = self.http_ping_once()
            if ttl != -1:
                log('HTTP TTL:', ttl)
                return True
        else:
            return False
    
    def push_ssr_param(self):
        if self.ssr_params_standby.empty():
            self.ping_all()
            ssr_params = sorted(self.ssr_supervisor.ssr_params, key=lambda p: p.ttl)
            for ssr_param in ssr_params:
                if ssr_param.ttl != -1:
                    self.ssr_params_standby.put(ssr_param)

        ssr_param = self.ssr_params_standby.get()
        self.ssr_supervisor.ssr.load(ssr_param)
        log('use:', ssr_param.remarks)
        time.sleep(1)
        self.check()

    def check(self):
        current_ssr_param = self.ssr_supervisor.ssr.ssr_param
        if current_ssr_param:
            self.ping(current_ssr_param)
            if current_ssr_param.ttl == -1 or not self.http_ping():
                self.push_ssr_param()
        else:
            self.push_ssr_param()
    
    def check_loop(self):
        while True:
            self.check()
            time.sleep(20)

    def run(self):
        spawn(target=self.check_loop)
