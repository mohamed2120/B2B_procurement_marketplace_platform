# Smart Container Detection

## Overview

The start scripts and Makefile now intelligently detect if containers are already running and skip unnecessary rebuilds/restarts.

## How It Works

### Detection Logic

1. **Checks running containers**: Counts how many containers are in "running" state
2. **Compares with total**: If at least 50% of services are running, assumes system is up
3. **Skips build/start**: If containers are running, skips the build and start process
4. **Continues with verification**: Still runs health checks and verification steps

### Benefits

- ✅ **Faster startup** - No unnecessary rebuilds when containers are already running
- ✅ **Saves resources** - Doesn't waste CPU/memory on rebuilds
- ✅ **Preserves state** - Keeps existing container state and data
- ✅ **Smart detection** - Only rebuilds if containers are actually stopped

## Usage

### Normal Start (Containers Not Running)

```bash
./start.sh
```

**Output:**
```
Checking if containers are already running...
Containers are not running. Starting services...

Step 1: Starting all services...
Building services...
Services started...
```

### Smart Start (Containers Already Running)

```bash
./start.sh
```

**Output:**
```
Checking if containers are already running...
✅ Containers are already running (12 services)
Skipping build/start. Using existing containers.
To rebuild, run: make dev-down && make dev-up

Step 2: Waiting for services to be ready...
Step 3: Verifying backend services are running...
...
```

### Force Rebuild

If you want to rebuild even when containers are running:

```bash
make dev-down    # Stop all containers
make dev-up      # Rebuild and start
```

Or:

```bash
docker compose -f docker-compose.all.yml down
./start.sh
```

## Implementation Details

### start.sh / start.ps1

- Checks if containers are running before calling `make dev-up`
- If running, skips the build/start step
- Still runs all verification and health checks

### Makefile (dev-up target)

- Checks container status before building
- Only builds if containers are not running
- Provides clear feedback about what's happening

## Detection Threshold

The system considers containers "running" if:
- At least 50% of configured services are in "running" state
- This prevents false positives from partial failures

## Manual Override

If you need to force a rebuild:

1. **Stop containers first:**
   ```bash
   make dev-down
   ```

2. **Then start:**
   ```bash
   ./start.sh
   ```

Or use Docker Compose directly:
```bash
docker compose -f docker-compose.all.yml down
docker compose -f docker-compose.all.yml up -d --build
```

## Troubleshooting

### Containers detected as running but services aren't responding

This can happen if containers are in a bad state. Force a restart:

```bash
make dev-down
./start.sh
```

### Want to always rebuild

Remove the smart detection by editing the scripts, or always run:

```bash
make dev-down && make dev-up
```

### Containers partially running

If some containers are running but not all, the script will:
- Detect that not all are running (< 50% threshold)
- Rebuild and restart everything
- This ensures consistency

## Summary

The smart container detection makes the start scripts more efficient by:
- ✅ Skipping unnecessary rebuilds
- ✅ Preserving running container state
- ✅ Still verifying everything is working
- ✅ Providing clear feedback

Just run `./start.sh` - it will automatically detect if containers are running and skip rebuilds when appropriate! 🚀
