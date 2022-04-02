from dataclasses import dataclass


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
