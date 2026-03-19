package model

import "time"

type RedemptionStatus string

const (
	RedemptionPending   RedemptionStatus = "pending"
	RedemptionApproved  RedemptionStatus = "approved"
	RedemptionFulfilled RedemptionStatus = "fulfilled"
	RedemptionRejected  RedemptionStatus = "rejected"
)

type Reward struct {
	ID         string    `json:"id"`
	CompanyID  string    `json:"company_id"`
	Name       string    `json:"name"`
	Description *string  `json:"description,omitempty"`
	CostCoins  int       `json:"cost_coins"`
	Stock      *int      `json:"stock,omitempty"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Redemption struct {
	ID           string           `json:"id"`
	MembershipID string           `json:"membership_id"`
	RewardID     string           `json:"reward_id"`
	CoinsSpent   int              `json:"coins_spent"`
	Status       RedemptionStatus `json:"status"`
	RedeemedAt   time.Time        `json:"redeemed_at"`
	FulfilledAt  *time.Time       `json:"fulfilled_at,omitempty"`
	Reward       *Reward          `json:"reward,omitempty"`
}

type CreateRewardRequest struct {
	Name        string  `json:"name" validate:"required"`
	Description *string `json:"description,omitempty"`
	CostCoins   int     `json:"cost_coins" validate:"required,min=1"`
	Stock       *int    `json:"stock,omitempty"`
}

type RedeemRequest struct {
	RewardID string `json:"reward_id" validate:"required"`
}
