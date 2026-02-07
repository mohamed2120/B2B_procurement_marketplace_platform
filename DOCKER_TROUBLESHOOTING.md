# Docker Troubleshooting Guide

## Common Issues

### Issue: "Cannot connect to Docker daemon after waiting"

This error occurs when Docker Desktop is not running or not fully started.

#### Quick Fixes:

1. **Start Docker Desktop**
   - Open Docker Desktop application
   - Wait for it to fully start (look for Docker icon in menu bar/system tray)
   - Wait 10-20 seconds after Docker Desktop shows "Running"

2. **Check Docker Desktop Status**
   - Open Docker Desktop
   - Check if it shows "Running" status
   - Look for any error messages

3. **Restart Docker Desktop**
   - Quit Docker Desktop completely
   - Wait 5 seconds
   - Start Docker Desktop again
   - Wait for it to fully initialize

4. **Check Docker Desktop Permissions**
   - macOS: System Settings > Privacy & Security > Allow Docker Desktop
   - Windows: Check Windows Defender/antivirus settings

#### Detailed Diagnostics:

The start scripts now provide detailed diagnostics:

```bash
./start.sh
```

This will show:
- ✅/❌ Docker Desktop process status
- ✅/❌ Docker socket existence
- ✅/❌ Docker CLI installation
- Docker version information

#### Manual Checks:

**macOS/Linux:**
```bash
# Check if Docker Desktop process is running
pgrep -f "Docker Desktop"

# Check Docker socket
ls -la /var/run/docker.sock
ls -la ~/.docker/run/docker.sock

# Test Docker connection
docker ps
```

**Windows (PowerShell):**
```powershell
# Check if Docker Desktop process is running
Get-Process -Name "Docker Desktop" -ErrorAction SilentlyContinue

# Test Docker connection
docker ps
```

#### If Docker Desktop Won't Start:

1. **Check System Requirements**
   - macOS: macOS 10.15 or later
   - Windows: Windows 10 64-bit or later with WSL 2

2. **Check Available Resources**
   - Docker Desktop needs at least 4GB RAM
   - Recommended: 8GB+ RAM allocated to Docker

3. **Check for Conflicts**
   - Close other virtualization software (VMware, VirtualBox, etc.)
   - Check if Hyper-V is enabled (Windows)

4. **Reinstall Docker Desktop**
   - Download latest version from docker.com
   - Uninstall current version
   - Install fresh copy

#### Alternative: Use Docker Without Desktop

If Docker Desktop is problematic, you can use Docker Engine directly:

**macOS:**
```bash
brew install docker
brew install docker-machine
```

**Linux:**
```bash
# Install Docker Engine
sudo apt-get update
sudo apt-get install docker.io
sudo systemctl start docker
sudo systemctl enable docker
```

### Issue: "Docker daemon is starting" (Long Wait)

The script waits up to 90 seconds for Docker to be ready. If it takes longer:

1. **Check Docker Desktop Logs**
   - macOS: `~/Library/Containers/com.docker.docker/Data/log/`
   - Windows: Check Docker Desktop logs in the app

2. **Increase Wait Time**
   - Edit `start.sh` or `start.ps1`
   - Change `max_wait=90` to a higher value

3. **Start Docker Desktop Manually First**
   - Open Docker Desktop
   - Wait until it's fully running
   - Then run `./start.sh`

### Issue: "Permission Denied" Errors

**macOS/Linux:**
```bash
# Add your user to docker group
sudo usermod -aG docker $USER
# Log out and log back in
```

**Windows:**
- Run PowerShell as Administrator
- Or ensure your user has Docker Desktop permissions

### Issue: Docker Desktop Keeps Crashing

1. **Check System Resources**
   - Free up RAM/disk space
   - Close unnecessary applications

2. **Reset Docker Desktop**
   - Docker Desktop > Settings > Troubleshoot > Reset to factory defaults

3. **Check Logs**
   - Docker Desktop > Settings > Troubleshoot > View logs

### Getting Help

If issues persist:

1. Check Docker Desktop documentation: https://docs.docker.com/desktop/
2. Check Docker Desktop logs for specific errors
3. Search Docker Desktop issues on GitHub
4. Contact Docker support

## Prevention Tips

1. **Always start Docker Desktop before running scripts**
2. **Wait for Docker Desktop to fully initialize** (10-20 seconds)
3. **Keep Docker Desktop updated** to latest version
4. **Allocate sufficient resources** to Docker Desktop (8GB+ RAM recommended)
5. **Don't run multiple Docker instances** simultaneously
