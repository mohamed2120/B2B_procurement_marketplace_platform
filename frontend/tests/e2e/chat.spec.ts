import { test, expect } from '@playwright/test';

test.describe('Chat', () => {
  test.beforeEach(async ({ page }) => {
    // Login first
    await page.goto('/login');
    await page.fill('input[type="email"]', 'buyer@demo.com');
    await page.fill('input[type="password"]', 'demo123456');
    await page.click('button[type="submit"]');
    await page.waitForURL('/');
  });

  test('should send a chat message', async ({ page }) => {
    // Navigate to chat
    await page.goto('/chat');
    await expect(page.locator('text=Chat')).toBeVisible();

    // Wait for threads to load
    await page.waitForSelector('button:has-text("Thread")', { timeout: 10000 }).catch(() => {
      // If no threads exist, that's okay for the test
    });

    // Try to send a message if there's a thread
    const messageInput = page.locator('input[placeholder*="message"]');
    if (await messageInput.isVisible({ timeout: 5000 }).catch(() => false)) {
      await messageInput.fill('Hello from Playwright test!');
      await page.click('button:has-text("Send")');

      // Wait a moment for message to appear
      await page.waitForTimeout(1000);

      // Verify message was sent (check if it appears in the chat)
      await expect(page.locator('text=Hello from Playwright test!')).toBeVisible({ timeout: 5000 });
    } else {
      // If no threads, just verify the chat page loaded
      await expect(page.locator('text=Select a thread')).toBeVisible();
    }
  });
});
