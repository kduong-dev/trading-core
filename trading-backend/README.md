# trading-backend

A Go monorepo housing the microservices and CLI tools that power the trading platform. Services are orchestrated via `trading-formation` and run behind a Traefik reverse proxy.

## Services

| Service | Port | Description |
|---|---|---|
| `authentication-service` | 9100 | JWT-based user authentication |
| `account-service` | dynamic | Broker account management and TastyTrade OAuth flow |
| `bot-service` | dynamic | Trading bot creation, configuration, and execution |
| `journal-service` | dynamic | Trading journal entry management |
| `reporting-service` | dynamic | Report generation backed by the storage service |
| `stock-screener` | dynamic | Stock screening with market sentiment (CNN Fear & Greed, Alpaca) |
| `storage-service` | dynamic | File storage (filesystem or pluggable backend) |

## CLI Tools

| Tool | Description |
|---|---|
| `backtest-cli` | Replay historical market data against a trading strategy; supports tune mode |
| `tasty-trade-cli` | Interactive CLI to browse TastyTrade accounts and stream live quotes |
| `token-secret-generator` | Generates a random 256-bit base64 secret for `TOKEN_SECRET` |

## Internal Packages

| Package | Description |
|---|---|
| `internal/tradingstrategy` | Trading strategy implementations: trend entry, breakout, oversold/overbought, ATR stop, position sizing, session guard |
| `internal/broker` | Broker client interfaces and implementations (TastyTrade, Alpaca) |
| `internal/eventsource` | Event sourcing infrastructure: log factory, Redis and in-memory backends, subscriptions |
| `internal/backtest` | Backtesting framework: replay engine, indicators, charting |
| `internal/auth` | JWT middleware and audience-scoped token validation |
| `internal/authz` | Request authorisation helpers |
| `internal/httpx` | HTTP response utilities including standardised error responses |
| `internal/config` | Environment variable helpers (`EnvString`, `EnvStringOrFatal`, base URL config) |
| `internal/contextx` | Context key/value helpers |
| `internal/fatal` | `fatal.OnError` and `fatal.Unless` for unrecoverable startup errors |
| `internal/logger` | Structured logging configuration |
| `internal/iterator` | Generic iterator interface |

## Tech Stack

- **Language:** Go 1.25
- **HTTP routing:** `gorilla/mux`
- **WebSocket:** `gorilla/websocket`
- **Auth:** `golang-jwt/jwt` (HS256)
- **Event store backends:** Redis (`go-redis`), in-memory (`miniredis` for tests)
- **Database:** PostgreSQL (`lib/pq`) — used by `authentication-service`
- **Error enrichment:** `ansel1/merry` with HTTP status codes
- **Testing:** `smartystreets/goconvey` (BDD style)

## Running Services Locally

Services are run via the Makefile in **proxy mode** — each run picks a free port and registers it with `trading-formation` so Traefik routes `<subdomain>.trading-core.local` to it.

```bash
make run-account-service
make run-authentication-service
make run-bot-service
make run-journal-service
make run-reporting-service
make run-stock-screener
make run-storage-service
```

Run `make help` for the full list of targets.

> Prerequisites: Redis on `localhost:6379`, PostgreSQL on `localhost:5432` (for `authentication-service`), and `trading-formation` running for proxy routing.

## Running Tests

```bash
go test ./...
```

Integration tests live in the sibling `integration-tests` repository.
