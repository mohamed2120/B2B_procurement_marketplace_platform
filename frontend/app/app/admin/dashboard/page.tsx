'use client';

import { getUser } from '@/lib/auth';
import Card from '@/components/ui/Card';
import Link from 'next/link';

export default function AdminDashboard() {
  const user = getUser();

  return (
    <div>
      <h1 className="text-3xl font-bold text-gray-900 mb-6">Admin Dashboard</h1>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <Card title="Company Verification">
          <p className="text-gray-600 mb-4">Review and approve companies</p>
          <Link href="/app/admin/company-verification">
            <button className="w-full bg-primary-600 text-white px-4 py-2 rounded-lg hover:bg-primary-700">
              Review Companies
            </button>
          </Link>
        </Card>

        <Card title="Catalog Approvals">
          <p className="text-gray-600 mb-4">Approve catalog items</p>
          <Link href="/app/admin/catalog-approvals">
            <button className="w-full bg-primary-600 text-white px-4 py-2 rounded-lg hover:bg-primary-700">
              Review Catalog
            </button>
          </Link>
        </Card>

        <Card title="Disputes">
          <p className="text-gray-600 mb-4">Manage disputes</p>
          <Link href="/app/admin/disputes">
            <button className="w-full bg-primary-600 text-white px-4 py-2 rounded-lg hover:bg-primary-700">
              View Disputes
            </button>
          </Link>
        </Card>

        <Card title="Subscriptions">
          <p className="text-gray-600 mb-4">Manage tenant subscriptions</p>
          <Link href="/app/admin/subscriptions">
            <button className="w-full bg-primary-600 text-white px-4 py-2 rounded-lg hover:bg-primary-700">
              Manage Subscriptions
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
