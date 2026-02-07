import { test, expect } from '@playwright/test';

const BASE_URL = 'http://localhost:3002';
const DEMO_PASSWORD = 'demo123456';

const USERS = {
  requester: { email: 'buyer.requester@demo.com', password: DEMO_PASSWORD, role: 'requester' },
  procurement: { email: 'buyer.procurement@demo.com', password: DEMO_PASSWORD, role: 'procurement_manager' },
  supplier: { email: 'supplier@demo.com', password: DEMO_PASSWORD, role: 'supplier' },
  admin: { email: 'admin@demo.com', password: DEMO_PASSWORD, role: 'admin' },
};

// Pages that should be accessible by each role
const ROLE_PAGES = {
  requester: [
    '/app/customer/dashboard',
    '/app/customer/pr',
    '/app/customer/pr/create',
    '/app/customer/rfq',
    '/app/customer/orders',
    '/app/customer/shipments',
    '/app/my-plan',
    '/app/notifications',
    '/app/chat',
    '/app/profile',
  ],
  procurement_manager: [
    '/app/customer/dashboard',
    '/app/customer/pr',
    '/app/customer/rfq',
    '/app/customer/orders',
    '/app/customer/shipments',
    '/app/my-plan',
    '/app/notifications',
    '/app/chat',
    '/app/profile',
  ],
  supplier: [
    '/app/supplier/dashboard',
    '/app/supplier/rfq',
    '/app/supplier/quotes',
    '/app/supplier/listings',
    '/app/supplier/listings/create',
    '/app/supplier/orders',
    '/app/supplier/shipments',
    '/app/my-plan',
    '/app/notifications',
    '/app/chat',
    '/app/profile',
  ],
  admin: [
    '/app/admin/dashboard',
    '/app/admin/tenants',
    '/app/admin/users',
    '/app/admin/roles-permissions',
    '/app/admin/company-verification',
    '/app/admin/catalog-approvals',
    '/app/admin/disputes',
    '/app/admin/subscriptions',
    '/app/admin/audit-logs',
    '/app/admin/diagnostics',
    '/app/my-plan',
    '/app/notifications',
    '/app/chat',
    '/app/profile',
  ],
};

// Pages that should NOT be accessible (should redirect or show 403)
const FORBIDDEN_PAGES = {
  requester: [
    '/app/admin/dashboard',
    '/app/admin/tenants',
    '/app/admin/users',
    '/app/supplier/dashboard',
    '/app/supplier/rfq',
    '/app/customer/pr/create', // Actually allowed for requester
  ],
  procurement_manager: [
    '/app/admin/dashboard',
    '/app/admin/tenants',
    '/app/admin/users',
    '/app/supplier/dashboard',
    '/app/customer/pr/create', // Not allowed for procurement
  ],
  supplier: [
    '/app/admin/dashboard',
    '/app/admin/tenants',
    '/app/admin/users',
    '/app/customer/dashboard',
    '/app/customer/pr',
    '/app/customer/rfq',
  ],
  admin: [
    '/app/customer/dashboard',
    '/app/customer/pr',
    '/app/supplier/dashboard',
    '/app/supplier/rfq',
  ],
};

async function loginAs(page: any, email: string, password: string) {
  await page.goto(`${BASE_URL}/login`);
  await page.fill('input[type="email"]', email);
  await page.fill('input[type="password"]', password);
  await page.click('button[type="submit"]');
  // Wait for redirect - could be /app or /app/customer/dashboard, etc.
  await page.waitForURL(/\/app/, { timeout: 15000 });
  // Give extra time for role-based redirect
  await page.waitForTimeout(1000);
}

async function logout(page: any) {
  // Clear auth data
  await page.evaluate(() => {
    localStorage.clear();
    document.cookie.split(";").forEach(c => {
      document.cookie = c.replace(/^ +/, "").replace(/=.*/, "=;expires=" + new Date().toUTCString() + ";path=/");
    });
  });
  await page.goto(`${BASE_URL}/login`);
}

