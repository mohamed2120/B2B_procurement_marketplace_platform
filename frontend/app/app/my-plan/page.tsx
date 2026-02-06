'use client';

import { useEffect, useState } from 'react';
import { apiClients } from '@/lib/api';
import Card from '@/components/ui/Card';

interface Plan {
  id: string;
  name: string;
  code: string;
  description: string;
  price: number;
  currency: string;
  billing_cycle: string;
  entitlements?: Entitlement[];
}

interface Entitlement {
  feature: string;
  limit: number;
  unit: string;
  description: string;
}

interface Subscription {
  id: string;
  tenant_id: string;
  plan_id: string;
  status: string;
  started_at: string;
  expires_at?: string;
  auto_renew: boolean;
  plan?: Plan;
}

interface Usage {
  users: number;
  rfqs: number;
  listings: number;
  storage_gb: number;
}

export default function MyPlan() {
  const [subscription, setSubscription] = useState<Subscription | null>(null);
  const [usage, setUsage] = useState<Usage>({ users: 0, rfqs: 0, listings: 0, storage_gb: 0 });
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchData();
  }, []);

  const fetchData = async () => {
    try {
      const subResponse = await apiClients.billing.get<Subscription>('/api/v1/subscriptions');
      setSubscription(subResponse.data);
      setUsage({ users: 5, rfqs: 12, listings: 45, storage_gb: 2.5 });
    } catch (error) {
      console.error('Failed to fetch plan data:', error);
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return <div>Loading...</div>;
  }

  if (!subscription || !subscription.plan) {
    return (
      <div>
        <h1 className="text-3xl font-bold text-gray-900 mb-6">My Plan</h1>
        <Card title="No Active Plan">
          <p className="text-gray-600 mb-4">You don't have an active subscription plan.</p>
        </Card>
      </div>
    );
  }

  const plan = subscription.plan;
  const renewalDate = subscription.expires_at
    ? new Date(subscription.expires_at).toLocaleDateString()
    : subscription.auto_renew
    ? 'Auto-renewal enabled'
    : 'No expiration';

  const getLimit = (feature: string): number => {
    const entitlement = plan.entitlements?.find((e) => e.feature === feature);
    return entitlement?.limit ?? -1;
  };

  const formatLimit = (limit: number): string => {
    if (limit === -1) return 'Unlimited';
    return limit.toString();
  };

  const getUsagePercentage = (feature: string, current: number): number => {
    const limit = getLimit(feature);
    if (limit === -1) return 0;
    return Math.min((current / limit) * 100, 100);
  };

  return (
    <div>
      <h1 className="text-3xl font-bold text-gray-900 mb-6">My Plan</h1>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div className="lg:col-span-2">
          <Card title="Current Plan">
            <div className="space-y-4">
              <div>
                <h3 className="text-2xl font-bold text-gray-900 mb-2">{plan.name}</h3>
                <p className="text-gray-600 mb-4">{plan.description}</p>
                <div className="flex items-baseline mb-4">
                  <span className="text-4xl font-bold text-gray-900">${plan.price}</span>
                  <span className="text-gray-600 ml-2">/{plan.billing_cycle === 'monthly' ? 'month' : 'year'}</span>
                </div>
              </div>

              <div className="border-t pt-4">
                <div className="flex justify-between items-center mb-2">
                  <span className="text-sm font-medium text-gray-700">Status</span>
                  <span className={`px-3 py-1 rounded-full text-sm font-semibold ${subscription.status === 'active' ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-800'}`}>
                    {subscription.status.charAt(0).toUpperCase() + subscription.status.slice(1)}
                  </span>
                </div>
                <div className="flex justify-between items-center mb-2">
                  <span className="text-sm font-medium text-gray-700">Renewal Date</span>
                  <span className="text-sm text-gray-600">{renewalDate}</span>
                </div>
                <div className="flex justify-between items-center">
                  <span className="text-sm font-medium text-gray-700">Auto-Renew</span>
                  <span className="text-sm text-gray-600">{subscription.auto_renew ? 'Yes' : 'No'}</span>
                </div>
              </div>

              <div className="border-t pt-4">
                <h4 className="font-semibold text-gray-900 mb-3">Plan Features</h4>
                {plan.entitlements && plan.entitlements.length > 0 ? (
                  <ul className="space-y-2">
                    {plan.entitlements.map((entitlement, idx) => (
                      <li key={idx} className="flex items-start">
                        <svg className="w-5 h-5 text-primary-600 mr-2 mt-0.5" fill="currentColor" viewBox="0 0 20 20">
                          <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
                        </svg>
                        <span className="text-gray-700">
                          {entitlement.description || `${entitlement.feature}: ${formatLimit(entitlement.limit)} ${entitlement.unit}`}
                        </span>
                      </li>
                    ))}
                  </ul>
                ) : (
                  <p className="text-gray-500 text-sm">No specific entitlements defined</p>
                )}
              </div>
            </div>
          </Card>
        </div>

        <div>
          <Card title="Usage">
            <div className="space-y-6">
              <div>
                <div className="flex justify-between items-center mb-2">
                  <span className="text-sm font-medium text-gray-700">Users</span>
                  <span className="text-sm text-gray-600">{usage.users} / {formatLimit(getLimit('users'))}</span>
                </div>
                <div className="w-full bg-gray-200 rounded-full h-2">
                  <div className="bg-primary-600 h-2 rounded-full" style={{ width: `${getUsagePercentage('users', usage.users)}%` }} />
                </div>
              </div>

              <div>
                <div className="flex justify-between items-center mb-2">
                  <span className="text-sm font-medium text-gray-700">RFQs</span>
                  <span className="text-sm text-gray-600">{usage.rfqs} / {formatLimit(getLimit('rfqs'))}</span>
                </div>
                <div className="w-full bg-gray-200 rounded-full h-2">
                  <div className="bg-primary-600 h-2 rounded-full" style={{ width: `${getUsagePercentage('rfqs', usage.rfqs)}%` }} />
                </div>
              </div>

              <div>
                <div className="flex justify-between items-center mb-2">
                  <span className="text-sm font-medium text-gray-700">Listings</span>
                  <span className="text-sm text-gray-600">{usage.listings} / {formatLimit(getLimit('listings'))}</span>
                </div>
                <div className="w-full bg-gray-200 rounded-full h-2">
                  <div className="bg-primary-600 h-2 rounded-full" style={{ width: `${getUsagePercentage('listings', usage.listings)}%` }} />
                </div>
              </div>

              <div>
                <div className="flex justify-between items-center mb-2">
                  <span className="text-sm font-medium text-gray-700">Storage</span>
                  <span className="text-sm text-gray-600">{usage.storage_gb.toFixed(1)} GB / {formatLimit(getLimit('storage_gb'))} GB</span>
                </div>
                <div className="w-full bg-gray-200 rounded-full h-2">
                  <div className="bg-primary-600 h-2 rounded-full" style={{ width: `${getUsagePercentage('storage_gb', usage.storage_gb)}%` }} />
                </div>
              </div>
            </div>
          </Card>
        </div>
      </div>
    </div>
  );
}
