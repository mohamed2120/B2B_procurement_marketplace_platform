'use client';

import { useState, useEffect } from 'react';
import { useSafeRouter } from '@/lib/useSafeRouter';
import { isAuthenticated, hasRole } from '@/lib/auth';
import { apiClients } from '@/lib/api';
import Header from '@/components/layout/Header';
import Sidebar from '@/components/layout/Sidebar';
import Card from '@/components/ui/Card';
import Button from '@/components/ui/Button';
import Link from 'next/link';
import { format } from 'date-fns';

interface PR {
  id: string;
  pr_number: string;
  title: string;
  description: string;
  status: string;
  created_at: string;
  total_amount?: number;
}

export default function PRListPage() {
  const router = useSafeRouter();
  const [prs, setPRs] = useState<PR[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!isAuthenticated() || !hasRole('buyer')) {
      router.push('/login');
      return;
    }

    fetchPRs();
  }, [router]);

  const fetchPRs = async () => {
    try {
      const response = await apiClients.procurement.get('/api/v1/purchase-requests?limit=100');
      setPRs(response.data.items || response.data || []);
    } catch (error) {
      console.error('Failed to fetch PRs:', error);
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
          <div className="flex justify-between items-center mb-6">
            <h1 className="text-3xl font-bold text-gray-900">Purchase Requests</h1>
            <Link href="/customer/prs/create">
              <Button>Create PR</Button>
            </Link>
          </div>

          <Card>
            {prs.length === 0 ? (
              <p className="text-gray-500 text-center py-8">No purchase requests found.</p>
            ) : (
              <div className="overflow-x-auto">
                <table className="min-w-full divide-y divide-gray-200">
                  <thead className="bg-gray-50">
                    <tr>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        PR Number
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Title
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Status
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Created
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Actions
                      </th>
                    </tr>
                  </thead>
                  <tbody className="bg-white divide-y divide-gray-200">
                    {prs.map((pr) => (
                      <tr key={pr.id}>
                        <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                          {pr.pr_number}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                          {pr.title}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap">
                          <span
                            className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${
                              pr.status === 'approved'
                                ? 'bg-green-100 text-green-800'
                                : pr.status === 'pending'
                                ? 'bg-yellow-100 text-yellow-800'
                                : 'bg-gray-100 text-gray-800'
                            }`}
                          >
                            {pr.status}
                          </span>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                          {format(new Date(pr.created_at), 'MMM dd, yyyy')}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">
                          <Link href={`/customer/prs/${pr.id}`}>
                            <button className="text-primary-600 hover:text-primary-900">View</button>
                          </Link>
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