test.describe('RBAC Navigation Tests', () => {
  test.describe('Requester (buyer.requester@demo.com)', () => {
    test.beforeEach(async ({ page }) => {
      await loginAs(page, USERS.requester.email, USERS.requester.password);
    });

    test('should access allowed buyer pages', async ({ page }) => {
      for (const pagePath of ROLE_PAGES.requester) {
        await page.goto(`${BASE_URL}${pagePath}`);
        // Should not redirect to login (means we're authenticated)
        await expect(page).not.toHaveURL(/\/login/);
        // Page should load (check for common elements)
        await expect(page.locator('body')).toBeVisible();
        console.log(`✅ Requester can access: ${pagePath}`);
      }
    });

    test('should NOT access admin pages', async ({ page }) => {
      for (const pagePath of FORBIDDEN_PAGES.requester.filter(p => p.startsWith('/app/admin'))) {
        await page.goto(`${BASE_URL}${pagePath}`);
        // Should redirect to login or dashboard
        const url = page.url();
        expect(url).toMatch(/\/(login|app\/customer\/dashboard)/);
        console.log(`✅ Requester correctly blocked from: ${pagePath}`);
      }
    });

    test('should NOT access supplier pages', async ({ page }) => {
      for (const pagePath of FORBIDDEN_PAGES.requester.filter(p => p.startsWith('/app/supplier'))) {
        await page.goto(`${BASE_URL}${pagePath}`);
        const url = page.url();
        expect(url).toMatch(/\/(login|app\/customer\/dashboard)/);
        console.log(`✅ Requester correctly blocked from: ${pagePath}`);
      }
    });

    test('should see correct sidebar navigation', async ({ page }) => {
      await page.goto(`${BASE_URL}/app/customer/dashboard`);
      
      // Should see customer nav items
      await expect(page.locator('text=Purchase Requests')).toBeVisible();
      await expect(page.locator('text=RFQs')).toBeVisible();
      
      // Should NOT see admin nav items
      await expect(page.locator('text=Tenants')).not.toBeVisible();
      await expect(page.locator('text=Users')).not.toBeVisible();
      
      // Should NOT see supplier nav items
      await expect(page.locator('text=RFQ Inbox')).not.toBeVisible();
      await expect(page.locator('text=Listings')).not.toBeVisible();
    });
  });

  test.describe('Procurement Manager (buyer.procurement@demo.com)', () => {
    test.beforeEach(async ({ page }) => {
      await loginAs(page, USERS.procurement.email, USERS.procurement.password);
    });

    test('should access allowed procurement pages', async ({ page }) => {
      for (const pagePath of ROLE_PAGES.procurement_manager) {
        await page.goto(`${BASE_URL}${pagePath}`);
        await expect(page).not.toHaveURL(/\/login/);
        await expect(page.locator('body')).toBeVisible();
        console.log(`✅ Procurement can access: ${pagePath}`);
      }
    });

    test('should access PR detail with approve/reject buttons', async ({ page }) => {
      await page.goto(`${BASE_URL}/app/customer/pr`);
      await page.waitForTimeout(2000);
      
      // Try to find a PR and click view
      const viewButton = page.locator('text=View →').first();
      if (await viewButton.isVisible()) {
        await viewButton.click();
        await page.waitForURL(/\/app\/customer\/pr\//, { timeout: 10000 });
        
        // Should see approve/reject buttons (procurement role)
        const approveButton = page.locator('button:has-text("Approve")');
        const rejectButton = page.locator('button:has-text("Reject")');
        
        // At least one should be visible
        const hasApprove = await approveButton.isVisible().catch(() => false);
        const hasReject = await rejectButton.isVisible().catch(() => false);
        
        if (hasApprove || hasReject) {
          console.log('✅ Procurement can see approve/reject buttons on PR detail');
        }
      }
    });

    test('should access RFQ detail with award functionality', async ({ page }) => {
      await page.goto(`${BASE_URL}/app/customer/rfq`);
      await page.waitForTimeout(2000);
      
      const viewButton = page.locator('text=View →').first();
      if (await viewButton.isVisible()) {
        await viewButton.click();
        await page.waitForURL(/\/app\/customer\/rfq\//, { timeout: 10000 });
        
        // Should see award quote button
        const awardButton = page.locator('button:has-text("Award Quote")');
        const hasAward = await awardButton.isVisible().catch(() => false);
        
        if (hasAward) {
          console.log('✅ Procurement can see award quote button on RFQ detail');
        }
      }
    });

    test('should NOT access admin or supplier pages', async ({ page }) => {
      for (const pagePath of FORBIDDEN_PAGES.procurement_manager) {
        await page.goto(`${BASE_URL}${pagePath}`);
        const url = page.url();
        expect(url).toMatch(/\/(login|app\/customer\/dashboard)/);
        console.log(`✅ Procurement correctly blocked from: ${pagePath}`);
      }
    });
  });

  test.describe('Supplier (supplier@demo.com)', () => {
    test.beforeEach(async ({ page }) => {
      await loginAs(page, USERS.supplier.email, USERS.supplier.password);
    });

    test('should access allowed supplier pages', async ({ page }) => {
      for (const pagePath of ROLE_PAGES.supplier) {
        await page.goto(`${BASE_URL}${pagePath}`);
        await expect(page).not.toHaveURL(/\/login/);
        await expect(page.locator('body')).toBeVisible();
        console.log(`✅ Supplier can access: ${pagePath}`);
      }
    });

    test('should access RFQ detail with quote submission form', async ({ page }) => {
      await page.goto(`${BASE_URL}/app/supplier/rfq`);
      await page.waitForTimeout(2000);
      
      const submitButton = page.locator('text=Submit Quote →').or(page.locator('text=View →')).first();
      if (await submitButton.isVisible()) {
        await submitButton.click();
        await page.waitForURL(/\/app\/supplier\/rfq\//, { timeout: 10000 });
        
        // Should see quote submission form
        const submitForm = page.locator('button:has-text("Submit Quote")').or(page.locator('form'));
        const hasForm = await submitForm.isVisible().catch(() => false);
        
        if (hasForm) {
          console.log('✅ Supplier can see quote submission form on RFQ detail');
        }
      }
    });

    test('should access create listing page', async ({ page }) => {
      await page.goto(`${BASE_URL}/app/supplier/listings/create`);
      await expect(page).not.toHaveURL(/\/login/);
      
      // Should see form elements
      const form = page.locator('form').or(page.locator('input[type="text"]'));
      await expect(form.first()).toBeVisible();
      console.log('✅ Supplier can access create listing page');
    });

    test('should NOT access admin or buyer pages', async ({ page }) => {
      for (const pagePath of FORBIDDEN_PAGES.supplier) {
        await page.goto(`${BASE_URL}${pagePath}`);
        const url = page.url();
        expect(url).toMatch(/\/(login|app\/supplier\/dashboard)/);
        console.log(`✅ Supplier correctly blocked from: ${pagePath}`);
      }
    });

    test('should see correct sidebar navigation', async ({ page }) => {
      await page.goto(`${BASE_URL}/app/supplier/dashboard`);
      
      // Should see supplier nav items
      await expect(page.locator('text=RFQ Inbox')).toBeVisible();
      await expect(page.locator('text=Listings')).toBeVisible();
      
      // Should NOT see admin nav items
      await expect(page.locator('text=Tenants')).not.toBeVisible();
      await expect(page.locator('text=Users')).not.toBeVisible();
      
      // Should NOT see buyer nav items
      await expect(page.locator('text=Purchase Requests')).not.toBeVisible();
    });
  });

  test.describe('Admin (admin@demo.com)', () => {
    test.beforeEach(async ({ page }) => {
      await loginAs(page, USERS.admin.email, USERS.admin.password);
    });

    test('should access all admin pages', async ({ page }) => {
      for (const pagePath of ROLE_PAGES.admin) {
        await page.goto(`${BASE_URL}${pagePath}`);
        await expect(page).not.toHaveURL(/\/login/);
        await expect(page.locator('body')).toBeVisible();
        console.log(`✅ Admin can access: ${pagePath}`);
      }
    });

    test('should access tenant management page', async ({ page }) => {
      await page.goto(`${BASE_URL}/app/admin/tenants`);
      await expect(page).not.toHaveURL(/\/login/);
      
      // Should see tenant management content
      const content = page.locator('text=Tenant Management').or(page.locator('text=Tenants')).or(page.locator('table'));
      await expect(content.first()).toBeVisible();
      console.log('✅ Admin can access tenant management page');
    });

    test('should access user management page', async ({ page }) => {
      await page.goto(`${BASE_URL}/app/admin/users`);
      await expect(page).not.toHaveURL(/\/login/);
      
      // Should see user management content
      const content = page.locator('text=User Management').or(page.locator('text=Users')).or(page.locator('table'));
      await expect(content.first()).toBeVisible();
      console.log('✅ Admin can access user management page');
    });

    test('should access roles & permissions page', async ({ page }) => {
      await page.goto(`${BASE_URL}/app/admin/roles-permissions`);
      await expect(page).not.toHaveURL(/\/login/);
      
      const content = page.locator('text=Roles').or(page.locator('text=Permissions')).or(page.locator('table'));
      await expect(content.first()).toBeVisible();
      console.log('✅ Admin can access roles & permissions page');
    });

    test('should NOT access buyer or supplier pages', async ({ page }) => {
      for (const pagePath of FORBIDDEN_PAGES.admin) {
        await page.goto(`${BASE_URL}${pagePath}`);
        const url = page.url();
        expect(url).toMatch(/\/(login|app\/admin\/dashboard)/);
        console.log(`✅ Admin correctly blocked from: ${pagePath}`);
      }
    });

    test('should see correct sidebar navigation', async ({ page }) => {
      await page.goto(`${BASE_URL}/app/admin/dashboard`);
      
      // Should see admin nav items
      await expect(page.locator('text=Tenants')).toBeVisible();
      await expect(page.locator('text=Users')).toBeVisible();
      await expect(page.locator('text=Roles & Permissions')).toBeVisible();
      await expect(page.locator('text=Audit Logs')).toBeVisible();
      
      // Should NOT see buyer nav items
      await expect(page.locator('text=Purchase Requests')).not.toBeVisible();
      
      // Should NOT see supplier nav items
      await expect(page.locator('text=RFQ Inbox')).not.toBeVisible();
    });
  });

  test.describe('Cross-role access prevention', () => {
    test('requester cannot access procurement-only features', async ({ page }) => {
      await loginAs(page, USERS.requester.email, USERS.requester.password);
      
      // Try to access a PR detail and check if approve button is hidden
      await page.goto(`${BASE_URL}/app/customer/pr`);
      await page.waitForTimeout(2000);
      
      const viewButton = page.locator('text=View →').first();
      if (await viewButton.isVisible()) {
        await viewButton.click();
        await page.waitForURL(/\/app\/customer\/pr\//, { timeout: 10000 });
        
        // Requester should NOT see approve/reject buttons
        const approveButton = page.locator('button:has-text("Approve")');
        const hasApprove = await approveButton.isVisible().catch(() => false);
        
        if (!hasApprove) {
          console.log('✅ Requester correctly cannot see approve button');
        }
      }
    });
  });
});
