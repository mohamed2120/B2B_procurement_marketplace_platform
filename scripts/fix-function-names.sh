#!/bin/bash
# Fix invalid function names with spaces

fix_page() {
  local file=$1
  local component_name=$2
  
  if [ ! -f "$file" ]; then
    return
  fi
  
  # Replace function name and page title
  sed -i '' "s/export default function [^P]*Page()/export default function ${component_name}Page()/g" "$file"
  sed -i '' "s/<h1[^>]*>.*<\/h1>/<h1 className=\"text-3xl font-bold text-gray-900\">${component_name//_/ }<\/h1>/g" "$file"
}

# Fix pages with spaces in function names
fix_page "frontend/app/app/customer/parts/page.tsx" "PartsLibrary"
fix_page "frontend/app/app/customer/equipment/page.tsx" "Equipment"
fix_page "frontend/app/app/customer/warehouse/page.tsx" "Warehouse"
fix_page "frontend/app/app/customer/reports/page.tsx" "Reports"
fix_page "frontend/app/app/supplier/store/page.tsx" "Store"
fix_page "frontend/app/app/supplier/services/page.tsx" "Services"
fix_page "frontend/app/app/supplier/inventory/page.tsx" "Inventory"
fix_page "frontend/app/app/supplier/reports/page.tsx" "Reports"
fix_page "frontend/app/app/admin/companies/page.tsx" "Companies"

echo "✅ Fixed function names"
