'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { hasRole } from '@/lib/auth';
import Card from '@/components/ui/Card';
import Button from '@/components/ui/Button';

export default function NotificationTemplatesPage() {
  const router = useRouter();
  const [templates, setTemplates] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!(hasRole('admin') || hasRole('super_admin'))) {
      router.push('/app');
      return;
    }
    // Mock data
    setTemplates([
      { id: '1', name: 'Order Confirmation', subject: 'Your order has been confirmed', type: 'email' },
      { id: '2', name: 'Shipment Update', subject: 'Your shipment status has been updated', type: 'email' },
      { id: '3', name: 'RFQ Notification', subject: 'New RFQ available', type: 'email' },
    ]);
    setLoading(false);
  }, []);

  if (loading) {
    return <div>Loading...</div>;
  }

  return (
    <div>
      <h1 className="text-3xl font-bold text-gray-900 mb-6">Notification Templates</h1>
      <Card>
        <p className="text-gray-600">Notification template management (TODO: Implement)</p>
      </Card>
    </div>
  );
}
