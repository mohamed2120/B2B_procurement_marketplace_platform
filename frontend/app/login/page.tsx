'use client';

import { useState, useEffect } from 'react';
import { useSafeRouter } from '@/lib/useSafeRouter';
import Link from 'next/link';
import { login, isAuthenticated, getUser } from '@/lib/auth';
import Cookies from 'js-cookie';
import Input from '@/components/ui/Input';
import Button from '@/components/ui/Button';
import Card from '@/components/ui/Card';
import PublicLayout from '@/components/layout/PublicLayout';

export default function LoginPage() {
  const router = useSafeRouter();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const [mounted, setMounted] = useState(false);
  const [alreadyLoggedIn, setAlreadyLoggedIn] = useState(false);
  const [currentUser, setCurrentUser] = useState<any>(null);

  // Check if already authenticated on mount
  useEffect(() => {
    setMounted(true);
    if (isAuthenticated()) {
      const user = getUser();
      setAlreadyLoggedIn(true);
      setCurrentUser(user);
    }
  }, []);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      await login(email, password);
      // Small delay to ensure localStorage is updated
      await new Promise(resolve => setTimeout(resolve, 150));
      // Redirect to app - AppRouterRedirect will handle role-based routing
      window.location.href = '/app';
    } catch (err: any) {
      const errorMessage = err.response?.data?.error || 
                          err.response?.data?.message || 
                          err.message || 
                          'Login failed. Please check your credentials.';
      setError(errorMessage);
    } finally {
      setLoading(false);
    }
  };

  const handleLogout = () => {
    // Clear auth data manually to avoid redirect
    Cookies.remove('auth_token');
    if (typeof window !== 'undefined') {
      localStorage.removeItem('auth_user');
    }
    setAlreadyLoggedIn(false);
    setCurrentUser(null);
    setEmail('');
    setPassword('');
    // Small delay to ensure state updates, then reload to show login form
    setTimeout(() => {
      window.location.reload();
    }, 100);
  };

  // If already logged in, show option to logout or continue
  if (mounted && alreadyLoggedIn) {
    return (
      <PublicLayout>
        <div className="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
          <Card className="w-full max-w-md">
            <h2 className="text-2xl font-bold text-center text-gray-900 mb-6">Already Logged In</h2>
            {currentUser && (
              <div className="mb-4 p-4 bg-blue-50 border border-blue-200 rounded-lg">
                <p className="text-sm text-gray-600 mb-1">Currently logged in as:</p>
                <p className="font-semibold text-gray-900">{currentUser.email}</p>
              </div>
            )}
            <div className="space-y-4">
              <Button
                onClick={() => window.location.href = '/app'}
                className="w-full"
              >
                Go to Portal
              </Button>
              <Button
                onClick={handleLogout}
                variant="outline"
                className="w-full"
              >
                Logout and Sign In as Different User
              </Button>
            </div>
          </Card>
        </div>
      </PublicLayout>
    );
  }

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
