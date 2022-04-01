from dataclasses import dataclass
import subprocess
from utils import log


@dataclass
class SSRParam:
    host: str
    port: int
    password: str
    method: str
    protocol: str
    obfs: str
    remarks: str = ''
    group: str = ''
    listen: str = '127.0.0.1'
    listen_port: int = 1080
    obfs_param: str = ''
    proto_param: str = ''
    udp_port: int = 0
    uot: bool = False
    ttl: int = 0


class SSR:
    p: subprocess.Popen
    ssr_param: SSRParam

    def __init__(self, ssr_supervisor):
        self.p = subprocess.Popen(['echo', '-n'])
        self.ssr_supervisor = ssr_supervisor
        self.ssr_param = None

    def start(self):
        c = self.ssr_param
        log(c)
        self.p = subprocess.Popen(['python3', 'shadowsocksr/shadowsocks/local.py',
                                   '-s', c.host,
                                   '-p', str(c.port),
                                   '-k', c.password,
                                   '-m', c.method,
                                   '-O', c.protocol,
                                   '-o', c.obfs,
                                   '-G', c.proto_param,
                                   '-g', c.obfs_param,
                                   '-b', c.listen,
                                   '-l', str(c.listen_port),
                                   ])

    def run(self):
        if self.ssr_param:
            self.start()

    def load(self, ssr_param: SSRParam):
        self.ssr_param = ssr_param
        self.ssr_param.listen = self.ssr_supervisor.listen
        self.ssr_param.listen_port = self.ssr_supervisor.listen_port
        if not self.p.poll():
            self.p.terminate()
            self.p.wait()
        self.run()
