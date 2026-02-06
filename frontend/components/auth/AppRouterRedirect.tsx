'use client';

import { useEffect, useState } from 'react';
import { usePathname } from 'next/navigation';
import { useSafeRouter } from '@/lib/useSafeRouter';
import { isAuthenticated, getUser, hasRole } from '@/lib/auth';

export default function AppRouterRedirect() {
  const [mounted, setMounted] = useState(false);
  const [redirecting, setRedirecting] = useState(false);
  const router = useSafeRouter();
  const pathname = usePathname();

  useEffect(() => {
    setMounted(true);
  }, []);

  useEffect(() => {
    if (!mounted || redirecting) return;
    
    if (!isAuthenticated()) {
      router.push('/login');
      return;
    }

    // Only redirect if we're at the base /app route
    if (pathname === '/app' || pathname === '/app/') {
      setRedirecting(true);
      
      // Wait a moment for user data to be available after login
      const timer = setTimeout(() => {
        const user = getUser();
        
        if (!user) {
          router.push('/login');
          return;
        }

        // Check roles and redirect
        if (hasRole('admin') || hasRole('super_admin')) {
          router.push('/app/admin/dashboard');
        } else if (hasRole('requester') || hasRole('procurement_manager') || hasRole('buyer')) {
          router.push('/app/customer/dashboard');
        } else if (hasRole('supplier')) {
          router.push('/app/supplier/dashboard');
        } else {
          // Default fallback
          router.push('/app/my-plan');
        }
      }, 200);
      
      return () => clearTimeout(timer);
    }
  }, [mounted, router, pathname, redirecting]);

  return null;
}
