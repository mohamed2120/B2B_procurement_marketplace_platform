'use client';

import { useState, useEffect } from 'react';
import { apiClients } from '@/lib/api';
import Card from '@/components/ui/Card';

export default function SupplierQuotes() {
  const [quotes, setQuotes] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchQuotes();
  }, []);

  const fetchQuotes = async () => {
    try {
      const response = await apiClients.procurement.get('/api/v1/quotes');
      setQuotes(response.data || []);
    } catch (error) {
      console.error('Failed to fetch quotes:', error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div>
      <h1 className="text-3xl font-bold text-gray-900 mb-6">My Quotes</h1>

      {loading ? (
        <div>Loading...</div>
      ) : quotes.length === 0 ? (
        <Card>
          <p className="text-gray-600">No quotes submitted yet.</p>
        </Card>
      ) : (
        <div className="space-y-4">
          {quotes.map((quote) => (
            <Card key={quote.id}>
              <div className="flex justify-between items-start">
                <div>
                  <h3 className="font-semibold text-lg">Quote #{quote.id.slice(0, 8)}</h3>
                  <p className="text-gray-600 text-sm mt-1">Status: {quote.status}</p>
                </div>
                <button className="text-primary-600 hover:text-primary-700">View →</button>
              </div>
            </Card>
          ))}
        </div>
      )}
    </div>
  );
}
