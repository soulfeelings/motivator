# Integrations Guide

Motivator supports real-time integrations with external tools via webhooks. When events happen in your tools (Jira ticket resolved, GitHub PR merged, etc.), Motivator automatically awards XP, coins, and badges to the right employees.

## Supported Providers

| Provider | Events | User Identification |
|---|---|---|
| **Jira** | issue_created, issue_updated, issue_resolved | Email from Jira user profile |
| **GitHub** | push, pr_opened, pr_merged, issue_opened, issue_closed | Email or GitHub login |
| **Salesforce** | deal_closed, lead_converted, opportunity_won | user_email field |
| **Zendesk** | ticket_solved, ticket_created, satisfaction_rated | Agent email |
| **Custom** | Any custom event name | email field in payload |

## Setup Guide

### Step 1: Create Integration

1. Go to **Admin Panel → Integrations → Add Integration**
2. Select your provider (e.g., Jira)
3. Give it a name (e.g., "Engineering Jira")
4. Click **Create** — you'll get a **Webhook URL**

### Step 2: Configure Event Mappings

Map external events to internal metrics:

| External Event | → | Internal Metric | Value |
|---|---|---|---|
| `issue_resolved` | → | `tickets_closed` | +1 |
| `pr_merged` | → | `prs_merged` | +1 |
| `deal_closed` | → | `deals_closed` | +1 |

This tells Motivator: "When Jira says an issue was resolved, add +1 to the `tickets_closed` metric for that user."

### Step 3: Paste Webhook URL in External Tool

#### Jira
1. Jira → **Settings → System → Webhooks**
2. Click **Create a Webhook**
3. Paste the Motivator webhook URL
4. Select events: Issue Created, Issue Updated
5. Save

#### GitHub
1. Repo → **Settings → Webhooks → Add webhook**
2. Payload URL: paste Motivator webhook URL
3. Content type: `application/json`
4. Select events: Push, Pull requests, Issues
5. Save

#### Salesforce
1. Setup → **Platform Events** or **Outbound Messages**
2. Configure to POST to Motivator webhook URL
3. Include `event` and `user_email` fields

#### Zendesk
1. Admin → **Extensions → Webhooks → Create webhook**
2. Endpoint URL: Motivator webhook URL
3. Request method: POST
4. Create a trigger that fires on ticket events

#### Custom
Send a POST request with this format:
```json
{
  "event": "your_event_name",
  "email": "user@company.com"
}
```

### Step 4: Create Achievement Rules

Now create achievements that use the metrics:

1. Go to **Achievements → New Achievement**
2. Set metric to `tickets_closed`, operator `>=`, threshold `10`
3. Set XP and coin rewards
4. Optionally link a badge

When an employee resolves their 10th ticket, they automatically get the achievement, XP, coins, and badge.

## How It Works Internally

```
Jira → POST webhook → Motivator Backend
                          ↓
                    Parse event (issue_resolved)
                          ↓
                    Find mapping (issue_resolved → tickets_closed)
                          ↓
                    Log event
                          ↓
                    Ready for achievement evaluation
```

## Monitoring

View all received events in **Admin Panel → Integrations → select integration → Recent Events**.

Each event shows:
- ✅ Green dot = processed successfully
- 🔴 Red dot = error (unmapped event, unknown user, etc.)

## Troubleshooting

| Problem | Solution |
|---|---|
| No events appearing | Check webhook URL is correct, check Jira/GitHub webhook delivery logs |
| Events show "no mapping" error | Add a mapping for that event type |
| Events show but no XP awarded | Create an achievement rule that uses the mapped metric |
| Wrong user getting credit | Check the user_field mapping matches the email format |
