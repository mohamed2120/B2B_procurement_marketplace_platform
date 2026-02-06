'use client';

import { useEffect, useState } from 'react';
import { useSafeRouter } from '@/lib/useSafeRouter';
import { isAuthenticated } from '@/lib/auth';

export default function ProtectedRoute({ children }: { children: React.ReactNode }) {
  const [mounted, setMounted] = useState(false);
  const router = useSafeRouter();

  useEffect(() => {
    setMounted(true);
  }, []);

  useEffect(() => {
    if (!mounted) return;
    if (!isAuthenticated()) {
      router.push('/login');
    }
  }, [mounted, router]);

  if (!mounted) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div>Loading...</div>
      </div>
    );
  }

  if (!isAuthenticated()) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div>Redirecting to login...</div>
      </div>
    );
  }

  return <>{children}</>;
}
