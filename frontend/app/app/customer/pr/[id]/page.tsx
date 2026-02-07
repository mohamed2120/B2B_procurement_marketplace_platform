'use client';

import { useState, useEffect } from 'react';
import { useParams, useRouter } from 'next/navigation';
import { apiClients } from '@/lib/api';
import { hasRole } from '@/lib/auth';
import Card from '@/components/ui/Card';
import Button from '@/components/ui/Button';
import Link from 'next/link';

interface PR {
  id: string;
  pr_number: string;
  title: string;
  description: string;
  status: string;
  priority: string;
  requested_by: string;
  created_at: string;
  items?: PRItem[];
}

interface PRItem {
  id: string;
  description: string;
  quantity: number;
  unit: string;
  specifications?: string;
  estimated_cost?: number;
}

export default function PRDetailPage() {
  const params = useParams();
  const router = useRouter();
  const prId = params.id as string;
  
  const [pr, setPR] = useState<PR | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [approving, setApproving] = useState(false);
  const [rejecting, setRejecting] = useState(false);
  
  const isProcurement = hasRole('procurement_manager') || hasRole('procurement');
  const isRequester = hasRole('requester');

  useEffect(() => {
    fetchPR();
  }, [prId]);

  const fetchPR = async () => {
    try {
      const response = await apiClients.procurement.get<PR>(`/api/v1/purchase-requests/${prId}`);
      setPR(response.data);
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to fetch PR');
    } finally {
      setLoading(false);
    }
  };

  const handleApprove = async () => {
    if (!pr) return;
    
    setApproving(true);
    setError('');

    try {
      await apiClients.procurement.post(`/api/v1/purchase-requests/${prId}/approve`, {
        approved: true,
      });
      
      // Refresh PR data
      await fetchPR();
      // Show success message
      alert('PR approved successfully!');
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to approve PR');
    } finally {
      setApproving(false);
    }
  };

  const handleReject = async () => {
    if (!pr) return;
    
    const reason = prompt('Please provide a reason for rejection:');
    if (!reason) {
      return;
    }

    setRejecting(true);
    setError('');

    try {
      await apiClients.procurement.post(`/api/v1/purchase-requests/${prId}/approve`, {
        approved: false,
        rejection_reason: reason,
      });
      
      // Refresh PR data
      await fetchPR();
      // Show success message
      alert('PR rejected successfully!');
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to reject PR');
    } finally {
      setRejecting(false);
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600 mx-auto mb-4"></div>
          <p className="text-gray-600">Loading PR...</p>
        </div>
      </div>
    );
  }

  if (error && !pr) {
    return (
      <div>
        <Link href="/app/customer/pr" className="text-primary-600 hover:text-primary-700 mb-4 inline-block">
          ← Back to PRs
        </Link>
        <Card>
          <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">
            {error}
          </div>
        </Card>
      </div>
    );
  }

  if (!pr) {
    return (
      <div>
        <Link href="/app/customer/pr" className="text-primary-600 hover:text-primary-700 mb-4 inline-block">
          ← Back to PRs
        </Link>
        <Card>
          <p className="text-gray-600">PR not found</p>
        </Card>
      </div>
    );
  }

  const canApprove = isProcurement && (pr.status === 'pending' || pr.status === 'submitted');
  const canView = isRequester || isProcurement;

  if (!canView) {
    return (
      <div>
        <Link href="/app/customer/pr" className="text-primary-600 hover:text-primary-700 mb-4 inline-block">
          ← Back to PRs
        </Link>
        <Card>
          <div className="bg-yellow-50 border border-yellow-200 text-yellow-700 px-4 py-3 rounded">
            You don't have permission to view this PR.
          </div>
        </Card>
      </div>
    );
  }

  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <div>
          <Link href="/app/customer/pr" className="text-primary-600 hover:text-primary-700 mb-2 inline-block">
            ← Back to PRs
          </Link>
          <h1 className="text-3xl font-bold text-gray-900">{pr.pr_number}</h1>
          <p className="text-gray-600 mt-1">{pr.title}</p>
        </div>
        <div className="text-right">
          <span className={`px-3 py-1 rounded-full text-sm font-semibold ${
            pr.status === 'approved' ? 'bg-green-100 text-green-800' :
            pr.status === 'rejected' ? 'bg-red-100 text-red-800' :
            pr.status === 'pending' ? 'bg-yellow-100 text-yellow-800' :
            'bg-gray-100 text-gray-800'
          }`}>
            {pr.status.toUpperCase()}
          </span>
        </div>
      </div>

      {error && (
        <div className="mb-4 bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">
          {error}
        </div>
      )}

      {/* PR Details */}
      <Card title="PR Details" className="mb-6">
        <div className="space-y-3">
          <div>
            <span className="text-sm font-medium text-gray-700">Description:</span>
            <p className="text-gray-900 mt-1">{pr.description}</p>
          </div>
          <div className="grid grid-cols-2 gap-4">
            <div>
              <span className="text-sm font-medium text-gray-700">Priority:</span>
              <p className="text-gray-900 mt-1">{pr.priority}</p>
            </div>
            <div>
              <span className="text-sm font-medium text-gray-700">Requested By:</span>
              <p className="text-gray-900 mt-1">{pr.requested_by}</p>
            </div>
            <div>
              <span className="text-sm font-medium text-gray-700">Created:</span>
              <p className="text-gray-900 mt-1">{new Date(pr.created_at).toLocaleDateString()}</p>
            </div>
          </div>
        </div>
      </Card>

      {/* PR Items */}
      {pr.items && pr.items.length > 0 && (
        <Card title="Requested Items" className="mb-6">
          <div className="overflow-x-auto">
            <table className="min-w-full divide-y divide-gray-200">
              <thead className="bg-gray-50">
                <tr>
                  <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Description</th>
                  <th className="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase">Quantity</th>
                  <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Unit</th>
                  <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Specifications</th>
                  {pr.items.some(item => item.estimated_cost) && (
                    <th className="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase">Est. Cost</th>
                  )}
                </tr>
              </thead>
              <tbody className="bg-white divide-y divide-gray-200">
                {pr.items.map((item) => (
                  <tr key={item.id}>
                    <td className="px-4 py-3 text-sm text-gray-900">{item.description}</td>
                    <td className="px-4 py-3 text-sm text-gray-900 text-right">{item.quantity}</td>
                    <td className="px-4 py-3 text-sm text-gray-900">{item.unit || 'pcs'}</td>
                    <td className="px-4 py-3 text-sm text-gray-600">{item.specifications || '-'}</td>
                    {pr.items!.some(i => i.estimated_cost) && (
                      <td className="px-4 py-3 text-sm text-gray-900 text-right">
                        {item.estimated_cost ? `$${item.estimated_cost.toFixed(2)}` : '-'}
                      </td>
                    )}
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </Card>
      )}

      {/* Approval Actions (Procurement only) */}
      {canApprove && (
        <Card title="Actions" className="mb-6">
          <div className="flex gap-4">
            <Button
              onClick={handleApprove}
              disabled={approving || rejecting}
              className="bg-green-600 hover:bg-green-700"
            >
              {approving ? 'Approving...' : 'Approve PR'}
            </Button>
            <Button
              onClick={handleReject}
              disabled={approving || rejecting}
              variant="secondary"
              className="bg-red-600 hover:bg-red-700 text-white"
            >
              {rejecting ? 'Rejecting...' : 'Reject PR'}
            </Button>
          </div>
        </Card>
      )}

      {/* Status Message */}
      {pr.status === 'approved' && (
        <Card className="bg-green-50 border-green-200">
          <p className="text-green-800 font-semibold">✓ This PR has been approved and is ready for RFQ creation.</p>
        </Card>
      )}
      {pr.status === 'rejected' && (
        <Card className="bg-red-50 border-red-200">
          <p className="text-red-800 font-semibold">✗ This PR has been rejected.</p>
        </Card>
      )}
    </div>
  );
}
