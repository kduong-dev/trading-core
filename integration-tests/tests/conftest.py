import os
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
    timeout_seconds: float


@pytest.fixture(scope="session")
def settings() -> ServiceSettings:
    return ServiceSettings(
        auth_url=os.getenv("AUTH_SERVICE_URL", "http://localhost:9100"),
        account_url=os.getenv("ACCOUNT_SERVICE_URL", "http://localhost:9000"),
        stock_url=os.getenv("STOCK_SCREENER_URL", "http://localhost:8080"),
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
    }

    def _request(service: str, method: str, path: str, **kwargs):
        base_url = service_urls[service].rstrip("/")
        url = f"{base_url}/{path.lstrip('/')}"
        timeout = kwargs.pop("timeout", settings.timeout_seconds)
        return http_session.request(method=method, url=url, timeout=timeout, **kwargs)

    return _request
