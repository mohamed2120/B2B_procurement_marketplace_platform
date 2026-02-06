import { test, expect } from '@playwright/test';

test.describe('Purchase Request', () => {
  test.beforeEach(async ({ page }) => {
    // Login first
    await page.goto('/login');
    await page.fill('input[type="email"]', 'buyer@demo.com');
    await page.fill('input[type="password"]', 'demo123456');
    await page.click('button[type="submit"]');
    await page.waitForURL('/');
  });

  test('should create a new PR', async ({ page }) => {
    // Navigate to PR list
    await page.goto('/customer/prs');
    await expect(page.locator('text=Purchase Requests')).toBeVisible();

    // Click create PR button
    await page.click('text=Create PR');

    // Fill in PR form
    await page.fill('input[placeholder*="Office Supplies"]', 'Test PR from Playwright');
    await page.fill('textarea', 'This is a test purchase request created by Playwright');

    // Add an item
    await page.fill('input[placeholder="Item description"]', 'Test Item');
    await page.fill('input[placeholder="Quantity"]', '10');
    await page.fill('input[placeholder="Unit Price"]', '25.50');

    // Submit form
    await page.click('button:has-text("Create PR")');

    // Should redirect to PR list
    await expect(page).toHaveURL('/customer/prs');
    
    // Should see the new PR in the list
    await expect(page.locator('text=Test PR from Playwright')).toBeVisible();
  });
});
