import axios, { AxiosInstance, AxiosError } from 'axios';
import Cookies from 'js-cookie';
import { getTenantID } from './tenant';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8001';

// Service URLs
const SERVICE_URLS = {
  identity: process.env.NEXT_PUBLIC_IDENTITY_SERVICE_URL || 'http://localhost:8001',
  company: process.env.NEXT_PUBLIC_COMPANY_SERVICE_URL || 'http://localhost:8002',
  catalog: process.env.NEXT_PUBLIC_CATALOG_SERVICE_URL || 'http://localhost:8003',
  procurement: process.env.NEXT_PUBLIC_PROCUREMENT_SERVICE_URL || 'http://localhost:8006',
  logistics: process.env.NEXT_PUBLIC_LOGISTICS_SERVICE_URL || 'http://localhost:8007',
  collaboration: process.env.NEXT_PUBLIC_COLLABORATION_SERVICE_URL || 'http://localhost:8008',
  notification: process.env.NEXT_PUBLIC_NOTIFICATION_SERVICE_URL || 'http://localhost:8009',
  billing: process.env.NEXT_PUBLIC_BILLING_SERVICE_URL || 'http://localhost:8010',
  marketplace: process.env.NEXT_PUBLIC_MARKETPLACE_SERVICE_URL || 'http://localhost:8005',
};

function createApiClient(baseURL: string): AxiosInstance {
  const client = axios.create({
    baseURL,
    headers: {
      'Content-Type': 'application/json',
    },
  });

  // Add auth token to requests
  client.interceptors.request.use((config) => {
    const token = Cookies.get('auth_token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }

    // Add tenant ID header
    const tenantID = getTenantID();
    if (tenantID) {
      config.headers['X-Tenant-ID'] = tenantID;
    }

    return config;
  });

  // Handle errors
  client.interceptors.response.use(
    (response) => response,
    (error: AxiosError) => {
      if (error.response?.status === 401) {
        // Unauthorized - clear token and redirect to login
        Cookies.remove('auth_token');
        if (typeof window !== 'undefined') {
          window.location.href = '/login';
        }
      }
      return Promise.reject(error);
    }
  );

  return client;
}

export const apiClients = {
  identity: createApiClient(SERVICE_URLS.identity),
  company: createApiClient(SERVICE_URLS.company),
  catalog: createApiClient(SERVICE_URLS.catalog),
  procurement: createApiClient(SERVICE_URLS.procurement),
  logistics: createApiClient(SERVICE_URLS.logistics),
  collaboration: createApiClient(SERVICE_URLS.collaboration),
  notification: createApiClient(SERVICE_URLS.notification),
  billing: createApiClient(SERVICE_URLS.billing),
  marketplace: createApiClient(SERVICE_URLS.marketplace),
};

export default apiClients;
