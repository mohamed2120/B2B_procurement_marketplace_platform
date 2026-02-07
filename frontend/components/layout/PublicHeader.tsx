'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { isAuthenticated } from '@/lib/auth';
import { useState, useEffect } from 'react';
import SearchBar from '@/components/search/SearchBar';

export default function PublicHeader() {
  const pathname = usePathname();
  const [authenticated, setAuthenticated] = useState(false);
  const [mounted, setMounted] = useState(false);

  useEffect(() => {
    setMounted(true);
    setAuthenticated(isAuthenticated());
  }, []);

  return (
    <header className="bg-white shadow-sm">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between items-center h-16">
          <div className="flex items-center">
            <Link href="/" className="text-2xl font-bold text-primary-600">
              B2B Marketplace
            </Link>
          </div>
          <div className="flex-1 max-w-md mx-4">
            <SearchBar />
          </div>
          <nav className="hidden md:flex space-x-8">
            <Link
              href="/"
              className={`${
                pathname === '/' ? 'text-primary-600' : 'text-gray-700'
              } hover:text-primary-600`}
            >
              Home
            </Link>
            <Link
              href="/how-it-works"
              className={`${
                pathname === '/how-it-works' ? 'text-primary-600' : 'text-gray-700'
              } hover:text-primary-600`}
            >
              How It Works
            </Link>
            <Link
              href="/pricing"
              className={`${
                pathname === '/pricing' ? 'text-primary-600' : 'text-gray-700'
              } hover:text-primary-600`}
            >
              Pricing
            </Link>
            <Link
              href="/search"
              className={`${
                pathname === '/search' ? 'text-primary-600' : 'text-gray-700'
              } hover:text-primary-600`}
            >
              Search
            </Link>
            <Link
              href="/contact"
              className={`${
                pathname === '/contact' ? 'text-primary-600' : 'text-gray-700'
              } hover:text-primary-600`}
            >
              Contact
            </Link>
          </nav>
          <div className="flex items-center space-x-4">
            {!mounted ? (
              // Show default state during SSR to match initial client render
              <>
                <Link
                  href="/login"
                  className="text-gray-700 hover:text-primary-600"
                >
                  Login
                </Link>
                <Link
                  href="/register"
                  className="bg-primary-600 text-white px-4 py-2 rounded-lg hover:bg-primary-700"
                >
                  Get Started
                </Link>
              </>
            ) : authenticated ? (
              <Link
                href="/app"
                className="bg-primary-600 text-white px-4 py-2 rounded-lg hover:bg-primary-700"
              >
                Go to Portal
              </Link>
            ) : (
              <>
                <Link
                  href="/login"
                  className="text-gray-700 hover:text-primary-600"
                >
                  Login
                </Link>
                <Link
                  href="/register"
                  className="bg-primary-600 text-white px-4 py-2 rounded-lg hover:bg-primary-700"
                >
                  Get Started
                </Link>
              </>
            )}
          </div>
        </div>
      </div>
    </header>
  );
}
