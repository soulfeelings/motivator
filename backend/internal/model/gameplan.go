package model

import "time"

type FlowData struct {
	Nodes []FlowNode `json:"nodes"`
	Edges []FlowEdge `json:"edges"`
}

type FlowNode struct {
	ID       string         `json:"id"`
	Type     string         `json:"type"`
	Position FlowPosition   `json:"position"`
	Data     map[string]any `json:"data"`
}

type FlowEdge struct {
	ID     string `json:"id"`
	Source string `json:"source"`
	Target string `json:"target"`
}

type FlowPosition struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type GamePlan struct {
	ID          string    `json:"id"`
	CompanyID   string    `json:"company_id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	FlowData    FlowData  `json:"flow_data"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateGamePlanRequest struct {
	Name        string   `json:"name" validate:"required"`
	Description *string  `json:"description,omitempty"`
	FlowData    FlowData `json:"flow_data"`
}

type UpdateGamePlanRequest struct {
	Name        *string   `json:"name,omitempty"`
	Description *string   `json:"description,omitempty"`
	FlowData    *FlowData `json:"flow_data,omitempty"`
}
