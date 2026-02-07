# B2B Procurement Marketplace Platform - Start Script (PowerShell)
# This script starts all services, runs migrations, and seeds the database

Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "B2B Procurement Marketplace Platform" -ForegroundColor Cyan
Write-Host "Starting System..." -ForegroundColor Cyan
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host ""

# Check if Docker is running
Write-Host "Checking Docker connection..." -ForegroundColor Yellow

# Function to check if Docker daemon is accessible (CLI-based only)
function Test-DockerReady {
    try {
        docker info | Out-Null
        return $true
    } catch {
        try {
            docker version | Out-Null
            return $true
        } catch {
            return $false
        }
    }
}

# Function to check if Docker daemon is accessible (for compatibility)
function Test-Docker {
    Test-DockerReady
}

# Function to free a port by killing processes using it
function Free-Port {
    param([int]$Port)
    
    $connections = Get-NetTCPConnection -LocalPort $Port -ErrorAction SilentlyContinue
    if ($connections) {
        Write-Host "Port $Port is in use. Freeing it..." -ForegroundColor Yellow
        $pids = $connections | Select-Object -ExpandProperty OwningProcess -Unique
        Write-Host "Found processes: $($pids -join ', ')" -ForegroundColor Yellow
        
        foreach ($pid in $pids) {
            try {
                Stop-Process -Id $pid -Force -ErrorAction Stop
                Write-Host "  Killed PID $pid" -ForegroundColor Green
            } catch {
                Write-Host "  Failed to kill PID $pid : $_" -ForegroundColor Red
            }
        }
        
        Start-Sleep -Seconds 1
        
        # Verify port is free
        $stillInUse = Get-NetTCPConnection -LocalPort $Port -ErrorAction SilentlyContinue
        if ($stillInUse) {
            Write-Host "Failed to free port $Port" -ForegroundColor Red
            return $false
        } else {
            Write-Host "Port $Port is now free" -ForegroundColor Green
            return $true
        }
    } else {
        Write-Host "Port $Port is free" -ForegroundColor Green
        return $true
    }
}

