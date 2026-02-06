import Cookies from 'js-cookie';
import { apiClients } from './api';
import { getTenantID } from './tenant';

export interface User {
  id: string;
  email: string;
  first_name: string;
  last_name: string;
  tenant_id: string;
  roles?: string[];
}

export interface LoginResponse {
  token: string;
  user: User;
  expires_at: string;
}

export interface LoginRequest {
  email: string;
  password: string;
  tenant_id: string;
}

const TOKEN_KEY = 'auth_token';
const USER_KEY = 'auth_user';

export async function login(email: string, password: string): Promise<LoginResponse> {
  const tenantID = getTenantID();
  
  const response = await apiClients.identity.post<LoginResponse>('/api/v1/auth/login', {
    email,
    password,
    tenant_id: tenantID,
  });

  // Store token and user
  Cookies.set(TOKEN_KEY, response.data.token, { expires: 7 }); // 7 days
  if (typeof window !== 'undefined') {
    localStorage.setItem(USER_KEY, JSON.stringify(response.data.user));
  }

  return response.data;
}

export function logout(): void {
  Cookies.remove(TOKEN_KEY);
  if (typeof window !== 'undefined') {
    localStorage.removeItem(USER_KEY);
    window.location.href = '/';
  }
}

export function getToken(): string | null {
  return Cookies.get(TOKEN_KEY) || null;
}

export function getUser(): User | null {
  if (typeof window === 'undefined') {
    return null;
  }

  const userStr = localStorage.getItem(USER_KEY);
  if (!userStr) {
    return null;
  }

  try {
    return JSON.parse(userStr) as User;
  } catch {
    return null;
  }
}

export function isAuthenticated(): boolean {
  return !!getToken();
}

export function hasRole(role: string): boolean {
  const user = getUser();
  if (!user || !user.roles) {
    return false;
  }
  return user.roles.includes(role) || user.roles.includes('admin') || user.roles.includes('super_admin');
}

export function hasAnyRole(roles: string[]): boolean {
  return roles.some(role => hasRole(role));
}

export async function validateToken(): Promise<boolean> {
  const token = getToken();
  if (!token) {
    return false;
  }

  try {
    await apiClients.identity.get('/api/v1/auth/validate');
    return true;
  } catch {
    logout();
    return false;
  }
}
