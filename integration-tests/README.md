# integration-tests

Python integration test boilerplate for the trading services.

## Prerequisites

- Python 3.10+
- trading-formation services running locally

## Setup

1. Create and activate a virtual environment.
2. Install dependencies:

```bash
pip install -e .
```

3. Copy env defaults and adjust if needed:

```bash
cp .env.example .env
```

On Windows:

```bat
copy .env.example .env
```

## Run tests

By default, live tests are skipped to avoid false failures if services are down.

```bash
pytest
```

Run live integration tests:

```bash
pytest --live
```

Run only smoke tests:

```bash
pytest --live -m smoke
```

## Environment variables

- AUTH_SERVICE_URL (default: http://localhost:9100)
- ACCOUNT_SERVICE_URL (default: http://localhost:9000)
- STOCK_SCREENER_URL (default: http://localhost:8080)
- TEST_TIMEOUT_SECONDS (default: 10)
