package model

import "time"

type InviteStatus string

const (
	InviteStatusPending  InviteStatus = "pending"
	InviteStatusAccepted InviteStatus = "accepted"
	InviteStatusExpired  InviteStatus = "expired"
	InviteStatusRevoked  InviteStatus = "revoked"
)

type Invite struct {
	ID         string       `json:"id"`
	CompanyID  string       `json:"company_id"`
	Email      string       `json:"email"`
	Role       Role         `json:"role"`
	Status     InviteStatus `json:"status"`
	InvitedBy  string       `json:"invited_by"`
	Token      string       `json:"token"`
	ExpiresAt  time.Time    `json:"expires_at"`
	AcceptedAt *time.Time   `json:"accepted_at,omitempty"`
	CreatedAt  time.Time    `json:"created_at"`
}

type CreateInviteRequest struct {
	Email string `json:"email" validate:"required,email"`
	Role  Role   `json:"role" validate:"required"`
}

type BatchInviteRequest struct {
	Invites []CreateInviteRequest `json:"invites" validate:"required,min=1"`
}
