from dataclasses import dataclass
import subprocess


@dataclass
class SSRParam:
    host: str
    port: int
    password: str
    method: str
    protocol: str
    obfs: str
    obfs_param: str
    proto_param: str
    listen: str
    listen_port: int
    remarks: str = ''
    group: str = ''
    udp_port: int = 0
    uot: bool = False


class SSR:
    p: subprocess.Popen
    ssr_param: SSRParam

    def __init__(self, ssr_param=None):
        self.ssr_param = ssr_param

    def start(self):
        c = self.ssr_param
        self.p = subprocess.Popen(['python3', 'shadowsocksr/shadowsocks/local.py', 
                                   '-s', c.host,
                                   '-p', c.port,
                                   '-k', c.password,
                                   '-m', c.method,
                                   '-O', c.protocol,
                                   '-o', c.obfs,
                                   '-G', c.proto_param,
                                   '-g', c.obfs_param,
                                   '-b', c.listen,
                                   '-l', c.listen_port,
                                   ])

    def run(self):
        if self.ssr_param:
            self.start()

    def load(self, ssr_param: SSRParam):
        self.ssr_param = ssr_param
        self.p.terminate()
        self.p.wait()
        self.run()
