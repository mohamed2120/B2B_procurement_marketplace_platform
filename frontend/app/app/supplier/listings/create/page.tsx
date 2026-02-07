'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { apiClients } from '@/lib/api';
import { getUser } from '@/lib/auth';
import Card from '@/components/ui/Card';
import Button from '@/components/ui/Button';
import Input from '@/components/ui/Input';
import Link from 'next/link';

interface ListingForm {
  title: string;
  description: string;
  part_number?: string;
  manufacturer?: string;
  category?: string;
  unit_price: number;
  currency: string;
  stock_quantity: number;
  min_order_quantity: number;
  lead_time_days: number;
  status: string;
}

export default function CreateListingPage() {
  const router = useRouter();
  const user = getUser();
  
  const [formData, setFormData] = useState<ListingForm>({
    title: '',
    description: '',
    part_number: '',
    manufacturer: '',
    category: '',
    unit_price: 0,
    currency: 'USD',
    stock_quantity: 0,
    min_order_quantity: 1,
    lead_time_days: 0,
    status: 'draft',
  });
  
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    // Validation
    if (!formData.title || !formData.description) {
      setError('Title and description are required');
      setLoading(false);
      return;
    }

    if (formData.unit_price <= 0) {
      setError('Unit price must be greater than 0');
      setLoading(false);
      return;
    }

    try {
      // Get supplier/tenant ID
      const supplierId = user?.tenant_id;
      if (!supplierId) {
        throw new Error('Supplier ID not found');
      }

      const listingData = {
        ...formData,
        supplier_id: supplierId,
      };

      await apiClients.marketplace.post('/api/v1/listings', listingData);
      
      // Success - redirect to listings
      router.push('/app/supplier/listings');
    } catch (err: any) {
      setError(err.response?.data?.error || err.message || 'Failed to create listing');
      setLoading(false);
    }
  };

  return (
    <div>
      <div className="mb-6">
        <Link href="/app/supplier/listings" className="text-primary-600 hover:text-primary-700 mb-2 inline-block">
          ← Back to Listings
        </Link>
        <h1 className="text-3xl font-bold text-gray-900">Create New Listing</h1>
      </div>

      {error && (
        <div className="mb-4 bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">
          {error}
        </div>
      )}

      <Card title="Listing Information">
        <form onSubmit={handleSubmit} className="space-y-6">
          {/* Basic Information */}
          <div className="space-y-4">
            <Input
              label="Title *"
              value={formData.title}
              onChange={(e) => setFormData(prev => ({ ...prev, title: e.target.value }))}
              required
              placeholder="Product name or title"
            />
            
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Description *</label>
              <textarea
                value={formData.description}
                onChange={(e) => setFormData(prev => ({ ...prev, description: e.target.value }))}
                rows={4}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-primary-500 focus:border-primary-500"
                required
                placeholder="Detailed product description, specifications, etc."
              />
            </div>

            <div className="grid grid-cols-2 gap-4">
              <Input
                label="Part Number"
                value={formData.part_number}
                onChange={(e) => setFormData(prev => ({ ...prev, part_number: e.target.value }))}
                placeholder="SKU or part number"
              />
              
              <Input
                label="Manufacturer"
                value={formData.manufacturer}
                onChange={(e) => setFormData(prev => ({ ...prev, manufacturer: e.target.value }))}
                placeholder="Manufacturer name"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Category</label>
              <select
                value={formData.category}
                onChange={(e) => setFormData(prev => ({ ...prev, category: e.target.value }))}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-primary-500 focus:border-primary-500"
              >
                <option value="">Select category</option>
                <option value="electronics">Electronics</option>
                <option value="mechanical">Mechanical</option>
                <option value="raw-materials">Raw Materials</option>
                <option value="components">Components</option>
                <option value="tools">Tools</option>
                <option value="other">Other</option>
              </select>
            </div>
          </div>

          {/* Pricing */}
          <div className="border-t pt-6">
            <h3 className="text-lg font-semibold text-gray-900 mb-4">Pricing & Inventory</h3>
            <div className="grid grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Currency</label>
                <select
                  value={formData.currency}
                  onChange={(e) => setFormData(prev => ({ ...prev, currency: e.target.value }))}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-primary-500 focus:border-primary-500"
                  required
                >
                  <option value="USD">USD</option>
                  <option value="EUR">EUR</option>
                  <option value="GBP">GBP</option>
                </select>
              </div>
              
              <Input
                label="Unit Price *"
                type="number"
                step="0.01"
                min="0"
                value={formData.unit_price || ''}
                onChange={(e) => setFormData(prev => ({ ...prev, unit_price: parseFloat(e.target.value) || 0 }))}
                required
                placeholder="0.00"
              />
            </div>

            <div className="grid grid-cols-3 gap-4 mt-4">
              <Input
                label="Stock Quantity *"
                type="number"
                min="0"
                value={formData.stock_quantity || ''}
                onChange={(e) => setFormData(prev => ({ ...prev, stock_quantity: parseInt(e.target.value) || 0 }))}
                required
              />
              
              <Input
                label="Min Order Quantity *"
                type="number"
                min="1"
                value={formData.min_order_quantity || ''}
                onChange={(e) => setFormData(prev => ({ ...prev, min_order_quantity: parseInt(e.target.value) || 1 }))}
                required
              />
              
              <Input
                label="Lead Time (days) *"
                type="number"
                min="0"
                value={formData.lead_time_days || ''}
                onChange={(e) => setFormData(prev => ({ ...prev, lead_time_days: parseInt(e.target.value) || 0 }))}
                required
              />
            </div>
          </div>

          {/* Status */}
          <div className="border-t pt-6">
            <label className="block text-sm font-medium text-gray-700 mb-1">Status</label>
            <select
              value={formData.status}
              onChange={(e) => setFormData(prev => ({ ...prev, status: e.target.value }))}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-primary-500 focus:border-primary-500"
            >
              <option value="draft">Draft (not visible to buyers)</option>
              <option value="active">Active (visible to buyers)</option>
            </select>
          </div>

          {/* Submit Buttons */}
          <div className="flex gap-4 pt-4 border-t">
            <Button type="submit" disabled={loading}>
              {loading ? 'Creating...' : 'Create Listing'}
            </Button>
            <Button type="button" variant="secondary" onClick={() => router.back()}>
              Cancel
            </Button>
          </div>
        </form>
      </Card>
    </div>
  );
}
