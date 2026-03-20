# API Reference

Base URL: `http://localhost:8080/api/v1`

All endpoints require `Authorization: Bearer <supabase_jwt>` unless noted.

## Authentication

Auth is handled by Supabase. The backend only validates JWTs.

1. Client calls Supabase Auth (signup/login) → gets JWT
2. Client sends JWT in `Authorization: Bearer <token>` header
3. Backend validates token and extracts `user_id` + `email`

## Company

| Method | Endpoint | Auth | Description |
|---|---|---|---|
| POST | `/companies` | User | Create company (caller becomes owner) |
| GET | `/companies/:id` | Member | Get company details |
| PATCH | `/companies/:id` | Owner/Admin | Update company |
| DELETE | `/companies/:id` | Owner | Delete company |

## Members

| Method | Endpoint | Auth | Description |
|---|---|---|---|
| GET | `/companies/:id/members` | Member | List members (paginated) |
| GET | `/companies/:id/members/:mid` | Member | Get member |
| GET | `/companies/:id/members/:mid/profile` | Member | Get profile with badges |
| POST | `/companies/:id/members/:mid/xp` | Admin+ | Award XP |
| PATCH | `/companies/:id/members/:mid` | Admin+ | Update role/name |
| DELETE | `/companies/:id/members/:mid` | Admin+ | Deactivate member |
| GET | `/me` | User | Current user + memberships |

## Badges

| Method | Endpoint | Auth | Description |
|---|---|---|---|
| GET | `/companies/:id/badges` | Member | List badges |
| POST | `/companies/:id/badges` | Admin+ | Create badge |
| DELETE | `/companies/:id/badges/:bid` | Admin+ | Delete badge |
| POST | `/companies/:id/members/:mid/badges` | Admin+ | Award badge to member |

## Achievements

| Method | Endpoint | Auth | Description |
|---|---|---|---|
| GET | `/companies/:id/achievements` | Member | List achievement rules |
| POST | `/companies/:id/achievements` | Admin+ | Create achievement rule |
| DELETE | `/companies/:id/achievements/:aid` | Admin+ | Delete achievement |
| GET | `/companies/:id/members/:mid/achievements` | Member | List member's completed achievements |
| POST | `/companies/:id/members/:mid/metrics` | Member | Report metric value (triggers evaluation) |

**Metric evaluation body:**
```json
{ "metric": "deals_closed", "value": 15 }
```

## Challenges (1v1)

| Method | Endpoint | Auth | Description |
|---|---|---|---|
| GET | `/companies/:id/challenges` | Member | List challenges |
| POST | `/companies/:id/challenges` | Member | Create challenge |
| POST | `/companies/:id/challenges/:cid/accept` | Opponent | Accept |
| POST | `/companies/:id/challenges/:cid/decline` | Opponent | Decline |
| POST | `/companies/:id/challenges/:cid/score` | Participant | Report score |

## Rewards

| Method | Endpoint | Auth | Description |
|---|---|---|---|
| GET | `/companies/:id/rewards` | Member | List rewards |
| POST | `/companies/:id/rewards` | Admin+ | Create reward |
| DELETE | `/companies/:id/rewards/:rid` | Admin+ | Delete reward |
| POST | `/companies/:id/rewards/redeem` | Member | Redeem reward (spends coins) |
| GET | `/companies/:id/rewards/redemptions` | Admin+ | List all redemptions |
| POST | `/companies/:id/rewards/redemptions/:rid/fulfill` | Admin+ | Mark as fulfilled |

## Teams

| Method | Endpoint | Auth | Description |
|---|---|---|---|
| GET | `/companies/:id/teams` | Member | List teams |
| POST | `/companies/:id/teams` | Admin+ | Create team |
| GET | `/companies/:id/teams/:tid` | Member | Get team |
| DELETE | `/companies/:id/teams/:tid` | Admin+ | Delete team |
| GET | `/companies/:id/teams/:tid/members` | Member | List team members |
| POST | `/companies/:id/teams/:tid/members` | Admin+ | Add member to team |
| DELETE | `/companies/:id/teams/:tid/members/:mid` | Admin+ | Remove from team |

## Team Battles

| Method | Endpoint | Auth | Description |
|---|---|---|---|
| GET | `/companies/:id/team-battles` | Member | List team battles |
| POST | `/companies/:id/team-battles` | Admin+ | Create team battle |
| POST | `/companies/:id/team-battles/:bid/score` | Member | Report team score |

