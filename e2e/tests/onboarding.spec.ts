import { test, expect } from '@playwright/test';

const BASE = 'https://admin-production-d9b7.up.railway.app';
const TEST_EMAIL = 'motivator.test.user@gmail.com';
const TEST_PASSWORD = 'TestPass123!';

test.describe('Onboarding Gate', () => {

  test('new user without company sees Setup page, not Dashboard', async ({ page }) => {
    await page.goto(BASE);
    await page.waitForLoadState('networkidle');

    // Login
    await page.fill('input[type="email"], input[placeholder*="email" i]', TEST_EMAIL);
    await page.fill('input[type="password"]', TEST_PASSWORD);
    await page.locator('button:has-text("Sign"), button:has-text("Log")').click();
    await page.waitForTimeout(3000);
    await page.waitForLoadState('networkidle');

    await page.screenshot({ path: 'screenshots/onboarding-gate.png', fullPage: false });

    // Should see the Setup page, not the Dashboard
    const pageText = await page.textContent('body');
    const isSetupPage = pageText?.includes('Create your company') || pageText?.includes('Set up your workspace');
    const isDashboard = pageText?.includes('Dashboard') && pageText?.includes('Companies');

    expect(isSetupPage).toBeTruthy();
    expect(isDashboard).toBeFalsy();
  });

  test('Setup page has no sidebar', async ({ page }) => {
    await page.goto(BASE);
    await page.waitForLoadState('networkidle');

    await page.fill('input[type="email"], input[placeholder*="email" i]', TEST_EMAIL);
    await page.fill('input[type="password"]', TEST_PASSWORD);
    await page.locator('button:has-text("Sign"), button:has-text("Log")').click();
    await page.waitForTimeout(3000);

    // Sidebar nav should NOT exist
    const sidebar = page.locator('aside, nav:has(a[href="/company"])');
    const sidebarCount = await sidebar.count();
    expect(sidebarCount).toBe(0);
  });

  test('Setup page has company creation form', async ({ page }) => {
    await page.goto(BASE);
    await page.waitForLoadState('networkidle');

    await page.fill('input[type="email"], input[placeholder*="email" i]', TEST_EMAIL);
    await page.fill('input[type="password"]', TEST_PASSWORD);
    await page.locator('button:has-text("Sign"), button:has-text("Log")').click();
    await page.waitForTimeout(3000);

    // Should have Company Name and Slug inputs
    const nameInput = page.locator('input[placeholder="Acme Corp"]');
    const slugInput = page.locator('input[placeholder="acme-corp"]');
    const createBtn = page.locator('button:has-text("Create")');

    await expect(nameInput).toBeVisible();
    await expect(slugInput).toBeVisible();
    await expect(createBtn).toBeVisible();

    await page.screenshot({ path: 'screenshots/onboarding-form.png', fullPage: false });
  });

  test('Setup page shows Motivator branding', async ({ page }) => {
    await page.goto(BASE);
    await page.waitForLoadState('networkidle');

    await page.fill('input[type="email"], input[placeholder*="email" i]', TEST_EMAIL);
    await page.fill('input[type="password"]', TEST_PASSWORD);
    await page.locator('button:has-text("Sign"), button:has-text("Log")').click();
    await page.waitForTimeout(3000);

    const pageText = await page.textContent('body');
    expect(pageText).toContain('Motivator');
    expect(pageText).toContain('Admin Panel');
  });

  test('auto-generates slug from company name', async ({ page }) => {
    await page.goto(BASE);
    await page.waitForLoadState('networkidle');

    await page.fill('input[type="email"], input[placeholder*="email" i]', TEST_EMAIL);
    await page.fill('input[type="password"]', TEST_PASSWORD);
    await page.locator('button:has-text("Sign"), button:has-text("Log")').click();
    await page.waitForTimeout(3000);

    // Type a company name
    const nameInput = page.locator('input[placeholder="Acme Corp"]');
    await nameInput.fill('My Test Company');

    // Slug should auto-generate
    const slugInput = page.locator('input[placeholder="acme-corp"]');
    const slugValue = await slugInput.inputValue();
    expect(slugValue).toBe('my-test-company');

    await page.screenshot({ path: 'screenshots/onboarding-autoslug.png', fullPage: false });
  });
});
