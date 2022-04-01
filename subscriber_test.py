import pytest
from ssr import SSRParam
from subscriber import ssr_params_from_subscription_url, ssr_param_from_ssr_url


def test_ssr_param_from_subscribe_url():
    url = 'ssr://MTI3LjAuMC4xOjgwODA6YXV0aF9hZXMxMjhfbWQ1OmFlcy0xMjgtY3RyOnBsYWluOmRHVnpkREl3TWpBLz9vYmZzcGFyYW09JnByb3RvcGFyYW09TXpBd09ERTZOREJ3WkVWd1REUTFkUSZyZW1hcmtzPTU0R3I1cGlmSURBeCZncm91cD1kR1Z6ZEE'
    ssr_param: SSRParam = ssr_param_from_ssr_url(url)
    assert ssr_param.host == '127.0.0.1'
    assert ssr_param.password == 'test2020'
    assert ssr_param.port == 8080
    assert ssr_param.protocol == 'auth_aes128_md5'
    assert ssr_param.method == 'aes-128-ctr'
    assert ssr_param.obfs == 'plain'
    assert ssr_param.obfs_param == ''
    assert ssr_param.proto_param == '30081:40pdEpL45u'
    assert ssr_param.remarks == '火星 01'
    assert ssr_param.group == 'test'


def test_ssr_param_from_url():
    ...
