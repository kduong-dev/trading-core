import uuid

import pytest
import requests as req

_FAKE_REPORT_ID = str(uuid.uuid4())


# --- auth guard tests ---

@pytest.mark.live
@pytest.mark.reporting
def test_enqueue_report_requires_authorization_header(require_live, request_api):
    response = request_api("reporting", "POST", "/reports/v1/jobs", json={})
    assert response.status_code == 401


@pytest.mark.live
@pytest.mark.reporting
def test_list_reports_requires_authorization_header(require_live, request_api):
    response = request_api("reporting", "GET", "/reports/v1/jobs")
    assert response.status_code == 401


@pytest.mark.live
@pytest.mark.reporting
def test_get_report_requires_authorization_header(require_live, request_api):
    response = request_api("reporting", "GET", f"/reports/v1/jobs/{_FAKE_REPORT_ID}")
    assert response.status_code == 401


@pytest.mark.live
@pytest.mark.reporting
def test_download_report_requires_authorization_header(require_live, request_api):
    response = request_api("reporting", "GET", f"/reports/v1/jobs/{_FAKE_REPORT_ID}/download")
    assert response.status_code == 401


@pytest.mark.live
@pytest.mark.reporting
def test_enqueue_report_rejects_invalid_token(require_live, request_api):
    headers = {"Authorization": "Bearer not-a-real-token"}
    response = request_api("reporting", "POST", "/reports/v1/jobs", json={}, headers=headers)
    assert response.status_code == 401


# --- validation tests ---

@pytest.mark.live
@pytest.mark.reporting
def test_enqueue_report_rejects_missing_kind(require_live, request_api, auth_headers):
    response = request_api("reporting", "POST", "/reports/v1/jobs", json={}, headers=auth_headers)
    assert response.status_code == 400


@pytest.mark.live
@pytest.mark.reporting
def test_list_reports_rejects_invalid_page_size(require_live, request_api, auth_headers):
    response = request_api(
        "reporting", "GET", "/reports/v1/jobs",
        params={"page_size": 0},
        headers=auth_headers,
    )
    assert response.status_code == 400


@pytest.mark.live
@pytest.mark.reporting
def test_list_reports_rejects_negative_page(require_live, request_api, auth_headers):
    response = request_api(
        "reporting", "GET", "/reports/v1/jobs",
        params={"page": -1},
        headers=auth_headers,
    )
    assert response.status_code == 400


# --- happy path tests ---

@pytest.fixture(scope="module")
def enqueued_report(request_api, auth_headers, require_live, settings):
    """Enqueues a report once per module and returns the report body."""
    reporting_url = settings.reporting_url.rstrip("/")
    resp = req.Session().post(
        f"{reporting_url}/reports/v1/jobs",
        json={"kind": "backtest", "name": "Integration Test Report"},
        headers=auth_headers,
        timeout=settings.timeout_seconds,
    )
    assert resp.status_code == 202, f"enqueue report failed: {resp.text}"
    return resp.json()


@pytest.mark.live
@pytest.mark.reporting
def test_enqueue_report(require_live, request_api, auth_headers):
    response = request_api(
        "reporting", "POST", "/reports/v1/jobs",
        json={"kind": "backtest", "name": "Test Report"},
        headers=auth_headers,
    )
    assert response.status_code == 202
    body = response.json()
    assert isinstance(body.get("id"), str)
    assert body.get("kind") == "backtest"
    assert body.get("status") == "pending"


@pytest.mark.live
@pytest.mark.reporting
def test_list_reports(require_live, request_api, auth_headers, enqueued_report):
    response = request_api("reporting", "GET", "/reports/v1/jobs", headers=auth_headers)
    assert response.status_code == 200
    body = response.json()
    assert isinstance(body, list)
    ids = [r["id"] for r in body]
    assert enqueued_report["id"] in ids


@pytest.mark.live
@pytest.mark.reporting
def test_get_report(require_live, request_api, auth_headers, enqueued_report):
    report_id = enqueued_report["id"]
    response = request_api("reporting", "GET", f"/reports/v1/jobs/{report_id}", headers=auth_headers)
    assert response.status_code == 200
    body = response.json()
    assert body.get("id") == report_id
    assert body.get("kind") == "backtest"


@pytest.mark.live
@pytest.mark.reporting
def test_get_report_not_found(require_live, request_api, auth_headers):
    response = request_api("reporting", "GET", f"/reports/v1/jobs/{_FAKE_REPORT_ID}", headers=auth_headers)
    assert response.status_code == 404


@pytest.mark.live
@pytest.mark.reporting
def test_download_report_conflicts_when_pending(require_live, request_api, auth_headers, enqueued_report):
    # Report is still pending — download should fail with 409.
    report_id = enqueued_report["id"]
    response = request_api("reporting", "GET", f"/reports/v1/jobs/{report_id}/download", headers=auth_headers)
    assert response.status_code == 409