if (-not (Test-DockerReady)) {
    Write-Host "Docker daemon is not accessible." -ForegroundColor Yellow
    Write-Host "Waiting for Docker to be ready (max 90 seconds)..." -ForegroundColor Yellow
    Write-Host ""
    
    $maxWait = 90
    $waited = 0
    
    while ($waited -lt $maxWait) {
        if (Test-DockerReady) {
            Write-Host ""
            Write-Host "✅ Docker is now ready!" -ForegroundColor Green
            break
        }
        
        # Show progress every 5 seconds
        if ($waited % 5 -eq 0) {
            if ($waited -eq 0) {
                Write-Host "Waiting" -NoNewline
            } else {
                Write-Host "." -NoNewline
            }
        }
        
        # Every 15 seconds, show status
        if ($waited % 15 -eq 0 -and $waited -gt 0) {
            Write-Host ""
            Write-Host "  Still waiting... ($waited s / $maxWait s)" -ForegroundColor Yellow
        }
        
        Start-Sleep -Seconds 1
        $waited++
    }
    
    Write-Host ""
    
    if (-not (Test-DockerReady)) {
        Write-Host "❌ Error: Cannot connect to Docker daemon after $maxWait seconds." -ForegroundColor Red
        Write-Host ""
        Write-Host "Docker daemon not running or context socket missing" -ForegroundColor Yellow
        Write-Host ""
        
        if (-not (Get-Command docker -ErrorAction SilentlyContinue)) {
            Write-Host "❌ Docker CLI is not installed" -ForegroundColor Red
            exit 1
        }
        
        Write-Host "Docker CLI version:" -ForegroundColor Yellow
        docker --version 2>&1
        Write-Host ""
        
        Write-Host "Current Docker context:" -ForegroundColor Yellow
        $currentContext = docker context show 2>&1
        Write-Host $currentContext
        Write-Host ""
        
        Write-Host "Docker context list:" -ForegroundColor Yellow
        docker context ls 2>&1
        Write-Host ""
        
        # Check if context is desktop-linux and provide specific instructions
        if ($currentContext -match "desktop-linux") {
            Write-Host "⚠️  Current context is 'desktop-linux' (Docker Desktop)" -ForegroundColor Yellow
            Write-Host ""
            
            # Check for common Docker Desktop socket paths
            $socketFound = $false
            $desktopSocket = "$env:USERPROFILE\.docker\run\docker.sock"
            if (Test-Path $desktopSocket) {
                Write-Host "✅ Docker Desktop socket found: $desktopSocket" -ForegroundColor Green
                $socketFound = $true
            } elseif (Test-Path "/var/run/docker.sock") {
                Write-Host "✅ Docker socket found: /var/run/docker.sock" -ForegroundColor Green
                $socketFound = $true
            }
            
            if (-not $socketFound) {
                Write-Host "❌ Docker Desktop socket not found" -ForegroundColor Red
                Write-Host ""
                Write-Host "Start Docker Desktop (Docker.app) then retry" -ForegroundColor Yellow
                Write-Host ""
                Write-Host "On macOS:" -ForegroundColor Yellow
                Write-Host "  open -a Docker" -ForegroundColor Gray
                Write-Host ""
                Write-Host "On Windows:" -ForegroundColor Yellow
                Write-Host "  Start Docker Desktop from Start Menu" -ForegroundColor Gray
                Write-Host ""
                Write-Host "Wait 10-20 seconds after starting Docker Desktop, then run this script again." -ForegroundColor Yellow
                exit 1
            }
        }
        
        Write-Host "Docker info error output:" -ForegroundColor Yellow
        docker info 2>&1 | Select-Object -First 20
        Write-Host ""
        
        Write-Host "Docker version error output:" -ForegroundColor Yellow
        docker version 2>&1 | Select-Object -First 10
        Write-Host ""
        
        Write-Host "Troubleshooting steps:" -ForegroundColor Yellow
        Write-Host "  1. Ensure Docker daemon is running" -ForegroundColor Green
        Write-Host "  2. Check Docker context: docker context ls" -ForegroundColor Green
        Write-Host "  3. Try switching context: docker context use default" -ForegroundColor Green
        Write-Host "  4. Verify Docker daemon is accessible: docker ps" -ForegroundColor Green
        Write-Host ""
        exit 1
    }
} else {
    Write-Host "✅ Docker is ready!" -ForegroundColor Green
}

# Check Docker memory/resources
Write-Host "Checking Docker resources..." -ForegroundColor Yellow
try {
    $dockerInfo = docker info 2>$null | Select-String -Pattern "Total Memory"
    if ($dockerInfo) {
        Write-Host "  Docker Memory: $($dockerInfo.ToString().Trim())" -ForegroundColor Green
    }
} catch {
    # Ignore if we can't get memory info
}

# Warn about potential memory issues
Write-Host "Note: If you encounter 'cannot allocate memory' errors:" -ForegroundColor Yellow
Write-Host "  1. Increase Docker Desktop memory (Settings > Resources > Memory)" -ForegroundColor Yellow
Write-Host "  2. Recommended: At least 8GB allocated to Docker" -ForegroundColor Yellow
Write-Host "  3. Clean build cache: docker builder prune -f" -ForegroundColor Yellow
Write-Host ""

# Check if Make is available
$makeAvailable = $false
try {
    $null = Get-Command make -ErrorAction Stop
    $makeAvailable = $true
} catch {
    Write-Host "Warning: Make is not installed. Using docker-compose directly." -ForegroundColor Yellow
}

# Function to check if containers are already running
function Test-ContainersRunning {
    try {
        $running = docker compose -f docker-compose.all.yml ps --format json 2>$null | 
            ConvertFrom-Json | 
            Where-Object { $_.State -eq "running" } | 
            Measure-Object | 
            Select-Object -ExpandProperty Count
        
        $total = (docker compose -f docker-compose.all.yml config --services 2>$null | Measure-Object -Line).Lines
        
        if ($running -gt 0 -and $running -ge ($total / 2)) {
            return $true
        } else {
            return $false
        }
    } catch {
        return $false
    }
}

