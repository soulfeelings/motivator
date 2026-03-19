package model

import "time"

type Team struct {
	ID          string    `json:"id"`
	CompanyID   string    `json:"company_id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	Color       string    `json:"color"`
	CreatedAt   time.Time `json:"created_at"`
	MemberCount int       `json:"member_count,omitempty"`
}

type TeamMember struct {
	ID           string    `json:"id"`
	TeamID       string    `json:"team_id"`
	MembershipID string    `json:"membership_id"`
	JoinedAt     time.Time `json:"joined_at"`
}

type TeamBattleStatus string

const (
	TeamBattlePending   TeamBattleStatus = "pending"
	TeamBattleActive    TeamBattleStatus = "active"
	TeamBattleCompleted TeamBattleStatus = "completed"
)

type TeamBattle struct {
	ID          string           `json:"id"`
	CompanyID   string           `json:"company_id"`
	TeamAID     string           `json:"team_a_id"`
	TeamBID     string           `json:"team_b_id"`
	Metric      string           `json:"metric"`
	Target      int              `json:"target"`
	TeamAScore  int              `json:"team_a_score"`
	TeamBScore  int              `json:"team_b_score"`
	Status      TeamBattleStatus `json:"status"`
	WinnerID    *string          `json:"winner_id,omitempty"`
	XPReward    int              `json:"xp_reward"`
	CoinReward  int              `json:"coin_reward"`
	Deadline    time.Time        `json:"deadline"`
	CreatedAt   time.Time        `json:"created_at"`
	CompletedAt *time.Time       `json:"completed_at,omitempty"`
}

type CreateTeamRequest struct {
	Name        string  `json:"name" validate:"required"`
	Description *string `json:"description,omitempty"`
	Color       string  `json:"color"`
}

type AddTeamMemberRequest struct {
	MembershipID string `json:"membership_id" validate:"required"`
}

type CreateTeamBattleRequest struct {
	TeamAID    string `json:"team_a_id" validate:"required"`
	TeamBID    string `json:"team_b_id" validate:"required"`
	Metric     string `json:"metric" validate:"required"`
	Target     int    `json:"target"`
	XPReward   int    `json:"xp_reward"`
	CoinReward int    `json:"coin_reward"`
}

type ReportTeamScoreRequest struct {
	TeamID string `json:"team_id" validate:"required"`
	Score  int    `json:"score" validate:"required"`
}
