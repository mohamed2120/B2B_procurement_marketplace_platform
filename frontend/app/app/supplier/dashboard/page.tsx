'use client';

import { getUser } from '@/lib/auth';
import Card from '@/components/ui/Card';
import Link from 'next/link';

export default function SupplierDashboard() {
  const user = getUser();

  return (
    <div>
      <h1 className="text-3xl font-bold text-gray-900 mb-6">Supplier Dashboard</h1>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <Card title="RFQ Inbox">
          <p className="text-gray-600 mb-4">View incoming RFQs</p>
          <Link href="/app/supplier/rfq">
            <button className="w-full bg-primary-600 text-white px-4 py-2 rounded-lg hover:bg-primary-700">
              View RFQs
            </button>
          </Link>
        </Card>

        <Card title="My Quotes">
          <p className="text-gray-600 mb-4">Manage your quotes</p>
          <Link href="/app/supplier/quotes">
            <button className="w-full bg-primary-600 text-white px-4 py-2 rounded-lg hover:bg-primary-700">
              View Quotes
            </button>
          </Link>
        </Card>

        <Card title="Listings">
          <p className="text-gray-600 mb-4">Manage your product listings</p>
          <Link href="/app/supplier/listings">
            <button className="w-full bg-primary-600 text-white px-4 py-2 rounded-lg hover:bg-primary-700">
              Manage Listings
            </button>
          </Link>
        </Card>

        <Card title="Orders">
          <p className="text-gray-600 mb-4">View customer orders</p>
          <Link href="/app/supplier/orders">
            <button className="w-full bg-primary-600 text-white px-4 py-2 rounded-lg hover:bg-primary-700">
              View Orders
            </button>
          </Link>
        </Card>

        <Card title="Shipments">
          <p className="text-gray-600 mb-4">Track shipments</p>
          <Link href="/app/supplier/shipments">
            <button className="w-full bg-primary-600 text-white px-4 py-2 rounded-lg hover:bg-primary-700">
              View Shipments
            </button>
          </Link>
        </Card>
      </div>

      <Card title="Welcome" className="mt-6">
        <p className="text-gray-600">
          Welcome back, {user?.first_name} {user?.last_name}!
        </p>
        <p className="text-sm text-gray-500 mt-2">
          Email: {user?.email} | Roles: {user?.roles?.join(', ') || 'None'}
        </p>
      </Card>
    </div>
  );
}
