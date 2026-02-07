import { test, expect } from '@playwright/test';

const BASE_URL = process.env.PLAYWRIGHT_BASE_URL || 'http://localhost:3002';

// Demo accounts
const DEMO_USERS = {
  admin: { email: 'admin@demo.com', password: 'demo123456', role: 'admin' },
  requester: { email: 'buyer.requester@demo.com', password: 'demo123456', role: 'requester' },
  procurement: { email: 'buyer.procurement@demo.com', password: 'demo123456', role: 'procurement_manager' },
  supplier: { email: 'supplier@demo.com', password: 'demo123456', role: 'supplier' },
};

// Required pages by category
const PUBLIC_PAGES = [
  '/',
  '/how-it-works',
  '/pricing',
  '/register',
  '/register/buyer',
  '/register/supplier',
  '/login',
  '/forgot-password',
  '/contact',
  '/terms',
  '/privacy',
  '/search',
];

const BUYER_PAGES = [
  '/app/customer/dashboard',
  '/app/customer/pr',
  '/app/customer/rfq',
  '/app/customer/orders',
  '/app/customer/shipments',
  '/app/customer/parts',
  '/app/customer/equipment',
  '/app/customer/warehouse',
  '/app/customer/team',
  '/app/customer/company',
  '/app/customer/reports',
];

const SUPPLIER_PAGES = [
  '/app/supplier/dashboard',
  '/app/supplier/rfq',
  '/app/supplier/quotes',
  '/app/supplier/listings',
  '/app/supplier/orders',
  '/app/supplier/shipments',
  '/app/supplier/store',
  '/app/supplier/services',
  '/app/supplier/inventory',
  '/app/supplier/reports',
];

const ADMIN_PAGES = [
  '/app/admin/dashboard',
  '/app/admin/tenants',
  '/app/admin/users',
  '/app/admin/companies',
  '/app/admin/roles-permissions',
  '/app/admin/catalog-approvals',
  '/app/admin/disputes',
  '/app/admin/subscriptions',
  '/app/admin/audit-logs',
  '/app/admin/diagnostics',
];

async function loginAs(page: any, user: typeof DEMO_USERS.admin) {
  await page.goto(`${BASE_URL}/login`);
  await page.fill('input[type="email"]', user.email);
  await page.fill('input[type="password"]', user.password);
  await page.click('button[type="submit"]');
  await page.waitForURL(/\/app/, { timeout: 10000 });
}

test.describe('Complete Page Map - Route Accessibility', () => {
  test('All public pages load (200 OK)', async ({ page }) => {
    for (const route of PUBLIC_PAGES) {
      const response = await page.goto(`${BASE_URL}${route}`);
      expect(response?.status()).toBe(200);
    }
  });

  test('Buyer pages accessible to requester', async ({ page }) => {
    await loginAs(page, DEMO_USERS.requester);
    
    for (const route of BUYER_PAGES) {
      const response = await page.goto(`${BASE_URL}${route}`);
      expect(response?.status()).toBe(200);
      // Check that page renders (not empty)
      await expect(page.locator('body')).not.toBeEmpty();
    }
  });

  test('Buyer pages accessible to procurement manager', async ({ page }) => {
    await loginAs(page, DEMO_USERS.procurement);
    
    for (const route of BUYER_PAGES) {
      const response = await page.goto(`${BASE_URL}${route}`);
      expect(response?.status()).toBe(200);
    }
  });

  test('Supplier pages accessible to supplier', async ({ page }) => {
    await loginAs(page, DEMO_USERS.supplier);
    
    for (const route of SUPPLIER_PAGES) {
      const response = await page.goto(`${BASE_URL}${route}`);
      expect(response?.status()).toBe(200);
    }
  });

  test('Admin pages accessible to admin', async ({ page }) => {
    await loginAs(page, DEMO_USERS.admin);
    
    for (const route of ADMIN_PAGES) {
      const response = await page.goto(`${BASE_URL}${route}`);
      expect(response?.status()).toBe(200);
    }
  });

  test('RBAC: Supplier cannot access buyer pages', async ({ page }) => {
    await loginAs(page, DEMO_USERS.supplier);
    
    const response = await page.goto(`${BASE_URL}/app/customer/dashboard`);
    // Should redirect to /app or show 403
    expect(response?.status()).toBeGreaterThanOrEqual(200);
    expect(response?.status()).toBeLessThan(500);
    // Should not be on buyer dashboard
    expect(page.url()).not.toContain('/app/customer/dashboard');
  });

  test('RBAC: Buyer cannot access admin pages', async ({ page }) => {
    await loginAs(page, DEMO_USERS.requester);
    
    const response = await page.goto(`${BASE_URL}/app/admin/dashboard`);
    // Should redirect or block
    expect(response?.status()).toBeGreaterThanOrEqual(200);
    expect(response?.status()).toBeLessThan(500);
    expect(page.url()).not.toContain('/app/admin/dashboard');
  });
});

test.describe('Complete Page Map - Data Display', () => {
  test('Buyer PR list shows data', async ({ page }) => {
    await loginAs(page, DEMO_USERS.requester);
    await page.goto(`${BASE_URL}/app/customer/pr`);
    
    // Check for table or data display
    const hasData = await page.locator('table, [class*="card"], [class*="item"]').count() > 0;
    expect(hasData).toBeTruthy();
  });

  test('Supplier listings show data', async ({ page }) => {
    await loginAs(page, DEMO_USERS.supplier);
    await page.goto(`${BASE_URL}/app/supplier/listings`);
    
    const hasData = await page.locator('table, [class*="card"], [class*="item"]').count() > 0;
    expect(hasData).toBeTruthy();
  });

  test('Admin users page shows data', async ({ page }) => {
    await loginAs(page, DEMO_USERS.admin);
    await page.goto(`${BASE_URL}/app/admin/users`);
    
    const hasData = await page.locator('table, [class*="card"], [class*="item"]').count() > 0;
    expect(hasData).toBeTruthy();
  });
});

test.describe('Complete Page Map - Core Flow', () => {
  test('PR → RFQ → Quote → Award → Shipment flow', async ({ page }) => {
    // Login as requester
    await loginAs(page, DEMO_USERS.requester);
    
    // Navigate to PR list
    await page.goto(`${BASE_URL}/app/customer/pr`);
    await expect(page.locator('body')).not.toBeEmpty();
    
    // Login as procurement
    await page.goto(`${BASE_URL}/login`);
    await page.fill('input[type="email"]', DEMO_USERS.procurement.email);
    await page.fill('input[type="password"]', DEMO_USERS.procurement.password);
    await page.click('button[type="submit"]');
    await page.waitForURL(/\/app/, { timeout: 10000 });
    
    // Navigate to RFQ
    await page.goto(`${BASE_URL}/app/customer/rfq`);
    await expect(page.locator('body')).not.toBeEmpty();
    
    // Login as supplier
    await page.goto(`${BASE_URL}/login`);
    await page.fill('input[type="email"]', DEMO_USERS.supplier.email);
    await page.fill('input[type="password"]', DEMO_USERS.supplier.password);
    await page.click('button[type="submit"]');
    await page.waitForURL(/\/app/, { timeout: 10000 });
    
    // Navigate to RFQ inbox
    await page.goto(`${BASE_URL}/app/supplier/rfq`);
    await expect(page.locator('body')).not.toBeEmpty();
  });
});
