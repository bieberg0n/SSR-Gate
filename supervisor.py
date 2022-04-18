from ssr import SSR
from config import Config
from checker import Checker
from app import WebServer
import otp


def start_supervisor():
    supervisor = otp.Supervisor()
    supervisor.start([
        Config,
        SSR,
        Checker,
        WebServer,
    ])


if __name__ == '__main__':
    start_supervisor()
