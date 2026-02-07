# Safe Port 3002 Freeing

## Overview

The port freeing script (`scripts/free-port.sh`) is now **safe and deterministic** - it only kills Node.js/Next.js processes and never kills system processes or Docker containers.

## Safety Features

### 1. Only Checks LISTENING Processes

The script uses `lsof -nP -iTCP:3002 -sTCP:LISTEN` to only find processes that are actually listening on port 3002, not just any process that might have the port open.

### 2. Allowed Process Whitelist

Only these processes can be killed:
- `node`
- `next`
- `npm`
- `yarn`
- `pnpm`

### 3. Safety Checks

If port 3002 is occupied by a process NOT in the allowed list:
- ❌ **DO NOT kill it**
- ✅ Print clear error message
- ✅ Show diagnostics
- ✅ Exit with error code

### 4. Diagnostics

The script shows:
- Docker containers using port 3002
- All processes listening on port 3002
- Which processes are safe to kill
- Which processes are NOT safe to kill

## Usage

### Automatic (npm scripts)
```bash
cd frontend
npm run dev
```
This automatically runs `free-port` before starting.

### Manual
```bash
cd frontend
bash scripts/free-port.sh
```

### Auto-kill (no prompts)
```bash
cd frontend
AUTO_KILL=true bash scripts/free-port.sh
```

## Example Output

### Port is Free
```
Checking port 3002...
✅ Port 3002 is free
```

### Port Used by Safe Process (Node.js)
```
Checking port 3002...
⚠️  Port 3002 is in use

Diagnostics:
--- Docker containers using port 3002 ---
  (none)

--- Processes listening on port 3002 ---
COMMAND   PID   USER   FD   TYPE DEVICE SIZE/OFF NODE NAME
node     12345  user   23u  IPv4  ...      0t0  TCP *:3002 (LISTEN)

  ✅ PID 12345 (node) - safe to kill

Found safe processes to kill: 12345
Kill these processes? (y/N): y
Killing processes on port 3002...
  Killing PID 12345...
✅ Port 3002 is now free
```

### Port Used by Unsafe Process (Chrome/Docker)
```
Checking port 3002...
⚠️  Port 3002 is in use

Diagnostics:
--- Docker containers using port 3002 ---
NAMES          PORTS
b2b-frontend   0.0.0.0:3002->3002/tcp

--- Processes listening on port 3002 ---
COMMAND   PID   USER   FD   TYPE DEVICE SIZE/OFF NODE NAME
com.docker 123  user   45u  IPv4  ...      0t0  TCP *:3002 (LISTEN)

  ❌ PID 123 (com.docker) - NOT safe to kill

❌ ERROR: Port 3002 is occupied by processes that cannot be safely killed:
   - PID 123: com.docker

Please close that application or stop the container using port 3002.

To see what's using the port:
  lsof -nP -iTCP:3002 -sTCP:LISTEN
  docker ps --format 'table {{.Names}}\t{{.Ports}}' | grep 3002
```

## Safety Guarantees

✅ **Will NOT kill:**
- Chrome/Browser processes
- Docker containers
- System processes
- Any process not in the whitelist

✅ **Will kill:**
- Old `node` processes
- Old `next` dev servers
- Stuck `npm` processes
- `yarn` or `pnpm` processes

## Integration

### start.sh
The `start.sh` script uses `free_port_3002()` which:
- Uses the safe `frontend/scripts/free-port.sh` script
- Exits with error if port cannot be safely freed
- Shows diagnostics before failing

### npm scripts
The `package.json` scripts:
- `npm run dev` - Automatically frees port 3002
- `npm run free-port` - Manual port freeing
- `npm run free-port:auto` - Auto-kill mode

## Troubleshooting

### Port occupied by Docker container
```bash
# Stop the container
docker compose -f ../docker-compose.all.yml stop frontend

# Or remove it
docker compose -f ../docker-compose.all.yml down
```

### Port occupied by unknown process
```bash
# See what's using it
lsof -nP -iTCP:3002 -sTCP:LISTEN

# Manually kill if it's safe (be careful!)
kill -9 <PID>
```

### Port still in use after kill
```bash
# Check for remaining processes
lsof -nP -iTCP:3002 -sTCP:LISTEN

# Check Docker containers
docker ps --format "table {{.Names}}\t{{.Ports}}" | grep 3002
```

## Summary

The port freeing script is now:
- ✅ **Safe** - Only kills Node.js processes
- ✅ **Deterministic** - Always uses port 3002
- ✅ **Helpful** - Shows diagnostics
- ✅ **Protective** - Never kills system processes or Docker
