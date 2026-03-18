package model

import "time"

type Role string

const (
	RoleOwner    Role = "owner"
	RoleAdmin    Role = "admin"
	RoleManager  Role = "manager"
	RoleEmployee Role = "employee"
)

type Membership struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	CompanyID   string    `json:"company_id"`
	Role        Role      `json:"role"`
	DisplayName *string   `json:"display_name,omitempty"`
	JobTitle    *string   `json:"job_title,omitempty"`
	IsActive    bool      `json:"is_active"`
	JoinedAt    time.Time `json:"joined_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UpdateMembershipRequest struct {
	Role        *Role   `json:"role,omitempty"`
	DisplayName *string `json:"display_name,omitempty"`
	JobTitle    *string `json:"job_title,omitempty"`
}
