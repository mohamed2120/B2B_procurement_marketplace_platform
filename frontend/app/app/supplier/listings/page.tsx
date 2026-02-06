'use client';

import { useState, useEffect } from 'react';
import { apiClients } from '@/lib/api';
import Card from '@/components/ui/Card';
import Button from '@/components/ui/Button';

export default function SupplierListings() {
  const [listings, setListings] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchListings();
  }, []);

  const fetchListings = async () => {
    try {
      const response = await apiClients.marketplace.get('/api/v1/listings');
      setListings(response.data || []);
    } catch (error) {
      console.error('Failed to fetch listings:', error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-3xl font-bold text-gray-900">My Listings</h1>
        <Button>Create Listing</Button>
      </div>

      {loading ? (
        <div>Loading...</div>
      ) : listings.length === 0 ? (
        <Card>
          <p className="text-gray-600 mb-4">No listings yet.</p>
          <Button>Create Your First Listing</Button>
        </Card>
      ) : (
        <div className="space-y-4">
          {listings.map((listing) => (
            <Card key={listing.id}>
              <div className="flex justify-between items-start">
                <div>
                  <h3 className="font-semibold text-lg">{listing.title || listing.name}</h3>
                  <p className="text-gray-600 text-sm mt-1">Status: {listing.status}</p>
                </div>
                <button className="text-primary-600 hover:text-primary-700">Edit →</button>
              </div>
            </Card>
          ))}
        </div>
      )}
    </div>
  );
}