## Tournaments

| Method | Endpoint | Auth | Description |
|---|---|---|---|
| GET | `/companies/:id/tournaments` | Member | List tournaments |
| POST | `/companies/:id/tournaments` | Admin+ | Create tournament |
| GET | `/companies/:id/tournaments/:tid` | Member | Get tournament |
| PATCH | `/companies/:id/tournaments/:tid/status` | Admin+ | Update status |
| DELETE | `/companies/:id/tournaments/:tid` | Admin+ | Delete |
| POST | `/companies/:id/tournaments/:tid/join` | Member | Join tournament |
| POST | `/companies/:id/tournaments/:tid/leave` | Member | Leave tournament |
| POST | `/companies/:id/tournaments/:tid/score` | Member | Submit score |
| GET | `/companies/:id/tournaments/:tid/standings` | Member | Get standings |
| POST | `/companies/:id/tournaments/:tid/complete` | Admin+ | Complete & award prizes |

## Game Plans

| Method | Endpoint | Auth | Description |
|---|---|---|---|
| GET | `/companies/:id/game-plans` | Member | List game plans |
| POST | `/companies/:id/game-plans` | Admin+ | Create game plan |
| GET | `/companies/:id/game-plans/:pid` | Member | Get game plan with flow data |
| PATCH | `/companies/:id/game-plans/:pid` | Admin+ | Update game plan |
| PUT | `/companies/:id/game-plans/:pid/flow` | Admin+ | Save flow data |
| POST | `/companies/:id/game-plans/:pid/activate` | Admin+ | Activate |
| POST | `/companies/:id/game-plans/:pid/deactivate` | Admin+ | Deactivate |
| DELETE | `/companies/:id/game-plans/:pid` | Admin+ | Delete |

## Leaderboard

| Method | Endpoint | Auth | Description |
|---|---|---|---|
| GET | `/companies/:id/leaderboard?limit=50` | Member | Ranked by XP |

## Integrations

| Method | Endpoint | Auth | Description |
|---|---|---|---|
| GET | `/companies/:id/integrations` | Admin+ | List integrations |
| POST | `/companies/:id/integrations` | Admin+ | Create integration |
| DELETE | `/companies/:id/integrations/:iid` | Admin+ | Delete |
| GET | `/companies/:id/integrations/:iid/mappings` | Admin+ | List event mappings |
| POST | `/companies/:id/integrations/:iid/mappings` | Admin+ | Create mapping |
| DELETE | `/companies/:id/integrations/:iid/mappings/:mid` | Admin+ | Delete mapping |
| GET | `/companies/:id/integrations/:iid/events` | Admin+ | Recent events log |
| POST | `/webhooks/inbound/:secret` | **Public** | Inbound webhook receiver |

## Webhooks (Slack/Teams)

| Method | Endpoint | Auth | Description |
|---|---|---|---|
| GET | `/companies/:id/webhooks` | Admin+ | List webhooks |
| POST | `/companies/:id/webhooks` | Admin+ | Create webhook |
| DELETE | `/companies/:id/webhooks/:wid` | Admin+ | Delete webhook |

## Analytics

| Method | Endpoint | Auth | Description |
|---|---|---|---|
| GET | `/companies/:id/analytics` | Admin+/Manager | Full dashboard data |

## Notifications

| Method | Endpoint | Auth | Description |
|---|---|---|---|
| POST | `/companies/:id/notifications/register` | Member | Register device token |
| POST | `/companies/:id/notifications/unregister` | Member | Unregister device token |

## Invites

| Method | Endpoint | Auth | Description |
|---|---|---|---|
| POST | `/companies/:id/invites` | Admin+ | Send invite |
| GET | `/companies/:id/invites` | Admin+ | List invites |
| DELETE | `/companies/:id/invites/:iid` | Admin+ | Revoke invite |
| POST | `/invites/:token/accept` | User | Accept invite |

## Standard Response Format

```json
{
  "success": true,
  "data": { ... }
}
```

Error:
```json
{
  "success": false,
  "error": "error message"
}
```

Paginated:
```json
{
  "success": true,
  "data": [ ... ],
  "meta": { "page": 1, "per_page": 20, "total": 42, "total_pages": 3 }
}
```
