'use client';

import { useEffect, useState } from 'react';
import { diagnostics } from '@/lib/api';
import { useSafeRouter } from '@/lib/useSafeRouter';

interface Metric {
  minute_ts: string;
  service_name: string;
  route: string;
  method: string;
  count_total: number;
  count_2xx: number;
  count_4xx: number;
  count_5xx: number;
  p95_ms: number | null;
  avg_ms: number | null;
}

export default function MetricsPage() {
  const router = useSafeRouter();
  const [metrics, setMetrics] = useState<Metric[]>([]);
  const [loading, setLoading] = useState(true);
  const [range, setRange] = useState('1h');

  useEffect(() => {
    fetchMetrics();
  }, [range]);

  const fetchMetrics = async () => {
    try {
      const res = await diagnostics.get(`/api/diagnostics/v1/metrics?range=${range}`);
      setMetrics(res.data);
    } catch (error) {
      console.error('Failed to fetch metrics:', error);
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return <div className="p-8">Loading metrics...</div>;
  }

  return (
    <div className="p-8">
      <div className="mb-6 flex justify-between items-center">
        <div>
          <h1 className="text-3xl font-bold text-gray-900 mb-2">API Metrics</h1>
          <p className="text-gray-600">View API performance and error rates</p>
        </div>
        <select
          value={range}
          onChange={(e) => setRange(e.target.value)}
          className="border rounded px-3 py-2"
        >
          <option value="1h">Last Hour</option>
          <option value="24h">Last 24 Hours</option>
          <option value="7d">Last 7 Days</option>
        </select>
      </div>

      <div className="bg-white rounded-lg shadow overflow-hidden">
        <table className="min-w-full divide-y divide-gray-200">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Time</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Service</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Route</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Method</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Total</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">2xx</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">4xx</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">5xx</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">P95 (ms)</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Avg (ms)</th>
            </tr>
          </thead>
          <tbody className="bg-white divide-y divide-gray-200">
            {metrics.map((metric, idx) => (
              <tr key={idx}>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  {new Date(metric.minute_ts).toLocaleString()}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  {metric.service_name}
                </td>
                <td className="px-6 py-4 text-sm text-gray-500 max-w-xs truncate">
                  {metric.route}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  {metric.method}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  {metric.count_total}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-green-600">
                  {metric.count_2xx}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-yellow-600">
                  {metric.count_4xx}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-red-600">
                  {metric.count_5xx}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  {metric.p95_ms || 'N/A'}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  {metric.avg_ms || 'N/A'}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
