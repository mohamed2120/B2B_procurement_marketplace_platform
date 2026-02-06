'use client';

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { isAuthenticated, getUser, hasRole } from '@/lib/auth';
import PortalLayout from '@/components/layout/PortalLayout';
import ProtectedRoute from '@/components/auth/ProtectedRoute';

export default function AppLayout({ children }: { children: React.ReactNode }) {
  const router = useRouter();

  useEffect(() => {
    if (!isAuthenticated()) {
      router.push('/login');
      return;
    }

    const user = getUser();
    if (!user) return;

    const path = window.location.pathname;
    if (path === '/app' || path === '/app/') {
      if (hasRole('admin') || hasRole('super_admin')) {
        router.push('/app/admin/dashboard');
      } else if (hasRole('requester') || hasRole('procurement_manager') || hasRole('buyer')) {
        router.push('/app/customer/dashboard');
      } else if (hasRole('supplier')) {
        router.push('/app/supplier/dashboard');
      } else {
        router.push('/app/my-plan');
      }
    }
  }, [router]);

  return (
    <ProtectedRoute>
      <PortalLayout>{children}</PortalLayout>
    </ProtectedRoute>
  );
}
