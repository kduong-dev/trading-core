import uuid

import pytest

_FAKE_ACCOUNT_ID = str(uuid.uuid4())


@pytest.mark.live
@pytest.mark.account
def test_get_balance_requires_authorization_header(require_live, request_api):
    response = request_api("account", "GET", f"/accounts/v1/accounts/{_FAKE_ACCOUNT_ID}/balances")
    assert response.status_code == 401


@pytest.mark.live
@pytest.mark.account
def test_get_balance_rejects_invalid_token(require_live, request_api):
    headers = {"Authorization": "Bearer not-a-real-token"}
    response = request_api("account", "GET", f"/accounts/v1/accounts/{_FAKE_ACCOUNT_ID}/balances", headers=headers)
    assert response.status_code == 401


@pytest.mark.live
@pytest.mark.account
def test_create_account_requires_authorization_header(require_live, request_api):
    response = request_api("account", "POST", "/accounts/v1/accounts", json={})
    assert response.status_code == 401


@pytest.mark.live
@pytest.mark.account
def test_list_accounts_requires_authorization_header(require_live, request_api):
    response = request_api("account", "GET", "/accounts/v1/accounts")
    assert response.status_code == 401


@pytest.mark.live
@pytest.mark.account
def test_get_account_requires_authorization_header(require_live, request_api):
    response = request_api("account", "GET", f"/accounts/v1/accounts/{_FAKE_ACCOUNT_ID}")
    assert response.status_code == 401


@pytest.mark.live
@pytest.mark.account
def test_start_broker_selection_requires_authorization_header(require_live, request_api):
    response = request_api(
        "account",
        "POST",
        f"/accounts/v1/accounts/{_FAKE_ACCOUNT_ID}/brokers",
        json={"broker_type": "tastytrade"},
    )
    assert response.status_code == 401


@pytest.mark.live
@pytest.mark.account
def test_get_pending_broker_selection_requires_authorization_header(require_live, request_api):
    response = request_api("account", "GET", f"/accounts/v1/accounts/{_FAKE_ACCOUNT_ID}/brokers")
    assert response.status_code == 401


@pytest.mark.live
@pytest.mark.account
def test_complete_broker_selection_requires_authorization_header(require_live, request_api):
    response = request_api(
        "account",
        "PUT",
        f"/accounts/v1/accounts/{_FAKE_ACCOUNT_ID}/brokers",
        json={},
    )
    assert response.status_code == 401
