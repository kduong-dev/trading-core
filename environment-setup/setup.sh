#!/usr/bin/env bash
# Installs all dependencies needed to run the trading-core stack on Linux/macOS.
# Picks the right package manager (apt or brew) automatically.

set -euo pipefail

have() { command -v "$1" >/dev/null 2>&1; }

if [[ "$OSTYPE" == "darwin"* ]]; then
    PM="brew"
    have brew || { echo "Install Homebrew first: https://brew.sh"; exit 1; }
elif have apt-get; then
    PM="apt"
    SUDO=$(have sudo && echo sudo || echo "")
    $SUDO apt-get update
else
    echo "Unsupported platform — install deps manually (docker, go, node, python3, make, ansible)."
    exit 1
fi

install() {
    local probe="$1"; shift
    if have "$probe"; then
        echo "[ok]   $probe already installed"
        return
    fi
    echo "[install] $probe ($*)"
    if [[ "$PM" == "brew" ]]; then
        brew install "$@"
    else
        $SUDO apt-get install -y "$@"
    fi
}

# --- system tools ---------------------------------------------------------
if [[ "$PM" == "brew" ]]; then
    install docker --cask docker
    install go go
    install node node
    install python3 python@3.12
    install make make
    install ansible-playbook ansible
    install git git
else
    install docker docker.io
    install go golang-go
    install node nodejs npm
    install python3 python3 python3-pip python3-venv
    install make make
    install ansible-playbook ansible
    install git git
fi

ROOT="$(cd "$(dirname "$0")/.." && pwd)"

# --- sibling repos --------------------------------------------------------
ORG="https://github.com/trading-core"
for r in trading-backend trading-frontend trading-formation integration-tests; do
    if [[ -d "$ROOT/$r/.git" ]]; then
        echo "[ok]   $r already cloned"
    else
        echo "[clone] $ORG/$r"
        git clone "$ORG/$r.git" "$ROOT/$r"
    fi
done

# --- frontend deps --------------------------------------------------------
if [[ -f "$ROOT/trading-frontend/package.json" ]]; then
    echo "[npm] installing frontend deps"
    (cd "$ROOT/trading-frontend" && npm install)
fi

# --- integration-tests deps ----------------------------------------------
if [[ -f "$ROOT/integration-tests/pyproject.toml" ]]; then
    echo "[pip] installing integration-tests deps"
    (cd "$ROOT/integration-tests" && python3 -m pip install --user -e .)
fi

# --- backend deps ---------------------------------------------------------
if [[ -f "$ROOT/trading-backend/go.mod" ]]; then
    echo "[go] downloading backend modules"
    (cd "$ROOT/trading-backend" && go mod download)
fi

echo
echo "Done. Make sure the Docker daemon is running before ./run-services.sh start."
