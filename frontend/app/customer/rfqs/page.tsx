'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { isAuthenticated } from '@/lib/auth';
import { apiClients } from '@/lib/api';
import Header from '@/components/layout/Header';
import Sidebar from '@/components/layout/Sidebar';
import Card from '@/components/ui/Card';
import { format } from 'date-fns';

interface RFQ {
  id: string;
  rfq_number: string;
  title: string;
  due_date: string;
  status: string;
  pr_id: string;
}

export default function RFQListPage() {
  const router = useRouter();
  const [rfqs, setRFQs] = useState<RFQ[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!isAuthenticated()) {
      router.push('/login');
      return;
    }

    fetchRFQs();
  }, [router]);

  const fetchRFQs = async () => {
    try {
      const response = await apiClients.procurement.get('/api/v1/rfqs?limit=100');
      setRFQs(response.data.items || response.data || []);
    } catch (error) {
      console.error('Failed to fetch RFQs:', error);
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return <div>Loading...</div>;
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <Header />
      <div className="flex">
        <Sidebar />
        <main className="flex-1 p-8">
          <h1 className="text-3xl font-bold text-gray-900 mb-6">RFQs</h1>

          <Card>
            {rfqs.length === 0 ? (
              <p className="text-gray-500 text-center py-8">No RFQs found.</p>
            ) : (
              <div className="overflow-x-auto">
                <table className="min-w-full divide-y divide-gray-200">
                  <thead className="bg-gray-50">
                    <tr>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">RFQ Number</th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Title</th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Due Date</th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Status</th>
                    </tr>
                  </thead>
                  <tbody className="bg-white divide-y divide-gray-200">
                    {rfqs.map((rfq) => (
                      <tr key={rfq.id}>
                        <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">{rfq.rfq_number}</td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm">{rfq.title}</td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm">
                          {format(new Date(rfq.due_date), 'MMM dd, yyyy')}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap">
                          <span className="px-2 py-1 text-xs rounded-full bg-blue-100 text-blue-800">
                            {rfq.status}
                          </span>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}
          </Card>
        </main>
      </div>
    </div>
  );
}
