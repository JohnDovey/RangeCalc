<#
.SYNOPSIS
    android-build.ps1 - Helper for building RangeCalcAndroid from CLI (Windows)
    (Windows counterpart of android-build.sh)

.DESCRIPTION
    Configuration cache is enabled (see gradle.properties) for much faster
    repeated builds after the first configuration.
      Override for one run: .\android-build.ps1 --no-configuration-cache ...

.EXAMPLE
    .\android-build.ps1                  # assembleDebug (default)
    .\android-build.ps1 assembleRelease
    .\android-build.ps1 installDebug
    .\android-build.ps1 clean
    .\android-build.ps1 tasks            # list available tasks
    .\android-build.ps1 assembleDebug installDebug
    .\android-build.ps1 --no-configuration-cache assembleDebug
    .\android-build.ps1 --warning-mode all :app:assembleDebug
#>

$ErrorActionPreference = "Stop"

# Activate Android dev environment (quietly)
$envScript = "F:\Android\source-android-dev.ps1"
if (Test-Path $envScript) {
    . $envScript -Quiet
}

# Move to project root (directory containing this script)
Set-Location $PSScriptRoot

# Default to assembleDebug if no arguments. All args (tasks + gradle flags like
# --no-configuration-cache, --warning-mode all, -x test, etc.) are passed through.
$gradleArgs = $args
if ($gradleArgs.Count -eq 0) {
    $gradleArgs = @("assembleDebug")
}

Write-Host "Building RangeCalcAndroid on Windows (F: drive)..."
Write-Host "   Args: $($gradleArgs -join ' ')"
Write-Host "   ANDROID_HOME=$env:ANDROID_HOME"
Write-Host ""

& .\gradlew.bat @gradleArgs
if ($LASTEXITCODE -ne 0) {
    exit $LASTEXITCODE
}

Write-Host ""
Write-Host "Build finished"
