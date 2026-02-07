'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { apiClients } from '@/lib/api';
import { hasRole } from '@/lib/auth';
import Card from '@/components/ui/Card';
import Button from '@/components/ui/Button';
import Input from '@/components/ui/Input';

interface Plan {
  id: string;
  name: string;
  code: string;
  description: string;
  price: number;
  currency: string;
  billing_cycle: string;
  is_active: boolean;
  entitlements?: Entitlement[];
}

interface Entitlement {
  id: string;
  feature: string;
  limit: number;
  unit: string;
  description: string;
}

export default function PlansPage() {
  const router = useRouter();
  const [plans, setPlans] = useState<Plan[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [editingPlan, setEditingPlan] = useState<Plan | null>(null);
  const [showFeaturesModal, setShowFeaturesModal] = useState(false);
  const [selectedPlan, setSelectedPlan] = useState<Plan | null>(null);

  useEffect(() => {
    if (!(hasRole('admin') || hasRole('super_admin'))) {
      router.push('/app');
      return;
    }
    fetchPlans();
  }, []);

  const fetchPlans = async () => {
    setLoading(true);
    setError('');
    try {
      const response = await apiClients.billing.get('/api/v1/plans');
      setPlans(response.data || []);
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to load plans');
      // Fallback to mock data
      setPlans([
        {
          id: '1',
          name: 'Basic',
          code: 'basic',
          description: 'Basic plan for small businesses',
          price: 99,
          currency: 'USD',
          billing_cycle: 'monthly',
          is_active: true,
          entitlements: [
            { id: '1', feature: 'users', limit: 10, unit: 'users', description: 'Maximum 10 users' },
            { id: '2', feature: 'orders', limit: 100, unit: 'orders/month', description: '100 orders per month' },
          ],
        },
        {
          id: '2',
          name: 'Professional',
          code: 'professional',
          description: 'Professional plan for growing businesses',
          price: 299,
          currency: 'USD',
          billing_cycle: 'monthly',
          is_active: true,
          entitlements: [
            { id: '3', feature: 'users', limit: 50, unit: 'users', description: 'Maximum 50 users' },
            { id: '4', feature: 'orders', limit: 500, unit: 'orders/month', description: '500 orders per month' },
            { id: '5', feature: 'storage', limit: 100, unit: 'GB', description: '100 GB storage' },
          ],
        },
        {
          id: '3',
          name: 'Enterprise',
          code: 'enterprise',
          description: 'Enterprise plan with unlimited features',
          price: 999,
          currency: 'USD',
          billing_cycle: 'monthly',
          is_active: true,
          entitlements: [
            { id: '6', feature: 'users', limit: -1, unit: 'users', description: 'Unlimited users' },
            { id: '7', feature: 'orders', limit: -1, unit: 'orders/month', description: 'Unlimited orders' },
            { id: '8', feature: 'storage', limit: -1, unit: 'GB', description: 'Unlimited storage' },
          ],
        },
      ]);
    } finally {
      setLoading(false);
    }
  };

  const handleCreatePlan = async (planData: Partial<Plan>) => {
    try {
      const response = await apiClients.billing.post('/api/v1/plans', planData);
      setPlans([...plans, response.data]);
      setShowCreateModal(false);
    } catch (err: any) {
      alert('Failed to create plan: ' + (err.response?.data?.error || err.message));
    }
  };

  const handleManageFeatures = (plan: Plan) => {
    setSelectedPlan(plan);
    setShowFeaturesModal(true);
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600 mx-auto mb-4"></div>
          <p className="text-gray-600">Loading...</p>
        </div>
      </div>
    );
  }

  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-3xl font-bold text-gray-900">Pricing Plans & Features</h1>
        <Button onClick={() => setShowCreateModal(true)}>Create New Plan</Button>
      </div>

      {error && (
        <div className="mb-4 bg-yellow-50 border border-yellow-200 text-yellow-700 px-4 py-3 rounded">
          <p className="font-semibold">Note</p>
          <p className="text-sm mt-1">Using mock data. API integration pending.</p>
        </div>
      )}

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {plans.map((plan) => (
          <Card key={plan.id} className="relative">
            <div className="flex justify-between items-start mb-4">
              <div>
                <h3 className="text-xl font-bold text-gray-900">{plan.name}</h3>
                <p className="text-sm text-gray-600 mt-1">{plan.code}</p>
              </div>
              <span className={`px-2 py-1 text-xs rounded ${plan.is_active ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-800'}`}>
                {plan.is_active ? 'Active' : 'Inactive'}
              </span>
            </div>

            <p className="text-gray-600 text-sm mb-4">{plan.description}</p>

            <div className="mb-4">
              <div className="text-3xl font-bold text-primary-600">
                ${plan.price}
                <span className="text-lg text-gray-600">/{plan.billing_cycle}</span>
              </div>
              <p className="text-sm text-gray-500">{plan.currency}</p>
            </div>

            <div className="mb-4">
              <h4 className="font-semibold text-sm text-gray-700 mb-2">Features:</h4>
              <ul className="space-y-1">
                {plan.entitlements?.slice(0, 3).map((ent) => (
                  <li key={ent.id} className="text-sm text-gray-600">
                    • {ent.feature}: {ent.limit === -1 ? 'Unlimited' : `${ent.limit} ${ent.unit}`}
                  </li>
                ))}
                {plan.entitlements && plan.entitlements.length > 3 && (
                  <li className="text-sm text-primary-600">+ {plan.entitlements.length - 3} more...</li>
                )}
              </ul>
            </div>

            <div className="flex space-x-2 mt-4">
              <Button size="sm" variant="secondary" onClick={() => handleManageFeatures(plan)}>
                Manage Features
              </Button>
              <Button size="sm" variant="secondary" onClick={() => setEditingPlan(plan)}>
                Edit
              </Button>
            </div>
          </Card>
        ))}
      </div>

      {/* Create Plan Modal */}
      {showCreateModal && (
        <CreatePlanModal
          onClose={() => setShowCreateModal(false)}
          onCreate={handleCreatePlan}
        />
      )}

      {/* Manage Features Modal */}
      {showFeaturesModal && selectedPlan && (
        <ManageFeaturesModal
          plan={selectedPlan}
          onClose={() => {
            setShowFeaturesModal(false);
            setSelectedPlan(null);
          }}
          onUpdate={fetchPlans}
        />
      )}
    </div>
  );
}

