# Changelog

## 2026-03-22 — Initial Build

Built a full B2B gamification platform from scratch in a single session.

### Monorepo with 5 services:
- `backend/` — Go/Fiber REST API (17 migrations, 15+ resource types)
- `game-server/` — Separate Go/Fiber service for the RTS game
- `admin/` — React admin panel with 18 pages
- `web/` — Phaser.js isometric RTS game
- `mobile/` — Flutter app with 6 screens

### MVP features (10):
- Company workspace, invites, RBAC (owner/admin/manager/employee)
- Employee profiles with XP, levels, coins
- Achievement engine with metric-based rules
- Badge system with XP/coin rewards
- Leaderboard ranked by XP
- Game Plan Builder (React Flow visual editor, n8n-style)
- 1v1 challenges with wager system
- Reward store with redemption flow
- Push notifications via FCM
- Command Center RTS mini-game (base building, army hiring, auto-battles)

### v1 features (5):
- Team vs team battles
- Seasonal tournaments with prize pools
- Slack/Teams webhook notifications
- Generic integrations (Jira, GitHub, Salesforce, Zendesk, Custom)
- Analytics dashboard

### Extras:
- Secret Motivator quest (anonymous positivity game)
- In-app docs viewer
- 5 documentation pages
- Deployed all 4 services to Railway

### Stats:
- ~15,000+ lines of code
- 17 database migrations, 30+ tables
- 100+ API endpoints
- 12 git commits pushed to GitHub
