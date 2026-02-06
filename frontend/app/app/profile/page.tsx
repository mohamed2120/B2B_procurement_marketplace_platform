'use client';

import { useState, useEffect } from 'react';
import { getUser } from '@/lib/auth';
import Card from '@/components/ui/Card';
import Input from '@/components/ui/Input';
import Button from '@/components/ui/Button';

export default function Profile() {
  const user = getUser();
  const [formData, setFormData] = useState({
    first_name: '',
    last_name: '',
    email: '',
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState(false);

  useEffect(() => {
    if (user) {
      setFormData({
        first_name: user.first_name || '',
        last_name: user.last_name || '',
        email: user.email || '',
      });
    }
  }, [user]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setSuccess(false);
    setLoading(true);

    try {
      // TODO: Implement profile update API call
      // For now, just show success message
      await new Promise(resolve => setTimeout(resolve, 500));
      setSuccess(true);
      
      // Update local storage user data
      if (typeof window !== 'undefined') {
        const updatedUser = {
          ...user,
          first_name: formData.first_name,
          last_name: formData.last_name,
        };
        localStorage.setItem('auth_user', JSON.stringify(updatedUser));
      }
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to update profile');
    } finally {
      setLoading(false);
    }
  };

  if (!user) {
    return (
      <div>
        <h1 className="text-3xl font-bold text-gray-900 mb-6">Profile</h1>
        <Card>
          <p className="text-gray-600">Please log in to view your profile.</p>
        </Card>
      </div>
    );
  }

  return (
    <div>
      <h1 className="text-3xl font-bold text-gray-900 mb-6">My Profile</h1>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div className="lg:col-span-2">
          <Card title="Personal Information">
            <form onSubmit={handleSubmit} className="space-y-4">
              <div className="grid grid-cols-2 gap-4">
                <Input
                  label="First Name"
                  type="text"
                  value={formData.first_name}
                  onChange={(e) => setFormData({ ...formData, first_name: e.target.value })}
                  required
                />
                <Input
                  label="Last Name"
                  type="text"
                  value={formData.last_name}
                  onChange={(e) => setFormData({ ...formData, last_name: e.target.value })}
                  required
                />
              </div>

              <Input
                label="Email"
                type="email"
                value={formData.email}
                disabled
                className="bg-gray-50"
              />
              <p className="text-xs text-gray-500">Email cannot be changed</p>

              {error && (
                <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">
                  {error}
                </div>
              )}

              {success && (
                <div className="bg-green-50 border border-green-200 text-green-700 px-4 py-3 rounded">
                  Profile updated successfully!
                </div>
              )}

              <div className="flex gap-4 pt-4">
                <Button type="submit" disabled={loading}>
                  {loading ? 'Saving...' : 'Save Changes'}
                </Button>
              </div>
            </form>
          </Card>
        </div>

        <div>
          <Card title="Account Information">
            <div className="space-y-4">
              <div>
                <span className="text-sm font-medium text-gray-700">User ID</span>
                <p className="text-sm text-gray-900 mt-1 font-mono">{user.id}</p>
              </div>
              <div>
                <span className="text-sm font-medium text-gray-700">Tenant ID</span>
                <p className="text-sm text-gray-900 mt-1 font-mono">{user.tenant_id}</p>
              </div>
              <div>
                <span className="text-sm font-medium text-gray-700">Roles</span>
                <p className="text-sm text-gray-900 mt-1">
                  {user.roles && user.roles.length > 0 ? user.roles.join(', ') : 'No roles assigned'}
                </p>
              </div>
            </div>
          </Card>
        </div>
      </div>
    </div>
  );
}