function CreatePlanModal({ onClose, onCreate }: { onClose: () => void; onCreate: (plan: Partial<Plan>) => void }) {
  const [formData, setFormData] = useState({
    name: '',
    code: '',
    description: '',
    price: '',
    currency: 'USD',
    billing_cycle: 'monthly',
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onCreate({
      ...formData,
      price: parseFloat(formData.price),
      is_active: true,
    });
  };

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <Card className="w-full max-w-md p-6">
        <h2 className="text-2xl font-bold mb-4">Create New Plan</h2>
        <form onSubmit={handleSubmit} className="space-y-4">
          <Input
            label="Plan Name"
            value={formData.name}
            onChange={(e) => setFormData({ ...formData, name: e.target.value })}
            required
          />
          <Input
            label="Code"
            value={formData.code}
            onChange={(e) => setFormData({ ...formData, code: e.target.value })}
            required
          />
          <Input
            label="Description"
            value={formData.description}
            onChange={(e) => setFormData({ ...formData, description: e.target.value })}
            required
          />
          <Input
            label="Price"
            type="number"
            step="0.01"
            value={formData.price}
            onChange={(e) => setFormData({ ...formData, price: e.target.value })}
            required
          />
          <div className="flex space-x-4">
            <Button type="submit">Create</Button>
            <Button type="button" variant="secondary" onClick={onClose}>
              Cancel
            </Button>
          </div>
        </form>
      </Card>
    </div>
  );
}

function ManageFeaturesModal({ plan, onClose, onUpdate }: { plan: Plan; onClose: () => void; onUpdate: () => void }) {
  const [features, setFeatures] = useState<Entitlement[]>(plan.entitlements || []);
  const [newFeature, setNewFeature] = useState({ feature: '', limit: '', unit: '', description: '' });

  const handleAddFeature = () => {
    if (!newFeature.feature) return;
    setFeatures([
      ...features,
      {
        id: Date.now().toString(),
        feature: newFeature.feature,
        limit: parseInt(newFeature.limit) || -1,
        unit: newFeature.unit,
        description: newFeature.description,
      },
    ]);
    setNewFeature({ feature: '', limit: '', unit: '', description: '' });
  };

  const handleRemoveFeature = (id: string) => {
    setFeatures(features.filter((f) => f.id !== id));
  };

  const handleSave = async () => {
    // TODO: Call API to update plan entitlements
    alert('Feature management API integration pending');
    onUpdate();
    onClose();
  };

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <Card className="w-full max-w-2xl p-6 max-h-[90vh] overflow-y-auto">
        <h2 className="text-2xl font-bold mb-4">Manage Features: {plan.name}</h2>

        <div className="space-y-4 mb-6">
          {features.map((feature) => (
            <div key={feature.id} className="flex items-center justify-between p-3 bg-gray-50 rounded">
              <div>
                <div className="font-semibold">{feature.feature}</div>
                <div className="text-sm text-gray-600">
                  Limit: {feature.limit === -1 ? 'Unlimited' : `${feature.limit} ${feature.unit}`}
                </div>
                {feature.description && (
                  <div className="text-xs text-gray-500 mt-1">{feature.description}</div>
                )}
              </div>
              <Button size="sm" variant="secondary" onClick={() => handleRemoveFeature(feature.id)}>
                Remove
              </Button>
            </div>
          ))}
        </div>

        <div className="border-t pt-4">
          <h3 className="font-semibold mb-3">Add New Feature</h3>
          <div className="grid grid-cols-2 gap-3">
            <Input
              placeholder="Feature name (e.g., users, orders)"
              value={newFeature.feature}
              onChange={(e) => setNewFeature({ ...newFeature, feature: e.target.value })}
            />
            <Input
              placeholder="Limit (-1 for unlimited)"
              type="number"
              value={newFeature.limit}
              onChange={(e) => setNewFeature({ ...newFeature, limit: e.target.value })}
            />
            <Input
              placeholder="Unit (e.g., users, GB)"
              value={newFeature.unit}
              onChange={(e) => setNewFeature({ ...newFeature, unit: e.target.value })}
            />
            <Input
              placeholder="Description"
              value={newFeature.description}
              onChange={(e) => setNewFeature({ ...newFeature, description: e.target.value })}
            />
          </div>
          <Button onClick={handleAddFeature} className="mt-3" size="sm">
            Add Feature
          </Button>
        </div>

        <div className="flex space-x-4 mt-6">
          <Button onClick={handleSave}>Save Changes</Button>
          <Button variant="secondary" onClick={onClose}>
            Cancel
          </Button>
        </div>
      </Card>
    </div>
  );
}
