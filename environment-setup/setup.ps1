# Installs all dependencies needed to run the trading-core stack on Windows.
# Requires winget (ships with Windows 10 1809+/Windows 11). Run from an elevated
# shell so winget can install machine-wide packages.

$ErrorActionPreference = "Stop"

function Have($cmd) { [bool](Get-Command $cmd -ErrorAction SilentlyContinue) }

function Install-Winget($id, $probe) {
    if (Have $probe) {
        Write-Host "[ok]   $probe already installed"
        return
    }
    Write-Host "[install] $id"
    winget install --id $id --silent --accept-source-agreements --accept-package-agreements
}

# --- system tools ---------------------------------------------------------
Install-Winget "Docker.DockerDesktop"   "docker"
Install-Winget "GoLang.Go"              "go"
Install-Winget "OpenJS.NodeJS.LTS"      "node"
Install-Winget "Python.Python.3.12"     "python"
Install-Winget "ezwinports.make"        "make"
Install-Winget "Git.Git"                "git"

# Refresh PATH so the rest of this script can see the freshly installed tools.
$env:Path = [Environment]::GetEnvironmentVariable("Path", "Machine") + ";" +
            [Environment]::GetEnvironmentVariable("Path", "User")

# --- ansible (via pipx) ---------------------------------------------------
if (-not (Have "pipx")) {
    Write-Host "[install] pipx"
    python -m pip install --user pipx
    python -m pipx ensurepath
    $env:Path = [Environment]::GetEnvironmentVariable("Path", "User") + ";" + $env:Path
}
if (-not (Have "ansible-playbook")) {
    Write-Host "[install] ansible"
    pipx install --include-deps ansible
}

# --- frontend deps --------------------------------------------------------
$root = Resolve-Path "$PSScriptRoot\.."

# --- sibling repos --------------------------------------------------------
$org = "https://github.com/trading-core"
$repos = @("trading-backend", "trading-frontend", "trading-formation", "integration-tests")
foreach ($r in $repos) {
    $dest = Join-Path $root $r
    if (Test-Path (Join-Path $dest ".git")) {
        Write-Host "[ok]   $r already cloned"
    } else {
        Write-Host "[clone] $org/$r"
        git clone "$org/$r.git" $dest
    }
}

$frontend = Join-Path $root "trading-frontend"
if (Test-Path (Join-Path $frontend "package.json")) {
    Write-Host "[npm] installing frontend deps"
    Push-Location $frontend
    npm install
    Pop-Location
}

# --- integration-tests deps ----------------------------------------------
$itests = Join-Path $root "integration-tests"
if (Test-Path (Join-Path $itests "pyproject.toml")) {
    Write-Host "[pip] installing integration-tests deps"
    Push-Location $itests
    python -m pip install --user -e .
    Pop-Location
}

# --- backend deps ---------------------------------------------------------
$backend = Join-Path $root "trading-backend"
if (Test-Path (Join-Path $backend "go.mod")) {
    Write-Host "[go] downloading backend modules"
    Push-Location $backend
    go mod download
    Pop-Location
}

Write-Host ""
Write-Host "Done. Open a new shell so PATH updates take effect, then start Docker Desktop."