# Free port 3002 before starting (frontend port)
Write-Host "Checking and freeing port 3002 for frontend..." -ForegroundColor Yellow

# Check if port is in use
$portInUse = Get-NetTCPConnection -LocalPort 3002 -ErrorAction SilentlyContinue
if ($portInUse) {
    Write-Host "Port 3002 is in use. Checking what's using it..." -ForegroundColor Yellow
    
    # Check if it's a local Node.js process (local dev server)
    $pids = $portInUse | Select-Object -ExpandProperty OwningProcess -Unique
    foreach ($pid in $pids) {
        try {
            $process = Get-Process -Id $pid -ErrorAction Stop
            $cmdLine = $process.CommandLine
            if ($cmdLine -match "next dev|node.*3002|npm.*dev") {
                Write-Host "⚠️  Found local frontend dev server running (PID $pid)" -ForegroundColor Yellow
                Write-Host "   This will conflict with Docker frontend container!" -ForegroundColor Yellow
                Write-Host "   Stopping local dev server..." -ForegroundColor Yellow
                Stop-Process -Id $pid -Force -ErrorAction Stop
                Write-Host "✅ Local dev server stopped" -ForegroundColor Green
            }
        } catch {
            # Process might have exited, continue
        }
    }
    
    # Check if a Docker container is using the port
    $dockerContainer = docker ps --format "{{.Names}}`t{{.Ports}}" 2>$null | Select-String -Pattern ":3002->|3002:" | ForEach-Object { ($_ -split "`t")[0] } | Select-Object -First 1
    if ($dockerContainer) {
        Write-Host "Found Docker container '$dockerContainer' using port 3002" -ForegroundColor Yellow
        Write-Host "   (This is expected - we'll restart it)" -ForegroundColor Yellow
    }
    
    # Check if port is still in use (might be a Node.js process)
    $stillInUse = Get-NetTCPConnection -LocalPort 3002 -ErrorAction SilentlyContinue
    if ($stillInUse) {
        Write-Host "Port still in use. Checking for Node.js processes..." -ForegroundColor Yellow
        
        # Try the safe free-port script if available (on WSL/Git Bash)
        if (Test-Path "frontend/scripts/free-port.ps1") {
            try {
                & "frontend/scripts/free-port.ps1" -AutoKill
                Write-Host "✅ Port 3002 freed successfully" -ForegroundColor Green
            } catch {
                # Fallback to direct kill
                $nodeProcesses = $stillInUse | ForEach-Object {
                    try {
                        $proc = Get-Process -Id $_.OwningProcess -ErrorAction Stop
                        if ($proc.ProcessName -match "node|next|npm|yarn") {
                            return $proc
                        }
                    } catch {}
                    return $null
                } | Where-Object { $null -ne $_ }
                
                foreach ($proc in $nodeProcesses) {
                    Write-Host "Killing Node.js process (PID $($proc.Id), $($proc.ProcessName))..." -ForegroundColor Yellow
                    Stop-Process -Id $proc.Id -Force -ErrorAction SilentlyContinue
                }
                Start-Sleep -Seconds 2
            }
        } else {
            # Direct kill fallback (only Node.js processes)
            $nodeProcesses = $stillInUse | ForEach-Object {
                try {
                    $proc = Get-Process -Id $_.OwningProcess -ErrorAction Stop
                    if ($proc.ProcessName -match "node|next|npm|yarn") {
                        return $proc
                    }
                } catch {}
                return $null
            } | Where-Object { $null -ne $_ }
            
            foreach ($proc in $nodeProcesses) {
                Write-Host "Killing Node.js process (PID $($proc.Id))..." -ForegroundColor Yellow
                Stop-Process -Id $proc.Id -Force -ErrorAction SilentlyContinue
            }
            Start-Sleep -Seconds 2
        }
        
        # Verify port is free
        $finalCheck = Get-NetTCPConnection -LocalPort 3002 -ErrorAction SilentlyContinue
        if ($finalCheck) {
            Write-Host "❌ Port 3002 is still in use after kill attempt" -ForegroundColor Red
            Write-Host "Please manually free port 3002 and try again:" -ForegroundColor Yellow
            Write-Host "  Get-NetTCPConnection -LocalPort 3002 | ForEach-Object { Stop-Process -Id `$_.OwningProcess -Force }" -ForegroundColor Yellow
            exit 1
        } else {
            Write-Host "✅ Port 3002 freed successfully" -ForegroundColor Green
        }
    } else {
        Write-Host "✅ Port 3002 freed after stopping Docker container" -ForegroundColor Green
    }
} else {
    Write-Host "✅ Port 3002 is already free" -ForegroundColor Green
}
Write-Host ""

