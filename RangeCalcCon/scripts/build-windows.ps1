# Build RangeCalcCon for Windows (amd64).
# Run from repo root or from scripts/:  .\scripts\build-windows.ps1
$ErrorActionPreference = "Stop"

$Root = Split-Path -Parent (Split-Path -Parent $MyInvocation.MyCommand.Path)
if (-not (Test-Path (Join-Path $Root "main.go"))) {
    $Root = Get-Location
}

$Version = "0.0.0"
$VersionFile = Join-Path $Root "VERSION"
if (Test-Path $VersionFile) {
    $Version = (Get-Content $VersionFile -Raw).Trim()
}

$Out = Join-Path $Root "rangecalccon.exe"
Write-Host "Building RangeCalcCon $Version → $Out"

Push-Location $Root
try {
    $env:GOOS = "windows"
    $env:GOARCH = "amd64"
    $env:CGO_ENABLED = "0"
    go build -ldflags "-s -w -X main.version=$Version" -o $Out .
    if ($LASTEXITCODE -ne 0) {
        throw "go build failed with exit code $LASTEXITCODE"
    }
} finally {
    Pop-Location
}

$item = Get-Item $Out
Write-Host ("OK: {0} ({1:N0} bytes)" -f $item.FullName, $item.Length)
