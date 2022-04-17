from typing import Tuple
import time
import socket
from queue import Queue
import requests

import otp
from otp import log, spawn
import subscriber
from ssr import SSR, SSRParam
from config import Config


class CheckerMethod:
    check = 'check'
    next = 'next'
    reload = 'reload'


class Checker(otp.Service):
    name: str = 'checker'
    methods = CheckerMethod

    def __init__(self):
        super(Checker, self).__init__()
        self.ssr_params_standby = Queue()

        methods = Checker.methods
        self.handle_map[methods.check] = self.check
        self.handle_map[methods.next] = self.push_ssr_param
        self.handle_map[methods.reload] = self.reload

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
        keyword = Config.get(Config.methods.keyword)
        auto_mode_flag = Config.get(Config.methods.auto_mode)
        if auto_mode_flag and keyword and keyword not in ssr_param.remarks:
            log(ssr_param.remarks, 'BAN')
            ssr_param.ttl = -1

        else:
            addr = (ssr_param.host, ssr_param.port)
            ttls = [self.ping_once(addr) for _ in range(3)]
            min_ttl = min(ttls)
            ssr_param.ttl = min_ttl
            log(ssr_param.remarks, 'TCP TTL:', ssr_param.ttl)

    def ping_all(self, ssr_params):
        for ssr_param in ssr_params:
            self.ping(ssr_param)
        return ssr_params
    
    def http_ping_once(self):
        old_time = time.time() * 1000
        listen_host = Config.get(Config.methods.listen_host)
        listen_port = Config.get(Config.methods.listen_port)
        try:
            r = requests.head('http://google.com',
                              proxies=dict(http=f'socks5h://{listen_host}:{listen_port}',
                                           https=f'socks5h://{listen_host}:{listen_port}'),
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
        is_auto_mode = Config.get(Config.methods.auto_mode)
        if not is_auto_mode:
            return

        elif self.ssr_params_standby.empty():
            url = Config.get(Config.methods.subscription_url)
            ssr_params = subscriber.ssr_params_from_subscription_url(url)
            ssr_params = self.ping_all(ssr_params)
            ssr_params = sorted(ssr_params, key=lambda p: p.ttl)
            Config.emit(Config.methods.set_ssr_params, ssr_params)

            for ssr_param in ssr_params:
                if ssr_param.ttl != -1:
                    self.ssr_params_standby.put(ssr_param)

        ssr_param = self.ssr_params_standby.get()
        SSR.emit(SSR.methods.set_param, ssr_param)
        time.sleep(1)
        Checker.emit(Checker.methods.check)

    def check(self):
        current_ssr_param = SSR.get(SSR.methods.param)
        if current_ssr_param:
            self.ping(current_ssr_param)
            if current_ssr_param.ttl == -1 or not self.http_ping():
                self.push_ssr_param()
        else:
            self.push_ssr_param()
    
    def check_loop(self):
        while True:
            Checker.emit(Checker.methods.check)
            time.sleep(20)

    def run(self):
        spawn(target=self.check_loop)
        super(Checker, self).run()
