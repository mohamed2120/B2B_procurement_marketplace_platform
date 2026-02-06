'use client';

import { useState, useEffect } from 'react';
import { apiClients } from '@/lib/api';
import Card from '@/components/ui/Card';

export default function SupplierOrders() {
  const [orders, setOrders] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchOrders();
  }, []);

  const fetchOrders = async () => {
    try {
      const response = await apiClients.procurement.get('/api/v1/orders');
      setOrders(response.data || []);
    } catch (error) {
      console.error('Failed to fetch orders:', error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div>
      <h1 className="text-3xl font-bold text-gray-900 mb-6">Customer Orders</h1>

      {loading ? (
        <div>Loading...</div>
      ) : orders.length === 0 ? (
        <Card>
          <p className="text-gray-600">No orders yet.</p>
        </Card>
      ) : (
        <div className="space-y-4">
          {orders.map((order) => (
            <Card key={order.id}>
              <div className="flex justify-between items-start">
                <div>
                  <h3 className="font-semibold text-lg">{order.po_number || order.id}</h3>
                  <p className="text-gray-600 text-sm mt-1">Status: {order.status}</p>
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
