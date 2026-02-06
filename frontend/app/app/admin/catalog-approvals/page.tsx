'use client';

import { useState, useEffect } from 'react';
import { apiClients } from '@/lib/api';
import Card from '@/components/ui/Card';

export default function CatalogApprovals() {
  const [items, setItems] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchItems();
  }, []);

  const fetchItems = async () => {
    try {
      const response = await apiClients.catalog.get('/api/v1/lib-parts?status=pending');
      setItems(response.data || []);
    } catch (error) {
      console.error('Failed to fetch catalog items:', error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div>
      <h1 className="text-3xl font-bold text-gray-900 mb-6">Catalog Approvals</h1>

      {loading ? (
        <div>Loading...</div>
      ) : items.length === 0 ? (
        <Card>
          <p className="text-gray-600">No pending catalog approvals.</p>
        </Card>
      ) : (
        <div className="space-y-4">
          {items.map((item) => (
            <Card key={item.id}>
              <div className="flex justify-between items-start">
                <div>
                  <h3 className="font-semibold text-lg">{item.name || item.part_number}</h3>
                  <p className="text-gray-600 text-sm mt-1">Status: {item.status}</p>
                </div>
                <div className="flex gap-2">
                  <button className="bg-green-600 text-white px-4 py-2 rounded-lg hover:bg-green-700">Approve</button>
                  <button className="bg-red-600 text-white px-4 py-2 rounded-lg hover:bg-red-700">Reject</button>
                </div>
              </div>
            </Card>
          ))}
        </div>
      )}
    </div>
  );
}
