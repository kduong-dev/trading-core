import os
import uuid
from dataclasses import dataclass

import pytest
import requests


def pytest_addoption(parser: pytest.Parser) -> None:
    parser.addoption(
        "--live",
        action="store_true",
        default=False,
        help="Run tests against live local services.",
    )


@dataclass(frozen=True)
class ServiceSettings:
    auth_url: str
    account_url: str
    stock_url: str
    bot_url: str
    reporting_url: str
    storage_url: str
    journal_url: str
    timeout_seconds: float


@pytest.fixture(scope="session")
def settings() -> ServiceSettings:
    return ServiceSettings(
        auth_url=os.getenv("AUTH_SERVICE_URL", "http://localhost:9100"),
        account_url=os.getenv("ACCOUNT_SERVICE_URL", "http://localhost:9000"),
        stock_url=os.getenv("STOCK_SCREENER_URL", "http://localhost:8080"),
        bot_url=os.getenv("BOT_SERVICE_URL", "http://localhost:8081"),
        reporting_url=os.getenv("REPORTING_SERVICE_URL", "http://localhost:8082"),
        storage_url=os.getenv("STORAGE_SERVICE_URL", "http://localhost:8083"),
        journal_url=os.getenv("JOURNAL_SERVICE_URL", "http://localhost:8084"),
        timeout_seconds=float(os.getenv("TEST_TIMEOUT_SECONDS", "10")),
    )


@pytest.fixture(scope="session")
def require_live(request: pytest.FixtureRequest) -> None:
    if not request.config.getoption("--live"):
        pytest.skip("Live integration tests are disabled. Re-run with --live.")


@pytest.fixture(scope="session")
def http_session() -> requests.Session:
    session = requests.Session()
    session.headers.update({"Accept": "application/json"})
    yield session
    session.close()


@pytest.fixture
def request_api(settings: ServiceSettings, http_session: requests.Session):
    service_urls = {
        "auth": settings.auth_url,
        "account": settings.account_url,
        "stock": settings.stock_url,
        "bot": settings.bot_url,
        "reporting": settings.reporting_url,
        "storage": settings.storage_url,
        "journal": settings.journal_url,
    }

    def _request(service: str, method: str, path: str, **kwargs):
        base_url = service_urls[service].rstrip("/")
        url = f"{base_url}/{path.lstrip('/')}"
        timeout = kwargs.pop("timeout", settings.timeout_seconds)
        return http_session.request(method=method, url=url, timeout=timeout, **kwargs)

    return _request


@pytest.fixture(scope="session")
def auth_token(require_live, settings, http_session):
    """Creates a test user and session once per suite. Returns the Bearer token."""
    email = f"it-{uuid.uuid4().hex[:8]}@example.com"
    password = "integration-test-pass-123"
    auth_url = settings.auth_url.rstrip("/")
    timeout = settings.timeout_seconds

    resp = http_session.post(
        f"{auth_url}/auth/v1/users",
        json={"email": email, "password": password},
        timeout=timeout,
    )
    assert resp.status_code == 200, f"create user failed: {resp.text}"

    resp = http_session.post(
        f"{auth_url}/auth/v1/sessions",
        json={"email": email, "password": password},
        timeout=timeout,
    )
    assert resp.status_code == 200, f"create session failed: {resp.text}"
    return resp.json()["access_token"]


@pytest.fixture(scope="session")
def auth_headers(auth_token):
    return {"Authorization": f"Bearer {auth_token}"}


@pytest.fixture(scope="session")
def account_id(require_live, settings, http_session, auth_headers):
    """Creates a test account once per suite. Returns the account_id."""
    account_url = settings.account_url.rstrip("/")
    resp = http_session.post(
        f"{account_url}/accounts/v1/accounts",
        json={"account_name": "Integration Test Account"},
        headers=auth_headers,
        timeout=settings.timeout_seconds,
    )
    assert resp.status_code == 201, f"create account failed: {resp.text}"
    return resp.json()["account_id"]
