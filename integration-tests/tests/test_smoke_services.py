import pytest


@pytest.mark.live
@pytest.mark.smoke
def test_auth_service_reachable(require_live, request_api):
    response = request_api("auth", "GET", "/__healthcheck__")
    assert response.status_code == 404


@pytest.mark.live
@pytest.mark.smoke
def test_account_service_reachable(require_live, request_api):
    response = request_api("account", "GET", "/__healthcheck__")
    assert response.status_code == 404


@pytest.mark.live
@pytest.mark.smoke
def test_stock_service_reachable(require_live, request_api):
    response = request_api("stock", "GET", "/__healthcheck__")
    assert response.status_code == 404
