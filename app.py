from flask import Flask, request, jsonify, send_from_directory
from flask_cors import CORS

from config import Config
from checker import Checker
from ssr import SSR
import otp
from otp import log


class WebServer(otp.Service):
    name = 'webserver'

    def __init__(self):
        super(WebServer, self).__init__()
        self.app = Flask(__name__)
        CORS(self.app)
        self.route()

    def get_app(self):
        return self.app

    def index(self):
        return send_from_directory('static/', 'index.html')

    def health(self):
        return 'SSR-GATE'

    def status(self):
        current_ssr_param = SSR.get(SSR.methods.param)
        status = dict(current_ssr_param=current_ssr_param)
        all_config = Config.get(Config.methods.all_config)
        return jsonify({**status, **all_config})

    def next(self):
        Checker.emit(Checker.methods.next)
        return ''

    def simulator(self):
        return 'c3NyOi8vTVRJM0xqQXVNQzR4T2pnd09EQTZZWFYwYUY5aFpYTXhNamhmYldRMU9tRmxjeTB4TWpndFkzUnlPbkJzWVdsdU9tUkhWbnBrUkVsM1RXcEJMejl2WW1aemNHRnlZVzA5Sm5CeWIzUnZjR0Z5WVcwOVRYcEJkMDlFUlRaT1JFSjNXa1ZXZDFSRVVURmtVU1p5WlcxaGNtdHpQVFUwUjNJMWNHbG1TVVJCZUNabmNtOTFjRDFrUjFaNlpFRQ'

    def post_subscription(self):
        url = request.json['url']
        Config.emit(Config.methods.set_subscription_url, url)
        Checker.emit(Checker.methods.reload)
        return jsonify(dict(url=url))

    def route(self):
        self.app.add_url_rule('/', methods=['GET'], view_func=self.index)
        self.app.add_url_rule('/api/health', methods=['GET'], view_func=self.health)
        self.app.add_url_rule('/api/status', methods=['GET'], view_func=self.status)
        self.app.add_url_rule('/api/next', methods=['POST'], view_func=self.next)
        self.app.add_url_rule('/api/subscription', methods=['POST'], view_func=self.post_subscription)
        self.app.add_url_rule('/api/simulator', methods=['GET'], view_func=self.simulator)

    def run_app(self):
        self.app.run(host='0.0.0.0')

    def run(self):
        otp.spawn(target=self.run_app)
        super(WebServer, self).run()
