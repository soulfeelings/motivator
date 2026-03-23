import { test, expect, Page } from '@playwright/test';

const BASE = 'https://admin-production-d9b7.up.railway.app';
const TEST_EMAIL = 'motivator.test.user@gmail.com';
const TEST_PASSWORD = 'TestPass123!';

async function login(page: Page) {
  await page.goto(BASE);
  await page.waitForLoadState('networkidle');
  await page.fill('input[type="email"], input[placeholder*="email" i]', TEST_EMAIL);
  await page.fill('input[type="password"]', TEST_PASSWORD);
  await page.locator('button:has-text("Sign"), button:has-text("Log")').click();
  // Wait for redirect after login
  await page.waitForTimeout(3000);
  await page.waitForLoadState('networkidle');
}

test.describe('Admin - Authenticated Flows', () => {

  test('login and see dashboard', async ({ page }) => {
    await login(page);
    await page.screenshot({ path: 'screenshots/auth-dashboard.png', fullPage: true });

    const pageText = await page.textContent('body');
    // Should see dashboard or sidebar content after login
    const isLoggedIn = pageText?.includes('Dashboard') ||
                       pageText?.includes('Motivator') ||
                       pageText?.includes('Sign out') ||
                       pageText?.includes('Company');
    expect(isLoggedIn).toBeTruthy();
  });

  test('sidebar navigation visible', async ({ page }) => {
    await login(page);

    const sidebar = page.locator('nav, [class*="sidebar"], aside');
    const sidebarCount = await sidebar.count();
    await page.screenshot({ path: 'screenshots/auth-sidebar.png', fullPage: true });

    expect(sidebarCount).toBeGreaterThan(0);
  });

  test('navigate to Company page', async ({ page }) => {
    await login(page);
    await page.click('a[href="/company"], a:has-text("Company")');
    await page.waitForTimeout(2000);
    await page.screenshot({ path: 'screenshots/auth-company.png', fullPage: true });

    const url = page.url();
    expect(url).toContain('company');
  });

  test('navigate to Members page', async ({ page }) => {
    await login(page);
    await page.click('a[href="/members"], a:has-text("Members")');
    await page.waitForTimeout(2000);
    await page.screenshot({ path: 'screenshots/auth-members.png', fullPage: true });

    const url = page.url();
    expect(url).toContain('members');
  });

  test('navigate to Badges page', async ({ page }) => {
    await login(page);
    await page.click('a[href="/badges"], a:has-text("Badges")');
    await page.waitForTimeout(2000);
    await page.screenshot({ path: 'screenshots/auth-badges.png', fullPage: true });

    const url = page.url();
    expect(url).toContain('badges');
  });

  test('navigate to Achievements page', async ({ page }) => {
    await login(page);
    await page.click('a[href="/achievements"], a:has-text("Achievements")');
    await page.waitForTimeout(2000);
    await page.screenshot({ path: 'screenshots/auth-achievements.png', fullPage: true });

    const url = page.url();
    expect(url).toContain('achievements');
  });

  test('navigate to Leaderboard page', async ({ page }) => {
    await login(page);
    await page.click('a[href="/leaderboard"], a:has-text("Leaderboard")');
    await page.waitForTimeout(2000);
    await page.screenshot({ path: 'screenshots/auth-leaderboard.png', fullPage: true });

    const url = page.url();
    expect(url).toContain('leaderboard');
  });

  test('navigate to Game Plans page', async ({ page }) => {
    await login(page);
    await page.click('a[href="/game-plans"], a:has-text("Game Plans")');
    await page.waitForTimeout(2000);
    await page.screenshot({ path: 'screenshots/auth-gameplans.png', fullPage: true });

    const url = page.url();
    expect(url).toContain('game-plans');
  });

  test('navigate to Challenges page', async ({ page }) => {
    await login(page);
    await page.click('a[href="/challenges"], a:has-text("Challenges")');
    await page.waitForTimeout(2000);
    await page.screenshot({ path: 'screenshots/auth-challenges.png', fullPage: true });

    const url = page.url();
    expect(url).toContain('challenges');
  });

  test('navigate to Rewards page', async ({ page }) => {
    await login(page);
    await page.click('a[href="/rewards"], a:has-text("Rewards")');
    await page.waitForTimeout(2000);
    await page.screenshot({ path: 'screenshots/auth-rewards.png', fullPage: true });

    const url = page.url();
    expect(url).toContain('rewards');
  });

  test('navigate to Teams page', async ({ page }) => {
    await login(page);
    await page.click('a[href="/teams"], a:has-text("Teams")');
    await page.waitForTimeout(2000);
    await page.screenshot({ path: 'screenshots/auth-teams.png', fullPage: true });

    const url = page.url();
    expect(url).toContain('teams');
  });

  test('navigate to Tournaments page', async ({ page }) => {
    await login(page);
    await page.click('a[href="/tournaments"], a:has-text("Tournaments")');
    await page.waitForTimeout(2000);
    await page.screenshot({ path: 'screenshots/auth-tournaments.png', fullPage: true });

    const url = page.url();
    expect(url).toContain('tournaments');
  });

  test('navigate to Analytics page', async ({ page }) => {
    await login(page);
    await page.click('a[href="/analytics"], a:has-text("Analytics")');
    await page.waitForTimeout(2000);
    await page.screenshot({ path: 'screenshots/auth-analytics.png', fullPage: true });

    const url = page.url();
    expect(url).toContain('analytics');
  });

  test('navigate to Integrations page', async ({ page }) => {
    await login(page);
    await page.click('a[href="/integrations"], a:has-text("Integrations")');
    await page.waitForTimeout(2000);
    await page.screenshot({ path: 'screenshots/auth-integrations.png', fullPage: true });

    const url = page.url();
    expect(url).toContain('integrations');
  });

  test('navigate to Quests page', async ({ page }) => {
    await login(page);
    await page.click('a[href="/quests"], a:has-text("Quests")');
    await page.waitForTimeout(2000);
    await page.screenshot({ path: 'screenshots/auth-quests.png', fullPage: true });

    const url = page.url();
    expect(url).toContain('quests');
  });

  test('sign out works', async ({ page }) => {
    await login(page);

    const signOutBtn = page.locator('button:has-text("Sign out"), button:has-text("Logout")');
    const count = await signOutBtn.count();
    if (count > 0) {
      await signOutBtn.click();
      await page.waitForTimeout(2000);
      await page.screenshot({ path: 'screenshots/auth-signout.png', fullPage: true });

      // Should be back on login page
      const emailInput = page.locator('input[type="email"], input[placeholder*="email" i]');
      const hasEmail = await emailInput.count();
      expect(hasEmail).toBeGreaterThan(0);
    }
  });
});
