from flask import Flask
from ssr import SSR


def index():
    return 'hello'


class SSRSupervisor:
    def __init__(self, debug=False):
        self.app = Flask(__name__)
        self.debug = debug
        self.route()
        self.ssr = SSR()

    def get_app(self):
        return self.app
    
    def route(self):
        self.app.add_url_rule('/', view_func=index)
        
    def run(self):
        self.app.run(debug=self.debug)


def main(debug=False):
    ssr_supervisor = SSRSupervisor(debug)
    ssr_supervisor.run()


if __name__ == '__main__':
    main(debug=True)
