package model

import "time"

type Webhook struct {
	ID        string    `json:"id"`
	CompanyID string    `json:"company_id"`
	Name      string    `json:"name"`
	URL       string    `json:"url"`
	Platform  string    `json:"platform"`
	Events    []string  `json:"events"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateWebhookRequest struct {
	Name     string   `json:"name" validate:"required"`
	URL      string   `json:"url" validate:"required"`
	Platform string   `json:"platform" validate:"required"`
	Events   []string `json:"events"`
}
