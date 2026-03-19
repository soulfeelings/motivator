package model

import "time"

type ChallengeStatus string

const (
	ChallengePending   ChallengeStatus = "pending"
	ChallengeActive    ChallengeStatus = "active"
	ChallengeCompleted ChallengeStatus = "completed"
	ChallengeDeclined  ChallengeStatus = "declined"
	ChallengeCancelled ChallengeStatus = "cancelled"
)

type Challenge struct {
	ID              string          `json:"id"`
	CompanyID       string          `json:"company_id"`
	ChallengerID    string          `json:"challenger_id"`
	OpponentID      string          `json:"opponent_id"`
	Metric          string          `json:"metric"`
	Target          int             `json:"target"`
	Wager           int             `json:"wager"`
	Status          ChallengeStatus `json:"status"`
	ChallengerScore int             `json:"challenger_score"`
	OpponentScore   int             `json:"opponent_score"`
	WinnerID        *string         `json:"winner_id,omitempty"`
	XPReward        int             `json:"xp_reward"`
	Deadline        time.Time       `json:"deadline"`
	CreatedAt       time.Time       `json:"created_at"`
	CompletedAt     *time.Time      `json:"completed_at,omitempty"`
}

type CreateChallengeRequest struct {
	OpponentID string `json:"opponent_id" validate:"required"`
	Metric     string `json:"metric" validate:"required"`
	Target     int    `json:"target" validate:"required,min=1"`
	Wager      int    `json:"wager"`
	XPReward   int    `json:"xp_reward"`
}

type ReportScoreRequest struct {
	Score int `json:"score" validate:"required"`
}
