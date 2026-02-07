'use client';

import { useState, useEffect } from 'react';
import { apiClients } from '@/lib/api';
import { hasRole } from '@/lib/auth';
import Card from '@/components/ui/Card';

interface Role {
  id: string;
  name: string;
  description: string;
  permissions: Permission[];
}

interface Permission {
  id: string;
  resource: string;
  action: string;
  description: string;
}

export default function RolesPermissionsPage() {
  const [roles, setRoles] = useState<Role[]>([]);
  const [permissions, setPermissions] = useState<Permission[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

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
    fetchData();
  }, []);

  const fetchData = async () => {
    try {
      setLoading(true);
      setError('');
      
      // TODO: Replace with actual API endpoints when available
      // try {
      //   const [rolesRes, permsRes] = await Promise.all([
      //     apiClients.identity.get('/api/v1/roles'),
      //     apiClients.identity.get('/api/v1/permissions'),
      //   ]);
      //   setRoles(rolesRes.data || []);
      //   setPermissions(permsRes.data || []);
      // } catch (apiError: any) {
      //   if (apiError.response?.status === 404 || apiError.response?.status === 501) {
      //     console.warn('Roles/permissions API not available, using mock data');
      //     setRoles(getMockRoles());
      //     setPermissions(getMockPermissions());
      //   } else {
      //     throw apiError;
      //   }
      // }
      
      // Using mock data for now
      setRoles(getMockRoles());
      setPermissions(getMockPermissions());
    } catch (err: any) {
      console.error('Failed to fetch roles/permissions:', err);
      setError(err.response?.data?.error || 'Failed to fetch data. Using mock data.');
      setRoles(getMockRoles());
      setPermissions(getMockPermissions());
    } finally {
      setLoading(false);
    }
  };

  // Mock data for development until API is available
  const getMockRoles = (): Role[] => {
    return [
      {
        id: '1',
        name: 'requester',
        description: 'Can create purchase requests and RFQs',
        permissions: [
          { id: '1', resource: 'pr', action: 'create', description: 'Create purchase requests' },
          { id: '2', resource: 'pr', action: 'read', description: 'View own purchase requests' },
          { id: '3', resource: 'rfq', action: 'create', description: 'Create RFQs' },
          { id: '4', resource: 'rfq', action: 'read', description: 'View own RFQs' },
        ],
      },
      {
        id: '2',
        name: 'procurement_manager',
        description: 'Can approve PRs, compare quotes, and award quotes',
        permissions: [
          { id: '1', resource: 'pr', action: 'read', description: 'View all purchase requests' },
          { id: '5', resource: 'pr', action: 'approve', description: 'Approve/reject purchase requests' },
          { id: '6', resource: 'quote', action: 'read', description: 'View all quotes' },
          { id: '7', resource: 'quote', action: 'award', description: 'Award quotes and create POs' },
          { id: '8', resource: 'po', action: 'create', description: 'Create purchase orders' },
        ],
      },
      {
        id: '3',
        name: 'supplier',
        description: 'Can manage listings and submit quotes',
        permissions: [
          { id: '9', resource: 'listing', action: 'create', description: 'Create product listings' },
          { id: '10', resource: 'listing', action: 'update', description: 'Update own listings' },
          { id: '11', resource: 'rfq', action: 'read', description: 'View RFQs' },
          { id: '12', resource: 'quote', action: 'create', description: 'Submit quotes' },
          { id: '13', resource: 'shipment', action: 'update', description: 'Update shipment status' },
        ],
      },
      {
        id: '4',
        name: 'admin',
        description: 'Platform administrator with full access',
        permissions: [
          { id: '14', resource: '*', action: '*', description: 'Full access to all resources' },
          { id: '15', resource: 'tenant', action: 'manage', description: 'Manage tenants' },
          { id: '16', resource: 'user', action: 'manage', description: 'Manage users' },
          { id: '17', resource: 'role', action: 'manage', description: 'Manage roles and permissions' },
        ],
      },
    ];
  };

  const getMockPermissions = (): Permission[] => {
    return [
      { id: '1', resource: 'pr', action: 'create', description: 'Create purchase requests' },
      { id: '2', resource: 'pr', action: 'read', description: 'View purchase requests' },
      { id: '3', resource: 'rfq', action: 'create', description: 'Create RFQs' },
      { id: '4', resource: 'rfq', action: 'read', description: 'View RFQs' },
      { id: '5', resource: 'pr', action: 'approve', description: 'Approve/reject purchase requests' },
      { id: '6', resource: 'quote', action: 'read', description: 'View quotes' },
      { id: '7', resource: 'quote', action: 'award', description: 'Award quotes' },
      { id: '8', resource: 'po', action: 'create', description: 'Create purchase orders' },
      { id: '9', resource: 'listing', action: 'create', description: 'Create product listings' },
      { id: '10', resource: 'listing', action: 'update', description: 'Update listings' },
      { id: '11', resource: 'rfq', action: 'read', description: 'View RFQs' },
      { id: '12', resource: 'quote', action: 'create', description: 'Submit quotes' },
      { id: '13', resource: 'shipment', action: 'update', description: 'Update shipment status' },
      { id: '14', resource: '*', action: '*', description: 'Full access' },
      { id: '15', resource: 'tenant', action: 'manage', description: 'Manage tenants' },
      { id: '16', resource: 'user', action: 'manage', description: 'Manage users' },
      { id: '17', resource: 'role', action: 'manage', description: 'Manage roles' },
    ];
  };

  // Get unique resources
  const resources = Array.from(new Set(permissions.map(p => p.resource))).sort();

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600 mx-auto mb-4"></div>
          <p className="text-gray-600">Loading roles and permissions...</p>
        </div>
      </div>
    );
  }

  return (
    <div>
      <div className="mb-6">
        <h1 className="text-3xl font-bold text-gray-900">Roles & Permissions Matrix</h1>
        <p className="text-gray-600 mt-2">View the RBAC (Role-Based Access Control) matrix for the platform</p>
      </div>

      {error && (
        <div className="mb-4 bg-yellow-50 border border-yellow-200 text-yellow-700 px-4 py-3 rounded">
          {error}
          <div className="text-sm mt-1">
            Note: Using mock data until roles/permissions API is implemented.
          </div>
        </div>
      )}

      {/* Roles Overview */}
      <Card title="Roles Overview" className="mb-6">
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
          {roles.map((role) => (
            <div key={role.id} className="border rounded-lg p-4">
              <h3 className="font-semibold text-gray-900 mb-1">{role.name}</h3>
              <p className="text-sm text-gray-600 mb-2">{role.description}</p>
              <p className="text-xs text-gray-500">
                {role.permissions.length} permission{role.permissions.length !== 1 ? 's' : ''}
              </p>
            </div>
          ))}
        </div>
      </Card>

      {/* Permissions Matrix */}
      <Card title="Permissions Matrix">
        <div className="overflow-x-auto">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Permission</th>
                {roles.map((role) => (
                  <th key={role.id} className="px-4 py-3 text-center text-xs font-medium text-gray-500 uppercase">
                    {role.name}
                  </th>
                ))}
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {permissions.map((permission) => (
                <tr key={permission.id}>
                  <td className="px-4 py-3 text-sm">
                    <div>
                      <span className="font-medium text-gray-900">
                        {permission.resource}.{permission.action}
                      </span>
                      <p className="text-xs text-gray-500 mt-1">{permission.description}</p>
                    </div>
                  </td>
                  {roles.map((role) => {
                    const hasPermission = role.permissions.some(
                      p => p.resource === permission.resource && p.action === permission.action
                    ) || role.permissions.some(p => p.resource === '*' && p.action === '*');
                    return (
                      <td key={role.id} className="px-4 py-3 text-center">
                        {hasPermission ? (
                          <span className="inline-flex items-center justify-center w-6 h-6 bg-green-100 text-green-800 rounded-full text-xs">
                            ✓
                          </span>
                        ) : (
                          <span className="inline-flex items-center justify-center w-6 h-6 bg-gray-100 text-gray-400 rounded-full text-xs">
                            -
                          </span>
                        )}
                      </td>
                    );
                  })}
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </Card>

      <div className="mt-4 text-sm text-gray-500">
        <p>Note: This is a read-only view. Role and permission management APIs are not yet implemented.</p>
        <p className="mt-1">To modify roles or permissions, backend API endpoints need to be created.</p>
      </div>
    </div>
  );
}
