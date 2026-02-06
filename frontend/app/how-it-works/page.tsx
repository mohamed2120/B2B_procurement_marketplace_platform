'use client';

import PublicLayout from '@/components/layout/PublicLayout';
import Link from 'next/link';

export default function HowItWorks() {
  return (
    <PublicLayout>
      
      <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        <h1 className="text-4xl font-bold text-gray-900 mb-8">How It Works</h1>
        
        <div className="space-y-12">
          <section>
            <h2 className="text-2xl font-semibold mb-4 text-primary-600">For Buyers</h2>
            <ol className="list-decimal list-inside space-y-4 text-gray-700">
              <li>Create a purchase request (PR) for items you need</li>
              <li>Submit RFQs to multiple suppliers</li>
              <li>Compare quotes and select the best supplier</li>
              <li>Approve purchase orders and track shipments</li>
              <li>Receive goods and complete transactions</li>
            </ol>
          </section>

          <section>
            <h2 className="text-2xl font-semibold mb-4 text-primary-600">For Suppliers</h2>
            <ol className="list-decimal list-inside space-y-4 text-gray-700">
              <li>Create your company profile and get verified</li>
              <li>List your products and services in the marketplace</li>
              <li>Receive RFQs from buyers</li>
              <li>Submit competitive quotes</li>
              <li>Fulfill orders and get paid securely</li>
            </ol>
          </section>

          <section>
            <h2 className="text-2xl font-semibold mb-4 text-primary-600">Key Features</h2>
            <ul className="list-disc list-inside space-y-2 text-gray-700">
              <li>Multi-tenant architecture for complete data isolation</li>
              <li>Role-based access control (RBAC) for security</li>
              <li>Real-time notifications and chat collaboration</li>
              <li>Escrow payment protection for both parties</li>
              <li>Comprehensive catalog and equipment management</li>
              <li>Advanced search and filtering capabilities</li>
            </ul>
          </section>
        </div>

        <div className="mt-12 text-center">
          <Link
            href="/register"
            className="bg-primary-600 text-white px-8 py-3 rounded-lg font-semibold hover:bg-primary-700 transition inline-block"
          >
            Get Started Today
          </Link>
        </div>
      </div>
    </PublicLayout>
  );
}
