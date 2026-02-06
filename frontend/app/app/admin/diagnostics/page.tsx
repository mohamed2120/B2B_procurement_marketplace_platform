'use client';

import { useEffect, useState } from 'react';
import { diagnostics } from '@/lib/api';
import { useSafeRouter } from '@/lib/useSafeRouter';
import Link from 'next/link';

interface Summary {
  unhealthy_services: number;
  incidents_24h: number;
  event_failures: number;
  critical_incidents: number;
}

export default function DiagnosticsDashboard() {
  const router = useSafeRouter();
  const [summary, setSummary] = useState<Summary | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchSummary();
    const interval = setInterval(fetchSummary, 30000); // Refresh every 30s
    return () => clearInterval(interval);
  }, []);

  const fetchSummary = async () => {
    try {
      const res = await diagnostics.get('/api/diagnostics/v1/summary');
      setSummary(res.data);
    } catch (error) {
      console.error('Failed to fetch diagnostics summary:', error);
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return <div className="p-8">Loading diagnostics...</div>;
  }

  return (
    <div className="p-8">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900 mb-2">Diagnostics Center</h1>
        <p className="text-gray-600">Platform health monitoring and debugging</p>
      </div>

      {/* Summary Cards */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        <Link href="/app/admin/diagnostics/services">
          <div className="bg-white rounded-lg shadow p-6 hover:shadow-lg transition">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-gray-600">Unhealthy Services</p>
                <p className="text-3xl font-bold text-red-600 mt-2">
                  {summary?.unhealthy_services || 0}
                </p>
              </div>
              <div className="text-red-500 text-4xl">⚠️</div>
            </div>
          </div>
        </Link>

        <Link href="/app/admin/diagnostics/incidents">
          <div className="bg-white rounded-lg shadow p-6 hover:shadow-lg transition">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-gray-600">Incidents (24h)</p>
                <p className="text-3xl font-bold text-orange-600 mt-2">
                  {summary?.incidents_24h || 0}
                </p>
              </div>
              <div className="text-orange-500 text-4xl">📊</div>
            </div>
          </div>
        </Link>

        <Link href="/app/admin/diagnostics/events">
          <div className="bg-white rounded-lg shadow p-6 hover:shadow-lg transition">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-gray-600">Event Failures</p>
                <p className="text-3xl font-bold text-yellow-600 mt-2">
                  {summary?.event_failures || 0}
                </p>
              </div>
              <div className="text-yellow-500 text-4xl">⚡</div>
            </div>
          </div>
        </Link>

        <div className="bg-white rounded-lg shadow p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-600">Critical Incidents</p>
              <p className="text-3xl font-bold text-red-700 mt-2">
                {summary?.critical_incidents || 0}
              </p>
            </div>
            <div className="text-red-700 text-4xl">🚨</div>
          </div>
        </div>
      </div>

      {/* Quick Links */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <Link href="/app/admin/diagnostics/services">
          <div className="bg-white rounded-lg shadow p-6 hover:shadow-lg transition cursor-pointer">
            <h3 className="text-xl font-semibold mb-2">Services</h3>
            <p className="text-gray-600">View service health and heartbeats</p>
          </div>
        </Link>

        <Link href="/app/admin/diagnostics/incidents">
          <div className="bg-white rounded-lg shadow p-6 hover:shadow-lg transition cursor-pointer">
            <h3 className="text-xl font-semibold mb-2">Incidents</h3>
            <p className="text-gray-600">View and resolve incidents</p>
          </div>
        </Link>

        <Link href="/app/admin/diagnostics/events">
          <div className="bg-white rounded-lg shadow p-6 hover:shadow-lg transition cursor-pointer">
            <h3 className="text-xl font-semibold mb-2">Event Failures</h3>
            <p className="text-gray-600">View and retry failed events</p>
          </div>
        </Link>

        <Link href="/app/admin/diagnostics/metrics">
          <div className="bg-white rounded-lg shadow p-6 hover:shadow-lg transition cursor-pointer">
            <h3 className="text-xl font-semibold mb-2">Metrics</h3>
            <p className="text-gray-600">View API metrics and performance</p>
          </div>
        </Link>
      </div>
    </div>
  );
}
