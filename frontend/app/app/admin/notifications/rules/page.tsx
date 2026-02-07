'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { hasRole } from '@/lib/auth';
import Card from '@/components/ui/Card';
import Button from '@/components/ui/Button';

export default function NotificationRulesPage() {
  const router = useRouter();
  const [rules, setRules] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!(hasRole('admin') || hasRole('super_admin'))) {
      router.push('/app');
      return;
    }
    // Mock data
    setRules([
      { id: '1', event: 'order.created', template: 'Order Confirmation', enabled: true },
      { id: '2', event: 'shipment.updated', template: 'Shipment Update', enabled: true },
      { id: '3', event: 'rfq.created', template: 'RFQ Notification', enabled: true },
    ]);
    setLoading(false);
  }, []);

  if (loading) {
    return <div>Loading...</div>;
  }

  return (
    <div>
      <h1 className="text-3xl font-bold text-gray-900 mb-6">Notification Rules</h1>
      <Card>
        <p className="text-gray-600">Notification rule management (TODO: Implement)</p>
      </Card>
    </div>
  );
}
