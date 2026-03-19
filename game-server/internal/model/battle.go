package model

import "time"

type Battle struct {
	ID           string         `json:"id"`
	AttackerID   string         `json:"attacker_id"`
	DefenderID   string         `json:"defender_id"`
	WinnerID     *string        `json:"winner_id,omitempty"`
	AttackerLost map[string]int `json:"attacker_lost"`
	DefenderLost map[string]int `json:"defender_lost"`
	ReplayData   []ReplayFrame  `json:"replay_data"`
	CoinsWon     int            `json:"coins_won"`
	XPWon        int            `json:"xp_won"`
	FoughtAt     time.Time      `json:"fought_at"`
}

type ReplayFrame struct {
	Tick    int            `json:"tick"`
	Events []ReplayEvent  `json:"events"`
}

type ReplayEvent struct {
	Type     string `json:"type"`
	UnitID   string `json:"unit_id,omitempty"`
	TargetID string `json:"target_id,omitempty"`
	Damage   int    `json:"damage,omitempty"`
	X        int    `json:"x,omitempty"`
	Y        int    `json:"y,omitempty"`
	Side     string `json:"side"`
}

type AttackRequest struct {
	DefenderBaseID string `json:"defender_base_id" validate:"required"`
}

type BaseOverview struct {
	Base      Base        `json:"base"`
	Buildings []BaseBuilding `json:"buildings"`
	Army      []ArmyUnit  `json:"army"`
}
