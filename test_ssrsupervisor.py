def test_a(client):
    rv = client.get('/')
    print(rv.data)
    assert b'hello' in rv.data
