'use client';

import { useState, useEffect } from 'react';
import { apiClients } from '@/lib/api';
import { hasRole } from '@/lib/auth';
import Card from '@/components/ui/Card';
import Button from '@/components/ui/Button';
import Input from '@/components/ui/Input';
import Link from 'next/link';
import { useRouter } from 'next/navigation';

interface User {
  id: string;
  email: string;
  first_name: string;
  last_name: string;
  tenant_id: string;
  tenant_name?: string;
  roles?: string[];
  is_active: boolean;
  is_verified: boolean;
  created_at: string;
  last_login_at?: string;
}

export default function AdminUsersPage() {
  const router = useRouter();
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [searchTerm, setSearchTerm] = useState('');

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
    fetchUsers();
  }, []);

  const fetchUsers = async () => {
    try {
      setLoading(true);
      setError('');
      
      // TODO: Replace with actual API endpoint when available
      // For now, try identity service endpoint
      try {
        const response = await apiClients.identity.get<User[]>('/api/v1/users');
        setUsers(response.data || []);
      } catch (apiError: any) {
        // If endpoint doesn't exist, use mock data
        if (apiError.response?.status === 404 || apiError.response?.status === 501) {
          console.warn('User management API not available, using mock data');
          setUsers(getMockUsers());
        } else {
          throw apiError;
        }
      }
    } catch (err: any) {
      console.error('Failed to fetch users:', err);
      setError(err.response?.data?.error || 'Failed to fetch users. Using mock data.');
      // Fallback to mock data
      setUsers(getMockUsers());
    } finally {
      setLoading(false);
    }
  };

  // Mock data for development until API is available
  const getMockUsers = (): User[] => {
    return [
      {
        id: '1',
        email: 'buyer.requester@demo.com',
        first_name: 'Buyer',
        last_name: 'Requester',
        tenant_id: '00000000-0000-0000-0000-000000000001',
        tenant_name: 'Demo Buyer Company',
        roles: ['requester'],
        is_active: true,
        is_verified: true,
        created_at: new Date().toISOString(),
        last_login_at: new Date().toISOString(),
      },
      {
        id: '2',
        email: 'buyer.procurement@demo.com',
        first_name: 'Buyer',
        last_name: 'Procurement',
        tenant_id: '00000000-0000-0000-0000-000000000001',
        tenant_name: 'Demo Buyer Company',
        roles: ['procurement_manager'],
        is_active: true,
        is_verified: true,
        created_at: new Date().toISOString(),
        last_login_at: new Date().toISOString(),
      },
      {
        id: '3',
        email: 'supplier@demo.com',
        first_name: 'Supplier',
        last_name: 'User',
        tenant_id: '00000000-0000-0000-0000-000000000002',
        tenant_name: 'Demo Supplier Company',
        roles: ['supplier'],
        is_active: true,
        is_verified: true,
        created_at: new Date().toISOString(),
        last_login_at: new Date().toISOString(),
      },
      {
        id: '4',
        email: 'admin@demo.com',
        first_name: 'Platform',
        last_name: 'Admin',
        tenant_id: '00000000-0000-0000-0000-000000000000',
        tenant_name: 'Platform',
        roles: ['admin', 'super_admin'],
        is_active: true,
        is_verified: true,
        created_at: new Date().toISOString(),
        last_login_at: new Date().toISOString(),
      },
    ];
  };

  const handleDeactivate = async (userId: string) => {
    if (!confirm('Are you sure you want to deactivate this user?')) {
      return;
    }

    try {
      // TODO: Replace with actual API endpoint when available
      // await apiClients.identity.put(`/api/v1/users/${userId}/deactivate`);
      alert('User deactivation API not yet implemented. This is a placeholder.');
      // Refresh users
      await fetchUsers();
    } catch (err: any) {
      alert(err.response?.data?.error || 'Failed to deactivate user');
    }
  };

  const handleActivate = async (userId: string) => {
    try {
      // TODO: Replace with actual API endpoint when available
      // await apiClients.identity.put(`/api/v1/users/${userId}/activate`);
      alert('User activation API not yet implemented. This is a placeholder.');
      // Refresh users
      await fetchUsers();
    } catch (err: any) {
      alert(err.response?.data?.error || 'Failed to activate user');
    }
  };

  const handleResetPassword = async (userId: string) => {
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

  const filteredUsers = users.filter(user => {
    if (!searchTerm) return true;
    const search = searchTerm.toLowerCase();
    return (
      user.email.toLowerCase().includes(search) ||
      user.first_name.toLowerCase().includes(search) ||
      user.last_name.toLowerCase().includes(search) ||
      user.tenant_name?.toLowerCase().includes(search) ||
      user.roles?.some(role => role.toLowerCase().includes(search))
    );
  });

  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-3xl font-bold text-gray-900">User Management</h1>
        <div className="text-sm text-gray-500">
          {users.length} user{users.length !== 1 ? 's' : ''}
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

      {/* Search */}
      <Card className="mb-6">
        <Input
          label="Search Users"
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
          placeholder="Search by email, name, company, or role..."
        />
      </Card>

      {/* Users Table */}
      {loading ? (
        <div className="flex items-center justify-center min-h-[400px]">
          <div className="text-center">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600 mx-auto mb-4"></div>
            <p className="text-gray-600">Loading users...</p>
          </div>
        </div>
      ) : filteredUsers.length === 0 ? (
        <Card>
          <p className="text-gray-600">No users found.</p>
        </Card>
      ) : (
        <Card>
          <div className="overflow-x-auto">
            <table className="min-w-full divide-y divide-gray-200">
              <thead className="bg-gray-50">
                <tr>
                  <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">User ID</th>
                  <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Email</th>
                  <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Name</th>
                  <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Company/Tenant</th>
                  <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Roles</th>
                  <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Status</th>
                  <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Created</th>
                  <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Last Login</th>
                  <th className="px-4 py-3 text-center text-xs font-medium text-gray-500 uppercase">Actions</th>
                </tr>
              </thead>
              <tbody className="bg-white divide-y divide-gray-200">
                {filteredUsers.map((user) => (
                  <tr key={user.id} className="hover:bg-gray-50">
                    <td className="px-4 py-3 text-sm text-gray-900 font-mono">
                      {user.id.slice(0, 8)}...
                    </td>
                    <td className="px-4 py-3 text-sm text-gray-900">{user.email}</td>
                    <td className="px-4 py-3 text-sm text-gray-900">
                      {user.first_name} {user.last_name}
                    </td>
                    <td className="px-4 py-3 text-sm text-gray-600">
                      {user.tenant_name || user.tenant_id.slice(0, 8)}
                    </td>
                    <td className="px-4 py-3 text-sm">
                      <div className="flex flex-wrap gap-1">
                        {user.roles?.map((role) => (
                          <span
                            key={role}
                            className="px-2 py-1 bg-blue-100 text-blue-800 text-xs rounded"
                          >
                            {role}
                          </span>
                        ))}
                      </div>
                    </td>
                    <td className="px-4 py-3 text-sm">
                      <div className="flex flex-col gap-1">
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
                    </td>
                    <td className="px-4 py-3 text-sm text-gray-600">
                      {new Date(user.created_at).toLocaleDateString()}
                    </td>
                    <td className="px-4 py-3 text-sm text-gray-600">
                      {user.last_login_at ? new Date(user.last_login_at).toLocaleDateString() : 'Never'}
                    </td>
                    <td className="px-4 py-3 text-sm text-center">
                      <div className="flex gap-2 justify-center">
                        <Link href={`/app/admin/users/${user.id}`}>
                          <Button size="sm" variant="secondary">View</Button>
                        </Link>
                        {user.is_active ? (
                          <Button
                            size="sm"
                            variant="secondary"
                            onClick={() => handleDeactivate(user.id)}
                            className="bg-red-50 text-red-700 hover:bg-red-100"
                          >
                            Deactivate
                          </Button>
                        ) : (
                          <Button
                            size="sm"
                            variant="secondary"
                            onClick={() => handleActivate(user.id)}
                            className="bg-green-50 text-green-700 hover:bg-green-100"
                          >
                            Activate
                          </Button>
                        )}
                        <Button
                          size="sm"
                          variant="secondary"
                          onClick={() => handleResetPassword(user.id)}
                          className="bg-blue-50 text-blue-700 hover:bg-blue-100"
                        >
                          Reset Pwd
                        </Button>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </Card>
      )}
    </div>
  );
}
