# trading-core

Trading platform monorepo.

| Directory | Contents |
|---|---|
| [backend](backend/) | Go microservices and CLI tools (auth, accounts, bots, journal, reporting, stock screener) |
| [frontend](frontend/) | Next.js web app |
| [integration-tests](integration-tests/) | Python integration tests against locally running services |

## Setup

`setup.sh` installs everything `trading-core` needs: Docker, Go, Node.js,
Python, and the per-project deps (`npm install`, `pip install -e`, `go mod download`).

> **Primary target: Linux (Debian/Ubuntu).** Windows and macOS users should
> read the platform notes below before running.

```bash
./setup.sh
```

The script is idempotent — re-running it only installs what's missing.

### Windows

Use [WSL 2](https://learn.microsoft.com/en-us/windows/wsl/install) with an
Ubuntu distro, then run `setup.sh` inside the WSL terminal. Docker Desktop with
the WSL 2 backend is the recommended way to run the compose stack on Windows.

```powershell
# one-time WSL setup (elevated PowerShell)
wsl --install -d Ubuntu
```

Then inside the WSL shell:

```bash
./setup.sh
```

### macOS

[Homebrew](https://brew.sh) is required. Install it first if you don't have it,
then run the script — it detects `brew` automatically and uses it instead of
`apt-get`.

```bash
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
./setup.sh
```

### What gets installed

| Tool | Why |
|---|---|
| Docker | runs Redis/Postgres and the compose stack |
| Go 1.27+ | builds [backend](backend/) services |
| Node.js (LTS) | builds & runs [frontend](frontend/) (Next.js) |
| Python 3.10+ | runs [integration-tests](integration-tests/) |
| Git | repo operations |

Packages come from `apt-get` on Debian/Ubuntu and Homebrew on macOS. The
script skips Go if any version is already installed, and apt's `golang-go` can
lag behind — check `go version` and install from [go.dev](https://go.dev/dl/)
if it's older than 1.27.

### After setup

1. Start Docker Desktop / the docker daemon.
2. Run a service: see [backend/README.md](backend/README.md) and
   [frontend/README.md](frontend/README.md).
