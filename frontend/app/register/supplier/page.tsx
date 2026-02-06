'use client';

import { useState } from 'react';
import { useSafeRouter } from '@/lib/useSafeRouter';
import Link from 'next/link';
import PublicLayout from '@/components/layout/PublicLayout';
import Input from '@/components/ui/Input';
import Button from '@/components/ui/Button';
import Card from '@/components/ui/Card';

export default function RegisterSupplier() {
  const router = useSafeRouter();
  const [formData, setFormData] = useState({
    companyName: '',
    email: '',
    password: '',
    firstName: '',
    lastName: '',
    phone: '',
    taxID: '',
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      // TODO: Implement registration API call
      // For now, redirect to login
      router.push('/login?registered=true');
    } catch (err: any) {
      setError(err.response?.data?.error || 'Registration failed');
    } finally {
      setLoading(false);
    }
  };

  return (
    <PublicLayout>
      <div className="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
        <Card className="w-full max-w-md">
          <h2 className="text-2xl font-bold text-center text-gray-900 mb-6">
            Register as Supplier Company
          </h2>

          <form onSubmit={handleSubmit} className="space-y-4">
            <Input
              label="Company Name"
              type="text"
              value={formData.companyName}
              onChange={(e) => setFormData({ ...formData, companyName: e.target.value })}
              required
            />

            <div className="grid grid-cols-2 gap-4">
              <Input
                label="First Name"
                type="text"
                value={formData.firstName}
                onChange={(e) => setFormData({ ...formData, firstName: e.target.value })}
                required
              />
              <Input
                label="Last Name"
                type="text"
                value={formData.lastName}
                onChange={(e) => setFormData({ ...formData, lastName: e.target.value })}
                required
              />
            </div>

            <Input
              label="Email"
              type="email"
              value={formData.email}
              onChange={(e) => setFormData({ ...formData, email: e.target.value })}
              required
            />

            <Input
              label="Phone"
              type="tel"
              value={formData.phone}
              onChange={(e) => setFormData({ ...formData, phone: e.target.value })}
            />

            <Input
              label="Tax ID"
              type="text"
              value={formData.taxID}
              onChange={(e) => setFormData({ ...formData, taxID: e.target.value })}
            />

            <Input
              label="Password"
              type="password"
              value={formData.password}
              onChange={(e) => setFormData({ ...formData, password: e.target.value })}
              required
              minLength={8}
            />

            {error && (
              <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">
                {error}
              </div>
            )}

            <Button type="submit" className="w-full" disabled={loading}>
              {loading ? 'Creating Account...' : 'Create Supplier Account'}
            </Button>
          </form>

          <div className="mt-6 text-center">
            <Link href="/register" className="text-sm text-primary-600 hover:text-primary-700">
              ← Back to registration options
            </Link>
          </div>
        </Card>
      </div>
    </PublicLayout>
  );
}
