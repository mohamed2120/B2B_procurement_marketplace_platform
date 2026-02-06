'use client';

import { useState, useEffect } from 'react';
import { apiClients } from '@/lib/api';
import Card from '@/components/ui/Card';

export default function AdminDisputes() {
  const [disputes, setDisputes] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchDisputes();
  }, []);

  const fetchDisputes = async () => {
    try {
      const response = await apiClients.collaboration.get('/api/v1/disputes');
      setDisputes(response.data || []);
    } catch (error) {
      console.error('Failed to fetch disputes:', error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div>
      <h1 className="text-3xl font-bold text-gray-900 mb-6">Disputes</h1>

      {loading ? (
        <div>Loading...</div>
      ) : disputes.length === 0 ? (
        <Card>
          <p className="text-gray-600">No disputes.</p>
        </Card>
      ) : (
        <div className="space-y-4">
          {disputes.map((dispute) => (
            <Card key={dispute.id}>
              <div className="flex justify-between items-start">
                <div>
                  <h3 className="font-semibold text-lg">{dispute.title}</h3>
                  <p className="text-gray-600 text-sm mt-1">Status: {dispute.status}</p>
                </div>
                <button className="text-primary-600 hover:text-primary-700">Review →</button>
              </div>
            </Card>
          ))}
        </div>
      )}
    </div>
  );
}
