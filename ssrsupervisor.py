import pickle
import os
from flask import Flask, request
from ssr import SSR, SSRParam
from subscriber import ssr_params_from_subscription_url
from checker import Checker
from utils import log


CONFIG_FILENAME = 'config.dat'


class SSRSupervisor:
    ssr_params: [SSRParam]

    def __init__(self, debug=False):
        self.app = Flask(__name__)
        self.debug = debug
        self.route()
        self.ssr = SSR(self)
        self.subscription_url = ''
        self.ssr_params = []
        self.listen = ''
        self.listen_port = 0
        self.keyword = '香港'
        self.read_config()
        self.checker = Checker(self)

    def get_app(self):
        return self.app

    def route(self):
        self.app.add_url_rule('/', view_func=self.index, methods=['GET'])
        self.app.add_url_rule('/next', view_func=self.next, methods=['POST'])
        self.app.add_url_rule('/subscription', view_func=self.post_subscription, methods=['POST'])

    def index(self):
        return 'hello'

    def next(self):
        self.checker.push_ssr_param()
        return '', 204

    def post_subscription(self):
        self.subscription_url = request.json['url']
        self.ssr_params = ssr_params_from_subscription_url(self.subscription_url)
        self.save_config()
        self.checker.reload()
        return '', 204

    def read_config(self):
        if not os.path.exists(CONFIG_FILENAME):
            return
        with open(CONFIG_FILENAME, 'rb') as f:
            config = pickle.load(f)
        self.subscription_url = config.get('subscription_url')
        self.ssr_params = config.get('ssr_params')
        self.listen = config.get('listen')
        self.listen_port = config.get('listen_port')

    def save_config(self):
        with open(CONFIG_FILENAME, 'wb') as f:
            config = dict(
                subscription_url=self.subscription_url,
                ssr_params=self.ssr_params,
                listen=self.listen,
                listen_port=self.listen_port,
            )
            f.write(pickle.dumps(config))

    def run(self):
        self.checker.run()
        self.app.run(debug=self.debug)


def main(debug=False):
    ssr_supervisor = SSRSupervisor(debug)
    ssr_supervisor.run()


if __name__ == '__main__':
    main(debug=False)
