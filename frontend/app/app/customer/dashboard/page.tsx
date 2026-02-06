'use client';

import { getUser, hasRole } from '@/lib/auth';
import Card from '@/components/ui/Card';
import Link from 'next/link';

export default function CustomerDashboard() {
  const user = getUser();
  const isRequester = hasRole('requester');
  const isProcurement = hasRole('procurement_manager');

  return (
    <div>
      <h1 className="text-3xl font-bold text-gray-900 mb-6">Customer Dashboard</h1>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {isRequester && (
          <>
            <Card title="Purchase Requests">
              <p className="text-gray-600 mb-4">Create and manage purchase requests</p>
              <Link href="/app/customer/pr">
                <button className="w-full bg-primary-600 text-white px-4 py-2 rounded-lg hover:bg-primary-700">
                  View PRs
                </button>
              </Link>
            </Card>

            <Card title="RFQs">
              <p className="text-gray-600 mb-4">View and manage RFQs</p>
              <Link href="/app/customer/rfq">
                <button className="w-full bg-primary-600 text-white px-4 py-2 rounded-lg hover:bg-primary-700">
                  View RFQs
                </button>
              </Link>
            </Card>
          </>
        )}

        {isProcurement && (
          <>
            <Card title="Approve PRs">
              <p className="text-gray-600 mb-4">Review and approve purchase requests</p>
              <Link href="/app/customer/pr">
                <button className="w-full bg-primary-600 text-white px-4 py-2 rounded-lg hover:bg-primary-700">
                  Review PRs
                </button>
              </Link>
            </Card>

            <Card title="Award Quotes">
              <p className="text-gray-600 mb-4">Compare and award quotes</p>
              <Link href="/app/customer/rfq">
                <button className="w-full bg-primary-600 text-white px-4 py-2 rounded-lg hover:bg-primary-700">
                  Manage Quotes
                </button>
              </Link>
            </Card>
          </>
        )}

        <Card title="Orders">
          <p className="text-gray-600 mb-4">Track your orders</p>
          <Link href="/app/customer/orders">
            <button className="w-full bg-primary-600 text-white px-4 py-2 rounded-lg hover:bg-primary-700">
              View Orders
            </button>
          </Link>
        </Card>

        <Card title="Shipments">
          <p className="text-gray-600 mb-4">Track shipments</p>
          <Link href="/app/customer/shipments">
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
