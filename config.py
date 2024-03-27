import os
import pickle

import otp
from otp import log
import subscriber


class ConfigMethod:
    all_config = 'all_config'
    listen_host = 'listen_host'
    listen_port = 'listen_port'
    set_listen_host = 'set_listen_host'
    set_listen_port = 'set_listen_port'
    subscription_url = 'subscription_url'
    set_subscription_url = 'set_subscription_url'
    ssr_params = 'ssr_params'
    set_ssr_params = 'set_ssr_params'
    keyword = 'keyword'
    set_keyword = 'set_keyword'
    auto_mode = 'auto_mode'
    set_auto_mode = 'set_auto_mode'
    update_subscription = 'update_subscription'
    current_ssr_param = 'current_ssr_param'
    set_current_ssr_param = 'set_current_ssr_param'


class Config(otp.Service):
    name: str = 'config'
    methods = ConfigMethod
    filename = 'config.dat'

    def __init__(self):
        super(Config, self).__init__()

        methods = Config.methods
        self.bind(methods.listen_host)
        self.bind(methods.listen_port)
        self.bind(methods.subscription_url)
        self.bind(methods.ssr_params)
        self.bind(methods.keyword)
        self.bind(methods.auto_mode)
        self.bind(methods.current_ssr_param)
        self.handle_map = {
            **self.handle_map,
            methods.set_listen_host: self.set_listen_host,
            methods.set_listen_port: self.set_listen_port,
            methods.set_subscription_url: self.set_subscription_url,
            methods.set_ssr_params: self.set_ssr_params,
            methods.set_keyword: self.set_keyword,
            methods.set_auto_mode: self.set_auto_mode,
            methods.set_current_ssr_param: self.set_current_ssr_param,
            methods.all_config: self.all_config,
            methods.update_subscription: self.update_subscription,
        }

        self.read_config_file()

    def read_config_file(self):
        if not os.path.exists(Config.filename):
            self.states['subscription_url'] = ''
            self.states['ssr_params'] = []
            self.states['listen_host'] = '127.0.0.1'
            self.states['listen_port'] = 1080
            self.states['keyword'] = '香港'
            self.states['auto_mode'] = False
            self.states['current_ssr_param'] = None
            return

        with open(Config.filename, 'rb') as f:
            config = pickle.load(f)
        # self.states['subscription_url'] = config.get('subscription_url')
        # self.states['ssr_params'] = config.get('ssr_params')
        # self.states['listen_host'] = config.get('listen_host')
        # self.states['listen_port'] = config.get('listen_port')
        # self.states['keyword'] = config.get('keyword')
        # self.states['auto_mode'] = config.get('auto_mode')
        # self.states['current_ssr_param'] = config.get('current_ssr_param')
        self.states = config
        log('current config:', self.states)

    def save_config_file(self):
        with open(Config.filename, 'wb') as f:
            config = self.states
            f.write(pickle.dumps(config))

    def set_listen_host(self, host):
        log('set listen:', host)
        self.states['listen_host'] = host
        self.save_config_file()

    def set_listen_port(self, port):
        log('set listen port:', port)
        self.states['listen_port'] = port
        self.save_config_file()

    def set_subscription_url(self, url):
        log('set subscription url:', url)
        self.states['subscription_url'] = url
        self.save_config_file()

    def set_ssr_params(self, ssr_params):
        self.states['ssr_params'] = ssr_params
        self.save_config_file()

    def set_keyword(self, keyword):
        self.states['keyword'] = keyword
        self.save_config_file()

    def set_auto_mode(self, mode):
        log('set_auto_mode', mode)
        self.states['auto_mode'] = mode
        self.save_config_file()

    def set_current_ssr_param(self, ssr_param):
        log('set_current_ssr_param', ssr_param)
        self.states['current_ssr_param'] = ssr_param
        self.save_config_file()

    def all_config(self, q):
        q.put(self.states)

    def update_subscription(self):
        url = self.states[ConfigMethod.subscription_url]
        log('update sub:', url)
        ssr_params = subscriber.ssr_params_from_subscription_url(url)
        self.states[ConfigMethod.ssr_params] = ssr_params
        log('updated')
