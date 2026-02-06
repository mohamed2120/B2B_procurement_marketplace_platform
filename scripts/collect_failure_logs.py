#!/usr/bin/env python3
"""
Collect logs for failing services.
"""
import subprocess
import requests
from typing import List, Tuple

SERVICES = [
    ("identity-service", "b2b-identity", 8001),
    ("company-service", "b2b-company", 8002),
    ("catalog-service", "b2b-catalog", 8003),
    ("equipment-service", "b2b-equipment", 8004),
    ("marketplace-service", "b2b-marketplace", 8005),
    ("procurement-service", "b2b-procurement", 8006),
    ("logistics-service", "b2b-logistics", 8007),
    ("collaboration-service", "b2b-collaboration", 8008),
    ("notification-service", "b2b-notification", 8009),
    ("billing-service", "b2b-billing", 8010),
    ("virtual-warehouse-service", "b2b-virtual-warehouse", 8011),
    ("search-indexer-service", "b2b-search-indexer", 8012),
    ("diagnostics-service", "b2b-diagnostics", 8013),
]

def check_service_ready(port: int) -> bool:
    """Check if service is ready."""
    try:
        response = requests.get(f"http://localhost:{port}/ready", timeout=2)
        return response.status_code == 200
    except:
        return False

def collect_logs(container_name: str, output_file: str):
    """Collect logs for a container."""
    try:
        result = subprocess.run(
            ["docker", "logs", container_name, "--tail", "100"],
            capture_output=True,
            text=True,
            timeout=10
        )
        with open(output_file, "w") as f:
            f.write(result.stdout)
            if result.stderr:
                f.write("\n\nSTDERR:\n")
                f.write(result.stderr)
        print(f"  ✅ Collected logs: {output_file}")
    except Exception as e:
        print(f"  ❌ Failed to collect logs for {container_name}: {e}")

def main():
    """Collect logs for all failing services."""
    print("Checking service status and collecting logs for failures...")
    print("")
    
    failed_services = []
    for service_name, container_name, port in SERVICES:
        if not check_service_ready(port):
            failed_services.append((service_name, container_name))
            print(f"❌ {service_name} is not ready")
    
    if not failed_services:
        print("✅ All services are ready - no logs to collect")
        return
    
    print("")
    print(f"Collecting logs for {len(failed_services)} failing service(s)...")
    print("")
    
    for service_name, container_name in failed_services:
        output_file = f"reports/logs/{service_name}.txt"
        print(f"Collecting logs for {service_name}...")
        collect_logs(container_name, output_file)
    
    print("")
    print(f"✅ Collected logs for {len(failed_services)} service(s)")
    print(f"   Logs saved to: reports/logs/")

if __name__ == "__main__":
    main()
