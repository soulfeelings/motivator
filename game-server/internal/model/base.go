package model

import "time"

type Base struct {
	ID           string         `json:"id"`
	MembershipID string         `json:"membership_id"`
	Name         string         `json:"name"`
	Level        int            `json:"level"`
	Layout       []BuildingSlot `json:"layout"`
	CoinsBalance int            `json:"coins_balance"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

type BuildingSlot struct {
	BuildingID string `json:"building_id"`
	GridX      int    `json:"grid_x"`
	GridY      int    `json:"grid_y"`
	Level      int    `json:"level"`
}

type BuildingType struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	Cost        int     `json:"cost"`
	BuildTime   int     `json:"build_time"`
	HP          int     `json:"hp"`
	Category    string  `json:"category"`
	Unlocks     *string `json:"unlocks,omitempty"`
}

type BaseBuilding struct {
	ID         string    `json:"id"`
	BaseID     string    `json:"base_id"`
	BuildingID string    `json:"building_id"`
	GridX      int       `json:"grid_x"`
	GridY      int       `json:"grid_y"`
	Level      int       `json:"level"`
	HP         int       `json:"hp"`
	BuiltAt    time.Time `json:"built_at"`
}

type BuildRequest struct {
	BuildingID string `json:"building_id" validate:"required"`
	GridX      int    `json:"grid_x"`
	GridY      int    `json:"grid_y"`
}

type DepositCoinsRequest struct {
	Amount int `json:"amount" validate:"required,min=1"`
}
