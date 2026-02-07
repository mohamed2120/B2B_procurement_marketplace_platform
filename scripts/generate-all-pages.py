#!/usr/bin/env python3
"""
Generate all missing pages for the B2B platform
"""
import os
import re

PAGES = [
    # PUBLIC
    ("/forgot-password", "ForgotPassword", None, None),
    ("/c/[companySlug]", "CompanyProfile", None, "company"),
    ("/p/[listingSlug]", "ProductListing", None, "marketplace"),
    ("/part/[manufacturer]/[partCode]", "PartDetail", None, "catalog"),
    ("/equipment/[brand]/[model]", "EquipmentDetail", None, "equipment"),
    
    # APP SHARED
    ("/app/settings", "Settings", None, None),
    ("/app/chat/[threadId]", "ChatThread", None, "collaboration"),
    ("/app/support", "Support", None, None),
    ("/app/support/new", "NewSupportTicket", None, None),
    ("/app/support/[ticketId]", "SupportTicket", None, None),
    
    # BUYER
    ("/app/customer/team", "Team", ["requester", "procurement_manager"], None),
    ("/app/customer/invite", "InviteTeam", ["requester", "procurement_manager"], None),
    ("/app/customer/roles", "TeamRoles", ["requester", "procurement_manager"], None),
    ("/app/customer/company", "Company", ["requester", "procurement_manager"], "company"),
    ("/app/customer/documents", "Documents", ["requester", "procurement_manager"], None),
    ("/app/customer/addresses", "Addresses", ["requester", "procurement_manager"], None),
    ("/app/customer/policies", "Policies", ["requester", "procurement_manager"], None),
    ("/app/customer/pr/[id]/edit", "EditPR", ["requester"], "procurement"),
    ("/app/customer/rfq/create", "CreateRFQ", ["requester", "procurement_manager"], "procurement"),
    ("/app/customer/rfq/[id]/quotes", "RFQQuotes", ["procurement_manager"], "procurement"),
    ("/app/customer/rfq/[id]/award", "AwardQuote", ["procurement_manager"], "procurement"),
    ("/app/customer/orders/[id]", "OrderDetail", ["requester", "procurement_manager"], "procurement"),
    ("/app/customer/shipments/[id]", "ShipmentDetail", ["requester", "procurement_manager"], "logistics"),
    ("/app/customer/parts", "PartsLibrary", ["requester", "procurement_manager"], "catalog"),
    ("/app/customer/parts/[id]", "PartDetail", ["requester", "procurement_manager"], "catalog"),
    ("/app/customer/equipment", "Equipment", ["requester", "procurement_manager"], "equipment"),
    ("/app/customer/equipment/[id]", "EquipmentDetail", ["requester", "procurement_manager"], "equipment"),
    ("/app/customer/equipment/bom/[equipmentId]", "BOM", ["requester", "procurement_manager"], "equipment"),
    ("/app/customer/equipment/compatibility", "Compatibility", ["requester", "procurement_manager"], "equipment"),
    ("/app/customer/warehouse", "Warehouse", ["requester", "procurement_manager"], None),
    ("/app/customer/warehouse/inventory", "Inventory", ["requester", "procurement_manager"], None),
    ("/app/customer/warehouse/shared", "SharedInventory", ["requester", "procurement_manager"], None),
    ("/app/customer/warehouse/transfers", "Transfers", ["requester", "procurement_manager"], None),
    ("/app/customer/warehouse/transfers/[id]", "TransferDetail", ["requester", "procurement_manager"], None),
    ("/app/customer/emergency", "Emergency", ["requester", "procurement_manager"], None),
    ("/app/customer/sell/listings", "SellListings", ["requester", "procurement_manager"], "marketplace"),
    ("/app/customer/sell/listings/create", "CreateSellListing", ["requester", "procurement_manager"], "marketplace"),
    ("/app/customer/auctions", "Auctions", ["requester", "procurement_manager"], None),
    ("/app/customer/auctions/create", "CreateAuction", ["requester", "procurement_manager"], None),
    ("/app/customer/auctions/[id]", "AuctionDetail", ["requester", "procurement_manager"], None),
    ("/app/customer/reports", "Reports", ["requester", "procurement_manager"], None),
    ("/app/customer/reports/spend", "SpendReport", ["requester", "procurement_manager"], None),
    ("/app/customer/reports/performance", "PerformanceReport", ["requester", "procurement_manager"], None),
    ("/app/customer/reports/late", "LateShipments", ["requester", "procurement_manager"], None),
    ("/app/customer/reports/inventory-aging", "InventoryAging", ["requester", "procurement_manager"], None),
    
    # SUPPLIER
    ("/app/supplier/company", "Company", ["supplier"], "company"),
    ("/app/supplier/documents", "Documents", ["supplier"], None),
    ("/app/supplier/team", "Team", ["supplier"], None),
    ("/app/supplier/invite", "InviteTeam", ["supplier"], None),
    ("/app/supplier/roles", "TeamRoles", ["supplier"], None),
    ("/app/supplier/store", "Store", ["supplier"], "marketplace"),
    ("/app/supplier/store/locations", "StoreLocations", ["supplier"], "marketplace"),
    ("/app/supplier/store/policies", "StorePolicies", ["supplier"], "marketplace"),
    ("/app/supplier/listings/[id]", "ListingDetail", ["supplier"], "marketplace"),
    ("/app/supplier/listings/[id]/edit", "EditListing", ["supplier"], "marketplace"),
    ("/app/supplier/services", "Services", ["supplier"], "marketplace"),
    ("/app/supplier/services/create", "CreateService", ["supplier"], "marketplace"),
    ("/app/supplier/services/[id]", "ServiceDetail", ["supplier"], "marketplace"),
    ("/app/supplier/inventory", "Inventory", ["supplier"], None),
    ("/app/supplier/inventory/import", "ImportInventory", ["supplier"], None),
    ("/app/supplier/pricing", "Pricing", ["supplier"], None),
    ("/app/supplier/rfq/[id]", "RFQDetail", ["supplier"], "procurement"),
    ("/app/supplier/quotes/[id]", "QuoteDetail", ["supplier"], "procurement"),
    ("/app/supplier/orders/[id]", "OrderDetail", ["supplier"], "procurement"),
    ("/app/supplier/shipments/[id]", "ShipmentDetail", ["supplier"], "logistics"),
    ("/app/supplier/ratings", "Ratings", ["supplier"], None),
    ("/app/supplier/performance", "Performance", ["supplier"], None),
    ("/app/supplier/reports", "Reports", ["supplier"], None),
    ("/app/supplier/reports/sales", "SalesReport", ["supplier"], None),
    ("/app/supplier/reports/win-rate", "WinRate", ["supplier"], None),
    ("/app/supplier/reports/on-time", "OnTimeDelivery", ["supplier"], None),
    
    # ADMIN
    ("/app/admin/companies", "Companies", ["admin", "super_admin"], "company"),
    ("/app/admin/subdomains", "Subdomains", ["admin", "super_admin"], None),
    ("/app/admin/catalog/parts", "CatalogParts", ["admin", "super_admin"], "catalog"),
    ("/app/admin/catalog/equipment", "CatalogEquipment", ["admin", "super_admin"], "equipment"),
    ("/app/admin/catalog/duplicates", "Duplicates", ["admin", "super_admin"], "catalog"),
    ("/app/admin/listings", "Listings", ["admin", "super_admin"], "marketplace"),
    ("/app/admin/stores", "Stores", ["admin", "super_admin"], "marketplace"),
    ("/app/admin/rfqs", "RFQs", ["admin", "super_admin"], "procurement"),
    ("/app/admin/orders", "Orders", ["admin", "super_admin"], "procurement"),
    ("/app/admin/shipments", "Shipments", ["admin", "super_admin"], "logistics"),
    ("/app/admin/chat-moderation", "ChatModeration", ["admin", "super_admin"], "collaboration"),
    ("/app/admin/plans", "Plans", ["admin", "super_admin"], "billing"),
    ("/app/admin/payments", "Payments", ["admin", "super_admin"], "billing"),
    ("/app/admin/notifications/templates", "NotificationTemplates", ["admin", "super_admin"], "notification"),
    ("/app/admin/notifications/rules", "NotificationRules", ["admin", "super_admin"], "notification"),
]

