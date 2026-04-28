# environment-setup

One-shot installer for everything `trading-core` needs: Docker, Go, Node.js,
Python, GNU Make, Ansible, the sibling repos (cloned from
`github.com/trading-core`), and the per-service deps (`npm install`, `pip
install -e`, `go mod download`).

## Usage

**Windows** (run an elevated PowerShell so winget can install machine-wide):

```powershell
.\setup.ps1
```

**Linux / macOS**:

```bash
./setup.sh
```

The scripts are idempotent — re-running them only installs what's missing.

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

Package sources:
- Windows → `winget` (Docker Desktop, Go, Node LTS, Python 3.12, ezwinports.make, Git) + `pipx` for Ansible
- macOS → Homebrew
- Linux → `apt-get` (Debian/Ubuntu)

## After setup

1. Start Docker Desktop / the docker daemon.
2. Decrypt or fill in `trading-formation/secrets.yml` (ansible-vault).
3. `cd trading-formation && ./run-services.sh render && ./run-services.sh start`.

See [trading-formation/README.md](../trading-formation/README.md) for the full
boot sequence.
