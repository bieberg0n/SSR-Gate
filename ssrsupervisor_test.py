def test_index(client):
    rv = client.get('/')
    assert b'hello' in rv.data


def test_subscription(client):
    ...
