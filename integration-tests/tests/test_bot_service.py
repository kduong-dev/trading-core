import uuid

import pytest

_FAKE_BOT_ID = str(uuid.uuid4())


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
