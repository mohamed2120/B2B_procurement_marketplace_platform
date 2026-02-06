'use client';

import { useEffect, useState } from 'react';
import { diagnostics } from '@/lib/api';
import { useSafeRouter } from '@/lib/useSafeRouter';

interface EventFailure {
  id: string;
  event_name: string;
  direction: string;
  service_name: string;
  payload_json: any;
  error_message: string | null;
  retry_count: number;
  last_attempt_at: string;
  status: string;
}

export default function EventFailuresPage() {
  const router = useSafeRouter();
  const [failures, setFailures] = useState<EventFailure[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchFailures();
    const interval = setInterval(fetchFailures, 10000);
    return () => clearInterval(interval);
  }, []);

  const fetchFailures = async () => {
    try {
      const res = await diagnostics.get('/api/diagnostics/v1/events/failures');
      setFailures(res.data);
    } catch (error) {
      console.error('Failed to fetch event failures:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleRetry = async (id: string) => {
    try {
      await diagnostics.post(`/api/diagnostics/v1/events/failures/${id}/retry`);
      fetchFailures();
    } catch (error) {
      console.error('Failed to retry event:', error);
      alert('Failed to retry event');
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'FAILED':
        return 'text-red-600 bg-red-100';
      case 'RETRYING':
        return 'text-yellow-600 bg-yellow-100';
      case 'RESOLVED':
        return 'text-green-600 bg-green-100';
      default:
        return 'text-gray-600 bg-gray-100';
    }
  };

  if (loading) {
    return <div className="p-8">Loading event failures...</div>;
  }

  return (
    <div className="p-8">
      <div className="mb-6">
        <h1 className="text-3xl font-bold text-gray-900 mb-2">Event Failures</h1>
        <p className="text-gray-600">View and retry failed event publish/consume operations</p>
      </div>

      <div className="bg-white rounded-lg shadow overflow-hidden">
        <table className="min-w-full divide-y divide-gray-200">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Event</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Direction</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Service</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Error</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Retries</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Last Attempt</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Status</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Actions</th>
            </tr>
          </thead>
          <tbody className="bg-white divide-y divide-gray-200">
            {failures.map((failure) => (
              <tr key={failure.id}>
                <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                  {failure.event_name}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  {failure.direction}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  {failure.service_name}
                </td>
                <td className="px-6 py-4 text-sm text-gray-500 max-w-xs truncate">
                  {failure.error_message || 'N/A'}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  {failure.retry_count}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  {new Date(failure.last_attempt_at).toLocaleString()}
                </td>
                <td className="px-6 py-4 whitespace-nowrap">
                  <span className={`px-2 py-1 text-xs font-semibold rounded-full ${getStatusColor(failure.status)}`}>
                    {failure.status}
                  </span>
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm">
                  {failure.status === 'FAILED' && (
                    <button
                      onClick={() => handleRetry(failure.id)}
                      className="text-blue-600 hover:text-blue-800"
                    >
                      Retry
                    </button>
                  )}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
