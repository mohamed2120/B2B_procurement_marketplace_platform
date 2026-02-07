# Service Version Tracking System

## Overview

The platform now includes a version tracking system that:
- ✅ Tracks versions for each service
- ✅ Compares source code versions with container versions
- ✅ Only rebuilds services that have changed
- ✅ Saves build time and resources

## How It Works

### Version Calculation

Each service's version is calculated using (in priority order):

1. **VERSION file** - If `services/<service>/VERSION` exists, uses that
2. **Git commit hash** - Last commit that modified the service directory
3. **File hash** - SHA256 hash of `main.go` and `Dockerfile`
4. **Modification time** - Last modified timestamp of service files
5. **Fallback** - Current timestamp

### Version Storage

Versions are stored in `.service-versions.json`:
```json
{
  "identity-service": {
    "source": "a1b2c3d4e5f6",
    "container": "a1b2c3d4e5f6"
  },
  "company-service": {
    "source": "f6e5d4c3b2a1",
    "container": "a1b2c3d4e5f6"
  }
}
```

### Docker Labels

Each built container includes a label:
```dockerfile
LABEL service.version="${SERVICE_VERSION}"
```

This allows checking container versions without reading the versions file.

## Usage

### Automatic (Recommended)

Just run the start script:
```bash
./start.sh
```

The script will:
1. Check all service versions
2. Compare source vs container versions
3. Only rebuild changed services
4. Update version tracking after build

### Manual Version Check

```bash
# Check which services need rebuild
bash scripts/check-versions.sh

# Get version for a specific service
bash scripts/get-service-version.sh identity-service

# Update container versions after manual build
bash scripts/update-container-versions.sh
```

### Force Rebuild All

```bash
make dev-down
make dev-up
```

### Rebuild Specific Service

```bash
# Rebuild just one service
docker compose -f docker-compose.all.yml build identity-service
docker compose -f docker-compose.all.yml up -d identity-service
bash scripts/update-container-versions.sh
```

## Example Output

```
Checking service versions...

  ✅ identity-service: SAME (a1b2c3d4)
  ✅ company-service: SAME (f6e5d4c3)
  🔄 catalog-service: CHANGED (a1b2c3d4 → b2c3d4e5)
  🔨 equipment-service: NEW (no version found)

Summary:
  New services: 1
  Changed services: 1
  Unchanged services: 2

🔨 Services that need rebuild:
  - equipment-service
  - catalog-service
```

## Setting Service Versions

### Option 1: VERSION File (Recommended)

Create `services/<service>/VERSION`:
```bash
echo "1.2.3" > services/identity-service/VERSION
```

### Option 2: Git-Based (Automatic)

The system automatically uses git commit hashes. Just commit your changes:
```bash
git add services/identity-service/
git commit -m "Update identity service"
```

### Option 3: Manual Version

Edit `.service-versions.json` directly (not recommended):
```json
{
  "identity-service": {
    "source": "1.2.3",
    "container": "1.2.3"
  }
}
```

## Updating Dockerfiles

All Dockerfiles should include version labels. To update all Dockerfiles:

```bash
bash scripts/update-all-dockerfiles.sh
```

This adds:
```dockerfile
ARG SERVICE_VERSION=unknown
LABEL service.version="${SERVICE_VERSION}"
```

## Benefits

1. **Faster builds** - Only rebuilds what changed
2. **Resource efficient** - Saves CPU/memory
3. **Deterministic** - Clear version tracking
4. **CI/CD friendly** - Easy to integrate into pipelines

## Troubleshooting

### Version not updating

1. Check if `.service-versions.json` exists
2. Verify Dockerfile has version labels
3. Run `bash scripts/update-container-versions.sh`

### All services showing as "NEW"

This is normal on first run. After first build, versions will be tracked.

### Version mismatch

If versions don't match:
1. Check git status - uncommitted changes affect version
2. Verify VERSION file format
3. Rebuild: `make dev-down && make dev-up`

### jq not installed

The scripts use `jq` for JSON manipulation. Install it:
```bash
# macOS
brew install jq

# Linux
sudo apt-get install jq

# Or use fallback (basic functionality)
```

## Integration with CI/CD

For CI/CD pipelines:

```bash
# Check if rebuild needed
if bash scripts/check-versions.sh; then
    # Rebuild changed services
    SERVICES=$(bash scripts/check-versions.sh --rebuild-list)
    docker compose build $SERVICES
    docker compose up -d $SERVICES
    bash scripts/update-container-versions.sh
fi
```

## Summary

The version tracking system makes builds smarter by:
- ✅ Only rebuilding what changed
- ✅ Tracking versions automatically
- ✅ Saving time and resources
- ✅ Providing clear feedback

Just run `./start.sh` and it handles everything automatically! 🚀
