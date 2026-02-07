#!/bin/bash

# Safe port freeing script for port 3002
# Only kills Node.js/Next.js processes, never kills system processes or Docker containers

PORT=3002

# Allowed process names that can be killed
ALLOWED_PROCESSES=("node" "next" "npm" "yarn" "pnpm")

echo "Checking port $PORT..."

# Check if port is in use (LISTENING only)
LISTENING_PIDS=$(lsof -nP -iTCP:$PORT -sTCP:LISTEN -t 2>/dev/null)

if [ -z "$LISTENING_PIDS" ]; then
    echo "✅ Port $PORT is free"
    exit 0
fi

echo "⚠️  Port $PORT is in use"
echo ""
echo "Diagnostics:"
echo "--- Docker containers using port $PORT ---"
docker ps --format "table {{.Names}}\t{{.Ports}}" 2>/dev/null | grep -E "3002|$PORT" || echo "  (none)"
echo ""
echo "--- Processes listening on port $PORT ---"
lsof -nP -iTCP:$PORT -sTCP:LISTEN 2>/dev/null || echo "  (none found)"
echo ""

# Check each process
SAFE_TO_KILL=()
UNSAFE_PROCESSES=()

for PID in $LISTENING_PIDS; do
    # Get process command name
    CMD=$(ps -p $PID -o comm= 2>/dev/null | xargs basename 2>/dev/null)
    
    if [ -z "$CMD" ]; then
        # Try alternative method
        CMD=$(ps -p $PID -o command= 2>/dev/null | awk '{print $1}' | xargs basename 2>/dev/null)
    fi
    
    if [ -z "$CMD" ]; then
        # Process might have exited, skip it
        continue
    fi
    
    # Check if command is in allowed list
    ALLOWED=false
    for allowed_cmd in "${ALLOWED_PROCESSES[@]}"; do
        if [ "$CMD" = "$allowed_cmd" ] || [[ "$CMD" == *"$allowed_cmd"* ]]; then
            ALLOWED=true
            break
        fi
    done
    
    if [ "$ALLOWED" = true ]; then
        SAFE_TO_KILL+=("$PID")
        echo "  ✅ PID $PID ($CMD) - safe to kill"
    else
        UNSAFE_PROCESSES+=("$PID:$CMD")
        echo "  ❌ PID $PID ($CMD) - NOT safe to kill"
    fi
done

echo ""

# If there are unsafe processes, don't proceed
if [ ${#UNSAFE_PROCESSES[@]} -gt 0 ]; then
    echo "❌ ERROR: Port $PORT is occupied by processes that cannot be safely killed:"
    for unsafe in "${UNSAFE_PROCESSES[@]}"; do
        IFS=':' read -r pid cmd <<< "$unsafe"
        echo "   - PID $pid: $cmd"
    done
    echo ""
    echo "Please close that application or stop the container using port $PORT."
    echo ""
    echo "To see what's using the port:"
    echo "  lsof -nP -iTCP:$PORT -sTCP:LISTEN"
    echo "  docker ps --format 'table {{.Names}}\t{{.Ports}}' | grep $PORT"
    echo ""
    exit 1
fi

# If no safe processes to kill, port might be free now
if [ ${#SAFE_TO_KILL[@]} -eq 0 ]; then
    # Double-check if port is still in use
    if lsof -nP -iTCP:$PORT -sTCP:LISTEN > /dev/null 2>&1; then
        echo "⚠️  Port $PORT is still in use but no safe processes found to kill"
        exit 1
    else
        echo "✅ Port $PORT is now free"
        exit 0
    fi
fi

# Kill safe processes
if [ "${AUTO_KILL:-false}" = "true" ]; then
    CONFIRM="y"
else
    echo "Found safe processes to kill: ${SAFE_TO_KILL[*]}"
    read -p "Kill these processes? (y/N): " CONFIRM
fi

if [[ "$CONFIRM" =~ ^[Yy]$ ]]; then
    echo "Killing processes on port $PORT..."
    for pid in "${SAFE_TO_KILL[@]}"; do
        echo "  Killing PID $pid..."
        kill -9 $pid 2>/dev/null || true
    done
    
    sleep 1
    
    # Verify port is free
    if lsof -nP -iTCP:$PORT -sTCP:LISTEN > /dev/null 2>&1; then
        echo "❌ Failed to free port $PORT"
        echo "Remaining processes:"
        lsof -nP -iTCP:$PORT -sTCP:LISTEN 2>/dev/null || echo "  (none found)"
        exit 1
    else
        echo "✅ Port $PORT is now free"
        exit 0
    fi
else
    echo "⚠️  Port $PORT is still in use. Exiting."
    exit 1
fi
