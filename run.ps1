<#
.SYNOPSIS
  Build and launch the whole MindSet Data solution: API server + edge agent + web UI.

.DESCRIPTION
  Builds bin/server.exe and bin/agent.exe, then opens three windows:
    1. API server   -> http://localhost:8080  (serves /api for the UI)
    2. Edge agent   -> connects to OPC-UA + MQTT (needs the simulator/broker up)
    3. Frontend dev -> http://localhost:5173  (the web UI, proxies /api to :8080)
  Finally it opens the UI in your browser.

.PARAMETER NoBuild
  Skip the Go build step and use the existing bin/*.exe.

.PARAMETER NoAgent
  Don't start the edge agent (use this when OPC-UA/MQTT aren't running).

.EXAMPLE
  .\run.ps1
  .\run.ps1 -NoAgent
  .\run.ps1 -NoBuild
#>
param(
  [switch]$NoBuild,
  [switch]$NoAgent
)

$ErrorActionPreference = 'Stop'
$root = $PSScriptRoot
Set-Location $root

Write-Host "== MindSet Data launcher ==" -ForegroundColor Cyan

# 1. Build the Go binaries
if (-not $NoBuild) {
  Write-Host "[build] Compiling server + agent..." -ForegroundColor Yellow
  & go build -o bin/server.exe ./cmd/server
  if ($LASTEXITCODE -ne 0) { throw "server build failed" }
  & go build -o bin/agent.exe ./cmd/agent
  if ($LASTEXITCODE -ne 0) { throw "agent build failed" }
  Write-Host "[build] OK -> bin/server.exe, bin/agent.exe" -ForegroundColor Green
}

# 2. Ensure frontend deps are installed
if (-not (Test-Path "$root/frontend/pipeline-builder/node_modules")) {
  Write-Host "[frontend] Installing dependencies (first run)..." -ForegroundColor Yellow
  Push-Location "$root/frontend/pipeline-builder"
  & npm install
  Pop-Location
}

# 3. Warn if no MQTT broker is listening on :1883 (needed for Run + auto-KG + agent)
$broker = Get-NetTCPConnection -LocalPort 1883 -State Listen -ErrorAction SilentlyContinue
if (-not $broker) {
  Write-Host "[warn] No MQTT broker detected on :1883 - 'Run', auto-KG and the agent need one." -ForegroundColor DarkYellow
}

# 4. Launch each component in its own window (so you can read logs / Ctrl+C individually)
Write-Host "[run] Starting API server on :8080..." -ForegroundColor Yellow
Start-Process powershell -ArgumentList '-NoExit', '-Command', "Set-Location '$root'; .\bin\server.exe"

if (-not $NoAgent) {
  Write-Host "[run] Starting edge agent..." -ForegroundColor Yellow
  Start-Process powershell -ArgumentList '-NoExit', '-Command', "Set-Location '$root'; .\bin\agent.exe"
} else {
  Write-Host "[run] Skipping edge agent (-NoAgent)." -ForegroundColor DarkGray
}

Write-Host "[run] Starting frontend dev server on :5173..." -ForegroundColor Yellow
Start-Process powershell -ArgumentList '-NoExit', '-Command', "Set-Location '$root/frontend/pipeline-builder'; npm run dev"

# 5. Give the dev server a moment, then open the browser
Start-Sleep -Seconds 4
Start-Process "http://localhost:5173"

Write-Host ""
Write-Host "All started. UI: http://localhost:5173   API: http://localhost:8080" -ForegroundColor Green
Write-Host "Close the spawned windows (or Ctrl+C in each) to stop." -ForegroundColor Green
