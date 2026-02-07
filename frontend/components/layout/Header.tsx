'use client';

import Link from 'next/link';
import { useSafeRouter } from '@/lib/useSafeRouter';
import { getUser, logout } from '@/lib/auth';
import { useState, useEffect } from 'react';
import SearchBar from '@/components/search/SearchBar';

export default function Header() {
  const [mounted, setMounted] = useState(false);
  const router = useSafeRouter();
  const [user, setUser] = useState(getUser());
  const [showMenu, setShowMenu] = useState(false);

  useEffect(() => {
    setMounted(true);
    setUser(getUser());
  }, []);

  const handleLogout = () => {
    logout();
    router.push('/login');
  };

  return (
    <header className="bg-white shadow-sm border-b">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between items-center h-16">
          <div className="flex items-center">
            <Link href="/app" className="text-xl font-bold text-primary-600">
              B2B Marketplace
            </Link>
          </div>

          <div className="flex-1 max-w-md mx-4">
            <SearchBar />
          </div>
          <nav className="hidden md:flex space-x-4">
            <Link href="/app" className="px-3 py-2 text-sm font-medium text-gray-700 hover:text-primary-600">
              Dashboard
            </Link>
            <Link href="/app/my-plan" className="px-3 py-2 text-sm font-medium text-gray-700 hover:text-primary-600">
              My Plan
            </Link>
            <Link href="/app/notifications" className="px-3 py-2 text-sm font-medium text-gray-700 hover:text-primary-600">
              Notifications
            </Link>
            <Link href="/app/chat" className="px-3 py-2 text-sm font-medium text-gray-700 hover:text-primary-600">
              Chat
            </Link>
          </nav>

          <div className="flex items-center space-x-4">
            {user && (
              <div className="relative">
                <button
                  onClick={() => setShowMenu(!showMenu)}
                  className="flex items-center space-x-2 text-sm text-gray-700 hover:text-primary-600"
                >
                  <span>{user.first_name} {user.last_name}</span>
                  <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
                  </svg>
                </button>

                {showMenu && (
                  <div className="absolute right-0 mt-2 w-48 bg-white rounded-md shadow-lg py-1 z-10">
                    <Link
                      href="/app/my-plan"
                      className="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
                      onClick={() => setShowMenu(false)}
                    >
                      My Plan
                    </Link>
                    <Link
                      href="/app/profile"
                      className="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
                      onClick={() => setShowMenu(false)}
                    >
                      Profile
                    </Link>
                    <button
                      onClick={handleLogout}
                      className="block w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
                    >
                      Logout
                    </button>
                  </div>
                )}
              </div>
            )}
          </div>
        </div>
      </div>
    </header>
  );
}
