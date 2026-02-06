'use client';

import { useState, useEffect } from 'react';
import { apiClients } from '@/lib/api';
import Card from '@/components/ui/Card';

export default function AdminSubscriptions() {
  const [subscriptions, setSubscriptions] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchSubscriptions();
  }, []);

  const fetchSubscriptions = async () => {
    try {
      // TODO: Implement admin endpoint to list all subscriptions
      setSubscriptions([]);
    } catch (error) {
      console.error('Failed to fetch subscriptions:', error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div>
      <h1 className="text-3xl font-bold text-gray-900 mb-6">Subscriptions</h1>

      {loading ? (
        <div>Loading...</div>
      ) : subscriptions.length === 0 ? (
        <Card>
          <p className="text-gray-600">No subscriptions found.</p>
          <p className="text-sm text-gray-500 mt-2">TODO: Implement admin subscription management endpoint</p>
        </Card>
      ) : (
        <div className="space-y-4">
          {subscriptions.map((sub) => (
            <Card key={sub.id}>
              <div className="flex justify-between items-start">
                <div>
                  <h3 className="font-semibold text-lg">Tenant: {sub.tenant_id}</h3>
                  <p className="text-gray-600 text-sm mt-1">Status: {sub.status}</p>
                </div>
                <button className="text-primary-600 hover:text-primary-700">Manage →</button>
              </div>
            </Card>
          ))}
        </div>
      )}
    </div>
  );
}
