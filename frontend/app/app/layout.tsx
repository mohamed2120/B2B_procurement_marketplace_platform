'use client';

import dynamic from 'next/dynamic';
import PortalLayout from '@/components/layout/PortalLayout';
import ProtectedRoute from '@/components/auth/ProtectedRoute';

// Dynamically import AppRouterRedirect to ensure it only loads client-side after hydration
const AppRouterRedirect = dynamic(
  () => import('@/components/auth/AppRouterRedirect'),
  { ssr: false }
);

export default function AppLayout({ children }: { children: React.ReactNode }) {
  return (
    <ProtectedRoute>
      <AppRouterRedirect />
      <PortalLayout>{children}</PortalLayout>
    </ProtectedRoute>
  );
}
