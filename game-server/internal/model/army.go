package model

type UnitType struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	Cost        int     `json:"cost"`
	HP          int     `json:"hp"`
	Attack      int     `json:"attack"`
	Defense     int     `json:"defense"`
	Speed       int     `json:"speed"`
	Category    string  `json:"category"`
}

type ArmyUnit struct {
	ID     string `json:"id"`
	BaseID string `json:"base_id"`
	UnitID string `json:"unit_id"`
	Count  int    `json:"count"`
}

type HireRequest struct {
	UnitID string `json:"unit_id" validate:"required"`
	Count  int    `json:"count" validate:"required,min=1"`
}
