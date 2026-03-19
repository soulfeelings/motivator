package model

import "time"

type DeviceToken struct {
	ID           string    `json:"id"`
	MembershipID string    `json:"membership_id"`
	Token        string    `json:"token"`
	Platform     string    `json:"platform"`
	CreatedAt    time.Time `json:"created_at"`
}

type RegisterTokenRequest struct {
	Token    string `json:"token" validate:"required"`
	Platform string `json:"platform" validate:"required"`
}

type Notification struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	Data  map[string]string `json:"data,omitempty"`
}
