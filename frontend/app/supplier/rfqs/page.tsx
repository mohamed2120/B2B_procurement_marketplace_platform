'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { isAuthenticated, hasRole } from '@/lib/auth';
import { apiClients } from '@/lib/api';
import Header from '@/components/layout/Header';
import Sidebar from '@/components/layout/Sidebar';
import Card from '@/components/ui/Card';
import Button from '@/components/ui/Button';
import Link from 'next/link';
import { format } from 'date-fns';

interface RFQ {
  id: string;
  rfq_number: string;
  title: string;
  due_date: string;
  status: string;
}

export default function SupplierRFQPage() {
  const router = useRouter();
  const [rfqs, setRFQs] = useState<RFQ[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!isAuthenticated() || !hasRole('supplier')) {
      router.push('/login');
      return;
    }

    fetchRFQs();
  }, [router]);

  const fetchRFQs = async () => {
    try {
      const response = await apiClients.procurement.get('/api/v1/rfqs?limit=100');
      setRFQs(response.data.items || response.data || []);
    } catch (error) {
      console.error('Failed to fetch RFQs:', error);
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return <div>Loading...</div>;
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <Header />
      <div className="flex">
        <Sidebar />
        <main className="flex-1 p-8">
          <h1 className="text-3xl font-bold text-gray-900 mb-6">RFQ Inbox</h1>

          <Card>
            {rfqs.length === 0 ? (
              <p className="text-gray-500 text-center py-8">No RFQs available.</p>
            ) : (
              <div className="space-y-4">
                {rfqs.map((rfq) => (
                  <div key={rfq.id} className="border rounded-lg p-4 hover:bg-gray-50">
                    <div className="flex justify-between items-start">
                      <div>
                        <h3 className="font-semibold text-lg">{rfq.rfq_number}</h3>
                        <p className="text-gray-600">{rfq.title}</p>
                        <p className="text-sm text-gray-500 mt-2">
                          Due: {format(new Date(rfq.due_date), 'MMM dd, yyyy')}
                        </p>
                      </div>
                      <Link href={`/supplier/rfqs/${rfq.id}/quote`}>
                        <Button size="sm">Submit Quote</Button>
                      </Link>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </Card>
        </main>
      </div>
    </div>
  );
}
