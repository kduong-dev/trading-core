import uuid

import pytest

_FAKE_UPLOAD_ID = str(uuid.uuid4())
_FAKE_FILE_ID = str(uuid.uuid4())


@pytest.mark.live
@pytest.mark.storage
def test_initialise_upload_requires_authorization_header(require_live, request_api):
    response = request_api("storage", "POST", "/storage/v1/uploads", json={})
    assert response.status_code == 401


@pytest.mark.live
@pytest.mark.storage
def test_upload_part_requires_authorization_header(require_live, request_api):
    response = request_api("storage", "PUT", f"/storage/v1/uploads/{_FAKE_UPLOAD_ID}/parts/1", data=b"chunk")
    assert response.status_code == 401


@pytest.mark.live
@pytest.mark.storage
def test_complete_upload_requires_authorization_header(require_live, request_api):
    response = request_api("storage", "POST", f"/storage/v1/uploads/{_FAKE_UPLOAD_ID}/complete", json={})
    assert response.status_code == 401


@pytest.mark.live
@pytest.mark.storage
def test_download_file_requires_authorization_header(require_live, request_api):
    response = request_api("storage", "GET", f"/storage/v1/files/{_FAKE_FILE_ID}")
    assert response.status_code == 401


@pytest.mark.live
@pytest.mark.storage
def test_initialise_upload_rejects_invalid_token(require_live, request_api):
    headers = {"Authorization": "Bearer not-a-real-token"}
    response = request_api("storage", "POST", "/storage/v1/uploads", json={}, headers=headers)
    assert response.status_code == 401


# --- validation tests ---

@pytest.mark.live
@pytest.mark.storage
def test_initialise_upload_rejects_missing_key(require_live, request_api, auth_headers):
    response = request_api("storage", "POST", "/storage/v1/uploads", json={}, headers=auth_headers)
    assert response.status_code == 400


@pytest.mark.live
@pytest.mark.storage
def test_download_file_not_found(require_live, request_api, auth_headers):
    response = request_api("storage", "GET", f"/storage/v1/files/{_FAKE_FILE_ID}", headers=auth_headers)
    assert response.status_code == 404


@pytest.mark.live
@pytest.mark.storage
def test_complete_upload_not_found(require_live, request_api, auth_headers):
    response = request_api(
        "storage", "POST", f"/storage/v1/uploads/{_FAKE_UPLOAD_ID}/complete",
        headers=auth_headers,
    )
    assert response.status_code == 404


# --- happy path tests ---

@pytest.fixture(scope="module")
def completed_upload(request_api, auth_headers, require_live, settings):
    """Runs a full upload flow once per module: initialise → upload part → complete. Returns file_id."""
    import requests as req
    storage_url = settings.storage_url.rstrip("/")
    session = req.Session()
    session.headers.update({"Accept": "application/json"})

    # Step 1: initialise upload
    resp = session.post(
        f"{storage_url}/storage/v1/uploads",
        json={"key": "integration-test.txt", "content_type": "text/plain"},
        headers=auth_headers,
        timeout=settings.timeout_seconds,
    )
    assert resp.status_code == 201, f"initialise upload failed: {resp.text}"
    upload_id = resp.json()["upload_id"]

    # Step 2: upload a single part
    content = b"hello from integration test"
    resp = session.put(
        f"{storage_url}/storage/v1/uploads/{upload_id}/parts/1",
        data=content,
        headers={**auth_headers, "Content-Type": "application/octet-stream"},
        timeout=settings.timeout_seconds,
    )
    assert resp.status_code == 200, f"upload part failed: {resp.text}"

    # Step 3: complete the upload
    resp = session.post(
        f"{storage_url}/storage/v1/uploads/{upload_id}/complete",
        headers=auth_headers,
        timeout=settings.timeout_seconds,
    )
    assert resp.status_code == 201, f"complete upload failed: {resp.text}"
    body = resp.json()
    return {"file_id": body["file_id"], "content": content}


@pytest.mark.live
@pytest.mark.storage
def test_initialise_upload(require_live, request_api, auth_headers):
    response = request_api(
        "storage", "POST", "/storage/v1/uploads",
        json={"key": "test.txt", "content_type": "text/plain"},
        headers=auth_headers,
    )
    assert response.status_code == 201
    body = response.json()
    assert isinstance(body.get("upload_id"), str)
    assert len(body["upload_id"]) > 0


@pytest.mark.live
@pytest.mark.storage
def test_upload_part(require_live, request_api, auth_headers, settings):
    import requests as req
    storage_url = settings.storage_url.rstrip("/")
    session = req.Session()

    # Need a fresh upload for this test
    resp = session.post(
        f"{storage_url}/storage/v1/uploads",
        json={"key": "part-test.txt", "content_type": "text/plain"},
        headers=auth_headers,
        timeout=settings.timeout_seconds,
    )
    assert resp.status_code == 201
    upload_id = resp.json()["upload_id"]

    resp = session.put(
        f"{storage_url}/storage/v1/uploads/{upload_id}/parts/1",
        data=b"test content for part upload",
        headers={**auth_headers, "Content-Type": "application/octet-stream"},
        timeout=settings.timeout_seconds,
    )
    assert resp.status_code == 200
    body = resp.json()
    assert body.get("part_number") == 1
    assert isinstance(body.get("size"), int)
    assert body["size"] > 0


@pytest.mark.live
@pytest.mark.storage
def test_complete_upload(require_live, completed_upload):
    assert isinstance(completed_upload["file_id"], str)
    assert len(completed_upload["file_id"]) > 0


@pytest.mark.live
@pytest.mark.storage
def test_download_file(require_live, request_api, auth_headers, completed_upload):
    file_id = completed_upload["file_id"]
    response = request_api("storage", "GET", f"/storage/v1/files/{file_id}", headers=auth_headers)
    assert response.status_code == 200
    assert response.content == completed_upload["content"]
