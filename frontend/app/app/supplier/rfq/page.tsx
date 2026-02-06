'use client';

import { useState, useEffect } from 'react';
import { apiClients } from '@/lib/api';
import Card from '@/components/ui/Card';

export default function SupplierRFQ() {
  const [rfqs, setRFQs] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchRFQs();
  }, []);

  const fetchRFQs = async () => {
    try {
      const response = await apiClients.procurement.get('/api/v1/rfqs');
      setRFQs(response.data || []);
    } catch (error) {
      console.error('Failed to fetch RFQs:', error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div>
      <h1 className="text-3xl font-bold text-gray-900 mb-6">RFQ Inbox</h1>

      {loading ? (
        <div>Loading...</div>
      ) : rfqs.length === 0 ? (
        <Card>
          <p className="text-gray-600">No RFQs in your inbox.</p>
        </Card>
      ) : (
        <div className="space-y-4">
          {rfqs.map((rfq) => (
            <Card key={rfq.id}>
              <div className="flex justify-between items-start">
                <div>
                  <h3 className="font-semibold text-lg">{rfq.rfq_number}</h3>
                  <p className="text-gray-600 text-sm mt-1">Status: {rfq.status}</p>
                </div>
                <button className="text-primary-600 hover:text-primary-700">Submit Quote →</button>
              </div>
            </Card>
          ))}
        </div>
      )}
    </div>
  );
}
