'use client';

import { useState, useEffect } from 'react';
import { apiClients } from '@/lib/api';
import { hasRole } from '@/lib/auth';
import Card from '@/components/ui/Card';
import Input from '@/components/ui/Input';
import Button from '@/components/ui/Button';

interface AuditLog {
  id: string;
  action: string;
  resource: string;
  user_id: string;
  user_email?: string;
  tenant_id?: string;
  timestamp: string;
  ip_address?: string;
  user_agent?: string;
  details?: any;
  status?: 'success' | 'failure';
}

export default function AuditLogsPage() {
  const [logs, setLogs] = useState<AuditLog[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [filters, setFilters] = useState({
    action: '',
    resource: '',
    user_id: '',
    start_date: '',
    end_date: '',
  });

  // Check if user has admin role
  if (!hasRole('admin') && !hasRole('super_admin')) {
    return (
      <div>
        <Card>
          <div className="bg-yellow-50 border border-yellow-200 text-yellow-700 px-4 py-3 rounded">
            You don't have permission to access this page.
          </div>
        </Card>
      </div>
    );
  }

  useEffect(() => {
    fetchLogs();
  }, []);

  const fetchLogs = async () => {
    try {
      setLoading(true);
      setError('');
      
      // TODO: Replace with actual API endpoint when available
      // try {
      //   const params = new URLSearchParams();
      //   if (filters.action) params.append('action', filters.action);
      //   if (filters.resource) params.append('resource', filters.resource);
      //   if (filters.user_id) params.append('user_id', filters.user_id);
      //   if (filters.start_date) params.append('start_date', filters.start_date);
      //   if (filters.end_date) params.append('end_date', filters.end_date);
      //   
      //   const response = await apiClients.identity.get(`/api/v1/audit-logs?${params.toString()}`);
      //   setLogs(response.data || []);
      // } catch (apiError: any) {
      //   if (apiError.response?.status === 404 || apiError.response?.status === 501) {
      //     console.warn('Audit logs API not available, using mock data');
      //     setLogs(getMockLogs());
      //   } else {
      //     throw apiError;
      //   }
      // }
      
      // Using mock data for now
      setLogs(getMockLogs());
    } catch (err: any) {
      console.error('Failed to fetch audit logs:', err);
      setError(err.response?.data?.error || 'Failed to fetch audit logs. Using mock data.');
      setLogs(getMockLogs());
    } finally {
      setLoading(false);
    }
  };

  // Mock audit logs for development until API is available
  const getMockLogs = (): AuditLog[] => {
    return [
      {
        id: '1',
        action: 'login',
        resource: 'auth',
        user_id: '1',
        user_email: 'buyer.requester@demo.com',
        tenant_id: '00000000-0000-0000-0000-000000000001',
        timestamp: new Date().toISOString(),
        ip_address: '192.168.1.1',
        user_agent: 'Mozilla/5.0...',
        status: 'success',
      },
      {
        id: '2',
        action: 'create',
        resource: 'pr',
        user_id: '1',
        user_email: 'buyer.requester@demo.com',
        tenant_id: '00000000-0000-0000-0000-000000000001',
        timestamp: new Date(Date.now() - 3600000).toISOString(),
        ip_address: '192.168.1.1',
        details: { pr_id: 'pr-123', title: 'Office Supplies' },
        status: 'success',
      },
      {
        id: '3',
        action: 'approve',
        resource: 'pr',
        user_id: '2',
        user_email: 'buyer.procurement@demo.com',
        tenant_id: '00000000-0000-0000-0000-000000000001',
        timestamp: new Date(Date.now() - 7200000).toISOString(),
        ip_address: '192.168.1.2',
        details: { pr_id: 'pr-123', approved: true },
        status: 'success',
      },
      {
        id: '4',
        action: 'submit',
        resource: 'quote',
        user_id: '3',
        user_email: 'supplier@demo.com',
        tenant_id: '00000000-0000-0000-0000-000000000002',
        timestamp: new Date(Date.now() - 10800000).toISOString(),
        ip_address: '192.168.1.3',
        details: { quote_id: 'quote-456', rfq_id: 'rfq-789' },
        status: 'success',
      },
      {
        id: '5',
        action: 'login',
        resource: 'auth',
        user_id: '4',
        user_email: 'admin@demo.com',
        tenant_id: '00000000-0000-0000-0000-000000000000',
        timestamp: new Date(Date.now() - 14400000).toISOString(),
        ip_address: '192.168.1.4',
        status: 'success',
      },
    ];
  };

  const handleFilterChange = (field: string, value: string) => {
    setFilters(prev => ({ ...prev, [field]: value }));
  };

  const handleApplyFilters = () => {
    fetchLogs();
  };

  const handleClearFilters = () => {
    setFilters({
      action: '',
      resource: '',
      user_id: '',
      start_date: '',
      end_date: '',
    });
    fetchLogs();
  };

  const filteredLogs = logs.filter(log => {
    if (filters.action && !log.action.toLowerCase().includes(filters.action.toLowerCase())) return false;
    if (filters.resource && !log.resource.toLowerCase().includes(filters.resource.toLowerCase())) return false;
    if (filters.user_id && log.user_id !== filters.user_id) return false;
    if (filters.start_date && new Date(log.timestamp) < new Date(filters.start_date)) return false;
    if (filters.end_date && new Date(log.timestamp) > new Date(filters.end_date)) return false;
    return true;
  });

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600 mx-auto mb-4"></div>
          <p className="text-gray-600">Loading audit logs...</p>
        </div>
      </div>
    );
  }

  return (
    <div>
      <div className="mb-6">
        <h1 className="text-3xl font-bold text-gray-900">Audit Logs</h1>
        <p className="text-gray-600 mt-2">View system activity and user actions</p>
      </div>

      {error && (
        <div className="mb-4 bg-yellow-50 border border-yellow-200 text-yellow-700 px-4 py-3 rounded">
          {error}
          <div className="text-sm mt-1">
            Note: Using mock data until audit logs API is implemented.
          </div>
        </div>
      )}

      {/* Filters */}
      <Card title="Filters" className="mb-6">
        <div className="grid grid-cols-1 md:grid-cols-3 lg:grid-cols-5 gap-4">
          <Input
            label="Action"
            value={filters.action}
            onChange={(e) => handleFilterChange('action', e.target.value)}
            placeholder="e.g., login, create"
          />
          <Input
            label="Resource"
            value={filters.resource}
            onChange={(e) => handleFilterChange('resource', e.target.value)}
            placeholder="e.g., auth, pr"
          />
          <Input
            label="User ID"
            value={filters.user_id}
            onChange={(e) => handleFilterChange('user_id', e.target.value)}
            placeholder="User ID"
          />
          <Input
            label="Start Date"
            type="date"
            value={filters.start_date}
            onChange={(e) => handleFilterChange('start_date', e.target.value)}
          />
          <Input
            label="End Date"
            type="date"
            value={filters.end_date}
            onChange={(e) => handleFilterChange('end_date', e.target.value)}
          />
        </div>
        <div className="flex gap-2 mt-4">
          <Button onClick={handleApplyFilters}>Apply Filters</Button>
          <Button onClick={handleClearFilters} variant="secondary">Clear</Button>
        </div>
      </Card>

      {/* Logs Table */}
      <Card title={`Audit Logs (${filteredLogs.length})`}>
        {filteredLogs.length === 0 ? (
          <p className="text-gray-600">No audit logs found.</p>
        ) : (
          <div className="overflow-x-auto">
            <table className="min-w-full divide-y divide-gray-200">
              <thead className="bg-gray-50">
                <tr>
                  <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Timestamp</th>
                  <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">User</th>
                  <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Action</th>
                  <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Resource</th>
                  <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Status</th>
                  <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">IP Address</th>
                  <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Details</th>
                </tr>
              </thead>
              <tbody className="bg-white divide-y divide-gray-200">
                {filteredLogs.map((log) => (
                  <tr key={log.id} className="hover:bg-gray-50">
                    <td className="px-4 py-3 text-sm text-gray-600">
                      {new Date(log.timestamp).toLocaleString()}
                    </td>
                    <td className="px-4 py-3 text-sm">
                      <div>
                        <div className="text-gray-900">{log.user_email || log.user_id}</div>
                        {log.tenant_id && (
                          <div className="text-xs text-gray-500">{log.tenant_id.slice(0, 8)}...</div>
                        )}
                      </div>
                    </td>
                    <td className="px-4 py-3 text-sm text-gray-900">{log.action}</td>
                    <td className="px-4 py-3 text-sm text-gray-600">{log.resource}</td>
                    <td className="px-4 py-3 text-sm">
                      <span className={`px-2 py-1 rounded text-xs font-semibold ${
                        log.status === 'success' ? 'bg-green-100 text-green-800' :
                        log.status === 'failure' ? 'bg-red-100 text-red-800' :
                        'bg-gray-100 text-gray-800'
                      }`}>
                        {log.status || 'unknown'}
                      </span>
                    </td>
                    <td className="px-4 py-3 text-sm text-gray-600">{log.ip_address || '-'}</td>
                    <td className="px-4 py-3 text-sm text-gray-600">
                      {log.details ? (
                        <details className="cursor-pointer">
                          <summary className="text-primary-600 hover:text-primary-700">View</summary>
                          <pre className="mt-2 text-xs bg-gray-50 p-2 rounded overflow-auto max-w-xs">
                            {JSON.stringify(log.details, null, 2)}
                          </pre>
                        </details>
                      ) : (
                        '-'
                      )}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </Card>

      <div className="mt-4 text-sm text-gray-500">
        <p>Note: Audit log API not yet implemented. Showing mock data.</p>
        <p className="mt-1">To implement real audit logging, backend API endpoints need to be created.</p>
      </div>
    </div>
  );
}
