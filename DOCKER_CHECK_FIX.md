# Docker Readiness Check Fix

## Problem Fixed

The previous Docker check looked for "Docker Desktop" process and specific socket paths, causing false failures on macOS even when Docker was working.

## Solution

Replaced all Docker checks with **CLI-based only** checks that determine readiness by whether the Docker CLI can reach the daemon.

## Changes Made

### 1. Created `scripts/check-docker.sh`

A standalone script that:
- Uses `docker info` (preferred) or `docker version` to check readiness
- Retries up to 90 seconds with progress indicators
- Shows diagnostics on failure:
  - `docker context ls`
  - `docker context show`
  - `docker info` error output
  - `docker version` error output
- **Does NOT** check for Docker Desktop process
- **Does NOT** check socket paths

### 2. Updated `start.sh`

- Removed all "Docker Desktop" process checks
- Removed socket path checks
- Uses `check-docker.sh` script
- Only checks if CLI can reach daemon

### 3. Updated `start.ps1`

- Removed all "Docker Desktop" process checks
- Uses `docker info` or `docker version` only
- Shows diagnostics on failure

### 4. Updated Makefile Targets

All targets now check Docker readiness first:
- `make dev-up` - Checks Docker before starting
- `make dev-up-search` - Checks Docker before starting
- `make verify` - Checks Docker as Step 0/9
- `make up-all` - Checks Docker before starting
- `make reset-all` - Checks Docker before starting

## How It Works

### Check Function
```bash
check_docker_ready() {
    docker info > /dev/null 2>&1 || docker version > /dev/null 2>&1
}
```

### Wait Function
- Retries up to 90 seconds
- Shows progress every 5 seconds
- Shows status every 15 seconds
- No process or socket checks

### Diagnostics on Failure
If Docker is not ready after 90 seconds:
1. Shows Docker CLI version
2. Shows Docker context list
3. Shows current Docker context
4. Shows `docker info` error
5. Shows `docker version` error
6. Provides troubleshooting steps

## Acceptance Criteria

✅ **If `docker ps` works, then `make dev-up` must not fail**
✅ **If `docker ps` works, then `make verify` must not fail**
✅ **No false failures due to Docker Desktop process checks**
✅ **No false failures due to socket path checks**

## Usage

### Standalone Check
```bash
bash scripts/check-docker.sh
```

### In Scripts
```bash
source scripts/check-docker.sh
if check_docker_ready; then
    echo "Docker is ready"
else
    wait_for_docker 90 || show_docker_diagnostics
fi
```

### In Makefile
```makefile
@bash scripts/check-docker.sh > /dev/null 2>&1 || (echo "❌ Docker not ready" && exit 1)
```

## Benefits

1. **No false failures** - Only checks if CLI can reach daemon
2. **Works on all platforms** - macOS, Linux, Windows
3. **Works with all Docker setups** - Docker Desktop, Docker Engine, remote Docker
4. **Better diagnostics** - Shows context and error information
5. **Simpler logic** - No process or socket detection needed

## Testing

Test that Docker check works:
```bash
# Should pass if docker ps works
docker ps > /dev/null 2>&1 && bash scripts/check-docker.sh && echo "✅ PASS"
```

Test Makefile targets:
```bash
# Should not fail if docker ps works
docker ps > /dev/null 2>&1 && make -n dev-up && echo "✅ PASS"
docker ps > /dev/null 2>&1 && make -n verify && echo "✅ PASS"
```

## Summary

The Docker readiness check now:
- ✅ Only uses CLI connectivity (`docker info` / `docker version`)
- ✅ Does NOT check for Docker Desktop process
- ✅ Does NOT check socket paths
- ✅ Retries up to 90 seconds
- ✅ Shows helpful diagnostics on failure
- ✅ Works on all platforms and Docker setups
