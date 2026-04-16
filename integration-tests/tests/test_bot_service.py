import uuid

import pytest

_FAKE_BOT_ID = str(uuid.uuid4())


# --- auth guard tests ---

@pytest.mark.live
@pytest.mark.bot
def test_create_bot_requires_authorization_header(require_live, request_api):
    response = request_api("bot", "POST", "/bots/v1/bots", json={})
    assert response.status_code == 401


@pytest.mark.live
@pytest.mark.bot
def test_list_bots_requires_authorization_header(require_live, request_api):
    response = request_api("bot", "GET", "/bots/v1/bots")
    assert response.status_code == 401


@pytest.mark.live
@pytest.mark.bot
def test_get_bot_requires_authorization_header(require_live, request_api):
    response = request_api("bot", "GET", f"/bots/v1/bots/{_FAKE_BOT_ID}")
    assert response.status_code == 401


@pytest.mark.live
@pytest.mark.bot
def test_stream_bot_events_requires_authorization_header(require_live, request_api):
    response = request_api("bot", "GET", f"/bots/v1/bots/{_FAKE_BOT_ID}/stream")
    assert response.status_code == 401


@pytest.mark.live
@pytest.mark.bot
def test_update_bot_requires_authorization_header(require_live, request_api):
    response = request_api("bot", "PATCH", f"/bots/v1/bots/{_FAKE_BOT_ID}", json={})
    assert response.status_code == 401


@pytest.mark.live
@pytest.mark.bot
def test_delete_bot_requires_authorization_header(require_live, request_api):
    response = request_api("bot", "DELETE", f"/bots/v1/bots/{_FAKE_BOT_ID}")
    assert response.status_code == 401


@pytest.mark.live
@pytest.mark.bot
def test_create_bot_rejects_invalid_token(require_live, request_api):
    headers = {"Authorization": "Bearer not-a-real-token"}
    response = request_api("bot", "POST", "/bots/v1/bots", json={}, headers=headers)
    assert response.status_code == 401


# --- happy path tests ---

@pytest.mark.live
@pytest.mark.bot
def test_list_bots_returns_list(require_live, request_api, auth_headers):
    response = request_api("bot", "GET", "/bots/v1/bots", headers=auth_headers)
    assert response.status_code == 200
    body = response.json()
    assert isinstance(body, list)


@pytest.mark.live
@pytest.mark.bot
def test_get_bot_not_found(require_live, request_api, auth_headers):
    response = request_api("bot", "GET", f"/bots/v1/bots/{_FAKE_BOT_ID}", headers=auth_headers)
    assert response.status_code == 404


@pytest.mark.live
@pytest.mark.bot
def test_delete_bot_not_found(require_live, request_api, auth_headers):
    response = request_api("bot", "DELETE", f"/bots/v1/bots/{_FAKE_BOT_ID}", headers=auth_headers)
    assert response.status_code == 404


@pytest.mark.live
@pytest.mark.bot
def test_update_bot_not_found(require_live, request_api, auth_headers):
    response = request_api(
        "bot", "PATCH", f"/bots/v1/bots/{_FAKE_BOT_ID}",
        json={"status": "stopped"},
        headers=auth_headers,
    )
    assert response.status_code == 404


@pytest.mark.live
@pytest.mark.bot
def test_create_bot_rejects_missing_fields(require_live, request_api, auth_headers):
    response = request_api("bot", "POST", "/bots/v1/bots", json={}, headers=auth_headers)
    assert response.status_code == 400


@pytest.mark.live
@pytest.mark.bot
def test_create_bot_rejects_invalid_symbol(require_live, request_api, auth_headers, account_id):
    response = request_api(
        "bot", "POST", "/bots/v1/bots",
        json={"account_id": account_id, "symbol": "invalid symbol!", "allocation_percent": 10},
        headers=auth_headers,
    )
    assert response.status_code == 400


@pytest.mark.live
@pytest.mark.bot
def test_create_bot_rejects_invalid_allocation(require_live, request_api, auth_headers, account_id):
    response = request_api(
        "bot", "POST", "/bots/v1/bots",
        json={"account_id": account_id, "symbol": "AAPL", "allocation_percent": 0},
        headers=auth_headers,
    )
    assert response.status_code == 400


@pytest.mark.live
@pytest.mark.bot
def test_create_bot_rejects_unlinked_account(require_live, request_api, auth_headers, account_id):
    # Account exists but has no broker linked — bot creation should fail.
    response = request_api(
        "bot", "POST", "/bots/v1/bots",
        json={"account_id": account_id, "symbol": "AAPL", "allocation_percent": 10},
        headers=auth_headers,
    )
    assert response.status_code == 400
