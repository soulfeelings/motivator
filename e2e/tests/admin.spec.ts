import { test, expect } from '@playwright/test';

test.describe('Admin Panel', () => {
  test.use({ project: { name: 'admin' } });

  const BASE = 'https://admin-production-d9b7.up.railway.app';

  test('loads login page', async ({ page }) => {
    await page.goto(BASE);
    await page.waitForLoadState('networkidle');
    await page.screenshot({ path: 'screenshots/admin-login.png', fullPage: true });

    // Should show login form or app content
    const body = await page.textContent('body');
    expect(body).toBeTruthy();
  });

  test('login page has email and password fields', async ({ page }) => {
    await page.goto(BASE);
    await page.waitForLoadState('networkidle');

    const emailInput = page.locator('input[type="email"], input[placeholder*="email" i]');
    const passwordInput = page.locator('input[type="password"]');

    const hasEmail = await emailInput.count();
    const hasPassword = await passwordInput.count();

    await page.screenshot({ path: 'screenshots/admin-login-fields.png', fullPage: true });

    // Login page should have email and password inputs
    expect(hasEmail).toBeGreaterThan(0);
    expect(hasPassword).toBeGreaterThan(0);
  });

  test('login page has sign in button', async ({ page }) => {
    await page.goto(BASE);
    await page.waitForLoadState('networkidle');

    const signInButton = page.locator('button:has-text("Sign"), button:has-text("Log")');
    const count = await signInButton.count();

    expect(count).toBeGreaterThan(0);
  });

  test('login page shows branding', async ({ page }) => {
    await page.goto(BASE);
    await page.waitForLoadState('networkidle');

    const pageText = await page.textContent('body');
    expect(pageText).toContain('Motivator');
  });

  test('invalid login shows error', async ({ page }) => {
    await page.goto(BASE);
    await page.waitForLoadState('networkidle');

    await page.fill('input[type="email"], input[placeholder*="email" i]', 'fake@test.com');
    await page.fill('input[type="password"]', 'wrongpassword123');

    const signInButton = page.locator('button:has-text("Sign"), button:has-text("Log")');
    await signInButton.click();

    // Wait for error message
    await page.waitForTimeout(3000);
    await page.screenshot({ path: 'screenshots/admin-login-error.png', fullPage: true });

    const pageText = await page.textContent('body');
    // Should show some error indication
    const hasError = pageText?.toLowerCase().includes('invalid') ||
                     pageText?.toLowerCase().includes('error') ||
                     pageText?.toLowerCase().includes('wrong') ||
                     pageText?.toLowerCase().includes('failed');
    expect(hasError).toBeTruthy();
  });
});
