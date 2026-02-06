'use client';

import { useState } from 'react';
import { useSafeRouter } from '@/lib/useSafeRouter';
import PublicLayout from '@/components/layout/PublicLayout';
import Link from 'next/link';

export default function Register() {
  const router = useSafeRouter();
  const [companyType, setCompanyType] = useState<'buyer' | 'supplier' | null>(null);
  const [formData, setFormData] = useState({
    companyName: '',
    email: '',
    password: '',
    firstName: '',
    lastName: '',
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    // TODO: Implement registration API call
    // For now, redirect to login
    router.push('/login');
  };

  return (
    <PublicLayout>
      <div className="max-w-2xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        <div className="bg-white rounded-lg shadow-md p-8">
          <h1 className="text-3xl font-bold text-gray-900 mb-6 text-center">Create Account</h1>
          
          <div className="space-y-4">
            <p className="text-gray-600 text-center mb-6">I want to register as:</p>
            <Link
              href="/register/buyer"
              className="block w-full bg-primary-600 text-white px-6 py-4 rounded-lg font-semibold hover:bg-primary-700 transition text-left"
            >
              <div className="font-bold text-lg mb-1">Buyer Company</div>
              <div className="text-sm text-primary-100">Create purchase requests and manage procurement</div>
            </Link>
            <Link
              href="/register/supplier"
              className="block w-full bg-primary-600 text-white px-6 py-4 rounded-lg font-semibold hover:bg-primary-700 transition text-left"
            >
              <div className="font-bold text-lg mb-1">Supplier Company</div>
              <div className="text-sm text-primary-100">List products and respond to RFQs</div>
            </Link>
          </div>

          {false && (
            <form onSubmit={handleSubmit} className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Company Name
                </label>
                <input
                  type="text"
                  required
                  value={formData.companyName}
                  onChange={(e) => setFormData({ ...formData, companyName: e.target.value })}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-primary-500 focus:border-primary-500"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  First Name
                </label>
                <input
                  type="text"
                  required
                  value={formData.firstName}
                  onChange={(e) => setFormData({ ...formData, firstName: e.target.value })}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-primary-500 focus:border-primary-500"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Last Name
                </label>
                <input
                  type="text"
                  required
                  value={formData.lastName}
                  onChange={(e) => setFormData({ ...formData, lastName: e.target.value })}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-primary-500 focus:border-primary-500"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Email
                </label>
                <input
                  type="email"
                  required
                  value={formData.email}
                  onChange={(e) => setFormData({ ...formData, email: e.target.value })}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-primary-500 focus:border-primary-500"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Password
                </label>
                <input
                  type="password"
                  required
                  value={formData.password}
                  onChange={(e) => setFormData({ ...formData, password: e.target.value })}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-primary-500 focus:border-primary-500"
                />
              </div>

              <div className="flex items-center space-x-2">
                <button
                  type="button"
                  onClick={() => setCompanyType(null)}
                  className="text-gray-600 hover:text-gray-800"
                >
                  ← Back
                </button>
              </div>

              <button
                type="submit"
                className="w-full bg-primary-600 text-white px-4 py-2 rounded-lg font-semibold hover:bg-primary-700 transition"
              >
                Create Account
              </button>

              <p className="text-center text-sm text-gray-600">
                Already have an account?{' '}
                <Link href="/login" className="text-primary-600 hover:text-primary-700">
                  Login
                </Link>
              </p>
            </form>
          )}
          <div className="mt-6 text-center">
            <Link href="/login" className="text-sm text-primary-600 hover:text-primary-700">
              Already have an account? Login here
            </Link>
          </div>
        </div>
      </div>
    </PublicLayout>
  );
}
