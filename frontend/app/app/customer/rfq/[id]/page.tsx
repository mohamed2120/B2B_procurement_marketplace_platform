'use client';

import { useState, useEffect } from 'react';
import { useParams, useRouter } from 'next/navigation';
import { apiClients } from '@/lib/api';
import { hasRole } from '@/lib/auth';
import Card from '@/components/ui/Card';
import Button from '@/components/ui/Button';
import Link from 'next/link';

interface RFQ {
  id: string;
  rfq_number: string;
  title: string;
  description: string;
  status: string;
  due_date: string;
  pr?: {
    id: string;
    pr_number: string;
    title: string;
    items?: PRItem[];
  };
  quotes?: Quote[];
}

interface PRItem {
  id: string;
  description: string;
  quantity: number;
  unit: string;
  specifications?: string;
}

interface Quote {
  id: string;
  quote_number: string;
  supplier_id: string;
  supplier_name?: string;
  status: string;
  total_amount: number;
  currency: string;
  valid_until: string;
  notes?: string;
  submitted_at: string;
  items?: QuoteItem[];
}

interface QuoteItem {
  id: string;
  pr_item_id: string;
  description: string;
  quantity: number;
  unit_price: number;
  total_price: number;
  lead_time: number;
}

export default function CustomerRFQDetail() {
  const params = useParams();
  const router = useRouter();
  const rfqId = params.id as string;
  
  const [rfq, setRFQ] = useState<RFQ | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [awarding, setAwarding] = useState<string | null>(null);
  const isProcurement = hasRole('procurement_manager');

  useEffect(() => {
    fetchRFQ();
  }, [rfqId]);

  const fetchRFQ = async () => {
    try {
      const response = await apiClients.procurement.get<RFQ>(`/api/v1/rfqs/${rfqId}`);
      setRFQ(response.data);
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to fetch RFQ');
    } finally {
      setLoading(false);
    }
  };

  const handleAwardQuote = async (quoteId: string) => {
    if (!rfq || !rfq.quotes) return;
    
    const quote = rfq.quotes.find(q => q.id === quoteId);
    if (!quote) return;

    setAwarding(quoteId);
    setError('');

    try {
      // Create PO from awarded quote
      const poData = {
        pr_id: rfq.pr?.id,
        rfq_id: rfq.id,
        quote_id: quoteId,
        supplier_id: quote.supplier_id,
        status: 'pending',
        total_amount: quote.total_amount,
        currency: quote.currency,
        payment_mode: 'DIRECT',
        payment_status: 'pending',
      };

      await apiClients.procurement.post('/api/v1/purchase-orders', poData);
      
      // Success - redirect to orders
      router.push('/app/customer/orders');
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to award quote');
      setAwarding(null);
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600 mx-auto mb-4"></div>
          <p className="text-gray-600">Loading RFQ...</p>
        </div>
      </div>
    );
  }

  if (error && !rfq) {
    return (
      <div>
        <Link href="/app/customer/rfq" className="text-primary-600 hover:text-primary-700 mb-4 inline-block">
          ← Back to RFQs
        </Link>
        <Card>
          <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">
            {error}
          </div>
        </Card>
      </div>
    );
  }

  if (!rfq) {
    return (
      <div>
        <Link href="/app/customer/rfq" className="text-primary-600 hover:text-primary-700 mb-4 inline-block">
          ← Back to RFQs
        </Link>
        <Card>
          <p className="text-gray-600">RFQ not found</p>
        </Card>
      </div>
    );
  }

  const quotes = rfq.quotes || [];

  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <div>
          <Link href="/app/customer/rfq" className="text-primary-600 hover:text-primary-700 mb-2 inline-block">
            ← Back to RFQs
          </Link>
          <h1 className="text-3xl font-bold text-gray-900">{rfq.rfq_number}</h1>
          <p className="text-gray-600 mt-1">{rfq.title}</p>
        </div>
        <div className="text-right">
          <span className={`px-3 py-1 rounded-full text-sm font-semibold ${
            rfq.status === 'closed' ? 'bg-gray-100 text-gray-800' :
            rfq.status === 'sent' ? 'bg-blue-100 text-blue-800' :
            'bg-yellow-100 text-yellow-800'
          }`}>
            {rfq.status.toUpperCase()}
          </span>
        </div>
      </div>

      {error && (
        <div className="mb-4 bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">
          {error}
        </div>
      )}

      {/* RFQ Details */}
      <Card title="RFQ Details" className="mb-6">
        <div className="space-y-3">
          <div>
            <span className="text-sm font-medium text-gray-700">Description:</span>
            <p className="text-gray-900 mt-1">{rfq.description}</p>
          </div>
          {rfq.pr && (
            <div>
              <span className="text-sm font-medium text-gray-700">Related PR:</span>
              <p className="text-gray-900 mt-1">{rfq.pr.pr_number} - {rfq.pr.title}</p>
            </div>
          )}
          <div>
            <span className="text-sm font-medium text-gray-700">Due Date:</span>
            <p className="text-gray-900 mt-1">{new Date(rfq.due_date).toLocaleDateString()}</p>
          </div>
        </div>
      </Card>

      {/* PR Items */}
      {rfq.pr?.items && rfq.pr.items.length > 0 && (
        <Card title="Requested Items" className="mb-6">
          <div className="overflow-x-auto">
            <table className="min-w-full divide-y divide-gray-200">
              <thead className="bg-gray-50">
                <tr>
                  <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Description</th>
                  <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Quantity</th>
                  <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Unit</th>
                  <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Specifications</th>
                </tr>
              </thead>
              <tbody className="bg-white divide-y divide-gray-200">
                {rfq.pr.items.map((item) => (
                  <tr key={item.id}>
                    <td className="px-4 py-3 text-sm text-gray-900">{item.description}</td>
                    <td className="px-4 py-3 text-sm text-gray-900">{item.quantity}</td>
                    <td className="px-4 py-3 text-sm text-gray-900">{item.unit || 'pcs'}</td>
                    <td className="px-4 py-3 text-sm text-gray-600">{item.specifications || '-'}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </Card>
      )}

      {/* Quotes Comparison */}
      <Card title={`Quotes (${quotes.length})`}>
        {quotes.length === 0 ? (
          <p className="text-gray-600">No quotes submitted yet.</p>
        ) : (
          <div className="space-y-4">
            <div className="overflow-x-auto">
              <table className="min-w-full divide-y divide-gray-200">
                <thead className="bg-gray-50">
                  <tr>
                    <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Supplier</th>
                    <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Quote #</th>
                    <th className="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase">Total Price</th>
                    <th className="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase">Lead Time</th>
                    <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Status</th>
                    <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Valid Until</th>
                    {isProcurement && (
                      <th className="px-4 py-3 text-center text-xs font-medium text-gray-500 uppercase">Action</th>
                    )}
                  </tr>
                </thead>
                <tbody className="bg-white divide-y divide-gray-200">
                  {quotes.map((quote) => {
                    // Calculate average lead time from items
                    const avgLeadTime = quote.items && quote.items.length > 0
                      ? Math.round(quote.items.reduce((sum, item) => sum + item.lead_time, 0) / quote.items.length)
                      : 0;

                    return (
                      <tr key={quote.id} className={quote.status === 'accepted' ? 'bg-green-50' : ''}>
                        <td className="px-4 py-3 text-sm text-gray-900">
                          {quote.supplier_name || `Supplier ${quote.supplier_id.slice(0, 8)}`}
                        </td>
                        <td className="px-4 py-3 text-sm text-gray-900">{quote.quote_number}</td>
                        <td className="px-4 py-3 text-sm text-gray-900 text-right font-semibold">
                          {quote.currency} {quote.total_amount.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}
                        </td>
                        <td className="px-4 py-3 text-sm text-gray-900 text-right">
                          {avgLeadTime} days
                        </td>
                        <td className="px-4 py-3 text-sm">
                          <span className={`px-2 py-1 rounded-full text-xs font-semibold ${
                            quote.status === 'accepted' ? 'bg-green-100 text-green-800' :
                            quote.status === 'rejected' ? 'bg-red-100 text-red-800' :
                            quote.status === 'expired' ? 'bg-gray-100 text-gray-800' :
                            'bg-blue-100 text-blue-800'
                          }`}>
                            {quote.status}
                          </span>
                        </td>
                        <td className="px-4 py-3 text-sm text-gray-600">
                          {new Date(quote.valid_until).toLocaleDateString()}
                        </td>
                        {isProcurement && (
                          <td className="px-4 py-3 text-center">
                            {quote.status === 'submitted' ? (
                              <Button
                                size="sm"
                                onClick={() => handleAwardQuote(quote.id)}
                                disabled={awarding === quote.id}
                              >
                                {awarding === quote.id ? 'Awarding...' : 'Award Quote'}
                              </Button>
                            ) : quote.status === 'accepted' ? (
                              <span className="text-sm text-green-600 font-semibold">Awarded</span>
                            ) : (
                              <span className="text-sm text-gray-400">-</span>
                            )}
                          </td>
                        )}
                      </tr>
                    );
                  })}
                </tbody>
              </table>
            </div>

            {/* Quote Details (expandable) */}
            {quotes.map((quote) => (
              <Card key={`details-${quote.id}`} title={`Quote Details: ${quote.quote_number}`} className="mt-4">
                {quote.notes && (
                  <div className="mb-4">
                    <span className="text-sm font-medium text-gray-700">Notes:</span>
                    <p className="text-gray-900 mt-1">{quote.notes}</p>
                  </div>
                )}
                {quote.items && quote.items.length > 0 && (
                  <div>
                    <span className="text-sm font-medium text-gray-700 mb-2 block">Line Items:</span>
                    <div className="overflow-x-auto">
                      <table className="min-w-full divide-y divide-gray-200">
                        <thead className="bg-gray-50">
                          <tr>
                            <th className="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Description</th>
                            <th className="px-4 py-2 text-right text-xs font-medium text-gray-500 uppercase">Quantity</th>
                            <th className="px-4 py-2 text-right text-xs font-medium text-gray-500 uppercase">Unit Price</th>
                            <th className="px-4 py-2 text-right text-xs font-medium text-gray-500 uppercase">Total</th>
                            <th className="px-4 py-2 text-right text-xs font-medium text-gray-500 uppercase">Lead Time</th>
                          </tr>
                        </thead>
                        <tbody className="bg-white divide-y divide-gray-200">
                          {quote.items.map((item) => (
                            <tr key={item.id}>
                              <td className="px-4 py-2 text-sm text-gray-900">{item.description}</td>
                              <td className="px-4 py-2 text-sm text-gray-900 text-right">{item.quantity}</td>
                              <td className="px-4 py-2 text-sm text-gray-900 text-right">
                                {quote.currency} {item.unit_price.toFixed(2)}
                              </td>
                              <td className="px-4 py-2 text-sm text-gray-900 text-right font-semibold">
                                {quote.currency} {item.total_price.toFixed(2)}
                              </td>
                              <td className="px-4 py-2 text-sm text-gray-900 text-right">{item.lead_time} days</td>
                            </tr>
                          ))}
                        </tbody>
                      </table>
                    </div>
                  </div>
                )}
              </Card>
            ))}
          </div>
        )}
      </Card>
    </div>
  );
}
