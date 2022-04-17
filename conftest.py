import pytest
import app


@pytest.fixture
def client():
    ssrs = ssrsupervisor.SSRSupervisor()
    app = ssrs.get_app()
    app.config['TESTING'] = True
    return app.test_client()