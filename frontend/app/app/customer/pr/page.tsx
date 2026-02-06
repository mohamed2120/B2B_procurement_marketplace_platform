'use client';

import { useState, useEffect } from 'react';
import { apiClients } from '@/lib/api';
import Card from '@/components/ui/Card';
import Link from 'next/link';
import Button from '@/components/ui/Button';

export default function CustomerPR() {
  const [prs, setPRs] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchPRs();
  }, []);

  const fetchPRs = async () => {
    try {
      const response = await apiClients.procurement.get('/api/v1/prs');
      setPRs(response.data);
    } catch (error) {
      console.error('Failed to fetch PRs:', error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-3xl font-bold text-gray-900">Purchase Requests</h1>
        <Link href="/app/customer/pr/create">
          <Button>Create PR</Button>
        </Link>
      </div>

      {loading ? (
        <div>Loading...</div>
      ) : prs.length === 0 ? (
        <Card>
          <p className="text-gray-600 mb-4">No purchase requests yet.</p>
          <Link href="/app/customer/pr/create">
            <Button>Create Your First PR</Button>
          </Link>
        </Card>
      ) : (
        <div className="space-y-4">
          {prs.map((pr) => (
            <Card key={pr.id}>
              <div className="flex justify-between items-start">
                <div>
                  <h3 className="font-semibold text-lg">{pr.title || pr.pr_number}</h3>
                  <p className="text-gray-600 text-sm mt-1">{pr.description}</p>
                  <div className="mt-2 flex gap-4 text-sm text-gray-500">
                    <span>Status: {pr.status}</span>
                    <span>Priority: {pr.priority}</span>
                  </div>
                </div>
                <Link href={`/app/customer/pr/${pr.id}`}>
                  <button className="text-primary-600 hover:text-primary-700">View →</button>
                </Link>
              </div>
            </Card>
          ))}
        </div>
      )}
    </div>
  );
}
