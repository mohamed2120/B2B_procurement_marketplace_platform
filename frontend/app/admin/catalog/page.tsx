'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { isAuthenticated, hasRole } from '@/lib/auth';
import { apiClients } from '@/lib/api';
import Header from '@/components/layout/Header';
import Sidebar from '@/components/layout/Sidebar';
import Card from '@/components/ui/Card';
import Button from '@/components/ui/Button';

interface Part {
  id: string;
  part_number: string;
  name: string;
  description: string;
  status: string;
}

export default function CatalogApprovalsPage() {
  const router = useRouter();
  const [parts, setParts] = useState<Part[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!isAuthenticated() || !hasRole('admin')) {
      router.push('/login');
      return;
    }

    fetchPendingParts();
  }, [router]);

  const fetchPendingParts = async () => {
    try {
      const response = await apiClients.catalog.get('/api/v1/parts?status=pending&limit=100');
      setParts(response.data.items || response.data || []);
    } catch (error) {
      console.error('Failed to fetch parts:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleApprove = async (partId: string) => {
    try {
      await apiClients.catalog.post(`/api/v1/parts/${partId}/approve`, {});
      fetchPendingParts();
    } catch (error) {
      console.error('Failed to approve part:', error);
    }
  };

  const handleReject = async (partId: string) => {
    if (!confirm('Are you sure you want to reject this part?')) return;
    try {
      await apiClients.catalog.post(`/api/v1/parts/${partId}/reject`, {
        reason: 'Rejected by admin',
      });
      fetchPendingParts();
    } catch (error) {
      console.error('Failed to reject part:', error);
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
          <h1 className="text-3xl font-bold text-gray-900 mb-6">Catalog Approvals</h1>

          <Card>
            {parts.length === 0 ? (
              <p className="text-gray-500 text-center py-8">No pending parts to approve.</p>
            ) : (
              <div className="space-y-4">
                {parts.map((part) => (
                  <div key={part.id} className="border rounded-lg p-4">
                    <div className="flex justify-between items-start">
                      <div className="flex-1">
                        <h3 className="font-semibold text-lg">{part.part_number}</h3>
                        <p className="text-gray-700">{part.name}</p>
                        <p className="text-sm text-gray-500 mt-2">{part.description}</p>
                      </div>
                      <div className="flex space-x-2">
                        <Button size="sm" onClick={() => handleApprove(part.id)}>
                          Approve
                        </Button>
                        <Button size="sm" variant="danger" onClick={() => handleReject(part.id)}>
                          Reject
                        </Button>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </Card>
        </main>
      </div>
    </div>
  );
}
