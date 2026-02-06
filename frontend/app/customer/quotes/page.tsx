'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { isAuthenticated } from '@/lib/auth';
import { apiClients } from '@/lib/api';
import Header from '@/components/layout/Header';
import Sidebar from '@/components/layout/Sidebar';
import Card from '@/components/ui/Card';
import Button from '@/components/ui/Button';
import { format } from 'date-fns';

interface Quote {
  id: string;
  quote_number: string;
  rfq_id: string;
  supplier_id: string;
  total_amount: number;
  status: string;
  submitted_at: string;
}

export default function QuotesPage() {
  const router = useRouter();
  const [quotes, setQuotes] = useState<Quote[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!isAuthenticated()) {
      router.push('/login');
      return;
    }

    fetchQuotes();
  }, [router]);

  const fetchQuotes = async () => {
    try {
      const response = await apiClients.procurement.get('/api/v1/quotes?limit=100');
      setQuotes(response.data.items || response.data || []);
    } catch (error) {
      console.error('Failed to fetch quotes:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleAward = async (quoteId: string) => {
    try {
      // Award quote and create PO
      await apiClients.procurement.post('/api/v1/quotes/' + quoteId + '/award', {});
      fetchQuotes();
    } catch (error) {
      console.error('Failed to award quote:', error);
    }
  };

  if (loading) {
    return <div>Loading...</div>;
  }

  // Group quotes by RFQ for comparison
  const quotesByRFQ = quotes.reduce((acc, quote) => {
    if (!acc[quote.rfq_id]) {
      acc[quote.rfq_id] = [];
    }
    acc[quote.rfq_id].push(quote);
    return acc;
  }, {} as Record<string, Quote[]>);

  return (
    <div className="min-h-screen bg-gray-50">
      <Header />
      <div className="flex">
        <Sidebar />
        <main className="flex-1 p-8">
          <h1 className="text-3xl font-bold text-gray-900 mb-6">Quotes</h1>

          <div className="space-y-6">
            {Object.entries(quotesByRFQ).map(([rfqId, rfqQuotes]) => (
              <Card key={rfqId} title={`RFQ: ${rfqQuotes[0].quote_number}`}>
                <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                  {rfqQuotes.map((quote) => (
                    <div key={quote.id} className="border rounded-lg p-4">
                      <div className="flex justify-between items-start mb-2">
                        <span className="font-semibold">{quote.quote_number}</span>
                        <span className="text-2xl font-bold text-primary-600">
                          ${quote.total_amount.toFixed(2)}
                        </span>
                      </div>
                      <p className="text-sm text-gray-500 mb-4">
                        Submitted: {format(new Date(quote.submitted_at), 'MMM dd, yyyy')}
                      </p>
                      <Button
                        size="sm"
                        className="w-full"
                        onClick={() => handleAward(quote.id)}
                        disabled={quote.status === 'awarded'}
                      >
                        {quote.status === 'awarded' ? 'Awarded' : 'Award Quote'}
                      </Button>
                    </div>
                  ))}
                </div>
              </Card>
            ))}
          </div>

          {quotes.length === 0 && (
            <Card>
              <p className="text-gray-500 text-center py-8">No quotes found.</p>
            </Card>
          )}
        </main>
      </div>
    </div>
  );
}
