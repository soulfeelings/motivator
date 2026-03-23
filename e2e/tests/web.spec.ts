import { test, expect } from '@playwright/test';

test.describe('Web - Command Center', () => {
  const BASE = 'https://web-production-8aabc.up.railway.app';

  test('loads game page', async ({ page }) => {
    await page.goto(BASE);
    await page.waitForLoadState('networkidle');
    await page.screenshot({ path: 'screenshots/web-game-load.png', fullPage: true });

    // Page should have the game container
    const title = await page.title();
    expect(title).toContain('Motivator');
  });

  test('has game container element', async ({ page }) => {
    await page.goto(BASE);
    await page.waitForLoadState('networkidle');

    const gameContainer = page.locator('#game-container');
    await expect(gameContainer).toBeVisible();
  });

  test('phaser canvas renders', async ({ page }) => {
    await page.goto(BASE);
    // Wait for Phaser to initialize
    await page.waitForTimeout(3000);
    await page.screenshot({ path: 'screenshots/web-game-canvas.png', fullPage: true });

    const canvas = page.locator('canvas');
    const count = await canvas.count();
    expect(count).toBeGreaterThan(0);
  });

  test('favicon is present', async ({ page }) => {
    await page.goto(BASE);
    await page.waitForLoadState('networkidle');

    const favicon = page.locator('link[rel="icon"]');
    const count = await favicon.count();
    expect(count).toBeGreaterThan(0);

    const href = await favicon.getAttribute('href');
    expect(href).toBe('/favicon.svg');
  });

  test('no console errors on load', async ({ page }) => {
    const errors: string[] = [];
    page.on('console', msg => {
      if (msg.type() === 'error') errors.push(msg.text());
    });

    await page.goto(BASE);
    await page.waitForTimeout(3000);

    // Filter out known non-critical errors
    const criticalErrors = errors.filter(e =>
      !e.includes('favicon') && !e.includes('404')
    );

    await page.screenshot({ path: 'screenshots/web-console-check.png', fullPage: true });

    // Should have no critical console errors
    expect(criticalErrors).toHaveLength(0);
  });
});
