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

        // Debug: Log user and roles
        console.log('AppRouterRedirect - User:', user);
        console.log('AppRouterRedirect - Roles:', user.roles);

        // Check roles and redirect
        if (hasRole('admin') || hasRole('super_admin')) {
          console.log('Redirecting admin to /app/admin/dashboard');
          router.push('/app/admin/dashboard');
        } else if (hasRole('requester') || hasRole('procurement_manager') || hasRole('buyer')) {
          console.log('Redirecting buyer to /app/customer/dashboard');
          router.push('/app/customer/dashboard');
        } else if (hasRole('supplier')) {
          console.log('Redirecting supplier to /app/supplier/dashboard');
          router.push('/app/supplier/dashboard');
        } else {
          // Default fallback
          console.log('No role detected, redirecting to /app/my-plan');
          router.push('/app/my-plan');
        }
      }, 500); // Increased timeout to ensure roles are loaded
      
      return () => clearTimeout(timer);
    }
  }, [mounted, router, pathname, redirecting]);

  return null;
}
