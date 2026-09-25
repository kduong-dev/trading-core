import pytest


# --- auth guard tests ---

@pytest.mark.live
@pytest.mark.journal
def test_upsert_entry_requires_authorization_header(require_live, request_api):
    response = request_api("journal", "PUT", "/journal/v1/entries/2026-04-22", json={})
    assert response.status_code == 401


@pytest.mark.live
@pytest.mark.journal
def test_list_entries_requires_authorization_header(require_live, request_api):
    response = request_api("journal", "GET", "/journal/v1/entries")
    assert response.status_code == 401


@pytest.mark.live
@pytest.mark.journal
def test_get_entry_requires_authorization_header(require_live, request_api):
    response = request_api("journal", "GET", "/journal/v1/entries/2026-04-22")
    assert response.status_code == 401


@pytest.mark.live
@pytest.mark.journal
def test_delete_entry_requires_authorization_header(require_live, request_api):
    response = request_api("journal", "DELETE", "/journal/v1/entries/2026-04-22")
    assert response.status_code == 401


@pytest.mark.live
@pytest.mark.journal
def test_upsert_entry_rejects_invalid_token(require_live, request_api):
    headers = {"Authorization": "Bearer not-a-real-token"}
    response = request_api("journal", "PUT", "/journal/v1/entries/2026-04-22", json={}, headers=headers)
    assert response.status_code == 401


# --- validation tests ---

@pytest.mark.live
@pytest.mark.journal
def test_upsert_entry_rejects_invalid_date(require_live, request_api, auth_headers):
    response = request_api(
        "journal", "PUT", "/journal/v1/entries/not-a-date",
        json={}, headers=auth_headers,
    )
    assert response.status_code == 400


@pytest.mark.live
@pytest.mark.journal
def test_upsert_entry_rejects_discipline_score_out_of_range(require_live, request_api, auth_headers):
    response = request_api(
        "journal", "PUT", "/journal/v1/entries/2026-04-22",
        json={"discipline_score": 42},
        headers=auth_headers,
    )
    assert response.status_code == 400


@pytest.mark.live
@pytest.mark.journal
def test_list_entries_rejects_invalid_page_size(require_live, request_api, auth_headers):
    response = request_api(
        "journal", "GET", "/journal/v1/entries",
        params={"page_size": 0},
        headers=auth_headers,
    )
    assert response.status_code == 400


@pytest.mark.live
@pytest.mark.journal
def test_list_entries_rejects_negative_page(require_live, request_api, auth_headers):
    response = request_api(
        "journal", "GET", "/journal/v1/entries",
        params={"page": -1},
        headers=auth_headers,
    )
    assert response.status_code == 400


@pytest.mark.live
@pytest.mark.journal
def test_list_entries_rejects_invalid_from_date(require_live, request_api, auth_headers):
    response = request_api(
        "journal", "GET", "/journal/v1/entries",
        params={"from": "nope"},
        headers=auth_headers,
    )
    assert response.status_code == 400


@pytest.mark.live
@pytest.mark.journal
def test_get_entry_rejects_invalid_date(require_live, request_api, auth_headers):
    response = request_api(
        "journal", "GET", "/journal/v1/entries/nope",
        headers=auth_headers,
    )
    assert response.status_code == 400


# --- happy path tests ---

@pytest.mark.live
@pytest.mark.journal
def test_upsert_entry(require_live, request_api, auth_headers):
    response = request_api(
        "journal", "PUT", "/journal/v1/entries/2026-04-22",
        json={
            "notes": "Clean breakout entry, let winner run",
            "tags": ["breakout", "AAPL"],
            "mood": "focused",
            "discipline_score": 8,
        },
        headers=auth_headers,
    )
    assert response.status_code == 202
    body = response.json()
    assert body["date"] == "2026-04-22"
    assert body["notes"] == "Clean breakout entry, let winner run"
    assert body["discipline_score"] == 8


@pytest.mark.live
@pytest.mark.journal
def test_get_entry(require_live, request_api, auth_headers):
    request_api(
        "journal", "PUT", "/journal/v1/entries/2026-04-23",
        json={"notes": "Sized down, avoided chop"},
        headers=auth_headers,
    )
    response = request_api(
        "journal", "GET", "/journal/v1/entries/2026-04-23",
        headers=auth_headers,
    )
    assert response.status_code == 200
    body = response.json()
    assert body["date"] == "2026-04-23"
    assert body["notes"] == "Sized down, avoided chop"


@pytest.mark.live
@pytest.mark.journal
def test_get_entry_not_found(require_live, request_api, auth_headers):
    response = request_api(
        "journal", "GET", "/journal/v1/entries/2030-01-01",
        headers=auth_headers,
    )
    assert response.status_code == 404


@pytest.mark.live
@pytest.mark.journal
def test_upsert_is_idempotent_on_same_date(require_live, request_api, auth_headers):
    date = "2026-04-24"
    request_api(
        "journal", "PUT", f"/journal/v1/entries/{date}",
        json={"notes": "first"}, headers=auth_headers,
    )
    request_api(
        "journal", "PUT", f"/journal/v1/entries/{date}",
        json={"notes": "second"}, headers=auth_headers,
    )
    response = request_api(
        "journal", "GET", f"/journal/v1/entries/{date}",
        headers=auth_headers,
    )
    assert response.status_code == 200
    assert response.json()["notes"] == "second"


@pytest.mark.live
@pytest.mark.journal
def test_list_entries_filters_by_date_range(require_live, request_api, auth_headers):
    for date in ["2026-05-01", "2026-05-02", "2026-05-03"]:
        request_api(
            "journal", "PUT", f"/journal/v1/entries/{date}",
            json={"notes": f"entry {date}"}, headers=auth_headers,
        )
    response = request_api(
        "journal", "GET", "/journal/v1/entries",
        params={"from": "2026-05-01", "to": "2026-05-02"},
        headers=auth_headers,
    )
    assert response.status_code == 200
    body = response.json()
    dates = [entry["date"] for entry in body["entries"]]
    assert "2026-05-01" in dates
    assert "2026-05-02" in dates
    assert "2026-05-03" not in dates


@pytest.mark.live
@pytest.mark.journal
def test_delete_entry(require_live, request_api, auth_headers):
    date = "2026-04-25"
    request_api(
        "journal", "PUT", f"/journal/v1/entries/{date}",
        json={"notes": "to delete"}, headers=auth_headers,
    )
    response = request_api(
        "journal", "DELETE", f"/journal/v1/entries/{date}",
        headers=auth_headers,
    )
    assert response.status_code == 204

    response = request_api(
        "journal", "GET", f"/journal/v1/entries/{date}",
        headers=auth_headers,
    )
    assert response.status_code == 404


@pytest.mark.live
@pytest.mark.journal
def test_delete_entry_not_found(require_live, request_api, auth_headers):
    response = request_api(
        "journal", "DELETE", "/journal/v1/entries/2030-02-02",
        headers=auth_headers,
    )
    assert response.status_code == 404
