'use client';

import { useState } from 'react';
import { useSafeRouter } from '@/lib/useSafeRouter';
import { isAuthenticated, hasRole } from '@/lib/auth';
import { apiClients } from '@/lib/api';
import Header from '@/components/layout/Header';
import Sidebar from '@/components/layout/Sidebar';
import Card from '@/components/ui/Card';
import Input from '@/components/ui/Input';
import Button from '@/components/ui/Button';

export default function CreatePRPage() {
  const router = useSafeRouter();
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [items, setItems] = useState([{ description: '', quantity: 1, unit_price: 0 }]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      await apiClients.procurement.post('/api/v1/purchase-requests', {
        title,
        description,
        status: 'draft',
        items: items.map(item => ({
          description: item.description,
          quantity: item.quantity,
          unit_price: item.unit_price,
        })),
      });

      router.push('/customer/prs');
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to create PR');
    } finally {
      setLoading(false);
    }
  };

  const addItem = () => {
    setItems([...items, { description: '', quantity: 1, unit_price: 0 }]);
  };

  const removeItem = (index: number) => {
    setItems(items.filter((_, i) => i !== index));
  };

  const updateItem = (index: number, field: string, value: any) => {
    const newItems = [...items];
    newItems[index] = { ...newItems[index], [field]: value };
    setItems(newItems);
  };

  return (
    <div className="min-h-screen bg-gray-50">
      <Header />
      <div className="flex">
        <Sidebar />
        <main className="flex-1 p-8">
          <h1 className="text-3xl font-bold text-gray-900 mb-6">Create Purchase Request</h1>

          <Card>
            <form onSubmit={handleSubmit} className="space-y-6">
              <Input
                label="Title"
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                required
                placeholder="Office Supplies Q1 2024"
              />

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Description</label>
                <textarea
                  value={description}
                  onChange={(e) => setDescription(e.target.value)}
                  required
                  rows={4}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500"
                  placeholder="Describe the purchase request..."
                />
              </div>

              <div>
                <div className="flex justify-between items-center mb-4">
                  <label className="block text-sm font-medium text-gray-700">Items</label>
                  <Button type="button" variant="outline" size="sm" onClick={addItem}>
                    + Add Item
                  </Button>
                </div>

                <div className="space-y-4">
                  {items.map((item, index) => (
                    <div key={index} className="grid grid-cols-12 gap-4 p-4 border rounded-lg">
                      <div className="col-span-5">
                        <Input
                          value={item.description}
                          onChange={(e) => updateItem(index, 'description', e.target.value)}
                          placeholder="Item description"
                          required
                        />
                      </div>
                      <div className="col-span-2">
                        <Input
                          type="number"
                          value={item.quantity}
                          onChange={(e) => updateItem(index, 'quantity', parseInt(e.target.value))}
                          placeholder="Quantity"
                          required
                          min="1"
                        />
                      </div>
                      <div className="col-span-2">
                        <Input
                          type="number"
                          step="0.01"
                          value={item.unit_price}
                          onChange={(e) => updateItem(index, 'unit_price', parseFloat(e.target.value))}
                          placeholder="Unit Price"
                          required
                          min="0"
                        />
                      </div>
                      <div className="col-span-2 flex items-end">
                        <span className="text-gray-700">
                          ${(item.quantity * item.unit_price).toFixed(2)}
                        </span>
                      </div>
                      <div className="col-span-1 flex items-end">
                        {items.length > 1 && (
                          <button
                            type="button"
                            onClick={() => removeItem(index)}
                            className="text-red-600 hover:text-red-800"
                          >
                            ✕
                          </button>
                        )}
                      </div>
                    </div>
                  ))}
                </div>
              </div>

              {error && (
                <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">
                  {error}
                </div>
              )}

              <div className="flex space-x-4">
                <Button type="submit" disabled={loading}>
                  {loading ? 'Creating...' : 'Create PR'}
                </Button>
                <Button
                  type="button"
                  variant="secondary"
                  onClick={() => router.back()}
                >
                  Cancel
                </Button>
              </div>
            </form>
          </Card>
        </main>
      </div>
    </div>
  );
}