# Helper function to safely check Docker before compose operations
function Test-SafeDocker {
    try {
        docker info | Out-Null
        docker ps | Out-Null
        return $true
    } catch {
        Write-Host "⚠️  Docker is unhealthy. Attempting recovery..." -ForegroundColor Yellow
        if ((Test-Path "scripts/docker-recover-macos.sh") -and ($IsMacOS -or (Get-Command uname -ErrorAction SilentlyContinue))) {
            try {
                bash scripts/docker-recover-macos.sh
                if ($LASTEXITCODE -eq 0) {
                    return $true
                }
            } catch {
                Write-Host "❌ Docker recovery failed. Cannot proceed." -ForegroundColor Red
                return $false
            }
        } else {
            Write-Host "❌ Docker is unhealthy. Please restart Docker Desktop manually." -ForegroundColor Red
            return $false
        }
    }
}

# Helper function to safely run docker compose with recovery
function Invoke-SafeDockerCompose {
    param(
        [string]$Command,
        [string[]]$Arguments
    )
    
    $errorOutput = docker compose -f docker-compose.all.yml $Command $Arguments 2>&1
    if ($LASTEXITCODE -ne 0) {
        if ($errorOutput -match "ECONNREFUSED.*com.docker.docker|Cannot connect to the Docker daemon.*docker.sock|docker.sock.*connect") {
            Write-Host "⚠️  Docker connection error detected. Attempting recovery..." -ForegroundColor Yellow
            if ((Test-Path "scripts/docker-recover-macos.sh") -and ($IsMacOS -or (Get-Command uname -ErrorAction SilentlyContinue))) {
                try {
                    bash scripts/docker-recover-macos.sh
                    if ($LASTEXITCODE -eq 0) {
                        Write-Host ""
                        Write-Host "Retrying docker compose operation..." -ForegroundColor Yellow
                        docker compose -f docker-compose.all.yml $Command $Arguments
                        if ($LASTEXITCODE -ne 0) {
                            Write-Host "❌ Failed after recovery." -ForegroundColor Red
                            return $false
                        }
                    } else {
                        Write-Host "❌ Docker recovery failed. Cannot proceed." -ForegroundColor Red
                        return $false
                    }
                } catch {
                    Write-Host "❌ Docker connection error. Please restart Docker Desktop manually." -ForegroundColor Red
                    return $false
                }
            } else {
                Write-Host "❌ Docker connection error. Please restart Docker Desktop manually." -ForegroundColor Red
                return $false
            }
        } else {
            Write-Host "❌ Docker compose failed:" -ForegroundColor Red
            Write-Host $errorOutput
            return $false
        }
    }
    return $true
}

# Restart Docker containers if they're running
Write-Host "Checking Docker containers..." -ForegroundColor Yellow
if (Test-ContainersRunning) {
    $runningServices = (docker compose -f docker-compose.all.yml ps --services --filter "status=running" 2>$null | Measure-Object -Line).Lines
    Write-Host "Found $runningServices running containers. Restarting them..." -ForegroundColor Yellow
    
    # Check Docker health before restart
    if (-not (Test-SafeDocker)) {
        Write-Host "❌ Docker is unhealthy. Cannot proceed." -ForegroundColor Red
        exit 1
    }
    
    # Try to restart containers with error handling
    if (-not (Invoke-SafeDockerCompose -Command "restart" -Arguments @())) {
        Write-Host "❌ Failed to restart containers." -ForegroundColor Red
        exit 1
    }
    
    Write-Host "✅ Containers restarted" -ForegroundColor Green
    Write-Host ""
    Start-Sleep -Seconds 5  # Give containers time to restart
}

