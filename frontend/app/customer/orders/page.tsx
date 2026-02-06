'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { isAuthenticated } from '@/lib/auth';
import { apiClients } from '@/lib/api';
import Header from '@/components/layout/Header';
import Sidebar from '@/components/layout/Sidebar';
import Card from '@/components/ui/Card';
import Link from 'next/link';
import { format } from 'date-fns';

interface Order {
  id: string;
  po_number: string;
  status: string;
  created_at: string;
  total_amount?: number;
}

interface Shipment {
  id: string;
  tracking_number: string;
  status: string;
  eta: string;
  po_id: string;
}

export default function OrdersPage() {
  const router = useRouter();
  const [orders, setOrders] = useState<Order[]>([]);
  const [shipments, setShipments] = useState<Shipment[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!isAuthenticated()) {
      router.push('/login');
      return;
    }

    fetchData();
  }, [router]);

  const fetchData = async () => {
    try {
      const [ordersRes, shipmentsRes] = await Promise.all([
        apiClients.procurement.get('/api/v1/purchase-orders?limit=100'),
        apiClients.logistics.get('/api/v1/shipments?limit=100'),
      ]);
      setOrders(ordersRes.data.items || ordersRes.data || []);
      setShipments(shipmentsRes.data.items || shipmentsRes.data || []);
    } catch (error) {
      console.error('Failed to fetch data:', error);
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
          <div className="flex justify-between items-center mb-6">
            <h1 className="text-3xl font-bold text-gray-900">Orders & Shipments</h1>
            <Link href="/customer/shipments">
              <button className="text-primary-600 hover:text-primary-800">View All Shipments</button>
            </Link>
          </div>

          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            <Card title="Purchase Orders">
              {orders.length === 0 ? (
                <p className="text-gray-500 text-center py-8">No orders found.</p>
              ) : (
                <div className="space-y-4">
                  {orders.map((order) => (
                    <div key={order.id} className="border rounded-lg p-4">
                      <div className="flex justify-between items-start">
                        <div>
                          <h3 className="font-semibold">{order.po_number}</h3>
                          <p className="text-sm text-gray-500">
                            {format(new Date(order.created_at), 'MMM dd, yyyy')}
                          </p>
                        </div>
                        <span
                          className={`px-2 py-1 text-xs rounded-full ${
                            order.status === 'completed'
                              ? 'bg-green-100 text-green-800'
                              : 'bg-yellow-100 text-yellow-800'
                          }`}
                        >
                          {order.status}
                        </span>
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </Card>

            <Card title="Shipments">
              {shipments.length === 0 ? (
                <p className="text-gray-500 text-center py-8">No shipments found.</p>
              ) : (
                <div className="space-y-4">
                  {shipments.map((shipment) => (
                    <div key={shipment.id} className="border rounded-lg p-4">
                      <div className="flex justify-between items-start">
                        <div>
                          <h3 className="font-semibold">{shipment.tracking_number}</h3>
                          <p className="text-sm text-gray-500">
                            ETA: {format(new Date(shipment.eta), 'MMM dd, yyyy')}
                          </p>
                        </div>
                        <span
                          className={`px-2 py-1 text-xs rounded-full ${
                            shipment.status === 'delivered'
                              ? 'bg-green-100 text-green-800'
                              : shipment.status === 'in_transit'
                              ? 'bg-blue-100 text-blue-800'
                              : 'bg-gray-100 text-gray-800'
                          }`}
                        >
                          {shipment.status}
                        </span>
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </Card>
          </div>
        </main>
      </div>
    </div>
  );
}
