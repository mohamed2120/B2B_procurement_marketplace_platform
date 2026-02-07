'use client';

import { useState, useEffect } from 'react';
import { useParams, useRouter } from 'next/navigation';
import { apiClients } from '@/lib/api';
import Card from '@/components/ui/Card';
import Button from '@/components/ui/Button';
import Link from 'next/link';

interface Tenant {
  id: string;
  name: string;
  domain?: string;
  company_type: string;
  status: string;
  subscription_tier?: string;
  created_at: string;
  verified_at?: string;
  user_count?: number;
  address?: string;
  phone?: string;
  email?: string;
  tax_id?: string;
}

interface TenantUser {
  id: string;
  email: string;
  first_name: string;
  last_name: string;
  roles: string[];
  is_active: boolean;
  last_login_at?: string;
}

export default function TenantDetailPage() {
  const params = useParams();
  const router = useRouter();
  const tenantId = params.id as string;
  
  const [tenant, setTenant] = useState<Tenant | null>(null);
  const [users, setUsers] = useState<TenantUser[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    if (tenantId) {
      fetchTenant();
      fetchTenantUsers();
    }
  }, [tenantId]);

  const fetchTenant = async () => {
    try {
      // TODO: Replace with actual API endpoint when available
      // const response = await apiClients.company.get(`/api/v1/admin/tenants/${tenantId}`);
      // setTenant(response.data);
      
      // Mock data
      setTenant({
        id: tenantId,
        name: 'Acme Corporation',
        domain: 'acme',
        company_type: 'buyer',
        status: 'active',
        subscription_tier: 'enterprise',
        created_at: '2024-01-15T10:00:00Z',
        verified_at: '2024-01-16T10:00:00Z',
        user_count: 25,
        address: '123 Business St, City, State 12345',
        phone: '+1-555-0123',
        email: 'contact@acme.com',
        tax_id: 'TAX-123456',
      });
    } catch (err: any) {
      setError(err.response?.data?.error || err.message || 'Failed to fetch tenant');
      console.error('Failed to fetch tenant:', err);
    } finally {
      setLoading(false);
    }
  };

  const fetchTenantUsers = async () => {
    try {
      // TODO: Replace with actual API endpoint when available
      // const response = await apiClients.company.get(`/api/v1/admin/tenants/${tenantId}/users`);
      // setUsers(response.data || []);
      
      // Mock data
      setUsers([
        {
          id: '1',
          email: 'admin@acme.com',
          first_name: 'John',
          last_name: 'Admin',
          roles: ['admin'],
          is_active: true,
          last_login_at: '2024-02-07T10:00:00Z',
        },
        {
          id: '2',
          email: 'buyer@acme.com',
          first_name: 'Jane',
          last_name: 'Buyer',
          roles: ['requester', 'procurement_manager'],
          is_active: true,
          last_login_at: '2024-02-07T09:00:00Z',
        },
      ]);
    } catch (err: any) {
      console.error('Failed to fetch tenant users:', err);
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'active':
        return 'bg-green-100 text-green-800';
      case 'pending_verification':
        return 'bg-yellow-100 text-yellow-800';
      case 'suspended':
        return 'bg-red-100 text-red-800';
      default:
        return 'bg-gray-100 text-gray-800';
    }
  };

  if (loading) {
    return (
      <Card>
        <div className="text-center py-8">Loading tenant details...</div>
      </Card>
    );
  }

  if (error || !tenant) {
    return (
      <Card className="bg-red-50 border-red-200">
        <p className="text-red-800">{error || 'Tenant not found'}</p>
        <Link href="/app/admin/tenants">
          <Button className="mt-4">Back to Tenants</Button>
        </Link>
      </Card>
    );
  }

  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <div>
          <Link href="/app/admin/tenants" className="text-primary-600 hover:text-primary-800 mb-2 inline-block">
            ← Back to Tenants
          </Link>
          <h1 className="text-3xl font-bold text-gray-900">{tenant.name}</h1>
        </div>
        <div className="flex gap-2">
          <Button variant="secondary">Edit</Button>
          <Button variant="secondary">Suspend</Button>
        </div>
      </div>

      {/* Tenant Details */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-6">
        <Card>
          <h2 className="text-xl font-semibold mb-4">Company Information</h2>
          <dl className="space-y-3">
            <div>
              <dt className="text-sm font-medium text-gray-500">Name</dt>
              <dd className="text-sm text-gray-900">{tenant.name}</dd>
            </div>
            <div>
              <dt className="text-sm font-medium text-gray-500">Domain</dt>
              <dd className="text-sm text-gray-900">{tenant.domain || '-'}</dd>
            </div>
            <div>
              <dt className="text-sm font-medium text-gray-500">Type</dt>
              <dd className="text-sm text-gray-900 capitalize">{tenant.company_type}</dd>
            </div>
            <div>
              <dt className="text-sm font-medium text-gray-500">Status</dt>
              <dd>
                <span className={`px-2 py-1 text-xs font-semibold rounded-full ${getStatusColor(tenant.status)}`}>
                  {tenant.status.replace('_', ' ')}
                </span>
              </dd>
            </div>
            <div>
              <dt className="text-sm font-medium text-gray-500">Subscription Tier</dt>
              <dd className="text-sm text-gray-900 capitalize">{tenant.subscription_tier || '-'}</dd>
            </div>
          </dl>
        </Card>

        <Card>
          <h2 className="text-xl font-semibold mb-4">Contact Information</h2>
          <dl className="space-y-3">
            {tenant.email && (
              <div>
                <dt className="text-sm font-medium text-gray-500">Email</dt>
                <dd className="text-sm text-gray-900">{tenant.email}</dd>
              </div>
            )}
            {tenant.phone && (
              <div>
                <dt className="text-sm font-medium text-gray-500">Phone</dt>
                <dd className="text-sm text-gray-900">{tenant.phone}</dd>
              </div>
            )}
            {tenant.address && (
              <div>
                <dt className="text-sm font-medium text-gray-500">Address</dt>
                <dd className="text-sm text-gray-900">{tenant.address}</dd>
              </div>
            )}
            {tenant.tax_id && (
              <div>
                <dt className="text-sm font-medium text-gray-500">Tax ID</dt>
                <dd className="text-sm text-gray-900">{tenant.tax_id}</dd>
              </div>
            )}
          </dl>
        </Card>
      </div>

      {/* Timestamps */}
      <Card className="mb-6">
        <h2 className="text-xl font-semibold mb-4">Timestamps</h2>
        <dl className="grid grid-cols-2 gap-4">
          <div>
            <dt className="text-sm font-medium text-gray-500">Created At</dt>
            <dd className="text-sm text-gray-900">
              {new Date(tenant.created_at).toLocaleString()}
            </dd>
          </div>
          {tenant.verified_at && (
            <div>
              <dt className="text-sm font-medium text-gray-500">Verified At</dt>
              <dd className="text-sm text-gray-900">
                {new Date(tenant.verified_at).toLocaleString()}
              </dd>
            </div>
          )}
        </dl>
      </Card>

      {/* Users */}
      <Card>
        <div className="flex justify-between items-center mb-4">
          <h2 className="text-xl font-semibold">Users ({users.length})</h2>
          <Link href={`/app/admin/users?tenant=${tenantId}`}>
            <Button size="sm">View All Users</Button>
          </Link>
        </div>
        {users.length === 0 ? (
          <p className="text-gray-600">No users found</p>
        ) : (
          <div className="overflow-x-auto">
            <table className="min-w-full divide-y divide-gray-200">
              <thead className="bg-gray-50">
                <tr>
                  <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Name</th>
                  <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Email</th>
                  <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Roles</th>
                  <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Status</th>
                  <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Last Login</th>
                  <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Actions</th>
                </tr>
              </thead>
              <tbody className="bg-white divide-y divide-gray-200">
                {users.map((user) => (
                  <tr key={user.id}>
                    <td className="px-4 py-3 whitespace-nowrap text-sm">
                      {user.first_name} {user.last_name}
                    </td>
                    <td className="px-4 py-3 whitespace-nowrap text-sm text-gray-500">{user.email}</td>
                    <td className="px-4 py-3 whitespace-nowrap text-sm text-gray-500">
                      {user.roles.join(', ')}
                    </td>
                    <td className="px-4 py-3 whitespace-nowrap">
                      <span className={`px-2 py-1 text-xs font-semibold rounded-full ${
                        user.is_active ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'
                      }`}>
                        {user.is_active ? 'Active' : 'Inactive'}
                      </span>
                    </td>
                    <td className="px-4 py-3 whitespace-nowrap text-sm text-gray-500">
                      {user.last_login_at ? new Date(user.last_login_at).toLocaleDateString() : 'Never'}
                    </td>
                    <td className="px-4 py-3 whitespace-nowrap text-sm font-medium">
                      <Link
                        href={`/app/admin/users/${user.id}`}
                        className="text-primary-600 hover:text-primary-900"
                      >
                        View →
                      </Link>
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
