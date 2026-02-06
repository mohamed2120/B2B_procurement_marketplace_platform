'use client';

import { useState, useEffect } from 'react';
import { useParams, useRouter } from 'next/navigation';
import { apiClients } from '@/lib/api';
import { getUser } from '@/lib/auth';
import Card from '@/components/ui/Card';
import Button from '@/components/ui/Button';
import Input from '@/components/ui/Input';
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
}

interface PRItem {
  id: string;
  description: string;
  quantity: number;
  unit: string;
  specifications?: string;
}

interface QuoteItem {
  pr_item_id: string;
  description: string;
  quantity: number;
  unit_price: number;
  total_price: number;
  lead_time: number;
}

interface QuoteForm {
  rfq_id: string;
  notes: string;
  currency: string;
  valid_until: string;
  items: QuoteItem[];
}

export default function SupplierRFQDetail() {
  const params = useParams();
  const router = useRouter();
  const rfqId = params.id as string;
  const user = getUser();
  
  const [rfq, setRFQ] = useState<RFQ | null>(null);
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState('');
  
  const [formData, setFormData] = useState<QuoteForm>({
    rfq_id: rfqId,
    notes: '',
    currency: 'USD',
    valid_until: '',
    items: [],
  });

  useEffect(() => {
    fetchRFQ();
  }, [rfqId]);

  useEffect(() => {
    // Initialize quote items from PR items when RFQ loads
    if (rfq?.pr?.items && formData.items.length === 0) {
      const items: QuoteItem[] = rfq.pr.items.map(item => ({
        pr_item_id: item.id,
        description: item.description,
        quantity: item.quantity,
        unit_price: 0,
        total_price: 0,
        lead_time: 0,
      }));
      setFormData(prev => ({ ...prev, items }));
      
      // Set default valid until (30 days from now)
      const validUntil = new Date();
      validUntil.setDate(validUntil.getDate() + 30);
      setFormData(prev => ({ ...prev, valid_until: validUntil.toISOString().split('T')[0] }));
    }
  }, [rfq]);

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

  const updateItem = (index: number, field: keyof QuoteItem, value: any) => {
    const newItems = [...formData.items];
    const item = { ...newItems[index] };
    
    if (field === 'unit_price') {
      item.unit_price = parseFloat(value) || 0;
      item.total_price = item.unit_price * item.quantity;
    } else if (field === 'lead_time') {
      item.lead_time = parseInt(value) || 0;
    } else {
      (item as any)[field] = value;
    }
    
    newItems[index] = item;
    setFormData(prev => ({ ...prev, items: newItems }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setSubmitting(true);

    // Validation
    if (!formData.valid_until) {
      setError('Valid until date is required');
      setSubmitting(false);
      return;
    }

    if (formData.items.some(item => item.unit_price <= 0)) {
      setError('All items must have a unit price greater than 0');
      setSubmitting(false);
      return;
    }

    try {
      // Get supplier ID from user's tenant (assuming supplier company = tenant)
      const supplierId = user?.tenant_id;
      if (!supplierId) {
        throw new Error('Supplier ID not found');
      }

      // Calculate total amount
      const totalAmount = formData.items.reduce((sum, item) => sum + item.total_price, 0);

      const quoteData = {
        rfq_id: rfqId,
        supplier_id: supplierId,
        status: 'submitted',
        total_amount: totalAmount,
        currency: formData.currency,
        valid_until: formData.valid_until,
        notes: formData.notes,
        items: formData.items,
      };

      await apiClients.procurement.post('/api/v1/quotes', quoteData);
      
      // Success - redirect to quotes
      router.push('/app/supplier/quotes');
    } catch (err: any) {
      setError(err.response?.data?.error || err.message || 'Failed to submit quote');
      setSubmitting(false);
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
        <Link href="/app/supplier/rfq" className="text-primary-600 hover:text-primary-700 mb-4 inline-block">
          ← Back to RFQ Inbox
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
        <Link href="/app/supplier/rfq" className="text-primary-600 hover:text-primary-700 mb-4 inline-block">
          ← Back to RFQ Inbox
        </Link>
        <Card>
          <p className="text-gray-600">RFQ not found</p>
        </Card>
      </div>
    );
  }

  const totalAmount = formData.items.reduce((sum, item) => sum + item.total_price, 0);

  return (
    <div>
      <div className="mb-6">
        <Link href="/app/supplier/rfq" className="text-primary-600 hover:text-primary-700 mb-2 inline-block">
          ← Back to RFQ Inbox
        </Link>
        <h1 className="text-3xl font-bold text-gray-900">{rfq.rfq_number}</h1>
        <p className="text-gray-600 mt-1">{rfq.title}</p>
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

      {/* Quote Submission Form */}
      <Card title="Submit Quote">
        <form onSubmit={handleSubmit} className="space-y-6">
          {/* Quote Header */}
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Currency</label>
              <select
                value={formData.currency}
                onChange={(e) => setFormData(prev => ({ ...prev, currency: e.target.value }))}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-primary-500 focus:border-primary-500"
                required
              >
                <option value="USD">USD</option>
                <option value="EUR">EUR</option>
                <option value="GBP">GBP</option>
              </select>
            </div>
            <Input
              label="Valid Until"
              type="date"
              value={formData.valid_until}
              onChange={(e) => setFormData(prev => ({ ...prev, valid_until: e.target.value }))}
              required
            />
          </div>

          {/* Notes */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Notes (Optional)</label>
            <textarea
              value={formData.notes}
              onChange={(e) => setFormData(prev => ({ ...prev, notes: e.target.value }))}
              rows={3}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-primary-500 focus:border-primary-500"
              placeholder="Add any additional notes or terms..."
            />
          </div>

          {/* Quote Items */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-3">Quote Items</label>
            <div className="overflow-x-auto">
              <table className="min-w-full divide-y divide-gray-200">
                <thead className="bg-gray-50">
                  <tr>
                    <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Description</th>
                    <th className="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase">Quantity</th>
                    <th className="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase">Unit Price</th>
                    <th className="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase">Total</th>
                    <th className="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase">Lead Time (days)</th>
                  </tr>
                </thead>
                <tbody className="bg-white divide-y divide-gray-200">
                  {formData.items.map((item, index) => (
                    <tr key={item.pr_item_id}>
                      <td className="px-4 py-3 text-sm text-gray-900">
                        {item.description}
                        {rfq.pr?.items?.find(pi => pi.id === item.pr_item_id)?.specifications && (
                          <div className="text-xs text-gray-500 mt-1">
                            Specs: {rfq.pr.items.find(pi => pi.id === item.pr_item_id)?.specifications}
                          </div>
                        )}
                      </td>
                      <td className="px-4 py-3 text-sm text-gray-900 text-right">{item.quantity}</td>
                      <td className="px-4 py-3 text-right">
                        <input
                          type="number"
                          step="0.01"
                          min="0"
                          value={item.unit_price || ''}
                          onChange={(e) => updateItem(index, 'unit_price', e.target.value)}
                          className="w-24 px-2 py-1 border border-gray-300 rounded text-sm text-right focus:ring-primary-500 focus:border-primary-500"
                          required
                        />
                      </td>
                      <td className="px-4 py-3 text-sm text-gray-900 text-right font-semibold">
                        {formData.currency} {item.total_price.toFixed(2)}
                      </td>
                      <td className="px-4 py-3 text-right">
                        <input
                          type="number"
                          min="0"
                          value={item.lead_time || ''}
                          onChange={(e) => updateItem(index, 'lead_time', e.target.value)}
                          className="w-20 px-2 py-1 border border-gray-300 rounded text-sm text-right focus:ring-primary-500 focus:border-primary-500"
                          required
                        />
                      </td>
                    </tr>
                  ))}
                </tbody>
                <tfoot className="bg-gray-50">
                  <tr>
                    <td colSpan={3} className="px-4 py-3 text-sm font-semibold text-gray-900 text-right">
                      Total Amount:
                    </td>
                    <td className="px-4 py-3 text-sm font-bold text-gray-900 text-right">
                      {formData.currency} {totalAmount.toFixed(2)}
                    </td>
                    <td></td>
                  </tr>
                </tfoot>
              </table>
            </div>
          </div>

          {/* Submit Button */}
          <div className="flex gap-4 pt-4 border-t">
            <Button type="submit" disabled={submitting}>
              {submitting ? 'Submitting...' : 'Submit Quote'}
            </Button>
            <Button type="button" variant="secondary" onClick={() => router.back()}>
              Cancel
            </Button>
          </div>
        </form>
      </Card>
    </div>
  );
}
