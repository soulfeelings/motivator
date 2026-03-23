import { defineConfig } from '@playwright/test';

export default defineConfig({
  testDir: './tests',
  timeout: 30000,
  retries: 0,
  use: {
    headless: true,
    screenshot: 'on',
    trace: 'on-first-retry',
    video: 'on',
    viewport: { width: 1280, height: 720 },
  },
  projects: [
    {
      name: 'admin',
      use: {
        baseURL: 'https://admin-production-d9b7.up.railway.app',
      },
    },
    {
      name: 'web',
      use: {
        baseURL: 'https://web-production-8aabc.up.railway.app',
      },
    },
  ],
  reporter: [['list'], ['html', { open: 'never', outputFolder: 'report' }]],
  outputDir: './screenshots',
});
