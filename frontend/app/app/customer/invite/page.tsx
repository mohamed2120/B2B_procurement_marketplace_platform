'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { apiClients } from '@/lib/api';
import { hasRole } from '@/lib/auth';
import Card from '@/components/ui/Card';
import Button from '@/components/ui/Button';

interface DataItem {
  id: string;
  [key: string]: any;
}

export default function Invite Team MemberPage() {
  const router = useRouter();
  const [data, setData] = useState<DataItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    if (!(hasRole('requester') || hasRole('procurement_manager'))) {
      router.push('/app');
      return;
    }
  }, []);

  useEffect(() => {
    if (!(hasRole('requester') || hasRole('procurement_manager'))) {
      router.push('/app');
      return;
    }
  }, []);

  useEffect(() => {
    fetchData();
  }, []);

  const fetchData = async () => {
    setLoading(true);
    setError('');
    // TODO: Implement API call
    setData([
      { id: '1', name: 'Sample Item 1', status: 'active' },
      { id: '2', name: 'Sample Item 2', status: 'pending' },
    ]);
    setLoading(false);
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
        <h1 className="text-3xl font-bold text-gray-900">Invite Team Member</h1>
        <Button onClick={() => alert('Create functionality (TODO)')}>Create New</Button>
      </div>

      {error && (
        <div className="mb-4 bg-yellow-50 border border-yellow-200 text-yellow-700 px-4 py-3 rounded">
          <p className="font-semibold">MVP Pending</p>
          <p className="text-sm mt-1">Using mock data. API integration pending.</p>
        </div>
      )}

      <Card>
        {data.length === 0 ? (
          <p className="text-gray-600 text-center py-8">No data available.</p>
        ) : (
          <div className="overflow-x-auto">
            <table className="min-w-full divide-y divide-gray-200">
              <thead className="bg-gray-50">
                <tr>
                  <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">ID</th>
                  <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Name</th>
                  <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Status</th>
                  <th className="px-4 py-3 text-center text-xs font-medium text-gray-500 uppercase">Actions</th>
                </tr>
              </thead>
              <tbody className="bg-white divide-y divide-gray-200">
                {data.map((item) => (
                  <tr key={item.id}>
                    <td className="px-4 py-3 text-sm text-gray-900">{item.id}</td>
                    <td className="px-4 py-3 text-sm text-gray-900">{item.name || 'N/A'}</td>
                    <td className="px-4 py-3 text-sm">
                      <span className="px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-green-100 text-green-800">
                        {item.status || 'active'}
                      </span>
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
