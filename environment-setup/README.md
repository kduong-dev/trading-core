# environment-setup

One-shot installer for everything `trading-core` needs: Docker, Go, Node.js,
Python, GNU Make, Ansible, the sibling repos (cloned from
`github.com/trading-core`), and the per-service deps (`npm install`, `pip
install -e`, `go mod download`).

> **Primary target: Linux (Debian/Ubuntu).** Windows and macOS users should
> read the platform notes below before running.

## Usage

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

## What gets installed

| Tool | Why |
|---|---|
| Docker | runs the compose stack ([trading-formation](../trading-formation/)) |
| Go 1.25+ | builds [trading-backend](../trading-backend/) services |
| Node.js (LTS) | builds & runs [trading-frontend](../trading-frontend/) (Next.js) |
| Python 3.10+ | runs [integration-tests](../integration-tests/) and hosts Ansible |
| Ansible | renders `.env` files and the host-mode `Makefile` |
| GNU Make | proxy-mode targets (`make run-<svc>`) |
| Git | repo operations |

Package source: `apt-get` (Debian/Ubuntu).

## After setup

1. Start Docker Desktop / the docker daemon.
2. Decrypt or fill in `trading-formation/secrets.yml` (ansible-vault).
3. `cd trading-formation && ./run-services.sh render && ./run-services.sh start`.

See [trading-formation/README.md](../trading-formation/README.md) for the full
boot sequence.
