import { test, expect } from '@playwright/test';

const BASE = 'https://admin-production-d9b7.up.railway.app';
const API = 'https://backend-production-880c.up.railway.app/api/v1';
const SUPABASE_URL = 'https://evfkxiphjhriwaozppsf.supabase.co';
const SUPABASE_KEY = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6ImV2Zmt4aXBoamhyaXdhb3pwcHNmIiwicm9sZSI6ImFub24iLCJpYXQiOjE3NjQwODgwMTcsImV4cCI6MjA3OTY2NDAxN30.Nqxq6azBpBq6qCq8aZ-DwEeH9E0eKlASpEbOf3Lgj9E';
const TEST_EMAIL = 'motivator.test.user@gmail.com';
const TEST_PASSWORD = 'TestPass123!';
const CID = 'aeedbc59-11df-42c6-b645-108d842f3437';

async function getToken(): Promise<string> {
  const res = await fetch(`${SUPABASE_URL}/auth/v1/token?grant_type=password`, {
    method: 'POST',
    headers: { 'apikey': SUPABASE_KEY, 'Content-Type': 'application/json' },
    body: JSON.stringify({ email: TEST_EMAIL, password: TEST_PASSWORD }),
  });
  const data = await res.json();
  return data.access_token;
}

async function apiCall(token: string, method: string, path: string, body?: any) {
  const res = await fetch(`${API}${path}`, {
    method,
    headers: { 'Authorization': `Bearer ${token}`, 'Content-Type': 'application/json' },
    body: body ? JSON.stringify(body) : undefined,
  });
  return res.json();
}

test.describe('Social Games — API Tests', () => {

  test('create, add questions, launch, and list a trivia game', async () => {
    const token = await getToken();

    // Create
    const created = await apiCall(token, 'POST', `/companies/${CID}/social-games`, {
      name: 'API Test Trivia', game_type: 'trivia', description: 'E2E test', duration_hours: 24,
    });
    expect(created.success).toBe(true);
    const gameId = created.data.id;
    expect(gameId).toBeTruthy();
    expect(created.data.status).toBe('draft');

    // Add questions
    const q1 = await apiCall(token, 'POST', `/companies/${CID}/social-games/${gameId}/questions`, {
      question: 'What is 2+2?', options: ['3', '4', '5', '6'], correct_index: 1,
    });
    expect(q1.success).toBe(true);

    const q2 = await apiCall(token, 'POST', `/companies/${CID}/social-games/${gameId}/questions`, {
      question: 'Capital of France?', options: ['London', 'Berlin', 'Paris', 'Madrid'], correct_index: 2,
    });
    expect(q2.success).toBe(true);

    // List questions
    const questions = await apiCall(token, 'GET', `/companies/${CID}/social-games/${gameId}/questions`);
    expect(questions.success).toBe(true);
    expect(questions.data.length).toBe(2);

    // Launch
    const launched = await apiCall(token, 'POST', `/companies/${CID}/social-games/${gameId}/launch`);
    expect(launched.success).toBe(true);

    // Verify status
    const game = await apiCall(token, 'GET', `/companies/${CID}/social-games/${gameId}`);
    expect(game.data.status).toBe('active');

    // List games
    const list = await apiCall(token, 'GET', `/companies/${CID}/social-games`);
    expect(list.success).toBe(true);
    expect(list.data.length).toBeGreaterThan(0);

    // Complete
    const completed = await apiCall(token, 'POST', `/companies/${CID}/social-games/${gameId}/complete`);
    expect(completed.success).toBe(true);

    // Results
    const results = await apiCall(token, 'GET', `/companies/${CID}/social-games/${gameId}/results`);
    expect(results.success).toBe(true);

    // Cleanup
    await apiCall(token, 'DELETE', `/companies/${CID}/social-games/${gameId}`);
  });

  test('create and launch a two truths game', async () => {
    const token = await getToken();

    const created = await apiCall(token, 'POST', `/companies/${CID}/social-games`, {
      name: 'API Test Two Truths', game_type: 'two_truths', description: 'E2E test', duration_hours: 24,
    });
    expect(created.success).toBe(true);
    const gameId = created.data.id;

    // Launch
    const launched = await apiCall(token, 'POST', `/companies/${CID}/social-games/${gameId}/launch`);
    expect(launched.success).toBe(true);

    // Submit entry
    const submitted = await apiCall(token, 'POST', `/companies/${CID}/social-games/${gameId}/submit`, {
      statements: ['I can juggle', 'I speak 3 languages', 'I climbed Everest'],
      lie_index: 2,
    });
    expect(submitted.success).toBe(true);

    // Complete
    await apiCall(token, 'POST', `/companies/${CID}/social-games/${gameId}/complete`);

    // Cleanup
    await apiCall(token, 'DELETE', `/companies/${CID}/social-games/${gameId}`);
  });

  test('create and launch a photo challenge', async () => {
    const token = await getToken();

    const created = await apiCall(token, 'POST', `/companies/${CID}/social-games`, {
      name: 'API Test Photo', game_type: 'photo_challenge', description: 'E2E test', duration_hours: 24,
    });
    expect(created.success).toBe(true);
    const gameId = created.data.id;

    // Launch
    const launched = await apiCall(token, 'POST', `/companies/${CID}/social-games/${gameId}/launch`);
    expect(launched.success).toBe(true);

    // Submit photo
    const submitted = await apiCall(token, 'POST', `/companies/${CID}/social-games/${gameId}/submit`, {
      content: 'https://example.com/test-photo.jpg',
    });
    expect(submitted.success).toBe(true);

    // Complete
    await apiCall(token, 'POST', `/companies/${CID}/social-games/${gameId}/complete`);

    // Cleanup
    await apiCall(token, 'DELETE', `/companies/${CID}/social-games/${gameId}`);
  });
});

