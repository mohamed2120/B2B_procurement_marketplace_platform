/**
 * Tenant resolution via subdomain
 * Reads from window.location.host to extract tenant subdomain
 */

export function getTenantFromSubdomain(): string | null {
  if (typeof window === 'undefined') {
    return null;
  }

  const host = window.location.host;
  const parts = host.split('.');

  // If we have more than 2 parts (e.g., tenant.localhost:3000 or tenant.example.com)
  // The first part is the subdomain
  if (parts.length > 2) {
    return parts[0];
  }

  // For localhost development, check if there's a subdomain pattern
  // e.g., tenant.localhost:3000
  if (host.includes('localhost') && parts.length > 1) {
    const subdomain = parts[0];
    if (subdomain !== 'localhost' && subdomain !== 'www') {
      return subdomain;
    }
  }

  // Default tenant for local development
  // In production, this would be null and require subdomain
  return process.env.NEXT_PUBLIC_DEFAULT_TENANT || 'demo';
}

export function getTenantID(): string {
  // For now, use a default tenant ID
  // In production, this would be resolved from the subdomain via API
  return '00000000-0000-0000-0000-000000000001';
}
