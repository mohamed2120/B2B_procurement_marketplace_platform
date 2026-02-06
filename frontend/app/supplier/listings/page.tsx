'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { isAuthenticated, hasRole } from '@/lib/auth';
import { apiClients } from '@/lib/api';
import Header from '@/components/layout/Header';
import Sidebar from '@/components/layout/Sidebar';
import Card from '@/components/ui/Card';
import Button from '@/components/ui/Button';
import Input from '@/components/ui/Input';
import Link from 'next/link';

interface Listing {
  id: string;
  title: string;
  description: string;
  price: number;
  stock_quantity: number;
  status: string;
}

export default function ListingsPage() {
  const router = useRouter();
  const [listings, setListings] = useState<Listing[]>([]);
  const [loading, setLoading] = useState(true);
  const [showCreate, setShowCreate] = useState(false);
  const [formData, setFormData] = useState({
    title: '',
    description: '',
    price: 0,
    stock_quantity: 0,
  });

  useEffect(() => {
    if (!isAuthenticated() || !hasRole('supplier')) {
      router.push('/login');
      return;
    }

    fetchListings();
  }, [router]);

  const fetchListings = async () => {
    try {
      const response = await apiClients.procurement.get('/api/v1/listings?limit=100');
      setListings(response.data.items || response.data || []);
    } catch (error) {
      console.error('Failed to fetch listings:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await apiClients.procurement.post('/api/v1/listings', formData);
      setShowCreate(false);
      setFormData({ title: '', description: '', price: 0, stock_quantity: 0 });
      fetchListings();
    } catch (error) {
      console.error('Failed to create listing:', error);
    }
  };

  const handleDelete = async (id: string) => {
    if (!confirm('Are you sure you want to delete this listing?')) return;
    try {
      await apiClients.procurement.delete(`/api/v1/listings/${id}`);
      fetchListings();
    } catch (error) {
      console.error('Failed to delete listing:', error);
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
            <h1 className="text-3xl font-bold text-gray-900">My Listings</h1>
            <Button onClick={() => setShowCreate(!showCreate)}>
              {showCreate ? 'Cancel' : '+ Create Listing'}
            </Button>
          </div>

          {showCreate && (
            <Card className="mb-6">
              <h2 className="text-xl font-semibold mb-4">Create New Listing</h2>
              <form onSubmit={handleCreate} className="space-y-4">
                <Input
                  label="Title"
                  value={formData.title}
                  onChange={(e) => setFormData({ ...formData, title: e.target.value })}
                  required
                />
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Description</label>
                  <textarea
                    value={formData.description}
                    onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                    required
                    rows={3}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg"
                  />
                </div>
                <div className="grid grid-cols-2 gap-4">
                  <Input
                    label="Price"
                    type="number"
                    step="0.01"
                    value={formData.price}
                    onChange={(e) => setFormData({ ...formData, price: parseFloat(e.target.value) })}
                    required
                  />
                  <Input
                    label="Stock Quantity"
                    type="number"
                    value={formData.stock_quantity}
                    onChange={(e) => setFormData({ ...formData, stock_quantity: parseInt(e.target.value) })}
                    required
                  />
                </div>
                <Button type="submit">Create Listing</Button>
              </form>
            </Card>
          )}

          <Card>
            {listings.length === 0 ? (
              <p className="text-gray-500 text-center py-8">No listings found. Create your first listing!</p>
            ) : (
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                {listings.map((listing) => (
                  <div key={listing.id} className="border rounded-lg p-4">
                    <h3 className="font-semibold text-lg mb-2">{listing.title}</h3>
                    <p className="text-sm text-gray-600 mb-4">{listing.description}</p>
                    <div className="flex justify-between items-center mb-4">
                      <span className="text-xl font-bold text-primary-600">${listing.price}</span>
                      <span className="text-sm text-gray-500">Stock: {listing.stock_quantity}</span>
                    </div>
                    <div className="flex space-x-2">
                      <Button size="sm" variant="outline" className="flex-1">
                        Edit
                      </Button>
                      <Button
                        size="sm"
                        variant="danger"
                        onClick={() => handleDelete(listing.id)}
                      >
                        Delete
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
