import pytest


@pytest.mark.live
@pytest.mark.stock
def test_most_actives_rejects_invalid_limit(require_live, request_api):
    response = request_api(
        "stock",
        "GET",
        "/stock-screener/v1/most-actives",
        params={"limit": "not-a-number"},
    )
    assert response.status_code >= 400


@pytest.mark.live
@pytest.mark.stock
def test_movers_rejects_invalid_limit(require_live, request_api):
    response = request_api(
        "stock",
        "GET",
        "/stock-screener/v1/movers",
        params={"limit": "not-a-number"},
    )
    assert response.status_code >= 400


@pytest.mark.live
@pytest.mark.stock
def test_news_rejects_invalid_limit(require_live, request_api):
    response = request_api(
        "stock",
        "GET",
        "/stock-screener/v1/news",
        params={"limit": "not-a-number"},
    )
    assert response.status_code >= 400


@pytest.mark.live
@pytest.mark.stock
@pytest.mark.xfail(reason="Requires valid upstream Alpaca configuration", strict=False)
def test_most_actives_happy_path(require_live, request_api):
    response = request_api(
        "stock",
        "GET",
        "/stock-screener/v1/most-actives",
        params={"limit": 1},
    )
    assert response.status_code == 200


@pytest.mark.live
@pytest.mark.stock
def test_fear_greed_requires_authorization_header(require_live, request_api):
    response = request_api("stock", "GET", "/stock-screener/v1/sentiments/fear-greed")
    assert response.status_code == 401


@pytest.mark.live
@pytest.mark.stock
def test_stock_bars_requires_authorization_header(require_live, request_api):
    response = request_api("stock", "GET", "/stock-screener/v1/stocks/AAPL/bars")
    assert response.status_code == 401


@pytest.mark.live
@pytest.mark.stock
def test_stock_snapshot_requires_authorization_header(require_live, request_api):
    response = request_api("stock", "GET", "/stock-screener/v1/stocks/AAPL/snapshot")
    assert response.status_code == 401
