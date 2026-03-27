import pytest


@pytest.mark.live
@pytest.mark.account
def test_balance_requires_authorization_header(require_live, request_api):
    response = request_api("account", "GET", "/accounts/v1/balance")
    assert response.status_code == 401


@pytest.mark.live
@pytest.mark.account
def test_balance_rejects_invalid_token(require_live, request_api):
    headers = {"Authorization": "Bearer not-a-real-token"}
    response = request_api("account", "GET", "/accounts/v1/balance", headers=headers)
    assert response.status_code == 401


@pytest.mark.live
@pytest.mark.account
def test_link_broker_requires_authorization_header(require_live, request_api):
    response = request_api(
        "account",
        "POST",
        "/accounts/v1/broker/link",
        json={"broker_type": "tastytrade", "broker_id": "12345678"},
    )
    assert response.status_code == 401
