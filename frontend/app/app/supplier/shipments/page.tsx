'use client';

import { useState, useEffect } from 'react';
import { apiClients } from '@/lib/api';
import Card from '@/components/ui/Card';

export default function SupplierShipments() {
  const [shipments, setShipments] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchShipments();
  }, []);

  const fetchShipments = async () => {
    try {
      const response = await apiClients.logistics.get('/api/v1/shipments');
      setShipments(response.data || []);
    } catch (error) {
      console.error('Failed to fetch shipments:', error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div>
      <h1 className="text-3xl font-bold text-gray-900 mb-6">Shipments</h1>

      {loading ? (
        <div>Loading...</div>
      ) : shipments.length === 0 ? (
        <Card>
          <p className="text-gray-600">No shipments yet.</p>
        </Card>
      ) : (
        <div className="space-y-4">
          {shipments.map((shipment) => (
            <Card key={shipment.id}>
              <div className="flex justify-between items-start">
                <div>
                  <h3 className="font-semibold text-lg">{shipment.tracking_number || shipment.id}</h3>
                  <p className="text-gray-600 text-sm mt-1">Status: {shipment.status}</p>
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
