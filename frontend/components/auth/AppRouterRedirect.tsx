'use client';

import { useEffect } from 'react';
import { useRouter, usePathname } from 'next/navigation';
import { isAuthenticated, getUser, hasRole } from '@/lib/auth';

export default function AppRouterRedirect() {
  const router = useRouter();
  const pathname = usePathname();

  useEffect(() => {
    if (!isAuthenticated()) {
      router.push('/login');
      return;
    }

    const user = getUser();
    if (!user) return;

    // Only redirect if we're at the base /app route
    if (pathname === '/app' || pathname === '/app/') {
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
  }, [router, pathname]);

  return null;
}
