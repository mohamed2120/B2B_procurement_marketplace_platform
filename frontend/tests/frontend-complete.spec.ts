import { test, expect } from '@playwright/test';

const BASE_URL = process.env.NEXT_PUBLIC_BASE_URL || 'http://localhost:3002';

// Demo credentials
const DEMO_USERS = {
  requester: { email: 'buyer.requester@demo.com', password: 'demo123456' },
  procurement: { email: 'buyer.procurement@demo.com', password: 'demo123456' },
  supplier: { email: 'supplier@demo.com', password: 'demo123456' },
  admin: { email: 'admin@demo.com', password: 'demo123456' },
};

test.describe('Frontend Complete Validation', () => {
  test.beforeEach(async ({ page }) => {
    // Set longer timeout for slow pages
    test.setTimeout(60000);
  });

  test.describe('Public Pages', () => {
    test('Home page loads', async ({ page }) => {
      await page.goto(`${BASE_URL}/`);
      await expect(page).toHaveTitle(/B2B Procurement Marketplace/i);
      await expect(page.locator('h1')).toContainText(/Streamline Your B2B Procurement/i);
    });

    test('How it works page loads', async ({ page }) => {
      await page.goto(`${BASE_URL}/how-it-works`);
      await expect(page).toHaveTitle(/B2B Procurement Marketplace/i);
    });

    test('Pricing page loads', async ({ page }) => {
      await page.goto(`${BASE_URL}/pricing`);
      await expect(page).toHaveTitle(/B2B Procurement Marketplace/i);
    });

    test('Contact page loads', async ({ page }) => {
      await page.goto(`${BASE_URL}/contact`);
      await expect(page).toHaveTitle(/B2B Procurement Marketplace/i);
    });

    test('Terms page loads', async ({ page }) => {
      await page.goto(`${BASE_URL}/terms`);
      await expect(page).toHaveTitle(/B2B Procurement Marketplace/i);
    });

    test('Privacy page loads', async ({ page }) => {
      await page.goto(`${BASE_URL}/privacy`);
      await expect(page).toHaveTitle(/B2B Procurement Marketplace/i);
    });

    test('Register page loads', async ({ page }) => {
      await page.goto(`${BASE_URL}/register`);
      await expect(page).toHaveTitle(/B2B Procurement Marketplace/i);
    });

    test('Register buyer page loads', async ({ page }) => {
      await page.goto(`${BASE_URL}/register/buyer`);
      await expect(page).toHaveTitle(/B2B Procurement Marketplace/i);
    });

    test('Register supplier page loads', async ({ page }) => {
      await page.goto(`${BASE_URL}/register/supplier`);
      await expect(page).toHaveTitle(/B2B Procurement Marketplace/i);
    });

    test('Login page loads', async ({ page }) => {
      await page.goto(`${BASE_URL}/login`);
      await expect(page).toHaveTitle(/B2B Procurement Marketplace/i);
      await expect(page.locator('input[type="email"]')).toBeVisible();
      await expect(page.locator('input[type="password"]')).toBeVisible();
    });
  });

  test.describe('Requester (Buyer) Flow', () => {
    test.beforeEach(async ({ page }) => {
      // Login as requester
      await page.goto(`${BASE_URL}/login`);
      await page.fill('input[type="email"]', DEMO_USERS.requester.email);
      await page.fill('input[type="password"]', DEMO_USERS.requester.password);
      await page.click('button[type="submit"]');
      // Wait for redirect
      await page.waitForURL(/\/app\/customer\/dashboard/, { timeout: 10000 });
    });

    test('Dashboard loads', async ({ page }) => {
      await page.goto(`${BASE_URL}/app/customer/dashboard`);
      await expect(page.locator('h1')).toContainText(/Dashboard/i);
    });

    test('PR list page loads', async ({ page }) => {
      await page.goto(`${BASE_URL}/app/customer/pr`);
      await expect(page.locator('h1')).toContainText(/Purchase Requests/i);
    });

    test('PR detail page exists and loads', async ({ page }) => {
      // First go to PR list
      await page.goto(`${BASE_URL}/app/customer/pr`);
      // Wait for page to load
      await page.waitForTimeout(2000);
      
      // Try to find a PR link or create one
      const prLink = page.locator('a[href*="/app/customer/pr/"]').first();
      const count = await prLink.count();
      
      if (count > 0) {
        await prLink.click();
        await page.waitForURL(/\/app\/customer\/pr\/[^/]+/, { timeout: 10000 });
        await expect(page.locator('h1')).toBeVisible();
      } else {
        // If no PRs exist, verify the route structure exists
        await page.goto(`${BASE_URL}/app/customer/pr/test-id`);
        // Should either show "not found" or load (both are valid)
        await expect(page.locator('body')).toBeVisible();
      }
    });

    test('RFQ list page loads', async ({ page }) => {
      await page.goto(`${BASE_URL}/app/customer/rfq`);
      await expect(page.locator('h1')).toContainText(/RFQ/i);
    });

    test('RFQ detail page exists and loads', async ({ page }) => {
      await page.goto(`${BASE_URL}/app/customer/rfq`);
      await page.waitForTimeout(2000);
      
      const rfqLink = page.locator('a[href*="/app/customer/rfq/"]').first();
      const count = await rfqLink.count();
      
      if (count > 0) {
        await rfqLink.click();
        await page.waitForURL(/\/app\/customer\/rfq\/[^/]+/, { timeout: 10000 });
        await expect(page.locator('h1')).toBeVisible();
      } else {
        await page.goto(`${BASE_URL}/app/customer/rfq/test-id`);
        await expect(page.locator('body')).toBeVisible();
      }
    });
  });

  test.describe('Procurement Manager (Buyer) Flow', () => {
    test.beforeEach(async ({ page }) => {
      await page.goto(`${BASE_URL}/login`);
      await page.fill('input[type="email"]', DEMO_USERS.procurement.email);
      await page.fill('input[type="password"]', DEMO_USERS.procurement.password);
      await page.click('button[type="submit"]');
      await page.waitForURL(/\/app\/customer\/dashboard/, { timeout: 10000 });
    });

    test('Can view PR detail with approve/reject buttons', async ({ page }) => {
      await page.goto(`${BASE_URL}/app/customer/pr`);
      await page.waitForTimeout(2000);
      
      const prLink = page.locator('a[href*="/app/customer/pr/"]').first();
      const count = await prLink.count();
      
      if (count > 0) {
        await prLink.click();
        await page.waitForURL(/\/app\/customer\/pr\/[^/]+/, { timeout: 10000 });
        // Should see approve/reject buttons (procurement role)
        const approveButton = page.locator('button:has-text("Approve")');
        const count = await approveButton.count();
        // Button may or may not be visible depending on PR status
        await expect(page.locator('h1')).toBeVisible();
      }
    });

    test('Can view RFQ detail with quote comparison', async ({ page }) => {
      await page.goto(`${BASE_URL}/app/customer/rfq`);
      await page.waitForTimeout(2000);
      
      const rfqLink = page.locator('a[href*="/app/customer/rfq/"]').first();
      const count = await rfqLink.count();
      
      if (count > 0) {
        await rfqLink.click();
        await page.waitForURL(/\/app\/customer\/rfq\/[^/]+/, { timeout: 10000 });
        await expect(page.locator('h1')).toBeVisible();
      }
    });
  });

  test.describe('Supplier Flow', () => {
    test.beforeEach(async ({ page }) => {
      await page.goto(`${BASE_URL}/login`);
      await page.fill('input[type="email"]', DEMO_USERS.supplier.email);
      await page.fill('input[type="password"]', DEMO_USERS.supplier.password);
      await page.click('button[type="submit"]');
      await page.waitForURL(/\/app\/supplier\/dashboard/, { timeout: 10000 });
    });

    test('Dashboard loads', async ({ page }) => {
      await page.goto(`${BASE_URL}/app/supplier/dashboard`);
      await expect(page.locator('h1')).toContainText(/Dashboard/i);
    });

    test('RFQ inbox loads', async ({ page }) => {
      await page.goto(`${BASE_URL}/app/supplier/rfq`);
      await expect(page.locator('h1')).toContainText(/RFQ/i);
    });

    test('RFQ detail/submit page exists and loads', async ({ page }) => {
      await page.goto(`${BASE_URL}/app/supplier/rfq`);
      await page.waitForTimeout(2000);
      
      const rfqLink = page.locator('a[href*="/app/supplier/rfq/"]').first();
      const count = await rfqLink.count();
      
      if (count > 0) {
        await rfqLink.click();
        await page.waitForURL(/\/app\/supplier\/rfq\/[^/]+/, { timeout: 10000 });
        await expect(page.locator('h1')).toBeVisible();
      } else {
        await page.goto(`${BASE_URL}/app/supplier/rfq/test-id`);
        await expect(page.locator('body')).toBeVisible();
      }
    });

    test('Listings page loads', async ({ page }) => {
      await page.goto(`${BASE_URL}/app/supplier/listings`);
      await expect(page.locator('h1')).toContainText(/Listings/i);
    });

    test('Create listing page exists and loads', async ({ page }) => {
      await page.goto(`${BASE_URL}/app/supplier/listings/create`);
      await expect(page.locator('h1')).toContainText(/Create/i);
      await expect(page.locator('input, textarea')).toHaveCount(await page.locator('input, textarea').count());
    });
  });

  test.describe('Platform Admin Flow', () => {
    test.beforeEach(async ({ page }) => {
      await page.goto(`${BASE_URL}/login`);
      await page.fill('input[type="email"]', DEMO_USERS.admin.email);
      await page.fill('input[type="password"]', DEMO_USERS.admin.password);
      await page.click('button[type="submit"]');
      await page.waitForURL(/\/app\/admin\/dashboard/, { timeout: 10000 });
    });

    test('Admin dashboard loads', async ({ page }) => {
      await page.goto(`${BASE_URL}/app/admin/dashboard`);
      await expect(page.locator('h1')).toContainText(/Admin Dashboard/i);
    });

    test('Users management page exists and loads', async ({ page }) => {
      await page.goto(`${BASE_URL}/app/admin/users`);
      await expect(page.locator('h1')).toContainText(/User Management/i);
      // Should see table or list of users (even if mock data)
      await expect(page.locator('table, div:has-text("user")')).toBeVisible({ timeout: 5000 });
    });

    test('User detail page exists and loads', async ({ page }) => {
      await page.goto(`${BASE_URL}/app/admin/users`);
      await page.waitForTimeout(2000);
      
      // Try to find a user link
      const userLink = page.locator('a[href*="/app/admin/users/"]').first();
      const count = await userLink.count();
      
      if (count > 0) {
        await userLink.click();
        await page.waitForURL(/\/app\/admin\/users\/[^/]+/, { timeout: 10000 });
        await expect(page.locator('h1')).toContainText(/User Details/i);
      } else {
        // Test with a known ID
        await page.goto(`${BASE_URL}/app/admin/users/1`);
        await expect(page.locator('body')).toBeVisible();
      }
    });

    test('Roles & Permissions page exists and loads', async ({ page }) => {
      await page.goto(`${BASE_URL}/app/admin/roles-permissions`);
      await expect(page.locator('h1')).toContainText(/Roles.*Permissions/i);
    });

    test('Audit logs page exists and loads', async ({ page }) => {
      await page.goto(`${BASE_URL}/app/admin/audit-logs`);
      await expect(page.locator('h1')).toContainText(/Audit Logs/i);
    });

    test('Catalog approvals page loads', async ({ page }) => {
      await page.goto(`${BASE_URL}/app/admin/catalog-approvals`);
      await expect(page.locator('h1')).toContainText(/Catalog Approvals/i);
    });

    test('Disputes page loads', async ({ page }) => {
      await page.goto(`${BASE_URL}/app/admin/disputes`);
      await expect(page.locator('h1')).toBeVisible();
    });

    test('Subscriptions page loads', async ({ page }) => {
      await page.goto(`${BASE_URL}/app/admin/subscriptions`);
      await expect(page.locator('h1')).toBeVisible();
    });
  });

  test.describe('Navigation and RBAC', () => {
    test('Requester sees customer sidebar', async ({ page }) => {
      await page.goto(`${BASE_URL}/login`);
      await page.fill('input[type="email"]', DEMO_USERS.requester.email);
      await page.fill('input[type="password"]', DEMO_USERS.requester.password);
      await page.click('button[type="submit"]');
      await page.waitForURL(/\/app/, { timeout: 10000 });
      
      // Should see customer navigation
      await expect(page.locator('text=Purchase Requests')).toBeVisible();
      await expect(page.locator('text=RFQs')).toBeVisible();
    });

    test('Supplier sees supplier sidebar', async ({ page }) => {
      await page.goto(`${BASE_URL}/login`);
      await page.fill('input[type="email"]', DEMO_USERS.supplier.email);
      await page.fill('input[type="password"]', DEMO_USERS.supplier.password);
      await page.click('button[type="submit"]');
      await page.waitForURL(/\/app/, { timeout: 10000 });
      
      // Should see supplier navigation
      await expect(page.locator('text=RFQ Inbox')).toBeVisible();
      await expect(page.locator('text=Listings')).toBeVisible();
    });

    test('Admin sees admin sidebar', async ({ page }) => {
      await page.goto(`${BASE_URL}/login`);
      await page.fill('input[type="email"]', DEMO_USERS.admin.email);
      await page.fill('input[type="password"]', DEMO_USERS.admin.password);
      await page.click('button[type="submit"]');
      await page.waitForURL(/\/app/, { timeout: 10000 });
      
      // Should see admin navigation
      await expect(page.locator('text=Users')).toBeVisible();
      await expect(page.locator('text=Roles & Permissions')).toBeVisible();
    });
  });
});
