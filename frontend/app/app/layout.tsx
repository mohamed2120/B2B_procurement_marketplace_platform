'use client';

import PortalLayout from '@/components/layout/PortalLayout';
import ProtectedRoute from '@/components/auth/ProtectedRoute';
import AppRouterRedirect from '@/components/auth/AppRouterRedirect';

export default function AppLayout({ children }: { children: React.ReactNode }) {
  return (
    <ProtectedRoute>
      <AppRouterRedirect />
      <PortalLayout>{children}</PortalLayout>
    </ProtectedRoute>
  );
}
