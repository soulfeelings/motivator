CREATE TABLE bases (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    membership_id UUID NOT NULL UNIQUE,
    name          VARCHAR(100) NOT NULL DEFAULT 'Base Alpha',
    level         INTEGER NOT NULL DEFAULT 1,
    layout        JSONB NOT NULL DEFAULT '[]',
    coins_balance INTEGER NOT NULL DEFAULT 0,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE building_types (
    id          VARCHAR(50) PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    description TEXT,
    cost        INTEGER NOT NULL,
    build_time  INTEGER NOT NULL DEFAULT 10,
    hp          INTEGER NOT NULL DEFAULT 100,
    category    VARCHAR(50) NOT NULL DEFAULT 'production',
    unlocks     VARCHAR(50)
);

INSERT INTO building_types (id, name, description, cost, build_time, hp, category, unlocks) VALUES
    ('hq',        'HQ',             'Command center. Upgrade to unlock more buildings.', 0,   0,  500, 'core',       NULL),
    ('barracks',  'Barracks',       'Train infantry units.',                             100, 30, 200, 'production', 'soldier'),
    ('factory',   'War Factory',    'Build vehicles.',                                   250, 60, 300, 'production', 'tank'),
    ('power',     'Power Plant',    'Provides power to your base.',                      75,  20, 150, 'utility',    NULL),
    ('turret',    'Gun Turret',     'Defends your base from attackers.',                 150, 45, 250, 'defense',    NULL),
    ('radar',     'Radar Station',  'Reveals opponent base layout before attack.',       200, 40, 150, 'utility',    NULL);

CREATE TABLE unit_types (
    id          VARCHAR(50) PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    description TEXT,
    cost        INTEGER NOT NULL,
    hp          INTEGER NOT NULL,
    attack      INTEGER NOT NULL,
    defense     INTEGER NOT NULL,
    speed       INTEGER NOT NULL DEFAULT 5,
    category    VARCHAR(50) NOT NULL DEFAULT 'infantry'
);

INSERT INTO unit_types (id, name, description, cost, hp, attack, defense, speed, category) VALUES
    ('soldier',  'Soldier',       'Basic infantry. Cheap but effective in numbers.', 10,  50,  15, 5,  8,  'infantry'),
    ('rpg',      'RPG Trooper',   'Anti-vehicle infantry.',                          25,  40,  30, 3,  6,  'infantry'),
    ('sniper',   'Sniper',        'Long range, high damage, fragile.',               35,  30,  45, 2,  4,  'infantry'),
    ('tank',     'Light Tank',    'Fast armored vehicle.',                            60,  120, 35, 20, 10, 'vehicle'),
    ('heavy',    'Heavy Tank',    'Slow but devastating.',                            100, 200, 50, 35, 4,  'vehicle'),
    ('apc',      'APC',           'Carries troops, boosts infantry defense.',         45,  150, 10, 30, 12, 'vehicle');

CREATE TABLE base_buildings (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    base_id     UUID NOT NULL REFERENCES bases(id) ON DELETE CASCADE,
    building_id VARCHAR(50) NOT NULL REFERENCES building_types(id),
    grid_x      INTEGER NOT NULL DEFAULT 0,
    grid_y      INTEGER NOT NULL DEFAULT 0,
    level       INTEGER NOT NULL DEFAULT 1,
    hp          INTEGER NOT NULL DEFAULT 100,
    built_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE army_units (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    base_id     UUID NOT NULL REFERENCES bases(id) ON DELETE CASCADE,
    unit_id     VARCHAR(50) NOT NULL REFERENCES unit_types(id),
    count       INTEGER NOT NULL DEFAULT 1,
    UNIQUE(base_id, unit_id)
);

CREATE TABLE battles (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    attacker_id   UUID NOT NULL REFERENCES bases(id),
    defender_id   UUID NOT NULL REFERENCES bases(id),
    winner_id     UUID REFERENCES bases(id),
    attacker_lost JSONB NOT NULL DEFAULT '{}',
    defender_lost JSONB NOT NULL DEFAULT '{}',
    replay_data   JSONB NOT NULL DEFAULT '[]',
    coins_won     INTEGER NOT NULL DEFAULT 0,
    xp_won        INTEGER NOT NULL DEFAULT 0,
    fought_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_bases_membership_id ON bases(membership_id);
CREATE INDEX idx_base_buildings_base_id ON base_buildings(base_id);
CREATE INDEX idx_army_units_base_id ON army_units(base_id);
CREATE INDEX idx_battles_attacker_id ON battles(attacker_id);
CREATE INDEX idx_battles_defender_id ON battles(defender_id);
