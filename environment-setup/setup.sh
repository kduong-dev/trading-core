#!/usr/bin/env bash
# Installs all dependencies needed to run the trading-core stack.
# Supports Linux (apt-get) and macOS (Homebrew). Windows users: run inside WSL 2.

#set -euo pipefail

have() { command -v "$1" >/dev/null 2>&1; }

if [[ "$OSTYPE" == "darwin"* ]]; then
    PM="brew"
    have brew || { echo "Install Homebrew first: https://brew.sh"; exit 1; }
elif have apt-get; then
    PM="apt"
    SUDO=$(have sudo && echo sudo || echo "")
    $SUDO apt-get update
else
    echo "Unsupported platform — install deps manually (docker, go, node, python3)."
    echo "Windows users: run this script inside a WSL 2 Ubuntu shell."
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
    install git git
else
    install docker docker.io
    install go golang-go
    install node nodejs npm
    install python3 python3 python3-pip python3-venv python3.12-venv python3-full
    install git git
fi

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

# --- frontend deps --------------------------------------------------------
if [[ -f "$ROOT/frontend/package.json" ]]; then
    echo "[npm] installing frontend deps"
    (cd "$ROOT/frontend" && npm install)
fi

# --- integration-tests deps ----------------------------------------------
if [[ -f "$ROOT/integration-tests/pyproject.toml" ]]; then
    echo "[pip] installing integration-tests deps"
    (cd "$ROOT/integration-tests" && python3 -m venv .venv && .venv/bin/pip install -e .)
fi

# --- backend deps ---------------------------------------------------------
if [[ -f "$ROOT/backend/go.mod" ]]; then
    echo "[go] downloading backend modules"
    (cd "$ROOT/backend" && go mod download)
fi

echo
echo "Done."
