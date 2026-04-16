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
