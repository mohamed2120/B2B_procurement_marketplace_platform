# Docker Desktop Automatic Recovery (macOS)

## Overview

The platform now includes automatic Docker Desktop recovery for macOS. When Docker becomes unhealthy, the system automatically attempts to restart it.

## When Recovery Triggers

Recovery is triggered when:
- `docker info` fails
- `docker ps` fails
- `docker compose` commands fail with:
  - `ECONNREFUSED .../com.docker.docker/Data/backend.sock`
  - `Cannot connect to the Docker daemon ... docker.sock`

## Recovery Process

### Step A: Soft Restart
1. Quit Docker Desktop application gracefully
2. Kill Docker Desktop processes:
   - `Docker Desktop`
   - `com.docker.backend`
   - `com.docker.supervisor`
3. Start Docker Desktop
4. Wait up to 180 seconds for Docker to become healthy (polling every 2 seconds)

### Step B: Hard Restart (if Step A fails)
1. Force kill all Docker processes:
   - `killall Docker`
   - `pkill -f com.docker`
2. Start Docker Desktop
3. Wait up to 180 seconds for Docker to become healthy

### If Both Steps Fail
- Prints clear error message
- Provides manual recovery instructions
- Exits with non-zero code

## Usage

### Safe Development Commands

```bash
# Start services with automatic Docker recovery
make dev-up-safe

# Stop services with automatic Docker recovery
make dev-down-safe
```

### Direct Script Usage

```bash
# Just recover Docker (if unhealthy)
bash scripts/docker-recover-macos.sh

# Safe startup
bash scripts/dev-up-safe.sh

# Safe shutdown
bash scripts/dev-down-safe.sh
```

## Integration

### dev-up-safe.sh
- Checks Docker health before starting
- Recovers Docker if unhealthy
- Checks for Docker compose connection errors
- Recovers and retries if needed
- Then proceeds with normal `make dev-up`

### dev-down-safe.sh
- Attempts to stop services
- If Docker connection error detected:
  - Recovers Docker
  - Retries stop operation once
- Handles both Docker and non-Docker errors gracefully

## Example Output

### Docker is Healthy
```
Checking Docker health...
✅ Docker is healthy. No recovery needed.
Docker is healthy. Proceeding with dev-up...
```

### Docker Recovery Needed
```
Checking Docker health...
⚠️  Docker is unhealthy. Attempting recovery...
Step A: Soft restart of Docker Desktop...
  Quitting Docker Desktop application...
  Killing Docker Desktop processes...
  Starting Docker Desktop...
  Waiting for Docker to become healthy...
✅ Docker recovered successfully (soft restart)
```

### Recovery Failed
```
❌ Docker Desktop failed to recover after soft and hard restart attempts.

Please try one of the following:
  1. Restart your Mac
  2. Manually restart Docker Desktop from Applications
  3. Reinstall Docker Desktop

To manually restart Docker Desktop:
  open -a Docker
```

## Platform Detection

The recovery script automatically detects macOS:
- On macOS (Darwin): Runs full recovery process
- On other platforms: Skips recovery (exits 0)

## Acceptance Criteria

✅ **`make dev-up-safe` automatically recovers Docker and continues**
- If Docker is unhealthy, recovery runs automatically
- After recovery, `make dev-up` proceeds normally

✅ **`make dev-down-safe` can stop stack even if Docker Desktop UI is broken**
- Detects Docker connection errors
- Recovers Docker if needed
- Retries stop operation after recovery
- Handles cases where recovery fails but still attempts cleanup

## Safety Features

- **Non-destructive**: Only kills Docker processes, never system processes
- **Graceful**: Tries soft restart before hard restart
- **Timeout protection**: Won't wait indefinitely (180s max per step)
- **Clear feedback**: Shows what's happening at each step
- **Platform-aware**: Only runs on macOS

## Troubleshooting

### Recovery keeps failing
1. Manually quit Docker Desktop from menu bar
2. Wait 10 seconds
3. Start Docker Desktop: `open -a Docker`
4. Wait 30 seconds
5. Try again: `make dev-up-safe`

### Docker Desktop won't start
1. Check Activity Monitor for stuck Docker processes
2. Force quit all Docker processes
3. Restart Mac
4. If still failing, reinstall Docker Desktop

### Recovery script not working
1. Verify you're on macOS: `uname` should show "Darwin"
2. Check script permissions: `chmod +x scripts/docker-recover-macos.sh`
3. Verify Docker Desktop is installed: `ls /Applications/Docker.app`

## Summary

The automatic Docker recovery system:
- ✅ Detects Docker health issues automatically
- ✅ Attempts soft and hard restarts
- ✅ Waits for Docker to become healthy
- ✅ Integrates seamlessly with dev-up/dev-down
- ✅ Works only on macOS (safe on other platforms)
- ✅ Provides clear feedback and error messages

Use `make dev-up-safe` and `make dev-down-safe` for automatic Docker recovery! 🚀
