'use client';

import { useRouter as useNextRouter } from 'next/navigation';
import { useEffect, useState } from 'react';

export function useSafeRouter() {
  const [mounted, setMounted] = useState(false);
  const router = useNextRouter();

  useEffect(() => {
    setMounted(true);
  }, []);

  if (!mounted) {
    // Return a safe router object that won't cause errors
    return {
      push: (url: string) => {
        if (typeof window !== 'undefined') {
          window.location.href = url;
        }
      },
      replace: (url: string) => {
        if (typeof window !== 'undefined') {
          window.location.replace(url);
        }
      },
      back: () => {
        if (typeof window !== 'undefined') {
          window.history.back();
        }
      },
      refresh: () => {
        if (typeof window !== 'undefined') {
          window.location.reload();
        }
      },
    } as ReturnType<typeof useNextRouter>;
  }

  return router;
}
