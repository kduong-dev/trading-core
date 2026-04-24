import uuid

import pytest

_FAKE_ACCOUNT_ID = str(uuid.uuid4())


# --- auth guard tests ---

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


# --- happy path tests ---

@pytest.mark.live
@pytest.mark.account
def test_create_account(require_live, request_api, auth_headers):
    response = request_api(
        "account",
        "POST",
        "/accounts/v1/accounts",
        json={"account_name": "My Test Account"},
        headers=auth_headers,
    )
    assert response.status_code == 201
    body = response.json()
    assert isinstance(body.get("account_id"), str)
    assert len(body["account_id"]) > 0
    assert body.get("account_name") == "My Test Account"


@pytest.mark.live
@pytest.mark.account
def test_list_accounts(require_live, request_api, auth_headers, account_id):
    response = request_api("account", "GET", "/accounts/v1/accounts", headers=auth_headers)
    assert response.status_code == 200
    accounts = response.json()
    assert isinstance(accounts, list)
    assert any(a["account_id"] == account_id for a in accounts)


@pytest.mark.live
@pytest.mark.account
def test_get_account(require_live, request_api, auth_headers, account_id):
    response = request_api("account", "GET", f"/accounts/v1/accounts/{account_id}", headers=auth_headers)
    assert response.status_code == 200
    body = response.json()
    assert body.get("account_id") == account_id


@pytest.mark.live
@pytest.mark.account
def test_get_account_not_found(require_live, request_api, auth_headers):
    response = request_api("account", "GET", f"/accounts/v1/accounts/{_FAKE_ACCOUNT_ID}", headers=auth_headers)
    assert response.status_code == 404


@pytest.mark.live
@pytest.mark.account
def test_get_balance_requires_linked_broker(require_live, request_api, auth_headers, account_id):
    # Account exists but has no broker linked — expect 400.
    response = request_api("account", "GET", f"/accounts/v1/accounts/{account_id}/balances", headers=auth_headers)
    assert response.status_code == 400


# --- daily PnL endpoint ---

@pytest.mark.live
@pytest.mark.account
def test_get_daily_pnl_requires_authorization_header(require_live, request_api):
    response = request_api(
        "account", "GET",
        f"/accounts/v1/accounts/{_FAKE_ACCOUNT_ID}/pnl/daily",
        params={"from": "2026-04-01", "to": "2026-04-30"},
    )
    assert response.status_code == 401


@pytest.mark.live
@pytest.mark.account
def test_get_daily_pnl_requires_from_and_to(require_live, request_api, auth_headers, account_id):
    response = request_api(
        "account", "GET",
        f"/accounts/v1/accounts/{account_id}/pnl/daily",
        headers=auth_headers,
    )
    assert response.status_code == 400


@pytest.mark.live
@pytest.mark.account
def test_get_daily_pnl_rejects_invalid_from(require_live, request_api, auth_headers, account_id):
    response = request_api(
        "account", "GET",
        f"/accounts/v1/accounts/{account_id}/pnl/daily",
        params={"from": "nope", "to": "2026-04-30"},
        headers=auth_headers,
    )
    assert response.status_code == 400


@pytest.mark.live
@pytest.mark.account
def test_get_daily_pnl_rejects_inverted_range(require_live, request_api, auth_headers, account_id):
    response = request_api(
        "account", "GET",
        f"/accounts/v1/accounts/{account_id}/pnl/daily",
        params={"from": "2026-04-30", "to": "2026-04-01"},
        headers=auth_headers,
    )
    assert response.status_code == 400


@pytest.mark.live
@pytest.mark.account
def test_get_daily_pnl_rejects_range_over_max(require_live, request_api, auth_headers, account_id):
    response = request_api(
        "account", "GET",
        f"/accounts/v1/accounts/{account_id}/pnl/daily",
        params={"from": "2024-01-01", "to": "2026-12-31"},
        headers=auth_headers,
    )
    assert response.status_code == 400


@pytest.mark.live
@pytest.mark.account
def test_get_daily_pnl_requires_linked_broker(require_live, request_api, auth_headers, account_id):
    # Account exists but has no broker linked — expect 400.
    response = request_api(
        "account", "GET",
        f"/accounts/v1/accounts/{account_id}/pnl/daily",
        params={"from": "2026-04-01", "to": "2026-04-30"},
        headers=auth_headers,
    )
    assert response.status_code == 400
