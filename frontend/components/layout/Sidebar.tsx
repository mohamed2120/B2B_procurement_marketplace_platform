'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { hasRole } from '@/lib/auth';
import { useState } from 'react';

interface NavItem {
  name: string;
  href: string;
  icon: string;
  roles?: string[];
}

interface NavSection {
  title: string;
  items: NavItem[];
  roles?: string[];
  collapsible?: boolean;
}

const customerSections: NavSection[] = [
  {
    title: 'Overview',
    items: [
      { name: 'Dashboard', href: '/app/customer/dashboard', icon: '🏠' },
    ],
  },
  {
    title: 'Procurement',
    items: [
      { name: 'Purchase Requests', href: '/app/customer/pr', icon: '📋' },
      { name: 'RFQs', href: '/app/customer/rfq', icon: '📝' },
      { name: 'Orders', href: '/app/customer/orders', icon: '📦' },
      { name: 'Shipments', href: '/app/customer/shipments', icon: '🚚' },
    ],
  },
  {
    title: 'Catalog',
    items: [
      { name: 'Parts Library', href: '/app/customer/parts', icon: '🔧' },
      { name: 'Equipment', href: '/app/customer/equipment', icon: '⚙️' },
    ],
  },
  {
    title: 'Warehouse',
    items: [
      { name: 'Overview', href: '/app/customer/warehouse', icon: '📦' },
      { name: 'Inventory', href: '/app/customer/warehouse/inventory', icon: '📊' },
      { name: 'Shared Inventory', href: '/app/customer/warehouse/shared', icon: '🤝' },
      { name: 'Transfers', href: '/app/customer/warehouse/transfers', icon: '🔄' },
    ],
  },
  {
    title: 'Company',
    items: [
      { name: 'Company Info', href: '/app/customer/company', icon: '🏢' },
      { name: 'Team', href: '/app/customer/team', icon: '👥' },
      { name: 'Documents', href: '/app/customer/documents', icon: '📄' },
      { name: 'Addresses', href: '/app/customer/addresses', icon: '📍' },
      { name: 'Policies', href: '/app/customer/policies', icon: '📋' },
    ],
  },
  {
    title: 'Selling',
    items: [
      { name: 'Listings', href: '/app/customer/sell/listings', icon: '🏪' },
      { name: 'Auctions', href: '/app/customer/auctions', icon: '🔨' },
    ],
  },
  {
    title: 'Reports',
    items: [
      { name: 'Overview', href: '/app/customer/reports', icon: '📊' },
      { name: 'Spend Report', href: '/app/customer/reports/spend', icon: '💰' },
      { name: 'Performance', href: '/app/customer/reports/performance', icon: '📈' },
      { name: 'Late Shipments', href: '/app/customer/reports/late', icon: '⚠️' },
      { name: 'Inventory Aging', href: '/app/customer/reports/inventory-aging', icon: '📅' },
    ],
  },
  {
    title: 'Emergency',
    items: [
      { name: 'Emergency', href: '/app/customer/emergency', icon: '🚨' },
    ],
  },
];

const supplierSections: NavSection[] = [
  {
    title: 'Overview',
    items: [
      { name: 'Dashboard', href: '/app/supplier/dashboard', icon: '🏠' },
    ],
  },
  {
    title: 'Opportunities',
    items: [
      { name: 'RFQ Inbox', href: '/app/supplier/rfq', icon: '📥' },
      { name: 'My Quotes', href: '/app/supplier/quotes', icon: '💵' },
      { name: 'Orders', href: '/app/supplier/orders', icon: '📦' },
      { name: 'Shipments', href: '/app/supplier/shipments', icon: '🚚' },
    ],
  },
  {
    title: 'Store',
    items: [
      { name: 'Listings', href: '/app/supplier/listings', icon: '🏪' },
      { name: 'Services', href: '/app/supplier/services', icon: '🛠️' },
      { name: 'Inventory', href: '/app/supplier/inventory', icon: '📊' },
      { name: 'Pricing', href: '/app/supplier/pricing', icon: '💰' },
      { name: 'Store Settings', href: '/app/supplier/store', icon: '⚙️' },
    ],
  },
  {
    title: 'Company',
    items: [
      { name: 'Company Info', href: '/app/supplier/company', icon: '🏢' },
      { name: 'Team', href: '/app/supplier/team', icon: '👥' },
      { name: 'Documents', href: '/app/supplier/documents', icon: '📄' },
    ],
  },
  {
    title: 'Performance',
    items: [
      { name: 'Ratings', href: '/app/supplier/ratings', icon: '⭐' },
      { name: 'Performance', href: '/app/supplier/performance', icon: '📈' },
    ],
  },
  {
    title: 'Reports',
    items: [
      { name: 'Overview', href: '/app/supplier/reports', icon: '📊' },
      { name: 'Sales', href: '/app/supplier/reports/sales', icon: '💰' },
      { name: 'Win Rate', href: '/app/supplier/reports/win-rate', icon: '🎯' },
      { name: 'On-Time Delivery', href: '/app/supplier/reports/on-time', icon: '⏰' },
    ],
  },
];

