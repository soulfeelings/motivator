# Product Requirements Document

## Motivator

### Overview
Motivator is a B2B gamification platform that transforms everyday work into an engaging game-like experience. Companies integrate Motivator to let employees earn achievements, compete in mini-games, and win real rewards — making routine tasks fun and driving productivity.

### Problem Statement
Work is boring. Repetitive tasks kill motivation, especially in sales, support, logistics, and operations. Traditional bonuses and KPI dashboards don't engage people emotionally. Employees disengage, churn rises, and performance drops. Companies need a way to make daily work feel rewarding in real-time — not just at the end of the quarter.

### Target Audience
- **B2B buyers**: HR leads, team managers, and C-level at mid-to-large companies (50-5000 employees)
- **End users**: Employees in performance-driven roles — sales reps, support agents, warehouse workers, delivery drivers

### User Personas

**Manager Maria**
Team lead at a sales org, 30 people. Wants to boost team energy without micromanaging. Needs dashboards, reward budgets, and the ability to set up challenges.

**Player Pete**
Sales rep. Competitive, wants to see how he ranks vs peers. Motivated by public recognition and real rewards (bonuses, gift cards, extra PTO).

**Admin Anna**
HR / operations. Sets up the company workspace, manages reward budgets, configures achievement rules, monitors engagement analytics.

### Features

#### MVP
- Company workspace setup (admin panel) — **TESTING**
- Employee profiles with XP, levels, and achievement badges — **TESTING**
- Achievement engine — define rules tied to work metrics (e.g. "Close 10 deals = Gold Closer badge") — **TESTING**
- Leaderboard — real-time ranking within teams / company — **TESTING**
- Game Plan Builder — visual drag-and-drop editor (n8n-style) for designing gamification flows: triggers, conditions, rewards, quest chains — **TESTING**
- 1v1 challenges — head-to-head mini-competitions based on work metrics — **TESTING**
- Reward store — earn coins from achievements, spend on real rewards — **TESTING**
- Push notifications for achievements, challenge invites, leaderboard changes — **TESTING**
- Mobile app for employees, web admin panel for managers — **TESTING**
- Command Center — pseudo-3D RTS mini-game (build base, hire army, auto-battle opponents using work-earned coins) — **TESTING**

#### v1
- Team vs team battles — **TESTING**
- Seasonal tournaments with prize pools
- Integration with CRMs, helpdesks, ERPs (Salesforce, Zendesk, SAP)
- Analytics dashboard — engagement, productivity correlation, ROI
- Slack / Teams notifications

### User Flows

**Onboarding (Admin)**
Sign up → Create company workspace → Invite employees → Configure achievement rules → Set reward budget → Go live

**Daily loop (Employee)**
Open app → See new achievements / challenges → Check leaderboard position → Accept a 1v1 challenge → Complete work tasks → Earn XP & coins → Spend coins in reward store

**Manager loop**
Open admin → View team leaderboard → Create a new challenge → Review reward spend → Check engagement stats

### Non-goals
- Not a project management tool (no tasks, sprints, tickets)
- Not a communication tool (no chat, no feed)
- Not an HR system (no payroll, no time tracking)
- No AI-generated goals — all achievement rules are configured by humans
