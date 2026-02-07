# Port Conflict Prevention Guide

## Problem
When starting the frontend, you may encounter:
```
Error: listen EADDRINUSE: address already in use :::3002
```

This happens when port 3002 is already in use by another process (local Next.js dev server, Docker container, or another application).

## Solutions

### Automatic Solution (Recommended)

The project now includes automatic port conflict handling:

#### 1. **npm scripts** (Automatic)
```bash
cd frontend
npm run dev
```
This automatically checks and frees port 3002 before starting.

#### 2. **Auto-kill mode** (No prompts)
```bash
cd frontend
npm run free-port:auto
npm run dev
```

#### 3. **Start scripts** (Automatic)
```bash
# Bash
./start.sh

# PowerShell
./start.ps1
```
These scripts automatically free port 3002 before starting services.

### Manual Solutions

#### Option 1: Use the free-port script
```bash
# Bash
cd frontend
./scripts/free-port.sh 3002

# PowerShell
cd frontend
.\scripts\free-port.ps1 3002
```

#### Option 2: Manual port check and kill
```bash
# Check what's using the port
lsof -ti:3002  # macOS/Linux
netstat -ano | findstr :3002  # Windows

# Kill the process
kill -9 $(lsof -ti:3002)  # macOS/Linux
Stop-Process -Id <PID> -Force  # Windows PowerShell
```

#### Option 3: Stop Docker container
```bash
docker compose -f docker-compose.all.yml stop frontend
```

#### Option 4: Use a different port
```bash
# In frontend/package.json, change:
"dev": "next dev -p 3003"  # Use port 3003 instead

# Update docker-compose.all.yml port mapping:
ports:
  - "3003:3002"
```

## Available Commands

### Frontend npm scripts:
- `npm run dev` - Start dev server (auto-frees port)
- `npm run dev:safe` - Same as dev (explicit safe mode)
- `npm run free-port` - Free port 3002 (with confirmation)
- `npm run free-port:auto` - Free port 3002 (auto-kill, no prompts)

### Standalone scripts:
- `frontend/scripts/free-port.sh` - Bash script (interactive)
- `frontend/scripts/free-port.ps1` - PowerShell script (interactive)

### Auto-kill mode:
Set environment variable for non-interactive mode:
```bash
export AUTO_KILL=true
npm run free-port
```

## Best Practices

1. **Always use npm scripts** - They handle port conflicts automatically
2. **Check before starting** - Run `npm run free-port` if unsure
3. **Use Docker consistently** - Either always use Docker OR always use local dev, not both
4. **Stop containers properly** - Use `docker compose down` when done

## Troubleshooting

### Port still in use after killing?
1. Wait a few seconds for the port to be released
2. Check if Docker container is still running: `docker ps`
3. Check for zombie processes: `ps aux | grep node`
4. Restart your terminal/IDE

### Multiple Next.js instances?
If you have multiple terminal windows running `npm run dev`, kill all of them:
```bash
pkill -f "next dev"  # macOS/Linux
Get-Process | Where-Object {$_.ProcessName -like "*node*"} | Stop-Process  # Windows
```

### Docker vs Local conflict?
Decide on one approach:
- **Docker**: Use `./start.sh` or `docker compose up`
- **Local**: Use `npm run dev` in frontend directory
- Don't run both simultaneously on the same port

## Summary

The easiest solution is to use:
```bash
cd frontend
npm run dev
```

This automatically handles port conflicts for you! 🎉