def to_camel_case(s):
    """Convert string to PascalCase"""
    return ''.join(word.capitalize() for word in re.split(r'[-_\s]+', s))

def get_api_client(api_type):
    if not api_type:
        return None
    mapping = {
        "company": "company",
        "catalog": "catalog",
        "equipment": "equipment",
        "marketplace": "marketplace",
        "procurement": "procurement",
        "logistics": "logistics",
        "collaboration": "collaboration",
        "billing": "billing",
        "notification": "notification",
    }
    return mapping.get(api_type)

def generate_page(route, component_name, roles, api_type):
    file_path = f"frontend/app{route}/page.tsx"
    os.makedirs(os.path.dirname(file_path), exist_ok=True)
    
    api_client = get_api_client(api_type)
    role_check = ""
    if roles:
        role_conditions = " || ".join([f"hasRole('{r}')" for r in roles])
        role_check = f"""
  useEffect(() => {
    if (!({role_conditions})) {{
      router.push('/app');
      return;
    }}
  }, []);"""
    
    api_call = ""
    if api_client:
        api_call = f"""
      try {{
        const response = await apiClients.{api_client}.get('/api/v1/...');
        setData(response.data || []);
      }} catch (err: any) {{
        setError(err.response?.data?.error || 'Failed to load data');
        // Fallback to mock data
        setData([{{ id: '1', name: 'Sample Item', status: 'active' }}]);
      }}"""
    else:
        api_call = """
      // TODO: Implement API call
      setData([
        { id: '1', name: 'Sample Item 1', status: 'active' },
        { id: '2', name: 'Sample Item 2', status: 'pending' },
      ]);"""
    
    # Escape curly braces for f-string
    role_check_escaped = role_check.replace('{', '{{').replace('}', '}}') if role_check else ""
    api_call_escaped = api_call.replace('{', '{{').replace('}', '}}')
    title = component_name.replace(/([A-Z])/g, r' $1').strip()
    
    content = f"""'use client';

import {{ useState, useEffect }} from 'react';
import {{ useRouter }} from 'next/navigation';
import {{ apiClients }} from '@/lib/api';
import {{ hasRole }} from '@/lib/auth';
import Card from '@/components/ui/Card';
import Button from '@/components/ui/Button';

interface DataItem {{
  id: string;
  [key: string]: any;
}}

export default function {component_name}Page() {{
  const router = useRouter();
  const [data, setData] = useState<DataItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
{role_check}

  useEffect(() => {{
    fetchData();
  }}, []);

  const fetchData = async () => {{
    setLoading(true);
    setError('');
{api_call}
    setLoading(false);
  }};

  if (loading) {{
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600 mx-auto mb-4"></div>
          <p className="text-gray-600">Loading...</p>
        </div>
      </div>
    );
  }}

  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-3xl font-bold text-gray-900">{title}</h1>
        <Button onClick={{() => alert('Create functionality (TODO)')}}>Create New</Button>
      </div>

      {{error && (
        <div className="mb-4 bg-yellow-50 border border-yellow-200 text-yellow-700 px-4 py-3 rounded">
          <p className="font-semibold">MVP Pending</p>
          <p className="text-sm mt-1">Using mock data. API integration pending.</p>
        </div>
      )}}

      <Card>
        {{data.length === 0 ? (
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
                {{data.map((item) => (
                  <tr key={{item.id}}>
                    <td className="px-4 py-3 text-sm text-gray-900">{{item.id}}</td>
                    <td className="px-4 py-3 text-sm text-gray-900">{{item.name || 'N/A'}}</td>
                    <td className="px-4 py-3 text-sm">
                      <span className="px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-green-100 text-green-800">
                        {{item.status || 'active'}}
                      </span>
                    </td>
                    <td className="px-4 py-3 text-center text-sm font-medium">
                      <Button size="sm" variant="secondary" onClick={{() => alert('View details (TODO)')}}>
                        View
                      </Button>
                    </td>
                  </tr>
                ))}}
              </tbody>
            </table>
          </div>
        )}}
      </Card>
    </div>
  );
}}
"""
    
    with open(file_path, 'w') as f:
        f.write(content)
    print(f"Created: {file_path}")

if __name__ == "__main__":
    for route, component_name, roles, api_type in PAGES:
        generate_page(route, component_name, roles, api_type)
    print(f"\n✅ Generated {len(PAGES)} pages")
