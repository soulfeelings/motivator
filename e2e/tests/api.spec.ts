import { test, expect } from '@playwright/test';

test.describe('Backend API', () => {
  const API = 'https://backend-production-880c.up.railway.app';

  test('health endpoint returns ok', async ({ request }) => {
    const response = await request.get(`${API}/health`);
    expect(response.status()).toBe(200);
    const body = await response.json();
    expect(body.status).toBe('ok');
  });

  test('swagger endpoint loads', async ({ request }) => {
    const response = await request.get(`${API}/swagger/index.html`);
    expect(response.status()).toBe(200);
  });

  test('protected endpoints return 401 without auth', async ({ request }) => {
    const endpoints = [
      { method: 'GET', path: '/api/v1/me' },
      { method: 'POST', path: '/api/v1/companies' },
    ];

    for (const ep of endpoints) {
      const response = ep.method === 'GET'
        ? await request.get(`${API}${ep.path}`)
        : await request.post(`${API}${ep.path}`);
      expect(response.status()).toBe(401);
    }
  });

  test('webhook inbound endpoint is publicly accessible', async ({ request }) => {
    const response = await request.post(`${API}/api/v1/webhooks/inbound/test-secret`, {
      data: { event: 'test' },
    });
    // Should NOT be 401 (auth error) — it should be 404 (invalid secret) or 200
    expect(response.status()).not.toBe(401);
  });
});

test.describe('Game Server API', () => {
  const GAME = 'https://game-server-production-7f61.up.railway.app';

  test('health endpoint returns ok', async ({ request }) => {
    const response = await request.get(`${GAME}/health`);
    expect(response.status()).toBe(200);
    const body = await response.json();
    expect(body.status).toBe('ok');
  });
});
