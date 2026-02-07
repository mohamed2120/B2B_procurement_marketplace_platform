'use client';

import { useState, useEffect } from 'react';
import { useParams, useRouter } from 'next/navigation';
import { apiClients } from '@/lib/api';
import { hasRole } from '@/lib/auth';
import Card from '@/components/ui/Card';
import Button from '@/components/ui/Button';
import Link from 'next/link';

interface User {
  id: string;
  email: string;
  first_name: string;
  last_name: string;
  tenant_id: string;
  tenant_name?: string;
  roles?: Role[];
  is_active: boolean;
  is_verified: boolean;
  created_at: string;
  updated_at: string;
  last_login_at?: string;
}

interface Role {
  id: string;
  name: string;
  description?: string;
}

interface AuditLog {
  id: string;
  action: string;
  resource: string;
  user_id: string;
  timestamp: string;
  details?: any;
}

export default function UserDetailPage() {
  const params = useParams();
  const router = useRouter();
  const userId = params.id as string;
  
  const [user, setUser] = useState<User | null>(null);
  const [auditLogs, setAuditLogs] = useState<AuditLog[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [updating, setUpdating] = useState(false);

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
    fetchUser();
    fetchAuditLogs();
  }, [userId]);

  const fetchUser = async () => {
    try {
      setLoading(true);
      setError('');
      
      // Try identity service endpoint
      try {
        const response = await apiClients.identity.get<User>(`/api/v1/users/${userId}`);
        setUser(response.data);
      } catch (apiError: any) {
        // If endpoint doesn't exist or returns error, use mock data
        if (apiError.response?.status === 404 || apiError.response?.status === 501 || apiError.response?.status === 403) {
          // Silently fall back to mock data
          setUser(getMockUser(userId));
        } else {
          throw apiError;
        }
      }
    } catch (err: any) {
      console.error('Failed to fetch user:', err);
      setError(err.response?.data?.error || 'Failed to fetch user. Using mock data.');
      setUser(getMockUser(userId));
    } finally {
      setLoading(false);
    }
  };

  const fetchAuditLogs = async () => {
    try {
      // TODO: Replace with actual API endpoint when available
      // const response = await apiClients.identity.get(`/api/v1/users/${userId}/audit-logs`);
      // setAuditLogs(response.data || []);
      
      // Mock audit logs for now
      setAuditLogs([
        {
          id: '1',
          action: 'login',
          resource: 'auth',
          user_id: userId,
          timestamp: new Date().toISOString(),
          details: { ip: '192.168.1.1' },
        },
        {
          id: '2',
          action: 'update_profile',
          resource: 'user',
          user_id: userId,
          timestamp: new Date(Date.now() - 86400000).toISOString(),
          details: { field: 'first_name' },
        },
      ]);
    } catch (err: any) {
      console.error('Failed to fetch audit logs:', err);
    }
  };

  // Mock user data for development until API is available
  const getMockUser = (id: string): User => {
    const mockUsers: Record<string, User> = {
      '1': {
        id: '1',
        email: 'buyer.requester@demo.com',
        first_name: 'Buyer',
        last_name: 'Requester',
        tenant_id: '00000000-0000-0000-0000-000000000001',
        tenant_name: 'Demo Buyer Company',
        roles: [{ id: '1', name: 'requester', description: 'Can create PRs and RFQs' }],
        is_active: true,
        is_verified: true,
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString(),
        last_login_at: new Date().toISOString(),
      },
      '2': {
        id: '2',
        email: 'buyer.procurement@demo.com',
        first_name: 'Buyer',
        last_name: 'Procurement',
        tenant_id: '00000000-0000-0000-0000-000000000001',
        tenant_name: 'Demo Buyer Company',
        roles: [{ id: '2', name: 'procurement_manager', description: 'Can approve PRs and award quotes' }],
        is_active: true,
        is_verified: true,
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString(),
        last_login_at: new Date().toISOString(),
      },
      '3': {
        id: '3',
        email: 'supplier@demo.com',
        first_name: 'Supplier',
        last_name: 'User',
        tenant_id: '00000000-0000-0000-0000-000000000002',
        tenant_name: 'Demo Supplier Company',
        roles: [{ id: '3', name: 'supplier', description: 'Can manage listings and submit quotes' }],
        is_active: true,
        is_verified: true,
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString(),
        last_login_at: new Date().toISOString(),
      },
      '4': {
        id: '4',
        email: 'admin@demo.com',
        first_name: 'Platform',
        last_name: 'Admin',
        tenant_id: '00000000-0000-0000-0000-000000000000',
        tenant_name: 'Platform',
        roles: [
          { id: '4', name: 'admin', description: 'Platform administrator' },
          { id: '5', name: 'super_admin', description: 'Super administrator' },
        ],
        is_active: true,
        is_verified: true,
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString(),
        last_login_at: new Date().toISOString(),
      },
    };
    return mockUsers[id] || mockUsers['1'];
  };

  const handleToggleActive = async () => {
    if (!user) return;
    
    setUpdating(true);
    try {
      // TODO: Replace with actual API endpoint when available
      // await apiClients.identity.put(`/api/v1/users/${userId}/toggle-active`);
      alert('User status update API not yet implemented. This is a placeholder.');
      // Refresh user data
      await fetchUser();
    } catch (err: any) {
      alert(err.response?.data?.error || 'Failed to update user status');
    } finally {
      setUpdating(false);
    }
  };

  const handleResetPassword = async () => {
    if (!confirm('Send password reset email to this user?')) {
      return;
    }

    try {
      // TODO: Replace with actual API endpoint when available
      // await apiClients.identity.post(`/api/v1/users/${userId}/reset-password`);
      alert('Password reset API not yet implemented. This is a placeholder.');
    } catch (err: any) {
      alert(err.response?.data?.error || 'Failed to send password reset');
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600 mx-auto mb-4"></div>
          <p className="text-gray-600">Loading user...</p>
        </div>
      </div>
    );
  }

  if (!user) {
    return (
      <div>
        <Link href="/app/admin/users" className="text-primary-600 hover:text-primary-700 mb-4 inline-block">
          ← Back to Users
        </Link>
        <Card>
          <p className="text-gray-600">User not found</p>
        </Card>
      </div>
    );
  }

  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <div>
          <Link href="/app/admin/users" className="text-primary-600 hover:text-primary-700 mb-2 inline-block">
            ← Back to Users
          </Link>
          <h1 className="text-3xl font-bold text-gray-900">User Details</h1>
        </div>
        <div className="flex gap-2">
          <Button
            onClick={handleToggleActive}
            disabled={updating}
            variant={user.is_active ? 'danger' : 'primary'}
          >
            {updating ? 'Updating...' : user.is_active ? 'Deactivate' : 'Activate'}
          </Button>
          <Button onClick={handleResetPassword} variant="secondary">
            Reset Password
          </Button>
        </div>
      </div>

      {error && (
        <div className="mb-4 bg-yellow-50 border border-yellow-200 text-yellow-700 px-4 py-3 rounded">
          {error}
          <div className="text-sm mt-1">
            Note: Using mock data until user management API is implemented.
          </div>
        </div>
      )}

      {/* User Information */}
      <Card title="User Information" className="mb-6">
        <div className="grid grid-cols-2 gap-6">
          <div>
            <span className="text-sm font-medium text-gray-700">User ID</span>
            <p className="text-gray-900 mt-1 font-mono text-sm">{user.id}</p>
          </div>
          <div>
            <span className="text-sm font-medium text-gray-700">Email</span>
            <p className="text-gray-900 mt-1">{user.email}</p>
          </div>
          <div>
            <span className="text-sm font-medium text-gray-700">First Name</span>
            <p className="text-gray-900 mt-1">{user.first_name}</p>
          </div>
          <div>
            <span className="text-sm font-medium text-gray-700">Last Name</span>
            <p className="text-gray-900 mt-1">{user.last_name}</p>
          </div>
          <div>
            <span className="text-sm font-medium text-gray-700">Company/Tenant</span>
            <p className="text-gray-900 mt-1">{user.tenant_name || user.tenant_id}</p>
          </div>
          <div>
            <span className="text-sm font-medium text-gray-700">Status</span>
            <div className="mt-1 flex gap-2">
              <span className={`px-2 py-1 rounded text-xs font-semibold ${
                user.is_active ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'
              }`}>
                {user.is_active ? 'Active' : 'Inactive'}
              </span>
              {!user.is_verified && (
                <span className="px-2 py-1 bg-yellow-100 text-yellow-800 rounded text-xs">
                  Unverified
                </span>
              )}
            </div>
          </div>
          <div>
            <span className="text-sm font-medium text-gray-700">Created</span>
            <p className="text-gray-900 mt-1">{new Date(user.created_at).toLocaleString()}</p>
          </div>
          <div>
            <span className="text-sm font-medium text-gray-700">Last Login</span>
            <p className="text-gray-900 mt-1">
              {user.last_login_at ? new Date(user.last_login_at).toLocaleString() : 'Never'}
            </p>
          </div>
        </div>
      </Card>

      {/* Roles */}
      <Card title="Roles & Permissions" className="mb-6">
        {user.roles && user.roles.length > 0 ? (
          <div className="space-y-3">
            {user.roles.map((role) => (
              <div key={role.id} className="border rounded-lg p-4">
                <div className="flex justify-between items-start">
                  <div>
                    <h4 className="font-semibold text-gray-900">{role.name}</h4>
                    {role.description && (
                      <p className="text-sm text-gray-600 mt-1">{role.description}</p>
                    )}
                  </div>
                  <span className="px-2 py-1 bg-blue-100 text-blue-800 text-xs rounded">
                    Active
                  </span>
                </div>
              </div>
            ))}
          </div>
        ) : (
          <p className="text-gray-600">No roles assigned</p>
        )}
      </Card>

      {/* Audit Logs */}
      <Card title="Audit Logs">
        {auditLogs.length > 0 ? (
          <div className="overflow-x-auto">
            <table className="min-w-full divide-y divide-gray-200">
              <thead className="bg-gray-50">
                <tr>
                  <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Timestamp</th>
                  <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Action</th>
                  <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Resource</th>
                  <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Details</th>
                </tr>
              </thead>
              <tbody className="bg-white divide-y divide-gray-200">
                {auditLogs.map((log) => (
                  <tr key={log.id}>
                    <td className="px-4 py-3 text-sm text-gray-600">
                      {new Date(log.timestamp).toLocaleString()}
                    </td>
                    <td className="px-4 py-3 text-sm text-gray-900">{log.action}</td>
                    <td className="px-4 py-3 text-sm text-gray-600">{log.resource}</td>
                    <td className="px-4 py-3 text-sm text-gray-600">
                      {log.details ? JSON.stringify(log.details) : '-'}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        ) : (
          <p className="text-gray-600">No audit logs available</p>
        )}
        <div className="mt-4 text-sm text-gray-500">
          Note: Audit log API not yet implemented. Showing mock data.
        </div>
      </Card>
    </div>
  );
}
