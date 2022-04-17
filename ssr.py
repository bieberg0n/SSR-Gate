import subprocess
import otp
from otp import log
from config import Config
from ssrparam import SSRParam


class SSRMethod:
    param = 'param'
    set_param = 'set-param'
    restart = 'restart'


class SSR(otp.Service):
    p: subprocess.Popen
    ssr_param: SSRParam
    name: str = 'ssr'
    methods = SSRMethod

    def __init__(self):
        super(SSR, self).__init__()
        self.p = subprocess.Popen(['echo', '-n'])

        methods = SSR.methods
        self.bind(methods.param)
        self.handle_map[methods.set_param] = self.set_param
        self.handle_map[methods.restart] = self.restart

    def start_ssr(self):
        c = self.states.get('param')
        c.listen = Config.get(Config.methods.listen_host)
        c.listen_port = Config.get(Config.methods.listen_port)
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

    def restart(self):
        if not self.p.poll():
            self.p.terminate()
            self.p.wait()
        self.start_ssr()

    def set_param(self, ssr_param: SSRParam):
        self.states['param'] = ssr_param
        self.restart()
