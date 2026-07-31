<#
.SYNOPSIS
    Windows alternative to: make wire-addons ADDONS_DIR=... ADDONS="..."

.DESCRIPTION
    Links external addon modules into app/business/ using directory junctions
    (mklink /J, no admin required), then regenerates Wire dependency injection code.

.PARAMETER AddonsDir
    Path to the external addons repository.
    Default: D:\code\AAAA\nova-factory-addons-be (or $env:NOVA_ADDONS_DIR)

.PARAMETER Addons
    Space/comma separated addon names to enable. "all" or empty = all addons.
    Default: all (or $env:NOVA_ADDONS)

.EXAMPLE
    .\wire-addons.ps1
    .\wire-addons.ps1 -AddonsDir "D:\code\AAAA\nova-factory-addons-be" -Addons "shop erp"
#>
param(
    [string]$AddonsDir,
    [string]$Addons
)

$ErrorActionPreference = "Stop"

# ---- resolve defaults ----
if (-not $AddonsDir) { $AddonsDir = $env:NOVA_ADDONS_DIR }
if (-not $AddonsDir) { $AddonsDir = "D:\code\AAAA\nova-factory-addons-be" }
if (-not $Addons)    { $Addons    = $env:NOVA_ADDONS }
if (-not $Addons)    { $Addons    = "all" }

Write-Host "============================================================"
Write-Host " Addon Directory : $AddonsDir"
Write-Host " Enabled Addons  : $Addons"
Write-Host "============================================================"

# ---- Step 1: link addons ----
Write-Host "`n[1/2] Linking addons..."
& go run ./tools/addonsync -addons-dir $AddonsDir -enabled $Addons -link
if ($LASTEXITCODE -ne 0) { throw "addonsync failed" }

# ---- Step 2: regenerate Wire ----
Write-Host "`n[2/2] Generating Wire..."
Push-Location app
try {
    if ($Addons -eq "all") {
        & wire gen -tags="ai iot"
    } else {
        & wire gen -tags="ai iot $Addons"
    }
    if ($LASTEXITCODE -ne 0) { throw "wire gen failed" }
} finally {
    Pop-Location
}

Write-Host "`n============================================================"
Write-Host " Done."
Write-Host "============================================================"