const adminSections: NavSection[] = [
  {
    title: 'Overview',
    items: [
      { name: 'Dashboard', href: '/app/admin/dashboard', icon: '🏠' },
    ],
  },
  {
    title: 'Tenants & Users',
    items: [
      { name: 'Tenants', href: '/app/admin/tenants', icon: '🏢' },
      { name: 'Companies', href: '/app/admin/companies', icon: '🏭' },
      { name: 'Users', href: '/app/admin/users', icon: '👥' },
      { name: 'Subdomains', href: '/app/admin/subdomains', icon: '🌐' },
    ],
  },
  {
    title: 'Access Control',
    items: [
      { name: 'Roles & Permissions', href: '/app/admin/roles-permissions', icon: '🔐' },
    ],
  },
  {
    title: 'Catalog',
    items: [
      { name: 'Parts', href: '/app/admin/catalog/parts', icon: '🔧' },
      { name: 'Equipment', href: '/app/admin/catalog/equipment', icon: '⚙️' },
      { name: 'Duplicates', href: '/app/admin/catalog/duplicates', icon: '🔍' },
      { name: 'Approvals', href: '/app/admin/catalog-approvals', icon: '✅' },
    ],
  },
  {
    title: 'Marketplace',
    items: [
      { name: 'Listings', href: '/app/admin/listings', icon: '🏪' },
      { name: 'Stores', href: '/app/admin/stores', icon: '🏬' },
    ],
  },
  {
    title: 'Procurement',
    items: [
      { name: 'RFQs', href: '/app/admin/rfqs', icon: '📝' },
      { name: 'Orders', href: '/app/admin/orders', icon: '📦' },
      { name: 'Shipments', href: '/app/admin/shipments', icon: '🚚' },
    ],
  },
  {
    title: 'Operations',
    items: [
      { name: 'Company Verification', href: '/app/admin/company-verification', icon: '✅' },
      { name: 'Disputes', href: '/app/admin/disputes', icon: '⚖️' },
      { name: 'Chat Moderation', href: '/app/admin/chat-moderation', icon: '💬' },
    ],
  },
  {
    title: 'Billing',
    items: [
      { name: 'Subscriptions', href: '/app/admin/subscriptions', icon: '💳' },
      { name: 'Plans', href: '/app/admin/plans', icon: '📋' },
      { name: 'Payments', href: '/app/admin/payments', icon: '💰' },
    ],
  },
  {
    title: 'Notifications',
    items: [
      { name: 'Templates', href: '/app/admin/notifications/templates', icon: '📧' },
      { name: 'Rules', href: '/app/admin/notifications/rules', icon: '⚙️' },
    ],
  },
  {
    title: 'System',
    items: [
      { name: 'Audit Logs', href: '/app/admin/audit-logs', icon: '📋' },
      { name: 'Diagnostics', href: '/app/admin/diagnostics', icon: '🔧' },
    ],
  },
];

const accountNavItems: NavItem[] = [
  { name: 'My Plan', href: '/app/my-plan', icon: '💳' },
  { name: 'Profile', href: '/app/profile', icon: '👤' },
  { name: 'Settings', href: '/app/settings', icon: '⚙️' },
  { name: 'Notifications', href: '/app/notifications', icon: '🔔' },
  { name: 'Chat', href: '/app/chat', icon: '💬' },
  { name: 'Support', href: '/app/support', icon: '🆘' },
];

export default function Sidebar() {
  const pathname = usePathname();
  const [expandedSections, setExpandedSections] = useState<Record<string, boolean>>({});

  const isActive = (href: string) => pathname === href || pathname?.startsWith(href + '/');

  const toggleSection = (title: string) => {
    setExpandedSections(prev => ({ ...prev, [title]: !prev[title] }));
  };

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

  const renderSection = (section: NavSection, defaultExpanded: boolean = true) => {
    const isExpanded = expandedSections[section.title] ?? defaultExpanded;
    const hasAccess = !section.roles || section.roles.some(role => hasRole(role));
    
    if (!hasAccess) return null;

    return (
      <div key={section.title} className="mb-4">
        {section.collapsible !== false ? (
          <button
            onClick={() => toggleSection(section.title)}
            className="w-full flex items-center justify-between px-2 py-1 text-xs font-semibold text-gray-500 uppercase tracking-wider hover:text-gray-700"
          >
            <span>{section.title}</span>
            <span>{isExpanded ? '▼' : '▶'}</span>
          </button>
        ) : (
          <h2 className="px-2 py-1 text-xs font-semibold text-gray-500 uppercase tracking-wider">
            {section.title}
          </h2>
        )}
        {isExpanded && (
          <div className="mt-1 space-y-1">
            {renderNavItems(section.items)}
          </div>
        )}
      </div>
    );
  };

  const isCustomer = hasRole('requester') || hasRole('procurement_manager') || hasRole('buyer');
  const isSupplier = hasRole('supplier');
  const isAdmin = hasRole('admin') || hasRole('super_admin');

  return (
    <aside className="w-64 bg-white border-r min-h-screen p-4 overflow-y-auto">
      <nav className="space-y-2">
        {isCustomer && (
          <div className="mb-6">
            <h2 className="text-xs font-semibold text-gray-500 uppercase tracking-wider mb-2 px-2">Customer</h2>
            {customerSections.map(section => renderSection(section))}
          </div>
        )}

        {isSupplier && (
          <div className="mb-6">
            <h2 className="text-xs font-semibold text-gray-500 uppercase tracking-wider mb-2 px-2">Supplier</h2>
            {supplierSections.map(section => renderSection(section))}
          </div>
        )}

        {isAdmin && (
          <div className="mb-6">
            <h2 className="text-xs font-semibold text-gray-500 uppercase tracking-wider mb-2 px-2">Admin</h2>
            {adminSections.map(section => renderSection(section))}
          </div>
        )}

        <div className="mb-6">
          <h2 className="text-xs font-semibold text-gray-500 uppercase tracking-wider mb-2 px-2">Account</h2>
          {renderNavItems(accountNavItems)}
        </div>
      </nav>
    </aside>
  );
}
