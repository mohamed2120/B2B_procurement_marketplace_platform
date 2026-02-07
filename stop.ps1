# B2B Procurement Marketplace Platform - Stop Script (PowerShell)
# This script stops all services and cleans up resources

Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "B2B Procurement Marketplace Platform" -ForegroundColor Cyan
Write-Host "Stopping System..." -ForegroundColor Cyan
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host ""

# Check if Docker is running
Write-Host "Checking Docker connection..." -ForegroundColor Yellow
try {
    docker info | Out-Null
    Write-Host "✅ Docker is accessible" -ForegroundColor Green
    $dockerAvailable = $true
} catch {
    Write-Host "⚠️  Docker is not accessible. Some cleanup steps may be skipped." -ForegroundColor Yellow
    $dockerAvailable = $false
}
Write-Host ""

# Stop local frontend dev server if running
Write-Host "Checking for local frontend dev server..." -ForegroundColor Yellow
$portInUse = Get-NetTCPConnection -LocalPort 3002 -ErrorAction SilentlyContinue
if ($portInUse) {
    $pids = $portInUse | Select-Object -ExpandProperty OwningProcess -Unique
    $localServerFound = $false
    
    foreach ($pid in $pids) {
        try {
            $process = Get-Process -Id $pid -ErrorAction Stop
            $cmdLine = $process.CommandLine
            if ($cmdLine -match "next dev|node.*3002|npm.*dev") {
                Write-Host "Found local frontend dev server (PID $pid). Stopping..." -ForegroundColor Yellow
                Stop-Process -Id $pid -Force -ErrorAction Stop
                Start-Sleep -Seconds 1
                Write-Host "✅ Local dev server stopped" -ForegroundColor Green
                $localServerFound = $true
            }
        } catch {
            # Process might have exited, continue
        }
    }
    
    if (-not $localServerFound) {
        Write-Host "Port 3002 is in use by non-Node.js process" -ForegroundColor Yellow
        Write-Host "   Skipping local dev server stop" -ForegroundColor Yellow
    }
} else {
    Write-Host "✅ No local frontend dev server running" -ForegroundColor Green
}
Write-Host ""

# Stop all Docker containers
Write-Host "Stopping all Docker containers..." -ForegroundColor Yellow
if ($dockerAvailable) {
    try {
        $runningJson = docker compose -f docker-compose.all.yml ps --format json 2>$null
        $runningCount = ($runningJson | ConvertFrom-Json | Where-Object { $_.State -eq "running" } | Measure-Object).Count
        
        if ($runningCount -gt 0) {
            Write-Host "Found $runningCount running container(s)" -ForegroundColor Yellow
            
            # Try graceful stop first
            try {
                docker compose -f docker-compose.all.yml stop 2>&1 | Out-Null
                Write-Host "✅ Containers stopped gracefully" -ForegroundColor Green
            } catch {
                # If graceful stop fails, try force stop
                Write-Host "Graceful stop failed. Attempting force stop..." -ForegroundColor Yellow
                docker compose -f docker-compose.all.yml down --remove-orphans 2>&1 | Out-Null
            }
            
            # Verify containers are stopped
            Start-Sleep -Seconds 2
            $stillRunningJson = docker compose -f docker-compose.all.yml ps --format json 2>$null
            $stillRunning = ($stillRunningJson | ConvertFrom-Json | Where-Object { $_.State -eq "running" } | Measure-Object).Count
            
            if ($stillRunning -gt 0) {
                Write-Host "⚠️  $stillRunning container(s) still running. Force stopping..." -ForegroundColor Yellow
                docker compose -f docker-compose.all.yml down --remove-orphans --timeout 5 2>&1 | Out-Null
            }
        } else {
            Write-Host "✅ No containers are running" -ForegroundColor Green
        }
    } catch {
        Write-Host "⚠️  Error stopping containers: $_" -ForegroundColor Yellow
    }
} else {
    Write-Host "⚠️  Docker is not accessible. Cannot stop containers." -ForegroundColor Yellow
}
Write-Host ""

# Optional: Remove containers (commented out by default)
# Uncomment the following lines if you want to remove containers on stop
# Write-Host "Removing containers..." -ForegroundColor Yellow
# docker compose -f docker-compose.all.yml down --remove-orphans 2>&1 | Out-Null
# Write-Host "✅ Containers removed" -ForegroundColor Green
# Write-Host ""

# Free port 3002
Write-Host "Verifying port 3002 is free..." -ForegroundColor Yellow
$stillInUse = Get-NetTCPConnection -LocalPort 3002 -ErrorAction SilentlyContinue
if ($stillInUse) {
    $nodeProcesses = $stillInUse | ForEach-Object {
        try {
            $proc = Get-Process -Id $_.OwningProcess -ErrorAction Stop
            if ($proc.ProcessName -match "node|next|npm|yarn") {
                return $proc
            }
        } catch {}
        return $null
    } | Where-Object { $null -ne $_ }
    
    if ($nodeProcesses) {
        Write-Host "Port 3002 is still in use. Cleaning up..." -ForegroundColor Yellow
        foreach ($proc in $nodeProcesses) {
            Write-Host "Killing Node.js process (PID $($proc.Id))..." -ForegroundColor Yellow
            Stop-Process -Id $proc.Id -Force -ErrorAction SilentlyContinue
        }
        Start-Sleep -Seconds 1
    }
}

# Verify port is free
$finalCheck = Get-NetTCPConnection -LocalPort 3002 -ErrorAction SilentlyContinue
if ($finalCheck) {
    Write-Host "⚠️  Port 3002 is still in use" -ForegroundColor Yellow
    Write-Host "   Run manually: Get-NetTCPConnection -LocalPort 3002" -ForegroundColor Yellow
} else {
    Write-Host "✅ Port 3002 is free" -ForegroundColor Green
}
Write-Host ""

# Summary
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "System stopped successfully!" -ForegroundColor Green
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "All services have been stopped."
Write-Host ""
Write-Host "To start again, run:"
Write-Host "  .\start.ps1"
Write-Host ""
Write-Host "To remove containers and volumes, run:"
Write-Host "  docker compose -f docker-compose.all.yml down -v"
Write-Host ""
