'use client';

import Link from 'next/link';
import PublicLayout from '@/components/layout/PublicLayout';

export default function Home() {

  return (
    <PublicLayout>
      {/* Hero Section */}
      <section className="bg-gradient-to-r from-primary-600 to-primary-800 text-white py-20">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center">
            <h1 className="text-4xl md:text-6xl font-bold mb-6">
              Streamline Your B2B Procurement
            </h1>
            <p className="text-xl md:text-2xl mb-8 text-primary-100">
              Connect buyers and suppliers in one powerful marketplace platform
            </p>
            <div className="flex justify-center space-x-4 flex-wrap gap-4">
              <Link
                href="/register/buyer"
                className="bg-white text-primary-600 px-8 py-3 rounded-lg font-semibold hover:bg-primary-50 transition"
              >
                Register as Buyer
              </Link>
              <Link
                href="/register/supplier"
                className="bg-white text-primary-600 px-8 py-3 rounded-lg font-semibold hover:bg-primary-50 transition"
              >
                Register as Supplier
              </Link>
              <Link
                href="/login"
                className="bg-transparent border-2 border-white text-white px-8 py-3 rounded-lg font-semibold hover:bg-white hover:text-primary-600 transition"
              >
                Login
              </Link>
            </div>
          </div>
        </div>
      </section>

      {/* Features Section */}
      <section className="py-20 bg-gray-50">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <h2 className="text-3xl font-bold text-center mb-12 text-gray-900">
            Why Choose Our Platform?
          </h2>
          <div className="grid md:grid-cols-3 gap-8">
            <div className="bg-white p-6 rounded-lg shadow">
              <div className="text-4xl mb-4">📦</div>
              <h3 className="text-xl font-semibold mb-2">Streamlined Procurement</h3>
              <p className="text-gray-600">
                Create purchase requests, manage RFQs, and track orders all in one place.
              </p>
            </div>
            <div className="bg-white p-6 rounded-lg shadow">
              <div className="text-4xl mb-4">🤝</div>
              <h3 className="text-xl font-semibold mb-2">Supplier Network</h3>
              <p className="text-gray-600">
                Connect with verified suppliers and access a global marketplace of products.
              </p>
            </div>
            <div className="bg-white p-6 rounded-lg shadow">
              <div className="text-4xl mb-4">⚡</div>
              <h3 className="text-xl font-semibold mb-2">Fast & Efficient</h3>
              <p className="text-gray-600">
                Automated workflows, real-time notifications, and intelligent matching.
              </p>
            </div>
          </div>
        </div>
      </section>

      {/* CTA Section */}
      <section className="py-20 bg-primary-600 text-white">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
          <h2 className="text-3xl font-bold mb-4">Ready to Transform Your Procurement?</h2>
          <p className="text-xl mb-8 text-primary-100">
            Join thousands of companies already using our platform
          </p>
          <Link
            href="/register"
            className="bg-white text-primary-600 px-8 py-3 rounded-lg font-semibold hover:bg-primary-50 transition inline-block"
          >
            Start Free Trial
          </Link>
        </div>
      </section>
    </PublicLayout>
  );
}
