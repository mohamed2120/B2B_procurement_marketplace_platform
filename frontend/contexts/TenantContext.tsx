'use client';

import { createContext, useContext, useState, useEffect, ReactNode } from 'react';

interface TenantContextType {
  tenant: string;
  tenantID: string;
}

const TenantContext = createContext<TenantContextType | undefined>(undefined);

export function TenantProvider({ children }: { children: ReactNode }) {
  const [tenant, setTenant] = useState<string>('demo');
  const [tenantID] = useState<string>('00000000-0000-0000-0000-000000000001');

  useEffect(() => {
    if (typeof window !== 'undefined') {
      const host = window.location.host;
      const parts = host.split('.');
      
      if (parts.length > 2) {
        setTenant(parts[0]);
      } else if (host.includes('localhost')) {
        const subdomain = parts[0];
        if (subdomain !== 'localhost' && subdomain !== 'www') {
          setTenant(subdomain);
        } else {
          setTenant('demo');
        }
      } else {
        setTenant('demo');
      }
    }
  }, []);

  return (
    <TenantContext.Provider value={{ tenant, tenantID }}>
      {children}
    </TenantContext.Provider>
  );
}

export function useTenant() {
  const context = useContext(TenantContext);
  if (!context) {
    throw new Error('useTenant must be used within TenantProvider');
  }
  return context;
}