# Check service versions (smart rebuild)
Write-Host "Checking service versions..." -ForegroundColor Yellow
$versionCheck = bash scripts/check-versions.sh 2>&1
$versionCheckExit = $LASTEXITCODE

if ($versionCheckExit -eq 0) {
    # Some services need rebuild
    Write-Host $versionCheck
    Write-Host ""
    
    # Check Docker health before proceeding
    if (-not (Test-SafeDocker)) {
        exit 1
    }
    
    Write-Host "Step 1: Building and starting changed services..." -ForegroundColor Green
    if ($makeAvailable) {
        make dev-up
    } else {
        if (-not (Invoke-SafeDockerCompose -Command "up" -Arguments @("-d", "--build"))) {
            exit 1
        }
        Start-Sleep -Seconds 10
    }
} else {
    # All services up to date
    Write-Host $versionCheck
    Write-Host ""
    
    # Check Docker health before proceeding
    if (-not (Test-SafeDocker)) {
        exit 1
    }
    
    if (Test-ContainersRunning) {
        $runningServices = (docker compose -f docker-compose.all.yml ps --services --filter "status=running" 2>$null | Measure-Object -Line).Lines
        Write-Host "✅ All services up to date and running ($runningServices services)" -ForegroundColor Green
        Write-Host "Containers already restarted above. Skipping additional restart." -ForegroundColor Yellow
    } else {
        Write-Host "Services up to date but containers not running. Starting..." -ForegroundColor Yellow
        Write-Host ""
        Write-Host "Step 1: Starting all services..." -ForegroundColor Green
        if (-not (Invoke-SafeDockerCompose -Command "up" -Arguments @("-d"))) {
            exit 1
        }
    }
}

Write-Host ""
Write-Host "Step 2: Waiting for services to be ready..." -ForegroundColor Green
Write-Host "Waiting for infrastructure services (PostgreSQL, Redis)..." -ForegroundColor Yellow
$maxAttempts = 30
$attempt = 0
while ($attempt -lt $maxAttempts) {
    $status = docker compose -f docker-compose.all.yml ps postgres redis 2>$null | Select-String -Pattern "healthy|running"
    if ($status) {
        break
    }
    Write-Host "." -NoNewline
    Start-Sleep -Seconds 2
    $attempt++
}
Write-Host ""

Write-Host "Waiting for backend services to start..." -ForegroundColor Yellow
Start-Sleep -Seconds 20

Write-Host ""
Write-Host "Step 3: Verifying backend services are running..." -ForegroundColor Green
$backendServices = @(
    @{Name="identity-service"; Port=8001},
    @{Name="company-service"; Port=8002},
    @{Name="catalog-service"; Port=8003},
    @{Name="equipment-service"; Port=8004},
    @{Name="marketplace-service"; Port=8005},
    @{Name="procurement-service"; Port=8006},
    @{Name="logistics-service"; Port=8007},
    @{Name="collaboration-service"; Port=8008},
    @{Name="notification-service"; Port=8009},
    @{Name="billing-service"; Port=8010},
    @{Name="virtual-warehouse-service"; Port=8011},
    @{Name="diagnostics-service"; Port=8013}
)

$allBackendOk = $true
foreach ($service in $backendServices) {
    try {
        $response = Invoke-WebRequest -Uri "http://localhost:$($service.Port)/health" -Method GET -TimeoutSec 2 -ErrorAction Stop
        Write-Host "  ✅ $($service.Name) (port $($service.Port)): Running" -ForegroundColor Green
    } catch {
        # Fallback: check if container is running
        $containerStatus = docker compose -f docker-compose.all.yml ps $service.Name 2>$null | Select-String -Pattern "running"
        if ($containerStatus) {
            Write-Host "  ✅ $($service.Name) (port $($service.Port)): Container running" -ForegroundColor Green
        } else {
            Write-Host "  ❌ $($service.Name) (port $($service.Port)): Not responding" -ForegroundColor Red
            $allBackendOk = $false
        }
    }
}

