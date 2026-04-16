import uuid

import pytest


@pytest.mark.live
@pytest.mark.auth
def test_create_user_rejects_short_password(require_live, request_api):
    payload = {
        "email": f"test-{uuid.uuid4().hex[:8]}@example.com",
        "password": "short",
    }
    response = request_api("auth", "POST", "/auth/v1/users", json=payload)
    assert response.status_code == 400


@pytest.mark.live
@pytest.mark.auth
def test_create_session_requires_email_and_password(require_live, request_api):
    response = request_api("auth", "POST", "/auth/v1/sessions", json={})
    assert response.status_code == 400


@pytest.mark.live
@pytest.mark.auth
def test_create_session_rejects_invalid_credentials(require_live, request_api):
    payload = {
        "email": f"missing-{uuid.uuid4().hex[:8]}@example.com",
        "password": "incorrect-password",
    }
    response = request_api("auth", "POST", "/auth/v1/sessions", json=payload)
    assert response.status_code == 401


@pytest.mark.live
@pytest.mark.auth
def test_refresh_session_requires_authorization_header(require_live, request_api):
    response = request_api("auth", "POST", "/auth/v1/sessions/refresh")
    assert response.status_code == 401


@pytest.mark.live
@pytest.mark.auth
def test_refresh_session_rejects_invalid_token(require_live, request_api):
    headers = {"Authorization": "Bearer not-a-real-token"}
    response = request_api("auth", "POST", "/auth/v1/sessions/refresh", headers=headers)
    assert response.status_code == 401


@pytest.mark.live
@pytest.mark.auth
def test_create_user_then_create_session(require_live, request_api):
    email = f"it-{uuid.uuid4().hex[:8]}@example.com"
    password = "test-pass-123"

    create_user_response = request_api(
        "auth",
        "POST",
        "/auth/v1/users",
        json={"email": email, "password": password},
    )
    assert create_user_response.status_code == 201

    create_session_response = request_api(
        "auth",
        "POST",
        "/auth/v1/sessions",
        json={"email": email, "password": password},
    )
    assert create_session_response.status_code == 200
    body = create_session_response.json()
    assert body.get("token_type") == "Bearer"
    assert isinstance(body.get("access_token"), str)
    assert len(body.get("access_token", "")) > 0
