param(
  [switch]$Help,
  [switch]$Stop,
  [switch]$NoBrowser,
  [switch]$BackendOnly,
  [switch]$FrontendOnly,
  [int]$BackendPort = 8888,
  [int]$FrontendPort = 8080
)

$ErrorActionPreference = "Stop"

function Show-Help {
  Write-Host @"
Medical Info Platform dev launcher

Usage:
  powershell.exe -ExecutionPolicy Bypass -File scripts\start-dev.ps1 [options]

Options:
  -Help             Show this help message and exit.
  -Stop             Stop frontend/backend processes that occupy the configured ports, then exit.
  -NoBrowser        Start services without opening the browser.
  -BackendOnly      Restart backend only.
  -FrontendOnly     Restart frontend only.
  -BackendPort      Backend port. Default: 8888.
  -FrontendPort     Frontend port. Default: 8080.

Examples:
  powershell.exe -ExecutionPolicy Bypass -File scripts\start-dev.ps1
  powershell.exe -ExecutionPolicy Bypass -File scripts\start-dev.ps1 -NoBrowser
  powershell.exe -ExecutionPolicy Bypass -File scripts\start-dev.ps1 -Stop
  powershell.exe -ExecutionPolicy Bypass -File scripts\start-dev.ps1 -BackendOnly
  powershell.exe -ExecutionPolicy Bypass -File scripts\start-dev.ps1 -FrontendOnly
  powershell.exe -ExecutionPolicy Bypass -File scripts\start-dev.ps1 -BackendPort 8888 -FrontendPort 8080

Behavior:
  1. Stops existing processes listening on the selected ports.
  2. Starts backend with: go run .
  3. Starts frontend with: npm.cmd run dev
  4. Writes logs to runtime_logs/.
  5. Opens http://localhost:<FrontendPort>/login unless -NoBrowser is set.
"@
}

if ($Help) {
  Show-Help
  exit 0
}

$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$ProjectRoot = Split-Path -Parent $ScriptDir
$BackendDir = Join-Path $ProjectRoot "backend"
$FrontendDir = Join-Path $ProjectRoot "frontend"
$LogDir = Join-Path $ProjectRoot "runtime_logs"

function Write-Step {
  param([string]$Message)
  Write-Host "[dev-start] $Message" -ForegroundColor Cyan
}

function Ensure-Directory {
  param([string]$Path)
  if (-not (Test-Path -LiteralPath $Path)) {
    New-Item -ItemType Directory -Path $Path | Out-Null
  }
}

function Get-PortProcessIds {
  param([int]$Port)

  $processIds = New-Object System.Collections.Generic.List[int]
  $netstatLines = netstat -ano -p tcp | Select-String -Pattern "LISTENING"
  foreach ($line in $netstatLines) {
    $parts = ($line.ToString().Trim() -split "\s+")
    if ($parts.Count -lt 5) {
      continue
    }

    $localAddress = $parts[1]
    $state = $parts[3]
    $processIdText = $parts[4]

    if ($state -ne "LISTENING") {
      continue
    }
    if (-not ($localAddress -match "[:.]$Port$")) {
      continue
    }
    if ($processIdText -match "^\d+$") {
      $processIds.Add([int]$processIdText)
    }
  }

  return $processIds | Sort-Object -Unique
}

function Stop-PortProcess {
  param(
    [int]$Port,
    [string]$Name
  )

  $processIds = @(Get-PortProcessIds -Port $Port)
  if ($processIds.Count -eq 0) {
    Write-Step "$Name port $Port is free."
    return
  }

  foreach ($processId in $processIds) {
    try {
      $process = Get-Process -Id $processId -ErrorAction Stop
      Write-Step "Stopping $Name process $($process.ProcessName)($processId) on port $Port..."
      Stop-Process -Id $processId -Force -ErrorAction Stop
    } catch {
      Write-Warning "Failed to stop process $processId on port ${Port}: $($_.Exception.Message)"
    }
  }
}

function Wait-Port {
  param(
    [int]$Port,
    [string]$Name,
    [int]$TimeoutSeconds = 45
  )

  $deadline = (Get-Date).AddSeconds($TimeoutSeconds)
  while ((Get-Date) -lt $deadline) {
    try {
      $client = New-Object System.Net.Sockets.TcpClient
      $asyncResult = $client.BeginConnect("127.0.0.1", $Port, $null, $null)
      if ($asyncResult.AsyncWaitHandle.WaitOne(1000)) {
        $client.EndConnect($asyncResult)
        $client.Close()
        Write-Step "$Name is listening on port $Port."
        return $true
      }
      $client.Close()
    } catch {
      Start-Sleep -Milliseconds 700
    }
  }

  Write-Warning "$Name did not become ready on port $Port within $TimeoutSeconds seconds."
  return $false
}

function Start-DevProcess {
  param(
    [string]$Name,
    [string]$WorkingDirectory,
    [string]$Command,
    [string]$LogPath
  )

  $escapedWorkingDirectory = $WorkingDirectory.Replace("'", "''")
  $escapedLogPath = $LogPath.Replace("'", "''")
  $innerCommand = "Set-Location '$escapedWorkingDirectory'; `$Host.UI.RawUI.WindowTitle = '$Name'; $Command 2>&1 | Tee-Object -FilePath '$escapedLogPath' -Append"

  Write-Step "Starting $Name..."
  Start-Process powershell.exe -WindowStyle Normal -ArgumentList @(
    "-NoExit",
    "-ExecutionPolicy", "Bypass",
    "-Command", $innerCommand
  ) | Out-Null
}

$startBackend = -not $FrontendOnly
$startFrontend = -not $BackendOnly

if ($BackendOnly -and $FrontendOnly) {
  $startBackend = $true
  $startFrontend = $true
}

Ensure-Directory -Path $LogDir

if ($startBackend) {
  Stop-PortProcess -Port $BackendPort -Name "backend"
}
if ($startFrontend) {
  Stop-PortProcess -Port $FrontendPort -Name "frontend"
}

if ($Stop) {
  Write-Step "Stop requested. Done."
  exit 0
}

$timestamp = Get-Date -Format "yyyyMMdd_HHmmss"

if ($startBackend) {
  $backendLog = Join-Path $LogDir "backend_$timestamp.log"
  Start-DevProcess -Name "medical-info backend" -WorkingDirectory $BackendDir -Command "`$env:SERVER_PORT='$BackendPort'; go run ." -LogPath $backendLog
  $backendReady = Wait-Port -Port $BackendPort -Name "backend" -TimeoutSeconds 90
  if (-not $backendReady -and $startFrontend) {
    Write-Warning "Backend is not ready. Frontend will not be started to avoid Vite proxy ECONNREFUSED errors."
    exit 1
  }
}

if ($startFrontend) {
  $frontendLog = Join-Path $LogDir "frontend_$timestamp.log"
  Start-DevProcess -Name "medical-info frontend" -WorkingDirectory $FrontendDir -Command "cmd /c npm.cmd run dev" -LogPath $frontendLog
  Wait-Port -Port $FrontendPort -Name "frontend" | Out-Null
}

if (-not $NoBrowser -and $startFrontend) {
  Write-Step "Opening login page..."
  Start-Process "http://localhost:$FrontendPort/login"
}

Write-Step "Done."
