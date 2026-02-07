#!/bin/bash
# Fix all generated pages with proper role checks

fix_page() {
  local file=$1
  local role_check=$2
  
  if [ ! -f "$file" ]; then
    return
  fi
  
  # Remove ROLE_CHECK placeholder
  sed -i '' '/ROLE_CHECK/d' "$file"
  
  # Add proper role check if provided
  if [ -n "$role_check" ]; then
    # Insert after useState declarations
    sed -i '' "/const \[error, setError\] = useState/a\\
\\
  useEffect(() => {\\
    if (!($role_check)) {\\
      router.push('/app');\\
      return;\\
    }\\
  }, []);" "$file"
  fi
}

# Fix buyer pages
fix_page "frontend/app/app/customer/team/page.tsx" "hasRole('requester') || hasRole('procurement_manager')"
fix_page "frontend/app/app/customer/invite/page.tsx" "hasRole('requester') || hasRole('procurement_manager')"
fix_page "frontend/app/app/customer/roles/page.tsx" "hasRole('requester') || hasRole('procurement_manager')"
fix_page "frontend/app/app/customer/company/page.tsx" "hasRole('requester') || hasRole('procurement_manager')"
fix_page "frontend/app/app/customer/documents/page.tsx" "hasRole('requester') || hasRole('procurement_manager')"
fix_page "frontend/app/app/customer/addresses/page.tsx" "hasRole('requester') || hasRole('procurement_manager')"
fix_page "frontend/app/app/customer/policies/page.tsx" "hasRole('requester') || hasRole('procurement_manager')"
fix_page "frontend/app/app/customer/pr/[id]/edit/page.tsx" "hasRole('requester')"
fix_page "frontend/app/app/customer/rfq/create/page.tsx" "hasRole('requester') || hasRole('procurement_manager')"
fix_page "frontend/app/app/customer/rfq/[id]/quotes/page.tsx" "hasRole('procurement_manager')"
fix_page "frontend/app/app/customer/rfq/[id]/award/page.tsx" "hasRole('procurement_manager')"

# Fix supplier pages
for file in frontend/app/app/supplier/{company,documents,team,invite,roles,store,listings,services,inventory,pricing,rfq,quotes,orders,shipments,ratings,performance,reports}*/page.tsx; do
  if [ -f "$file" ]; then
    fix_page "$file" "hasRole('supplier')"
  fi
done

# Fix admin pages
for file in frontend/app/app/admin/{companies,subdomains,catalog,listings,stores,rfqs,orders,shipments,chat-moderation,plans,payments,notifications}*/page.tsx; do
  if [ -f "$file" ]; then
    fix_page "$file" "hasRole('admin') || hasRole('super_admin')"
  fi
done

echo "✅ Fixed all pages"
