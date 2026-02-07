# Frontend Setup Guide

## Important: Choose ONE Method

The frontend can run in **two ways**, but **NOT both at the same time**:

1. **Docker Container** (default) - Frontend runs inside Docker
2. **Local Development** - Frontend runs on your machine with `npm run dev`

**⚠️ Running both will cause port conflicts on port 3002!**

## Option 1: Docker Container (Recommended for Full Stack)

This is the default when using `./start.sh` or `make dev-up`.

### Start Everything (including frontend in Docker):
```bash
./start.sh
# or
make dev-up
```

### Access:
- Frontend: http://localhost:3002
- All backend services are available

### Stop:
```bash
make dev-down
# or
docker compose -f docker-compose.all.yml down
```

## Option 2: Local Development (Recommended for Frontend Development)

Use this when you're actively developing the frontend and want hot reload.

### Step 1: Stop Docker Frontend Container
```bash
docker stop b2b-frontend
# or
docker compose -f docker-compose.all.yml stop frontend
```

### Step 2: Start Backend Services Only (without frontend)
```bash
# Start all services except frontend
docker compose -f docker-compose.all.yml up -d --scale frontend=0
# Or manually start each service
make dev-up
docker compose -f docker-compose.all.yml stop frontend
```

### Step 3: Start Frontend Locally
```bash
cd frontend
npm install  # First time only
npm run dev
```

### Access:
- Frontend: http://localhost:3002 (running locally)
- Backend services: http://localhost:8001, etc. (running in Docker)

### Benefits:
- ✅ Hot reload on code changes
- ✅ Faster development cycle
- ✅ Better debugging experience
- ✅ No need to rebuild Docker image

## Option 3: Exclude Frontend from Docker Compose

You can create a custom docker-compose override to exclude frontend:

### Create `docker-compose.override.yml`:
```yaml
services:
  frontend:
    profiles:
      - "disabled"
```

### Then start without frontend:
```bash
docker compose -f docker-compose.all.yml --profile disabled up -d
```

## Troubleshooting

### Port 3002 Already in Use

**Error:** `Error: listen EADDRINUSE: address already in use :::3002`

**Solution:**
1. Check what's using the port:
   ```bash
   lsof -nP -iTCP:3002 -sTCP:LISTEN
   ```

2. If it's Docker:
   ```bash
   docker stop b2b-frontend
   ```

3. If it's a local Node.js process:
   ```bash
   # Kill Node.js processes on port 3002
   lsof -nP -iTCP:3002 -sTCP:LISTEN -t | xargs kill -9
   # Or use the safe script
   cd frontend && npm run free-port:auto
   ```

### Frontend Not Loading

1. **Check if Docker container is running:**
   ```bash
   docker ps | grep frontend
   ```

2. **Check if local dev server is running:**
   ```bash
   ps aux | grep "next dev"
   ```

3. **Check logs:**
   ```bash
   # Docker logs
   docker compose -f docker-compose.all.yml logs frontend
   
   # Local dev (check terminal where you ran npm run dev)
   ```

### "Missing HTML Tags" Warning

This is a **false positive** in Next.js 14.2.5 dev mode. The layout file is correct. You can safely ignore this warning.

## Recommended Workflow

### For Backend Development:
```bash
./start.sh  # Everything in Docker
```

### For Frontend Development:
```bash
# Terminal 1: Start backend services
docker compose -f docker-compose.all.yml up -d
docker compose -f docker-compose.all.yml stop frontend

# Terminal 2: Start frontend locally
cd frontend
npm run dev
```

### For Full Stack Development:
```bash
# Use Docker for everything (slower but consistent)
./start.sh
```

## Summary

| Method | Command | Port | Hot Reload | Best For |
|--------|---------|------|------------|----------|
| Docker | `./start.sh` | 3002 | ❌ | Full stack, testing |
| Local | `cd frontend && npm run dev` | 3002 | ✅ | Frontend development |

**Remember:** Only use ONE method at a time!
