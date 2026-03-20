# Command Center — Game Server

The Command Center is a pseudo-3D RTS mini-game where employees use coins earned from work to build bases, hire armies, and battle opponents.

## Architecture

- **Game Server** (`game-server/`, port 8081) — separate Go/Fiber service
- **Game Client** (`web/`, port 5173) — Phaser.js isometric game
- **Flutter** — WebView embeds the game client

All three share the same Supabase database and JWT authentication.

## Buildings

| Building | Cost | HP | Category | Unlocks |
|---|---|---|---|---|
| HQ | Free | 500 | Core | — |
| Barracks | 100 | 200 | Production | Soldier |
| War Factory | 250 | 300 | Production | Tank |
| Power Plant | 75 | 150 | Utility | — |
| Gun Turret | 150 | 250 | Defense | +20 def/level |
| Radar Station | 200 | 150 | Utility | — |

## Units

| Unit | Cost | HP | Attack | Defense | Speed | Type |
|---|---|---|---|---|---|---|
| Soldier | 10 | 50 | 15 | 5 | 8 | Infantry |
| RPG Trooper | 25 | 40 | 30 | 3 | 6 | Infantry |
| Sniper | 35 | 30 | 45 | 2 | 4 | Infantry |
| Light Tank | 60 | 120 | 35 | 20 | 10 | Vehicle |
| Heavy Tank | 100 | 200 | 50 | 35 | 4 | Vehicle |
| APC | 45 | 150 | 10 | 30 | 12 | Vehicle |

## Battle System

Battles are **auto-resolved** with tick-based simulation:

1. Each tick, units attack a random living enemy
2. Damage = attacker.attack - defender.defense/2 + random(0-4)
3. Defender turrets add bonus damage split across defending units
4. Battle runs up to 30 ticks or until one side is eliminated
5. Winner = side with more surviving units

**Rewards:** Winner gets 50 coins + 100 XP. Losses are applied to both armies.

**Replay data** is stored as JSON frames for animated playback in the game client.

## API Endpoints

Base URL: `http://localhost:8081/api/v1/game`

| Method | Endpoint | Description |
|---|---|---|
| GET | `/base?membership_id=X` | Get or create base |
| GET | `/bases` | List all bases |
| GET | `/bases/:id` | Get base overview |
| POST | `/bases/:id/build` | Build a building |
| POST | `/bases/:id/hire` | Hire units |
| POST | `/bases/:id/deposit` | Deposit coins from work |
| POST | `/bases/:id/attack` | Attack another base |
| GET | `/bases/:id/battles` | Battle history |
| GET | `/battles/:id` | Get battle with replay data |
| GET | `/building-types` | Available buildings |
| GET | `/unit-types` | Available units |

## Coin Flow

```
Work achievements → Coins earned → Deposit to game base → Build & hire → Battle → Win more coins
```

Coins in the game are separate from the main Motivator coins. Players choose how many to deposit into the game.
