'use client';

import { usePathname as useNextPathname } from 'next/navigation';
import { useEffect, useState } from 'react';

export function useSafePathname() {
  const [mounted, setMounted] = useState(false);
  const [pathname, setPathname] = useState<string>('');

  useEffect(() => {
    setMounted(true);
  }, []);

  useEffect(() => {
    if (mounted) {
      try {
        const nextPathname = useNextPathname();
        setPathname(nextPathname || window.location.pathname);
      } catch {
        setPathname(window.location.pathname);
      }
    }
  }, [mounted]);

  if (!mounted) {
    if (typeof window !== 'undefined') {
      return window.location.pathname;
    }
    return '';
  }

  try {
    return useNextPathname() || '';
  } catch {
    return typeof window !== 'undefined' ? window.location.pathname : '';
  }
}
