'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import { login } from '@/lib/auth';
import Input from '@/components/ui/Input';
import Button from '@/components/ui/Button';
import Card from '@/components/ui/Card';
import PublicLayout from '@/components/layout/PublicLayout';

export default function LoginPage() {
  const router = useRouter();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      await login(email, password);
      router.push('/app');
    } catch (err: any) {
      setError(err.response?.data?.error || 'Login failed. Please check your credentials.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <PublicLayout>
      <div className="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
        <Card className="w-full max-w-md">
        <h2 className="text-2xl font-bold text-center text-gray-900 mb-6">Sign in to your account</h2>

        <form onSubmit={handleSubmit} className="space-y-4">
          <Input
            label="Email address"
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
            placeholder="admin@demo.com"
          />

          <Input
            label="Password"
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
            placeholder="demo123456"
          />

          {error && (
            <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">
              {error}
            </div>
          )}

          <Button type="submit" className="w-full" disabled={loading}>
            {loading ? 'Signing in...' : 'Sign in'}
          </Button>
        </form>

        <div className="mt-6 pt-6 border-t">
          <p className="text-sm text-gray-600 text-center mb-3">Demo Accounts (password: demo123456):</p>
          <div className="grid grid-cols-2 gap-2">
            <button
              onClick={() => {
                setEmail('buyer.requester@demo.com');
                setPassword('demo123456');
              }}
              className="text-xs bg-gray-100 hover:bg-gray-200 px-3 py-2 rounded text-left"
            >
              <div className="font-semibold">Requester</div>
              <div className="text-gray-600">buyer.requester@demo.com</div>
            </button>
            <button
              onClick={() => {
                setEmail('buyer.procurement@demo.com');
                setPassword('demo123456');
              }}
              className="text-xs bg-gray-100 hover:bg-gray-200 px-3 py-2 rounded text-left"
            >
              <div className="font-semibold">Procurement</div>
              <div className="text-gray-600">buyer.procurement@demo.com</div>
            </button>
            <button
              onClick={() => {
                setEmail('supplier@demo.com');
                setPassword('demo123456');
              }}
              className="text-xs bg-gray-100 hover:bg-gray-200 px-3 py-2 rounded text-left"
            >
              <div className="font-semibold">Supplier</div>
              <div className="text-gray-600">supplier@demo.com</div>
            </button>
            <button
              onClick={() => {
                setEmail('admin@demo.com');
                setPassword('demo123456');
              }}
              className="text-xs bg-gray-100 hover:bg-gray-200 px-3 py-2 rounded text-left"
            >
              <div className="font-semibold">Platform Admin</div>
              <div className="text-gray-600">admin@demo.com</div>
            </button>
          </div>
        </div>
        <div className="mt-4 text-center">
          <Link href="/register" className="text-sm text-primary-600 hover:text-primary-700">
            Don't have an account? Register here
          </Link>
        </div>
      </Card>
    </div>
    </PublicLayout>
  );
}
