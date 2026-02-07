#!/bin/bash

# Simple batch page generator
# This creates basic MVP pages for all missing routes

generate_page() {
  local route=$1
  local title=$2
  local roles=$3
  local file_path="frontend/app${route}/page.tsx"
  
  mkdir -p "$(dirname "$file_path")"
  
  # Build role check
  local role_check=""
  if [ -n "$roles" ]; then
    role_check="  useEffect(() => {
    if (!($roles)) {
      router.push('/app');
      return;
    }
  }, []);"
  fi
  
  cat > "$file_path" << 'PAGEEOF'
'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { apiClients } from '@/lib/api';
import { hasRole } from '@/lib/auth';
import Card from '@/components/ui/Card';
import Button from '@/components/ui/Button';

interface DataItem {
  id: string;
  [key: string]: any;
}

export default function PAGETITLEPage() {
  const router = useRouter();
  const [data, setData] = useState<DataItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
ROLE_CHECK

  useEffect(() => {
    fetchData();
  }, []);

  const fetchData = async () => {
    setLoading(true);
    setError('');
    // TODO: Implement API call
    setData([
      { id: '1', name: 'Sample Item 1', status: 'active' },
      { id: '2', name: 'Sample Item 2', status: 'pending' },
    ]);
    setLoading(false);
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600 mx-auto mb-4"></div>
          <p className="text-gray-600">Loading...</p>
        </div>
      </div>
    );
  }

  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-3xl font-bold text-gray-900">PAGETITLE</h1>
        <Button onClick={() => alert('Create functionality (TODO)')}>Create New</Button>
      </div>

      {error && (
        <div className="mb-4 bg-yellow-50 border border-yellow-200 text-yellow-700 px-4 py-3 rounded">
          <p className="font-semibold">MVP Pending</p>
          <p className="text-sm mt-1">Using mock data. API integration pending.</p>
        </div>
      )}

      <Card>
        {data.length === 0 ? (
          <p className="text-gray-600 text-center py-8">No data available.</p>
        ) : (
          <div className="overflow-x-auto">
            <table className="min-w-full divide-y divide-gray-200">
              <thead className="bg-gray-50">
                <tr>
                  <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">ID</th>
                  <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Name</th>
                  <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Status</th>
                  <th className="px-4 py-3 text-center text-xs font-medium text-gray-500 uppercase">Actions</th>
                </tr>
              </thead>
              <tbody className="bg-white divide-y divide-gray-200">
                {data.map((item) => (
                  <tr key={item.id}>
                    <td className="px-4 py-3 text-sm text-gray-900">{item.id}</td>
                    <td className="px-4 py-3 text-sm text-gray-900">{item.name || 'N/A'}</td>
                    <td className="px-4 py-3 text-sm">
                      <span className="px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-green-100 text-green-800">
                        {item.status || 'active'}
                      </span>
                    </td>
                    <td className="px-4 py-3 text-center text-sm font-medium">
                      <Button size="sm" variant="secondary" onClick={() => alert('View details (TODO)')}>
                        View
                      </Button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </Card>
    </div>
  );
}
PAGEEOF

  # Replace placeholders
  sed -i '' "s/PAGETITLE/$title/g" "$file_path"
  if [ -n "$role_check" ]; then
    sed -i '' "s/ROLE_CHECK/$role_check/" "$file_path"
  else
    sed -i '' "/ROLE_CHECK/d" "$file_path"
  fi
  
  echo "Created: $file_path"
}

# Generate pages
generate_page "/forgot-password" "Forgot Password" ""
generate_page "/c/[companySlug]" "Company Profile" ""
generate_page "/p/[listingSlug]" "Product Listing" ""
generate_page "/part/[manufacturer]/[partCode]" "Part Detail" ""
generate_page "/equipment/[brand]/[model]" "Equipment Detail" ""

generate_page "/app/settings" "Settings" ""
generate_page "/app/chat/[threadId]" "Chat Thread" ""
generate_page "/app/support" "Support" ""
generate_page "/app/support/new" "New Support Ticket" ""
generate_page "/app/support/[ticketId]" "Support Ticket" ""

# Buyer pages
generate_page "/app/customer/team" "Team" "hasRole('requester') || hasRole('procurement_manager')"
generate_page "/app/customer/invite" "Invite Team Member" "hasRole('requester') || hasRole('procurement_manager')"
generate_page "/app/customer/roles" "Team Roles" "hasRole('requester') || hasRole('procurement_manager')"
generate_page "/app/customer/company" "Company" "hasRole('requester') || hasRole('procurement_manager')"
generate_page "/app/customer/documents" "Documents" "hasRole('requester') || hasRole('procurement_manager')"
generate_page "/app/customer/addresses" "Addresses" "hasRole('requester') || hasRole('procurement_manager')"
generate_page "/app/customer/policies" "Policies" "hasRole('requester') || hasRole('procurement_manager')"
generate_page "/app/customer/pr/[id]/edit" "Edit PR" "hasRole('requester')"
generate_page "/app/customer/rfq/create" "Create RFQ" "hasRole('requester') || hasRole('procurement_manager')"
generate_page "/app/customer/rfq/[id]/quotes" "RFQ Quotes" "hasRole('procurement_manager')"
generate_page "/app/customer/rfq/[id]/award" "Award Quote" "hasRole('procurement_manager')"
generate_page "/app/customer/orders/[id]" "Order Detail" "hasRole('requester') || hasRole('procurement_manager')"
generate_page "/app/customer/shipments/[id]" "Shipment Detail" "hasRole('requester') || hasRole('procurement_manager')"
generate_page "/app/customer/parts" "Parts Library" "hasRole('requester') || hasRole('procurement_manager')"
generate_page "/app/customer/parts/[id]" "Part Detail" "hasRole('requester') || hasRole('procurement_manager')"
generate_page "/app/customer/equipment" "Equipment" "hasRole('requester') || hasRole('procurement_manager')"
generate_page "/app/customer/equipment/[id]" "Equipment Detail" "hasRole('requester') || hasRole('procurement_manager')"
generate_page "/app/customer/equipment/bom/[equipmentId]" "BOM" "hasRole('requester') || hasRole('procurement_manager')"
generate_page "/app/customer/equipment/compatibility" "Compatibility" "hasRole('requester') || hasRole('procurement_manager')"
generate_page "/app/customer/warehouse" "Warehouse" "hasRole('requester') || hasRole('procurement_manager')"
generate_page "/app/customer/warehouse/inventory" "Inventory" "hasRole('requester') || hasRole('procurement_manager')"
generate_page "/app/customer/warehouse/shared" "Shared Inventory" "hasRole('requester') || hasRole('procurement_manager')"
generate_page "/app/customer/warehouse/transfers" "Transfers" "hasRole('requester') || hasRole('procurement_manager')"
generate_page "/app/customer/warehouse/transfers/[id]" "Transfer Detail" "hasRole('requester') || hasRole('procurement_manager')"
generate_page "/app/customer/emergency" "Emergency" "hasRole('requester') || hasRole('procurement_manager')"
generate_page "/app/customer/sell/listings" "Sell Listings" "hasRole('requester') || hasRole('procurement_manager')"
generate_page "/app/customer/sell/listings/create" "Create Sell Listing" "hasRole('requester') || hasRole('procurement_manager')"
generate_page "/app/customer/auctions" "Auctions" "hasRole('requester') || hasRole('procurement_manager')"
generate_page "/app/customer/auctions/create" "Create Auction" "hasRole('requester') || hasRole('procurement_manager')"
generate_page "/app/customer/auctions/[id]" "Auction Detail" "hasRole('requester') || hasRole('procurement_manager')"
generate_page "/app/customer/reports" "Reports" "hasRole('requester') || hasRole('procurement_manager')"
generate_page "/app/customer/reports/spend" "Spend Report" "hasRole('requester') || hasRole('procurement_manager')"
generate_page "/app/customer/reports/performance" "Performance Report" "hasRole('requester') || hasRole('procurement_manager')"
generate_page "/app/customer/reports/late" "Late Shipments" "hasRole('requester') || hasRole('procurement_manager')"
generate_page "/app/customer/reports/inventory-aging" "Inventory Aging" "hasRole('requester') || hasRole('procurement_manager')"

# Supplier pages  
generate_page "/app/supplier/company" "Company" "hasRole('supplier')"
generate_page "/app/supplier/documents" "Documents" "hasRole('supplier')"
generate_page "/app/supplier/team" "Team" "hasRole('supplier')"
generate_page "/app/supplier/invite" "Invite Team Member" "hasRole('supplier')"
generate_page "/app/supplier/roles" "Team Roles" "hasRole('supplier')"
generate_page "/app/supplier/store" "Store" "hasRole('supplier')"
generate_page "/app/supplier/store/locations" "Store Locations" "hasRole('supplier')"
generate_page "/app/supplier/store/policies" "Store Policies" "hasRole('supplier')"
generate_page "/app/supplier/listings/[id]" "Listing Detail" "hasRole('supplier')"
generate_page "/app/supplier/listings/[id]/edit" "Edit Listing" "hasRole('supplier')"
generate_page "/app/supplier/services" "Services" "hasRole('supplier')"
generate_page "/app/supplier/services/create" "Create Service" "hasRole('supplier')"
generate_page "/app/supplier/services/[id]" "Service Detail" "hasRole('supplier')"
generate_page "/app/supplier/inventory" "Inventory" "hasRole('supplier')"
generate_page "/app/supplier/inventory/import" "Import Inventory" "hasRole('supplier')"
generate_page "/app/supplier/pricing" "Pricing" "hasRole('supplier')"
generate_page "/app/supplier/rfq/[id]" "RFQ Detail" "hasRole('supplier')"
generate_page "/app/supplier/quotes/[id]" "Quote Detail" "hasRole('supplier')"
generate_page "/app/supplier/orders/[id]" "Order Detail" "hasRole('supplier')"
generate_page "/app/supplier/shipments/[id]" "Shipment Detail" "hasRole('supplier')"
generate_page "/app/supplier/ratings" "Ratings" "hasRole('supplier')"
generate_page "/app/supplier/performance" "Performance" "hasRole('supplier')"
generate_page "/app/supplier/reports" "Reports" "hasRole('supplier')"
generate_page "/app/supplier/reports/sales" "Sales Report" "hasRole('supplier')"
generate_page "/app/supplier/reports/win-rate" "Win Rate" "hasRole('supplier')"
generate_page "/app/supplier/reports/on-time" "On-Time Delivery" "hasRole('supplier')"

# Admin pages
generate_page "/app/admin/companies" "Companies" "hasRole('admin') || hasRole('super_admin')"
generate_page "/app/admin/subdomains" "Subdomains" "hasRole('admin') || hasRole('super_admin')"
generate_page "/app/admin/catalog/parts" "Catalog Parts" "hasRole('admin') || hasRole('super_admin')"
generate_page "/app/admin/catalog/equipment" "Catalog Equipment" "hasRole('admin') || hasRole('super_admin')"
generate_page "/app/admin/catalog/duplicates" "Duplicates" "hasRole('admin') || hasRole('super_admin')"
generate_page "/app/admin/listings" "Listings" "hasRole('admin') || hasRole('super_admin')"
generate_page "/app/admin/stores" "Stores" "hasRole('admin') || hasRole('super_admin')"
generate_page "/app/admin/rfqs" "RFQs" "hasRole('admin') || hasRole('super_admin')"
generate_page "/app/admin/orders" "Orders" "hasRole('admin') || hasRole('super_admin')"
generate_page "/app/admin/shipments" "Shipments" "hasRole('admin') || hasRole('super_admin')"
generate_page "/app/admin/chat-moderation" "Chat Moderation" "hasRole('admin') || hasRole('super_admin')"
generate_page "/app/admin/plans" "Plans" "hasRole('admin') || hasRole('super_admin')"
generate_page "/app/admin/payments" "Payments" "hasRole('admin') || hasRole('super_admin')"
generate_page "/app/admin/notifications/templates" "Notification Templates" "hasRole('admin') || hasRole('super_admin')"
generate_page "/app/admin/notifications/rules" "Notification Rules" "hasRole('admin') || hasRole('super_admin')"

echo ""
echo "✅ Generated all pages"
