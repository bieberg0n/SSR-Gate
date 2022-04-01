from urllib.parse import parse_qs
from pprint import pprint as log
import requests
from ssr import SSRParam
from utils import b64decode


def ssr_param_from_ssr_url(ssr_url):
    url = ssr_url
    if not url.startswith('ssr://'):
        return None
    
    b64content = url[6:]
    b64content += (len(b64content) % 4) * '='
    content = b64decode(b64content)
    
    parts = content.split('/?')
    left_parts = parts[0].split(':')
    right_parts = parse_qs(content)

    ssr_param = SSRParam(
        host=left_parts[0],
        port=int(left_parts[1]),
        protocol=left_parts[2],
        method=left_parts[3],
        obfs=left_parts[4],
        password=b64decode(left_parts[5]),
    )

    for key, values in right_parts.items():
        value = values[0]
        if key == "obfsparam":
            ssr_param.obfs_param = b64decode(value)
        elif key == 'protoparam':
            ssr_param.proto_param = b64decode(value)
        elif key == 'remarks':
            ssr_param.remarks = b64decode(value)
        elif key == 'group':
            ssr_param.group = b64decode(value)
        elif key == 'udpport':
            ssr_param.udp_port = int(value)
        elif key == 'uot':
            ssr_param.uot = value != '0'

    return ssr_param


def ssr_params_from_subscription_url(url):
    r = requests.get(url)
    data = r.text
    ssr_urls = b64decode(data).split('\n')
    ssr_urls = [url for url in ssr_urls if url]
    ssr_params = [ssr_param_from_ssr_url(ssr_url) for ssr_url in ssr_urls]
    return [i for i in ssr_params if i]
