# Safe port freeing script for port 3002 (PowerShell)
# Only kills Node.js/Next.js processes, never kills system processes or Docker containers

param(
    [int]$Port = 3002
)

# Allowed process names that can be killed
$allowedProcesses = @("node", "next", "npm", "yarn", "pnpm")

Write-Host "Checking port $Port..." -ForegroundColor Cyan

# Check if port is in use (LISTENING only)
$listeningConnections = Get-NetTCPConnection -LocalPort $Port -State Listen -ErrorAction SilentlyContinue

if (-not $listeningConnections) {
    Write-Host "✅ Port $Port is free" -ForegroundColor Green
    exit 0
}

Write-Host "⚠️  Port $Port is in use" -ForegroundColor Yellow
Write-Host ""
Write-Host "Diagnostics:" -ForegroundColor Yellow
Write-Host "--- Docker containers using port $Port ---" -ForegroundColor Cyan
try {
    docker ps --format "table {{.Names}}`t{{.Ports}}" 2>$null | Select-String -Pattern "$Port|3002" | ForEach-Object { Write-Host $_.Line }
} catch {
    Write-Host "  (docker command not available or no containers found)" -ForegroundColor Gray
}
Write-Host ""
Write-Host "--- Processes listening on port $Port ---" -ForegroundColor Cyan
$listeningConnections | ForEach-Object {
    $proc = Get-Process -Id $_.OwningProcess -ErrorAction SilentlyContinue
    if ($proc) {
        Write-Host "  PID $($_.OwningProcess): $($proc.ProcessName) - $($proc.Path)" -ForegroundColor Gray
    }
}
Write-Host ""

# Check each process
$safeToKill = @()
$unsafeProcesses = @()

foreach ($conn in $listeningConnections) {
    $pid = $conn.OwningProcess
    $proc = Get-Process -Id $pid -ErrorAction SilentlyContinue
    
    if (-not $proc) {
        continue
    }
    
    $cmdName = $proc.ProcessName.ToLower()
    $isAllowed = $false
    
    foreach ($allowed in $allowedProcesses) {
        if ($cmdName -eq $allowed -or $cmdName -like "*$allowed*") {
            $isAllowed = $true
            break
        }
    }
    
    if ($isAllowed) {
        $safeToKill += $pid
        Write-Host "  ✅ PID $pid ($($proc.ProcessName)) - safe to kill" -ForegroundColor Green
    } else {
        $unsafeProcesses += [PSCustomObject]@{PID = $pid; Name = $proc.ProcessName; Path = $proc.Path}
        Write-Host "  ❌ PID $pid ($($proc.ProcessName)) - NOT safe to kill" -ForegroundColor Red
    }
}

Write-Host ""

# If there are unsafe processes, don't proceed
if ($unsafeProcesses.Count -gt 0) {
    Write-Host "❌ ERROR: Port $Port is occupied by processes that cannot be safely killed:" -ForegroundColor Red
    foreach ($unsafe in $unsafeProcesses) {
        Write-Host "   - PID $($unsafe.PID): $($unsafe.Name)" -ForegroundColor Red
        if ($unsafe.Path) {
            Write-Host "     Path: $($unsafe.Path)" -ForegroundColor Gray
        }
    }
    Write-Host ""
    Write-Host "Please close that application or stop the container using port $Port." -ForegroundColor Yellow
    Write-Host ""
    Write-Host "To see what's using the port:" -ForegroundColor Yellow
    Write-Host "  Get-NetTCPConnection -LocalPort $Port -State Listen | Get-Process" -ForegroundColor Gray
    Write-Host "  docker ps --format 'table {{.Names}}`t{{.Ports}}' | Select-String '$Port'" -ForegroundColor Gray
    Write-Host ""
    exit 1
}

# If no safe processes to kill, port might be free now
if ($safeToKill.Count -eq 0) {
    $stillInUse = Get-NetTCPConnection -LocalPort $Port -State Listen -ErrorAction SilentlyContinue
    if ($stillInUse) {
        Write-Host "⚠️  Port $Port is still in use but no safe processes found to kill" -ForegroundColor Yellow
        exit 1
    } else {
        Write-Host "✅ Port $Port is now free" -ForegroundColor Green
        exit 0
    }
}

# Kill safe processes
$autoKill = $env:AUTO_KILL -eq "true"

if (-not $autoKill) {
    Write-Host "Found safe processes to kill: $($safeToKill -join ', ')" -ForegroundColor Yellow
    $confirm = Read-Host "Kill these processes? (y/N)"
} else {
    $confirm = "y"
}

if ($confirm -match "^[Yy]$") {
    Write-Host "Killing processes on port $Port..." -ForegroundColor Yellow
    foreach ($pid in $safeToKill) {
        try {
            Write-Host "  Killing PID $pid..." -ForegroundColor Yellow
            Stop-Process -Id $pid -Force -ErrorAction Stop
        } catch {
            Write-Host "  Failed to kill PID $pid : $_" -ForegroundColor Red
        }
    }
    
    Start-Sleep -Seconds 1
    
    # Verify port is free
    $stillInUse = Get-NetTCPConnection -LocalPort $Port -State Listen -ErrorAction SilentlyContinue
    if ($stillInUse) {
        Write-Host "❌ Failed to free port $Port" -ForegroundColor Red
        Write-Host "Remaining processes:" -ForegroundColor Yellow
        $stillInUse | ForEach-Object {
            $proc = Get-Process -Id $_.OwningProcess -ErrorAction SilentlyContinue
            if ($proc) {
                Write-Host "  PID $($_.OwningProcess): $($proc.ProcessName)" -ForegroundColor Gray
            }
        }
        exit 1
    } else {
        Write-Host "✅ Port $Port is now free" -ForegroundColor Green
        exit 0
    }
} else {
    Write-Host "⚠️  Port $Port is still in use. Exiting." -ForegroundColor Yellow
    exit 1
}
