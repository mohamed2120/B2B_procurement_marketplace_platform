import type { Metadata } from 'next';
import { Inter } from 'next/font/google';
import './globals.css';
import { TenantProvider } from '@/contexts/TenantContext';
import ErrorBoundary from '@/components/ErrorBoundary';

const inter = Inter({ subsets: ['latin'] });

export const metadata: Metadata = {
  title: 'B2B Procurement Marketplace',
  description: 'Enterprise B2B Procurement Marketplace Platform',
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en" suppressHydrationWarning>
      <body className={inter.className} suppressHydrationWarning>
        <ErrorBoundary>
          <TenantProvider>
            {children}
          </TenantProvider>
        </ErrorBoundary>
      </body>
    </html>
  );
}
