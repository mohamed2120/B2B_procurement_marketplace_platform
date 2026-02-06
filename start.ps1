# B2B Procurement Marketplace Platform - Start Script (PowerShell)
# This script starts all services, runs migrations, and seeds the database

Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "B2B Procurement Marketplace Platform" -ForegroundColor Cyan
Write-Host "Starting System..." -ForegroundColor Cyan
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host ""

# Check if Docker is running
try {
    docker info | Out-Null
} catch {
    Write-Host "Error: Docker is not running. Please start Docker Desktop first." -ForegroundColor Red
    exit 1
}

# Check if Make is available
$makeAvailable = $false
try {
    $null = Get-Command make -ErrorAction Stop
    $makeAvailable = $true
} catch {
    Write-Host "Warning: Make is not installed. Using docker-compose directly." -ForegroundColor Yellow
}

Write-Host "Step 1: Starting all services (infrastructure + backend + frontend)..." -ForegroundColor Green
if ($makeAvailable) {
    make up-all
} else {
    docker compose -f docker-compose.all.yml up -d --build
    Start-Sleep -Seconds 10
}

Write-Host ""
Write-Host "Step 2: Waiting for services to be ready..." -ForegroundColor Green
Start-Sleep -Seconds 15

Write-Host ""
Write-Host "Step 3: Running database migrations..." -ForegroundColor Green
if ($makeAvailable) {
    $migrateOutput = make migrate-all 2>&1 | Out-String
    if ($migrateOutput -match "already exists") {
        Write-Host "Some migrations already applied (this is normal)." -ForegroundColor Yellow
    } else {
        Write-Host "Migrations completed." -ForegroundColor Green
    }
} else {
    Write-Host "Please run migrations manually: make migrate-all" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "Step 4: Seeding database with demo data..." -ForegroundColor Green
Push-Location services/identity-service
go run cmd/seed/main.go
Pop-Location
Write-Host "Identity service seeded." -ForegroundColor Green

Write-Host ""
Write-Host "Step 5: Checking service health..." -ForegroundColor Green
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
Write-Host "Access the application:"
Write-Host "  Frontend:  http://localhost:3000"
Write-Host "  API Docs:  http://localhost:8001/health"
Write-Host ""
Write-Host "Demo Accounts (password: demo123456):"
Write-Host "  - Platform Admin: admin@demo.com"
Write-Host "  - Requester: buyer.requester@demo.com"
Write-Host "  - Procurement: buyer.procurement@demo.com"
Write-Host "  - Supplier: supplier@demo.com"
Write-Host ""
Write-Host "Useful commands:"
Write-Host "  make logs-all      - View all service logs"
Write-Host "  make health-check  - Check service health"
Write-Host "  make down-all      - Stop all services"
Write-Host ""