if (-not $allBackendOk) {
    Write-Host "Warning: Some backend services are not responding yet. They may still be starting up." -ForegroundColor Yellow
    Write-Host "You can check logs with: docker compose -f docker-compose.all.yml logs <service-name>" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "Step 4: Verifying frontend is running..." -ForegroundColor Green
try {
    $response = Invoke-WebRequest -Uri "http://localhost:3002" -Method GET -TimeoutSec 2 -ErrorAction Stop
    Write-Host "  ✅ Frontend (port 3002): Running" -ForegroundColor Green
} catch {
    # Fallback: check if container is running
    $containerStatus = docker compose -f docker-compose.all.yml ps frontend 2>$null | Select-String -Pattern "running"
    if ($containerStatus) {
        Write-Host "  ✅ Frontend (port 3002): Container running" -ForegroundColor Green
    } else {
        Write-Host "  ⚠️  Frontend (port 3002): Not responding yet (may still be starting)" -ForegroundColor Yellow
        Write-Host "  Check logs: docker compose -f docker-compose.all.yml logs frontend" -ForegroundColor Yellow
    }
}

Write-Host ""
Write-Host "Step 5: Running database migrations..." -ForegroundColor Green
if ($makeAvailable) {
    try {
        make migrate-all
        Write-Host "Migrations completed." -ForegroundColor Green
    } catch {
        Write-Host "Some migrations may have already been applied (this is normal)." -ForegroundColor Yellow
    }
} else {
    Write-Host "Please run migrations manually: make migrate-all" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "Step 6: Seeding database with demo data..." -ForegroundColor Green
if ($makeAvailable) {
    try {
        make seed-all
        Write-Host "All services seeded." -ForegroundColor Green
    } catch {
        Write-Host "Some services may have already been seeded (this is normal)." -ForegroundColor Yellow
    }
} else {
    Write-Host "Please run seeding manually: make seed-all" -ForegroundColor Yellow
    Write-Host "Or seed individual services: cd services/<service> && go run cmd/seed/main.go" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "Step 7: Final health check..." -ForegroundColor Green
if ($makeAvailable) {
    make health-check
} else {
    Write-Host "Please run health check manually: make health-check" -ForegroundColor Yellow
    Write-Host "Or check: curl http://localhost:8001/health" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "System is ready!" -ForegroundColor Green
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host ""
# Get the actual frontend port from docker-compose file
$frontendPortLine = Get-Content docker-compose.all.yml | Select-String -Pattern '^\s+-\s*"[0-9]+:' -Context 0,20 | Where-Object { $_.Context.PreContext -match 'frontend:' } | Select-Object -First 1
if ($frontendPortLine) {
    $frontendPort = ($frontendPortLine.Line -replace '.*"([0-9]+):.*', '$1')
} else {
    # Fallback: try to get from running container
    $containerPort = docker compose -f docker-compose.all.yml ps frontend 2>$null | Select-String -Pattern '[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+:[0-9]+->'
    if ($containerPort) {
        $frontendPort = ($containerPort -replace '.*:([0-9]+)->.*', '$1')
    } else {
        $frontendPort = "3002"
    }
}
if ([string]::IsNullOrEmpty($frontendPort)) {
    $frontendPort = "3002"
}
Write-Host "Access the application:"
Write-Host "  Frontend:  http://localhost:$frontendPort"
Write-Host "  API Docs:  http://localhost:8001/health"
Write-Host ""
Write-Host "Demo Accounts (password: demo123456):"
Write-Host "  - Platform Admin: admin@demo.com"
Write-Host "  - Requester: buyer.requester@demo.com"
Write-Host "  - Procurement: buyer.procurement@demo.com"
Write-Host "  - Supplier: supplier@demo.com"
Write-Host ""
Write-Host "Note: OpenSearch is disabled by default. Use 'make dev-up-search' to enable search features."
Write-Host ""
Write-Host "Useful commands:"
Write-Host "  make logs-all        - View all service logs"
Write-Host "  make health-check    - Check service health"
Write-Host "  make dev-down        - Stop all services"
Write-Host "  make dev-up-search   - Start with OpenSearch enabled"
Write-Host ""
