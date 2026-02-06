'use client';

import PublicLayout from '@/components/layout/PublicLayout';
import Link from 'next/link';

export default function Pricing() {
  const plans = [
    {
      name: 'Starter',
      price: '$99',
      period: 'per month',
      description: 'Perfect for small businesses',
      features: [
        'Up to 10 users',
        '50 RFQs per month',
        '100 listings',
        '10GB storage',
        'Email support',
      ],
      cta: 'Start Free Trial',
    },
    {
      name: 'Professional',
      price: '$299',
      period: 'per month',
      description: 'For growing companies',
      features: [
        'Up to 50 users',
        'Unlimited RFQs',
        '500 listings',
        '100GB storage',
        'Priority support',
        'Advanced analytics',
      ],
      cta: 'Start Free Trial',
      popular: true,
    },
    {
      name: 'Enterprise',
      price: 'Custom',
      period: '',
      description: 'For large organizations',
      features: [
        'Unlimited users',
        'Unlimited RFQs',
        'Unlimited listings',
        'Unlimited storage',
        'Dedicated support',
        'Custom integrations',
        'SLA guarantee',
      ],
      cta: 'Contact Sales',
    },
  ];

  return (
    <PublicLayout>
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        <div className="text-center mb-12">
          <h1 className="text-4xl font-bold text-gray-900 mb-4">Pricing Plans</h1>
          <p className="text-xl text-gray-600">
            Choose the plan that fits your business needs
          </p>
        </div>

        <div className="grid md:grid-cols-3 gap-8">
          {plans.map((plan) => (
            <div
              key={plan.name}
              className={`bg-white rounded-lg shadow-lg p-8 ${
                plan.popular ? 'ring-2 ring-primary-600 transform scale-105' : ''
              }`}
            >
              {plan.popular && (
                <div className="bg-primary-600 text-white text-sm font-semibold px-3 py-1 rounded-full inline-block mb-4">
                  Most Popular
                </div>
              )}
              <h3 className="text-2xl font-bold text-gray-900 mb-2">{plan.name}</h3>
              <p className="text-gray-600 mb-4">{plan.description}</p>
              <div className="mb-6">
                <span className="text-4xl font-bold text-gray-900">{plan.price}</span>
                {plan.period && <span className="text-gray-600 ml-2">{plan.period}</span>}
              </div>
              <ul className="space-y-3 mb-8">
                {plan.features.map((feature, idx) => (
                  <li key={idx} className="flex items-start">
                    <svg className="w-5 h-5 text-primary-600 mr-2 mt-0.5" fill="currentColor" viewBox="0 0 20 20">
                      <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
                    </svg>
                    <span className="text-gray-700">{feature}</span>
                  </li>
                ))}
              </ul>
              <Link
                href={plan.name === 'Enterprise' ? '/contact' : '/register'}
                className={`block w-full text-center py-3 rounded-lg font-semibold transition ${
                  plan.popular
                    ? 'bg-primary-600 text-white hover:bg-primary-700'
                    : 'bg-gray-100 text-gray-900 hover:bg-gray-200'
                }`}
              >
                {plan.cta}
              </Link>
            </div>
          ))}
        </div>

        <div className="mt-12 text-center">
          <p className="text-gray-600 mb-4">All plans include a 14-day free trial. No credit card required.</p>
          <Link
            href="/contact"
            className="text-primary-600 hover:text-primary-700 font-semibold"
          >
            Contact us for custom enterprise pricing
          </Link>
        </div>
      </div>
    </PublicLayout>
  );
}
