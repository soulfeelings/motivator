# Admin Guide

## First-Time Setup

1. **Sign up** via Supabase Auth (email + password)
2. **Create company** — Admin Panel → Company → Create Company
3. **Invite employees** — Invites → enter email + role → Send
4. **Create badges** — Badges → define badges with XP/coin rewards
5. **Create achievements** — Achievements → define metric rules
6. **Set up integrations** — Integrations → connect Jira/GitHub/etc.
7. **Create rewards** — Rewards → add items employees can redeem

## Admin Panel Navigation

| Page | Purpose |
|---|---|
| **Dashboard** | Quick stats overview |
| **Company** | Edit company name, slug |
| **Members** | View all members with XP/level/coins, manage roles |
| **Badges** | Create/delete badges with XP and coin rewards |
| **Achievements** | Define metric-based rules that auto-award on completion |
| **Leaderboard** | Ranked member list by XP |
| **Game Plans** | Visual drag-and-drop flow editor for gamification rules |
| **Teams** | Create teams, assign members, start team battles |
| **Challenges** | View 1v1 challenges between members |
| **Rewards** | Create reward items, manage redemptions |
| **Tournaments** | Seasonal competitions with prize pools |
| **Analytics** | Engagement, performance, and ROI metrics |
| **Integrations** | Connect Jira, GitHub, Salesforce, Zendesk |
| **Webhooks** | Slack/Teams notifications for events |
| **Invites** | Invite new employees to the workspace |

## Roles

| Role | Can do |
|---|---|
| **Owner** | Everything + delete company |
| **Admin** | Everything except delete company |
| **Manager** | View analytics + all member features |
| **Employee** | View leaderboard, join challenges, redeem rewards |

## Game Plan Builder

The visual flow editor (like n8n) lets you design gamification flows:

1. **Trigger** (amber) — metric event that starts the flow
2. **Condition** (violet) — check if a threshold is met
3. **Action** (emerald) — award XP, coins, or badges

Drag nodes from the sidebar, connect them, and save. Activate the plan to make it live.

## Achievement Engine

Achievements auto-trigger when metrics are reported:

1. External tool sends webhook → metric is recorded
2. Achievement engine checks all rules for that metric
3. If condition is met → awards XP, coins, badge (if linked)
4. Push notification sent to employee

### Example Achievement

- **Name:** Gold Closer
- **Metric:** `deals_closed`
- **Operator:** `>=`
- **Threshold:** 10
- **Rewards:** +50 XP, +25 coins, "Gold Closer" badge

## Tournaments

Seasonal competitions for the whole company:

1. Create tournament with metric, dates, and prize pool
2. Open registration → employees join
3. Activate → scores are submitted
4. Complete → top 3 get XP/coin prizes, ranks calculated
