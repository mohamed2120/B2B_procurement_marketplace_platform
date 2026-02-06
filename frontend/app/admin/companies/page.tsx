'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { isAuthenticated, hasRole } from '@/lib/auth';
import { apiClients } from '@/lib/api';
import Header from '@/components/layout/Header';
import Sidebar from '@/components/layout/Sidebar';
import Card from '@/components/ui/Card';
import Button from '@/components/ui/Button';

interface Company {
  id: string;
  name: string;
  legal_name: string;
  status: string;
  verification_status: string;
}

export default function CompanyVerificationPage() {
  const router = useRouter();
  const [companies, setCompanies] = useState<Company[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!isAuthenticated() || !hasRole('admin')) {
      router.push('/login');
      return;
    }

    fetchCompanies();
  }, [router]);

  const fetchCompanies = async () => {
    try {
      const response = await apiClients.company.get('/api/v1/companies?limit=100');
      setCompanies(response.data.items || response.data || []);
    } catch (error) {
      console.error('Failed to fetch companies:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleApprove = async (companyId: string) => {
    try {
      await apiClients.company.post(`/api/v1/companies/${companyId}/approve`, {});
      fetchCompanies();
    } catch (error) {
      console.error('Failed to approve company:', error);
    }
  };

  if (loading) {
    return <div>Loading...</div>;
  }

  const pendingCompanies = companies.filter(c => c.status === 'pending');

  return (
    <div className="min-h-screen bg-gray-50">
      <Header />
      <div className="flex">
        <Sidebar />
        <main className="flex-1 p-8">
          <h1 className="text-3xl font-bold text-gray-900 mb-6">Company Verification</h1>

          <Card>
            {pendingCompanies.length === 0 ? (
              <p className="text-gray-500 text-center py-8">No pending companies to verify.</p>
            ) : (
              <div className="space-y-4">
                {pendingCompanies.map((company) => (
                  <div key={company.id} className="border rounded-lg p-4">
                    <div className="flex justify-between items-start">
                      <div>
                        <h3 className="font-semibold text-lg">{company.name}</h3>
                        <p className="text-gray-600">{company.legal_name}</p>
                        <p className="text-sm text-gray-500 mt-2">
                          Status: <span className="font-medium">{company.status}</span>
                        </p>
                      </div>
                      <Button size="sm" onClick={() => handleApprove(company.id)}>
                        Approve
                      </Button>
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
