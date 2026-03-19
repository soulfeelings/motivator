package model

import "time"

type MetricOperator string

const (
	OpGTE MetricOperator = "gte"
	OpLTE MetricOperator = "lte"
	OpEQ  MetricOperator = "eq"
	OpGT  MetricOperator = "gt"
	OpLT  MetricOperator = "lt"
)

type Achievement struct {
	ID          string         `json:"id"`
	CompanyID   string         `json:"company_id"`
	Name        string         `json:"name"`
	Description *string        `json:"description,omitempty"`
	Metric      string         `json:"metric"`
	Operator    MetricOperator `json:"operator"`
	Threshold   int            `json:"threshold"`
	BadgeID     *string        `json:"badge_id,omitempty"`
	XPReward    int            `json:"xp_reward"`
	CoinReward  int            `json:"coin_reward"`
	IsActive    bool           `json:"is_active"`
	CreatedAt   time.Time      `json:"created_at"`
}

type MemberAchievement struct {
	ID            string       `json:"id"`
	MembershipID  string       `json:"membership_id"`
	AchievementID string       `json:"achievement_id"`
	CompletedAt   time.Time    `json:"completed_at"`
	Achievement   *Achievement `json:"achievement,omitempty"`
}

type CreateAchievementRequest struct {
	Name        string         `json:"name" validate:"required"`
	Description *string        `json:"description,omitempty"`
	Metric      string         `json:"metric" validate:"required"`
	Operator    MetricOperator `json:"operator" validate:"required"`
	Threshold   int            `json:"threshold" validate:"required"`
	BadgeID     *string        `json:"badge_id,omitempty"`
	XPReward    int            `json:"xp_reward"`
	CoinReward  int            `json:"coin_reward"`
}

// EvaluateMetricRequest is sent by external systems to report a metric value.
type EvaluateMetricRequest struct {
	Metric string `json:"metric" validate:"required"`
	Value  int    `json:"value" validate:"required"`
}
