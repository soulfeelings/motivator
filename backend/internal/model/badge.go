package model

import "time"

type Badge struct {
	ID          string    `json:"id"`
	CompanyID   string    `json:"company_id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	IconURL     *string   `json:"icon_url,omitempty"`
	XPReward    int       `json:"xp_reward"`
	CoinReward  int       `json:"coin_reward"`
	CreatedAt   time.Time `json:"created_at"`
}

type MemberBadge struct {
	ID           string    `json:"id"`
	MembershipID string    `json:"membership_id"`
	BadgeID      string    `json:"badge_id"`
	AwardedAt    time.Time `json:"awarded_at"`
	Badge        *Badge    `json:"badge,omitempty"`
}

type CreateBadgeRequest struct {
	Name        string  `json:"name" validate:"required"`
	Description *string `json:"description,omitempty"`
	IconURL     *string `json:"icon_url,omitempty"`
	XPReward    int     `json:"xp_reward"`
	CoinReward  int     `json:"coin_reward"`
}

type AwardBadgeRequest struct {
	BadgeID string `json:"badge_id" validate:"required"`
}
