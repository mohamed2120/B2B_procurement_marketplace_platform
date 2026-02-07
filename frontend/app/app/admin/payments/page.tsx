'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { apiClients } from '@/lib/api';
import { hasRole } from '@/lib/auth';
import Card from '@/components/ui/Card';
import Button from '@/components/ui/Button';

interface Payment {
  id: string;
  tenant_id: string;
  amount: number;
  currency: string;
  status: string;
  payment_method: string;
  created_at: string;
}

export default function PaymentsPage() {
  const router = useRouter();
  const [payments, setPayments] = useState<Payment[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    if (!(hasRole('admin') || hasRole('super_admin'))) {
      router.push('/app');
      return;
    }
    fetchPayments();
  }, []);

  const fetchPayments = async () => {
    setLoading(true);
    setError('');
    try {
      const response = await apiClients.billing.get('/api/v1/billing/v1/payments');
      setPayments(response.data || []);
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to load payments');
      // Mock data
      setPayments([
        {
          id: '1',
          tenant_id: '00000000-0000-0000-0000-000000000001',
          amount: 299.00,
          currency: 'USD',
          status: 'succeeded',
          payment_method: 'credit_card',
          created_at: '2024-01-15T10:00:00Z',
        },
        {
          id: '2',
          tenant_id: '00000000-0000-0000-0000-000000000001',
          amount: 99.00,
          currency: 'USD',
          status: 'pending',
          payment_method: 'bank_transfer',
          created_at: '2024-01-20T10:00:00Z',
        },
      ]);
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600 mx-auto mb-4"></div>
          <p className="text-gray-600">Loading...</p>
        </div>
      </div>
    );
  }

  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-3xl font-bold text-gray-900">Payments</h1>
      </div>

      {error && (
        <div className="mb-4 bg-yellow-50 border border-yellow-200 text-yellow-700 px-4 py-3 rounded">
          <p className="font-semibold">Note</p>
          <p className="text-sm mt-1">Using mock data. API integration pending.</p>
        </div>
      )}

      <Card>
        {payments.length === 0 ? (
          <p className="text-gray-600 text-center py-8">No payments found.</p>
        ) : (
          <div className="overflow-x-auto">
            <table className="min-w-full divide-y divide-gray-200">
              <thead className="bg-gray-50">
                <tr>
                  <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">ID</th>
                  <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Tenant</th>
                  <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Amount</th>
                  <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Status</th>
                  <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Method</th>
                  <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Date</th>
                  <th className="px-4 py-3 text-center text-xs font-medium text-gray-500 uppercase">Actions</th>
                </tr>
              </thead>
              <tbody className="bg-white divide-y divide-gray-200">
                {payments.map((payment) => (
                  <tr key={payment.id}>
                    <td className="px-4 py-3 text-sm text-gray-900">{payment.id.slice(0, 8)}...</td>
                    <td className="px-4 py-3 text-sm text-gray-900">{payment.tenant_id.slice(0, 8)}...</td>
                    <td className="px-4 py-3 text-sm text-gray-900">
                      {payment.currency} {payment.amount.toFixed(2)}
                    </td>
                    <td className="px-4 py-3 text-sm">
                      <span className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${
                        payment.status === 'succeeded' ? 'bg-green-100 text-green-800' :
                        payment.status === 'pending' ? 'bg-yellow-100 text-yellow-800' :
                        'bg-red-100 text-red-800'
                      }`}>
                        {payment.status}
                      </span>
                    </td>
                    <td className="px-4 py-3 text-sm text-gray-900">{payment.payment_method}</td>
                    <td className="px-4 py-3 text-sm text-gray-900">
                      {new Date(payment.created_at).toLocaleDateString()}
                    </td>
                    <td className="px-4 py-3 text-center text-sm font-medium">
                      <Button size="sm" variant="secondary" onClick={() => alert('View details (TODO)')}>
                        View
                      </Button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </Card>
    </div>
  );
}
