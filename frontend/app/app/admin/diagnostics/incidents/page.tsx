'use client';

import { useEffect, useState } from 'react';
import { diagnostics } from '@/lib/api';
import { useSafeRouter } from '@/lib/useSafeRouter';

interface Incident {
  id: string;
  severity: string;
  service_name: string;
  category: string;
  error_code: string | null;
  title: string;
  details_json: any;
  request_id: string | null;
  tenant_id: string | null;
  user_id: string | null;
  occurred_at: string;
  resolved_at: string | null;
  resolution_notes: string | null;
}

export default function IncidentsPage() {
  const router = useSafeRouter();
  const [incidents, setIncidents] = useState<Incident[]>([]);
  const [loading, setLoading] = useState(true);
  const [filters, setFilters] = useState({
    severity: '',
    category: '',
    service_name: '',
    resolved: '',
  });

  useEffect(() => {
    fetchIncidents();
  }, [filters]);

  const fetchIncidents = async () => {
    try {
      const params = new URLSearchParams();
      if (filters.severity) params.append('severity', filters.severity);
      if (filters.category) params.append('category', filters.category);
      if (filters.service_name) params.append('service_name', filters.service_name);
      if (filters.resolved) params.append('resolved', filters.resolved);

      const res = await diagnostics.get(`/api/diagnostics/v1/incidents?${params.toString()}`);
      setIncidents(res.data);
    } catch (error) {
      console.error('Failed to fetch incidents:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleResolve = async (id: string, notes: string) => {
    try {
      await diagnostics.post(`/api/diagnostics/v1/incidents/${id}/resolve`, { notes });
      fetchIncidents();
    } catch (error) {
      console.error('Failed to resolve incident:', error);
      alert('Failed to resolve incident');
    }
  };

  const getSeverityColor = (severity: string) => {
    switch (severity) {
      case 'CRITICAL':
        return 'text-red-700 bg-red-100';
      case 'ERROR':
        return 'text-red-600 bg-red-50';
      case 'WARN':
        return 'text-yellow-600 bg-yellow-50';
      default:
        return 'text-blue-600 bg-blue-50';
    }
  };

  if (loading) {
    return <div className="p-8">Loading incidents...</div>;
  }

  return (
    <div className="p-8">
      <div className="mb-6">
        <h1 className="text-3xl font-bold text-gray-900 mb-2">Incidents</h1>
        <p className="text-gray-600">View and manage platform incidents</p>
      </div>

      {/* Filters */}
      <div className="bg-white rounded-lg shadow p-4 mb-6">
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
          <select
            value={filters.severity}
            onChange={(e) => setFilters({ ...filters, severity: e.target.value })}
            className="border rounded px-3 py-2"
          >
            <option value="">All Severities</option>
            <option value="CRITICAL">Critical</option>
            <option value="ERROR">Error</option>
            <option value="WARN">Warning</option>
            <option value="INFO">Info</option>
          </select>

          <select
            value={filters.category}
            onChange={(e) => setFilters({ ...filters, category: e.target.value })}
            className="border rounded px-3 py-2"
          >
            <option value="">All Categories</option>
            <option value="DB">Database</option>
            <option value="AUTH">Authentication</option>
            <option value="EVENT">Event</option>
            <option value="API">API</option>
            <option value="FILE">File</option>
            <option value="SEARCH">Search</option>
          </select>

          <input
            type="text"
            placeholder="Service name"
            value={filters.service_name}
            onChange={(e) => setFilters({ ...filters, service_name: e.target.value })}
            className="border rounded px-3 py-2"
          />

          <select
            value={filters.resolved}
            onChange={(e) => setFilters({ ...filters, resolved: e.target.value })}
            className="border rounded px-3 py-2"
          >
            <option value="">All</option>
            <option value="false">Unresolved</option>
            <option value="true">Resolved</option>
          </select>
        </div>
      </div>

      {/* Incidents Table */}
      <div className="bg-white rounded-lg shadow overflow-hidden">
        <table className="min-w-full divide-y divide-gray-200">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Severity</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Service</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Category</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Title</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Error Code</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Occurred</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Status</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Actions</th>
            </tr>
          </thead>
          <tbody className="bg-white divide-y divide-gray-200">
            {incidents.map((incident) => (
              <tr key={incident.id}>
                <td className="px-6 py-4 whitespace-nowrap">
                  <span className={`px-2 py-1 text-xs font-semibold rounded-full ${getSeverityColor(incident.severity)}`}>
                    {incident.severity}
                  </span>
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  {incident.service_name}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  {incident.category}
                </td>
                <td className="px-6 py-4 text-sm text-gray-900">
                  {incident.title}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  {incident.error_code || 'N/A'}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  {new Date(incident.occurred_at).toLocaleString()}
                </td>
                <td className="px-6 py-4 whitespace-nowrap">
                  {incident.resolved_at ? (
                    <span className="text-green-600">✓ Resolved</span>
                  ) : (
                    <span className="text-red-600">● Open</span>
                  )}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm">
                  {!incident.resolved_at && (
                    <button
                      onClick={() => {
                        const notes = prompt('Resolution notes:');
                        if (notes) handleResolve(incident.id, notes);
                      }}
                      className="text-blue-600 hover:text-blue-800"
                    >
                      Resolve
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
