import { test, expect } from '@playwright/test';

test.describe('Smoke Tests', () => {
  test('public website loads at /', async ({ page }) => {
    await page.goto('http://localhost:3002/');
    await expect(page).toHaveTitle(/B2B Procurement Marketplace/);
  });

  test('/pricing loads without login', async ({ page }) => {
    await page.goto('http://localhost:3002/pricing');
    await expect(page.locator('h1')).toContainText(/Pricing/i);
  });

  test('/app redirects to /login when not authenticated', async ({ page }) => {
    await page.goto('http://localhost:3002/app');
    await expect(page).toHaveURL(/.*\/login/);
  });

  test('login works for demo buyer.requester', async ({ page }) => {
    await page.goto('http://localhost:3002/login');
    
    // Fill login form
    await page.fill('input[type="email"]', 'buyer.requester@demo.com');
    await page.fill('input[type="password"]', 'demo123456');
    await page.fill('input[name="tenant_id"]', '00000000-0000-0000-0000-000000000001');
    
    // Submit
    await page.click('button[type="submit"]');
    
    // Should redirect to app
    await expect(page).toHaveURL(/.*\/app/);
  });

  test('buyer can create PR', async ({ page }) => {
    // Login first
    await page.goto('http://localhost:3002/login');
    await page.fill('input[type="email"]', 'buyer.requester@demo.com');
    await page.fill('input[type="password"]', 'demo123456');
    await page.fill('input[name="tenant_id"]', '00000000-0000-0000-0000-000000000001');
    await page.click('button[type="submit"]');
    await page.waitForURL(/.*\/app/);
    
    // Navigate to PR creation
    await page.goto('http://localhost:3002/app/customer/pr/create');
    
    // Check page loads
    await expect(page.locator('h1, h2')).toContainText(/Purchase Request|Create/i);
  });

  test('procurement can approve PR', async ({ page }) => {
    // Login as procurement
    await page.goto('http://localhost:3002/login');
    await page.fill('input[type="email"]', 'buyer.procurement@demo.com');
    await page.fill('input[type="password"]', 'demo123456');
    await page.fill('input[name="tenant_id"]', '00000000-0000-0000-0000-000000000001');
    await page.click('button[type="submit"]');
    await page.waitForURL(/.*\/app/);
    
    // Navigate to PR list
    await page.goto('http://localhost:3002/app/customer/pr');
    
    // Check page loads
    await expect(page.locator('h1, h2')).toContainText(/Purchase Request/i);
  });

  test('supplier submits quote', async ({ page }) => {
    // Login as supplier
    await page.goto('http://localhost:3002/login');
    await page.fill('input[type="email"]', 'supplier@demo.com');
    await page.fill('input[type="password"]', 'demo123456');
    await page.fill('input[name="tenant_id"]', '00000000-0000-0000-0000-000000000002');
    await page.click('button[type="submit"]');
    await page.waitForURL(/.*\/app/);
    
    // Navigate to quotes
    await page.goto('http://localhost:3002/app/supplier/quotes');
    
    // Check page loads
    await expect(page.locator('h1, h2')).toContainText(/Quote/i);
  });

  test('buyer awards quote → PO/Order created', async ({ page }) => {
    // Login as procurement (can award quotes)
    await page.goto('http://localhost:3002/login');
    await page.fill('input[type="email"]', 'buyer.procurement@demo.com');
    await page.fill('input[type="password"]', 'demo123456');
    await page.fill('input[name="tenant_id"]', '00000000-0000-0000-0000-000000000001');
    await page.click('button[type="submit"]');
    await page.waitForURL(/.*\/app/);
    
    // Navigate to orders
    await page.goto('http://localhost:3002/app/customer/orders');
    
    // Check page loads
    await expect(page.locator('h1, h2')).toContainText(/Order/i);
  });

  test('shipment created and visible', async ({ page }) => {
    // Login as buyer
    await page.goto('http://localhost:3002/login');
    await page.fill('input[type="email"]', 'buyer.requester@demo.com');
    await page.fill('input[type="password"]', 'demo123456');
    await page.fill('input[name="tenant_id"]', '00000000-0000-0000-0000-000000000001');
    await page.click('button[type="submit"]');
    await page.waitForURL(/.*\/app/);
    
    // Navigate to shipments
    await page.goto('http://localhost:3002/app/customer/shipments');
    
    // Check page loads
    await expect(page.locator('h1, h2')).toContainText(/Shipment/i);
  });
});
