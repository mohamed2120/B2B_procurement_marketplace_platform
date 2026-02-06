'use client';

import ErrorBoundary from './ErrorBoundary';
import { TenantProvider } from '@/contexts/TenantContext';

export default function ErrorBoundaryWrapper({ children }: { children: React.ReactNode }) {
  return (
    <ErrorBoundary>
      <TenantProvider>
        {children}
      </TenantProvider>
    </ErrorBoundary>
  );
}
