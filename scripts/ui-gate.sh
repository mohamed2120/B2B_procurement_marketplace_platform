#!/bin/bash

# UI Gate Script - Validates all required pages exist
# Fails CI if any required page is missing

set -e

REQUIRED_PAGES=(
  # Public
  "frontend/app/page.tsx"
  "frontend/app/how-it-works/page.tsx"
  "frontend/app/pricing/page.tsx"
  "frontend/app/register/page.tsx"
  "frontend/app/register/buyer/page.tsx"
  "frontend/app/register/supplier/page.tsx"
  "frontend/app/login/page.tsx"
  "frontend/app/forgot-password/page.tsx"
  "frontend/app/contact/page.tsx"
  "frontend/app/terms/page.tsx"
  "frontend/app/privacy/page.tsx"
  "frontend/app/search/page.tsx"
  
  # Buyer
  "frontend/app/app/customer/dashboard/page.tsx"
  "frontend/app/app/customer/pr/page.tsx"
  "frontend/app/app/customer/rfq/page.tsx"
  "frontend/app/app/customer/orders/page.tsx"
  "frontend/app/app/customer/shipments/page.tsx"
  "frontend/app/app/customer/parts/page.tsx"
  "frontend/app/app/customer/equipment/page.tsx"
  "frontend/app/app/customer/warehouse/page.tsx"
  "frontend/app/app/customer/team/page.tsx"
  "frontend/app/app/customer/company/page.tsx"
  "frontend/app/app/customer/reports/page.tsx"
  
  # Supplier
  "frontend/app/app/supplier/dashboard/page.tsx"
  "frontend/app/app/supplier/rfq/page.tsx"
  "frontend/app/app/supplier/quotes/page.tsx"
  "frontend/app/app/supplier/listings/page.tsx"
  "frontend/app/app/supplier/orders/page.tsx"
  "frontend/app/app/supplier/shipments/page.tsx"
  "frontend/app/app/supplier/store/page.tsx"
  "frontend/app/app/supplier/services/page.tsx"
  "frontend/app/app/supplier/inventory/page.tsx"
  "frontend/app/app/supplier/reports/page.tsx"
  
  # Admin
  "frontend/app/app/admin/dashboard/page.tsx"
  "frontend/app/app/admin/tenants/page.tsx"
  "frontend/app/app/admin/users/page.tsx"
  "frontend/app/app/admin/companies/page.tsx"
  "frontend/app/app/admin/roles-permissions/page.tsx"
  "frontend/app/app/admin/catalog-approvals/page.tsx"
  "frontend/app/app/admin/disputes/page.tsx"
  "frontend/app/app/admin/subscriptions/page.tsx"
  "frontend/app/app/admin/audit-logs/page.tsx"
  "frontend/app/app/admin/diagnostics/page.tsx"
)

MISSING_PAGES=()

echo "🔍 Checking for required pages..."
echo ""

for page in "${REQUIRED_PAGES[@]}"; do
  if [ ! -f "$page" ]; then
    MISSING_PAGES+=("$page")
    echo "❌ Missing: $page"
  else
    echo "✅ Found: $page"
  fi
done

echo ""

if [ ${#MISSING_PAGES[@]} -gt 0 ]; then
  echo "❌ UI GATE FAILED: ${#MISSING_PAGES[@]} required page(s) missing:"
  for page in "${MISSING_PAGES[@]}"; do
    echo "   - $page"
  done
  exit 1
else
  echo "✅ UI GATE PASSED: All required pages exist"
  exit 0
fi
