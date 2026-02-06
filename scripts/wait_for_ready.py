#!/usr/bin/env python3
"""
Wait for all services to be ready by polling /ready endpoints.
"""
import time
import sys
import requests
from typing import List, Tuple

SERVICES = [
    ("identity-service", 8001),
    ("company-service", 8002),
    ("catalog-service", 8003),
    ("equipment-service", 8004),
    ("marketplace-service", 8005),
    ("procurement-service", 8006),
    ("logistics-service", 8007),
    ("collaboration-service", 8008),
    ("notification-service", 8009),
    ("billing-service", 8010),
    ("virtual-warehouse-service", 8011),
    ("search-indexer-service", 8012),
    ("diagnostics-service", 8013),
]

MAX_WAIT = 300  # 5 minutes
POLL_INTERVAL = 5  # 5 seconds
TIMEOUT = 2  # 2 seconds per request

def check_ready(service_name: str, port: int) -> Tuple[bool, str]:
    """Check if a service is ready."""
    try:
        url = f"http://localhost:{port}/ready"
        response = requests.get(url, timeout=TIMEOUT)
        if response.status_code == 200:
            return True, "ready"
        else:
            return False, f"status {response.status_code}"
    except requests.exceptions.ConnectionError:
        return False, "connection refused"
    except requests.exceptions.Timeout:
        return False, "timeout"
    except Exception as e:
        return False, str(e)

def wait_for_all_ready():
    """Wait for all services to be ready."""
    start_time = time.time()
    ready_services = set()
    
    print(f"Waiting for {len(SERVICES)} services to be ready (max {MAX_WAIT}s)...")
    
    while time.time() - start_time < MAX_WAIT:
        all_ready = True
        status_lines = []
        
        for service_name, port in SERVICES:
            if service_name in ready_services:
                status_lines.append(f"  ✅ {service_name}: ready")
                continue
            
            is_ready, status = check_ready(service_name, port)
            if is_ready:
                ready_services.add(service_name)
                status_lines.append(f"  ✅ {service_name}: ready")
            else:
                all_ready = False
                status_lines.append(f"  ⏳ {service_name}: {status}")
        
        # Print status
        print("\r" + "\n".join(status_lines), end="", flush=True)
        
        if all_ready:
            print("\n")
            print(f"✅ All {len(SERVICES)} services are ready!")
            return True
        
        time.sleep(POLL_INTERVAL)
    
    print("\n")
    print("❌ Timeout waiting for services to be ready")
    print("\nNot ready services:")
    for service_name, port in SERVICES:
        if service_name not in ready_services:
            is_ready, status = check_ready(service_name, port)
            print(f"  ❌ {service_name} (port {port}): {status}")
    
    return False

if __name__ == "__main__":
    success = wait_for_all_ready()
    sys.exit(0 if success else 1)
