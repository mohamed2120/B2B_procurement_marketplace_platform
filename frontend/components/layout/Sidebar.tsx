'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { hasRole } from '@/lib/auth';

interface NavItem {
  name: string;
  href: string;
  icon: string;
  roles?: string[];
}

const customerNavItems: NavItem[] = [
  { name: 'Dashboard', href: '/app/customer/dashboard', icon: '🏠' },
  { name: 'Purchase Requests', href: '/app/customer/pr', icon: '📋' },
  { name: 'RFQs', href: '/app/customer/rfq', icon: '📝' },
  { name: 'Orders', href: '/app/customer/orders', icon: '📦' },
  { name: 'Shipments', href: '/app/customer/shipments', icon: '🚚' },
];

const accountNavItems: NavItem[] = [
  { name: 'My Plan', href: '/app/my-plan', icon: '💳' },
  { name: 'Notifications', href: '/app/notifications', icon: '🔔' },
  { name: 'Chat', href: '/app/chat', icon: '💬' },
];

const supplierNavItems: NavItem[] = [
  { name: 'Dashboard', href: '/app/supplier/dashboard', icon: '🏠' },
  { name: 'RFQ Inbox', href: '/app/supplier/rfq', icon: '📥' },
  { name: 'My Quotes', href: '/app/supplier/quotes', icon: '💵' },
  { name: 'Listings', href: '/app/supplier/listings', icon: '🏪' },
  { name: 'Orders', href: '/app/supplier/orders', icon: '📦' },
  { name: 'Shipments', href: '/app/supplier/shipments', icon: '🚚' },
];

const adminNavItems: NavItem[] = [
  { name: 'Dashboard', href: '/app/admin/dashboard', icon: '🏠' },
  { name: 'Company Verification', href: '/app/admin/company-verification', icon: '🏢' },
  { name: 'Catalog Approvals', href: '/app/admin/catalog-approvals', icon: '📚' },
  { name: 'Disputes', href: '/app/admin/disputes', icon: '⚖️' },
  { name: 'Subscriptions', href: '/app/admin/subscriptions', icon: '💳' },
  { name: 'Diagnostics', href: '/app/admin/diagnostics', icon: '🔧' },
];

export default function Sidebar() {
  const pathname = usePathname();

  const isActive = (href: string) => pathname === href || pathname?.startsWith(href + '/');

  const renderNavItems = (items: NavItem[]) => {
    return items
      .filter(item => !item.roles || item.roles.some(role => hasRole(role)))
      .map((item) => (
        <Link
          key={item.href}
          href={item.href}
          className={`flex items-center space-x-3 px-4 py-2 rounded-lg transition-colors ${
            isActive(item.href)
              ? 'bg-primary-50 text-primary-700 font-medium'
              : 'text-gray-700 hover:bg-gray-50'
          }`}
        >
          <span className="text-xl">{item.icon}</span>
          <span>{item.name}</span>
        </Link>
      ));
  };

  const isCustomer = hasRole('requester') || hasRole('procurement_manager') || hasRole('buyer');
  const isSupplier = hasRole('supplier');
  const isAdmin = hasRole('admin') || hasRole('super_admin');

  return (
    <aside className="w-64 bg-white border-r min-h-screen p-4">
      <nav className="space-y-2">
        {isCustomer && (
          <div className="mb-6">
            <h2 className="text-xs font-semibold text-gray-500 uppercase tracking-wider mb-2">Customer</h2>
            {renderNavItems(customerNavItems)}
          </div>
        )}

        {isSupplier && (
          <div className="mb-6">
            <h2 className="text-xs font-semibold text-gray-500 uppercase tracking-wider mb-2">Supplier</h2>
            {renderNavItems(supplierNavItems)}
          </div>
        )}

        {isAdmin && (
          <div className="mb-6">
            <h2 className="text-xs font-semibold text-gray-500 uppercase tracking-wider mb-2">Admin</h2>
            {renderNavItems(adminNavItems)}
          </div>
        )}

        <div className="mb-6">
          <h2 className="text-xs font-semibold text-gray-500 uppercase tracking-wider mb-2">Account</h2>
          {renderNavItems(accountNavItems)}
        </div>
      </nav>
    </aside>
  );
}