test.describe('Social Games — UI Tests', () => {

  test('admin can see Social Games page with game list', async ({ page }) => {
    await page.goto(BASE);
    await page.waitForLoadState('networkidle');
    await page.fill('input[type="email"]', TEST_EMAIL);
    await page.fill('input[type="password"]', TEST_PASSWORD);
    await page.locator('button:has-text("Sign")').click();
    await page.waitForTimeout(4000);

    await page.click('a:has-text("Social Games")');
    await page.waitForTimeout(2000);

    const pageText = await page.textContent('body');
    expect(pageText).toContain('Social Games');
    expect(pageText).toContain('New Game');

    await page.screenshot({ path: 'screenshots/social-games-list.png', fullPage: false });
  });

  test('admin can open game type picker', async ({ page }) => {
    await page.goto(BASE);
    await page.waitForLoadState('networkidle');
    await page.fill('input[type="email"]', TEST_EMAIL);
    await page.fill('input[type="password"]', TEST_PASSWORD);
    await page.locator('button:has-text("Sign")').click();
    await page.waitForTimeout(4000);

    await page.click('a:has-text("Social Games")');
    await page.waitForTimeout(2000);
    await page.locator('button:has-text("New Game")').click();
    await page.waitForTimeout(500);

    const pageText = await page.textContent('body');
    expect(pageText).toContain('Trivia');
    expect(pageText).toContain('Photo Challenge');
    expect(pageText).toContain('Two Truths');

    await page.screenshot({ path: 'screenshots/social-games-picker.png', fullPage: false });
  });

  test('trivia play page loads with questions', async ({ page }) => {
    await page.goto(BASE);
    await page.waitForLoadState('networkidle');
    await page.fill('input[type="email"]', TEST_EMAIL);
    await page.fill('input[type="password"]', TEST_PASSWORD);
    await page.locator('button:has-text("Sign")').click();
    await page.waitForTimeout(4000);

    // Find active trivia game
    const token = await getToken();
    const list = await apiCall(token, 'GET', `/companies/${CID}/social-games`);
    const triviaGame = list.data?.find((g: any) => g.game_type === 'trivia' && g.status === 'active');

    if (triviaGame) {
      await page.goto(`${BASE}/play/${triviaGame.id}`);
      await page.waitForTimeout(3000);

      const pageText = await page.textContent('body');
      expect(pageText).toContain('Motivator');
      expect(pageText).toContain('Question');

      await page.screenshot({ path: 'screenshots/social-games-play-trivia.png', fullPage: false });
    }
  });
});
