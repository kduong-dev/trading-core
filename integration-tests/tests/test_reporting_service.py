import uuid

import pytest

_FAKE_REPORT_ID = str(uuid.uuid4())


@pytest.mark.live
@pytest.mark.reporting
def test_enqueue_report_requires_authorization_header(require_live, request_api):
    response = request_api("reporting", "POST", "/reports/v1/reports", json={})
    assert response.status_code == 401


@pytest.mark.live
@pytest.mark.reporting
def test_list_reports_requires_authorization_header(require_live, request_api):
    response = request_api("reporting", "GET", "/reports/v1/reports")
    assert response.status_code == 401


@pytest.mark.live
@pytest.mark.reporting
def test_get_report_requires_authorization_header(require_live, request_api):
    response = request_api("reporting", "GET", f"/reports/v1/reports/{_FAKE_REPORT_ID}")
    assert response.status_code == 401


@pytest.mark.live
@pytest.mark.reporting
def test_download_report_requires_authorization_header(require_live, request_api):
    response = request_api("reporting", "GET", f"/reports/v1/reports/{_FAKE_REPORT_ID}/download")
    assert response.status_code == 401


@pytest.mark.live
@pytest.mark.reporting
def test_enqueue_report_rejects_invalid_token(require_live, request_api):
    headers = {"Authorization": "Bearer not-a-real-token"}
    response = request_api("reporting", "POST", "/reports/v1/reports", json={}, headers=headers)
    assert response.status_code == 401
