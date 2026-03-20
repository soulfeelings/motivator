# Getting Started

## Prerequisites

- **Go** 1.24+ — [install](https://go.dev/dl/)
- **Node.js** 20+ — [install](https://nodejs.org/)
- **Flutter** 3.x — [install](https://docs.flutter.dev/get-started/install)
- **Supabase** account — [supabase.com](https://supabase.com)
- **psql** (PostgreSQL client) — comes with Postgres or `brew install libpq`

## Project Structure

```
Motivator/
├── backend/        → Main API (Go/Fiber, port 8080)
├── game-server/    → Game API (Go/Fiber, port 8081)
├── admin/          → Admin panel (React/Vite, port 3001)
├── web/            → RTS game client (Phaser.js/Vite, port 5173)
├── mobile/         → Employee app (Flutter)
├── docs/           → Documentation
└── Makefile        → Root orchestrator
```

## 1. Clone & Install

```bash
git clone git@github.com:soulfeelings/motivator.git
cd motivator

# Install all dependencies
make backend-install
make admin-install
make web-install
make game-install
cd mobile && flutter pub get && cd ..
```

## 2. Supabase Setup

### Create a project
1. Go to [supabase.com/dashboard](https://supabase.com/dashboard)
2. Create a new project
3. Note down: **Project URL**, **Anon Key**, **JWT Secret**, **Database connection string** (use Session Pooler if on IPv4)

### Run migrations

```bash
# Set your connection string
export DB_URL="postgresql://postgres.YOUR_REF:YOUR_PASSWORD@aws-X-REGION.pooler.supabase.com:5432/postgres"

# Run all backend migrations (001-016)
for f in backend/migrations/*.up.sql; do psql "$DB_URL" -f "$f"; done

# Run game-server migration
psql "$DB_URL" -f game-server/migrations/001_create_game_tables.up.sql
```

### Create a test user
1. Supabase Dashboard → **Authentication** → **Users** → **Add user**
2. Enter email + password — this will be your admin account

## 3. Environment Variables

### backend/.env
```
PORT=8080
SUPABASE_URL=https://YOUR_REF.supabase.co
DATABASE_URL=postgresql://postgres.YOUR_REF:PASSWORD@aws-X-REGION.pooler.supabase.com:5432/postgres
SUPABASE_JWT_SECRET=your-jwt-secret
```

### game-server/.env
```
PORT=8081
SUPABASE_URL=https://YOUR_REF.supabase.co
DATABASE_URL=postgresql://postgres.YOUR_REF:PASSWORD@aws-X-REGION.pooler.supabase.com:5432/postgres
SUPABASE_JWT_SECRET=your-jwt-secret
```

### admin/.env
```
VITE_SUPABASE_URL=https://YOUR_REF.supabase.co
VITE_SUPABASE_ANON_KEY=your-anon-key
```

## 4. Run Everything

Open 4 terminals:

```bash
# Terminal 1 — Backend API
make backend-dev

# Terminal 2 — Game Server
make game-dev

# Terminal 3 — Admin Panel
make admin-dev

# Terminal 4 — Game Client
make web-dev
```

### Mobile (separate)
```bash
cd mobile
flutter run
```

## 5. First Login

1. Open admin panel at `http://localhost:3001`
2. Sign in with the user you created in Supabase Auth
3. Create a company (Company → Create Company)
4. You're now the owner — all features are unlocked

## What's Running

| Service | URL | Purpose |
|---|---|---|
| Backend API | `http://localhost:8080` | Main REST API |
| Game Server | `http://localhost:8081` | RTS game API |
| Admin Panel | `http://localhost:3001` | Management dashboard |
| Game Client | `http://localhost:5173` | Phaser.js game |
| Swagger | `http://localhost:8080/swagger/` | API docs |
