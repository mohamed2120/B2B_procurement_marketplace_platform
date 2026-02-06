'use client';

import { useState, useEffect } from 'react';
import { apiClients } from '@/lib/api';
import Card from '@/components/ui/Card';

export default function CompanyVerification() {
  const [companies, setCompanies] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchCompanies();
  }, []);

  const fetchCompanies = async () => {
    try {
      const response = await apiClients.company.get('/api/v1/companies?status=pending');
      setCompanies(response.data || []);
    } catch (error) {
      console.error('Failed to fetch companies:', error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div>
      <h1 className="text-3xl font-bold text-gray-900 mb-6">Company Verification</h1>

      {loading ? (
        <div>Loading...</div>
      ) : companies.length === 0 ? (
        <Card>
          <p className="text-gray-600">No pending company verifications.</p>
        </Card>
      ) : (
        <div className="space-y-4">
          {companies.map((company) => (
            <Card key={company.id}>
              <div className="flex justify-between items-start">
                <div>
                  <h3 className="font-semibold text-lg">{company.name}</h3>
                  <p className="text-gray-600 text-sm mt-1">Status: {company.status}</p>
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
