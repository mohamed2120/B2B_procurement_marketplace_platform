import { test, expect } from '@playwright/test';

test.describe('Authentication', () => {
  test('should login with demo account', async ({ page }) => {
    await page.goto('/login');

    // Fill in login form
    await page.fill('input[type="email"]', 'buyer@demo.com');
    await page.fill('input[type="password"]', 'demo123456');

    // Submit form
    await page.click('button[type="submit"]');

    // Should redirect to dashboard
    await expect(page).toHaveURL('/');
    await expect(page.locator('text=Dashboard')).toBeVisible();
  });

  test('should show error on invalid credentials', async ({ page }) => {
    await page.goto('/login');

    await page.fill('input[type="email"]', 'invalid@demo.com');
    await page.fill('input[type="password"]', 'wrongpassword');
    await page.click('button[type="submit"]');

    // Should show error message
    await expect(page.locator('text=/Login failed|invalid credentials/i')).toBeVisible();
  });
});
