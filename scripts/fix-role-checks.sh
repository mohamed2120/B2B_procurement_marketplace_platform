#!/bin/bash
# Fix role checks in all generated pages

# Buyer pages - requester or procurement_manager
for file in frontend/app/app/customer/{team,invite,roles,company,documents,addresses,policies,warehouse,warehouse/*,emergency,sell/*,auctions,auctions/*,reports,reports/*,parts,parts/*,equipment,equipment/*}/page.tsx; do
  if [ -f "$file" ]; then
    sed -i '' 's/\/\/ RBAC check.*/if (!(hasRole('\''requester'\'') || hasRole('\''procurement_manager'\''))) {\n      router.push('\''\/app'\'');\n      return;\n    }/' "$file"
  fi
done

# Buyer pages - requester only
for file in frontend/app/app/customer/pr/*/edit/page.tsx; do
  if [ -f "$file" ]; then
    sed -i '' 's/\/\/ RBAC check.*/if (!hasRole('\''requester'\'')) {\n      router.push('\''\/app'\'');\n      return;\n    }/' "$file"
  fi
done

# Buyer pages - procurement_manager only
for file in frontend/app/app/customer/rfq/*/quotes/page.tsx frontend/app/app/customer/rfq/*/award/page.tsx; do
  if [ -f "$file" ]; then
    sed -i '' 's/\/\/ RBAC check.*/if (!hasRole('\''procurement_manager'\'')) {\n      router.push('\''\/app'\'');\n      return;\n    }/' "$file"
  fi
done

# Supplier pages
for file in frontend/app/app/supplier/*/page.tsx frontend/app/app/supplier/*/*/page.tsx; do
  if [ -f "$file" ]; then
    sed -i '' 's/\/\/ RBAC check.*/if (!hasRole('\''supplier'\'')) {\n      router.push('\''\/app'\'');\n      return;\n    }/' "$file"
  fi
done

# Admin pages
for file in frontend/app/app/admin/*/page.tsx frontend/app/app/admin/*/*/page.tsx; do
  if [ -f "$file" ]; then
    sed -i '' 's/\/\/ RBAC check.*/if (!(hasRole('\''admin'\'') || hasRole('\''super_admin'\''))) {\n      router.push('\''\/app'\'');\n      return;\n    }/' "$file"
  fi
done

echo "✅ Fixed role checks in all pages"
