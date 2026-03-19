package model

import "time"

type TournamentStatus string

const (
	TournamentDraft        TournamentStatus = "draft"
	TournamentRegistration TournamentStatus = "registration"
	TournamentActive       TournamentStatus = "active"
	TournamentCompleted    TournamentStatus = "completed"
	TournamentCancelled    TournamentStatus = "cancelled"
)

type Tournament struct {
	ID              string           `json:"id"`
	CompanyID       string           `json:"company_id"`
	Name            string           `json:"name"`
	Description     *string          `json:"description,omitempty"`
	Season          *string          `json:"season,omitempty"`
	Metric          string           `json:"metric"`
	PrizePool       int              `json:"prize_pool"`
	XPPrizes        []int            `json:"xp_prizes"`
	CoinPrizes      []int            `json:"coin_prizes"`
	Status          TournamentStatus `json:"status"`
	MaxParticipants *int             `json:"max_participants,omitempty"`
	StartsAt        time.Time        `json:"starts_at"`
	EndsAt          time.Time        `json:"ends_at"`
	CreatedAt       time.Time        `json:"created_at"`
	CompletedAt     *time.Time       `json:"completed_at,omitempty"`
	ParticipantCount int             `json:"participant_count,omitempty"`
}

type TournamentParticipant struct {
	ID           string    `json:"id"`
	TournamentID string    `json:"tournament_id"`
	MembershipID string    `json:"membership_id"`
	Score        int       `json:"score"`
	Rank         *int      `json:"rank,omitempty"`
	JoinedAt     time.Time `json:"joined_at"`
}

type CreateTournamentRequest struct {
	Name            string   `json:"name" validate:"required"`
	Description     *string  `json:"description,omitempty"`
	Season          *string  `json:"season,omitempty"`
	Metric          string   `json:"metric" validate:"required"`
	PrizePool       int      `json:"prize_pool"`
	XPPrizes        []int    `json:"xp_prizes"`
	CoinPrizes      []int    `json:"coin_prizes"`
	MaxParticipants *int     `json:"max_participants,omitempty"`
	StartsAt        string   `json:"starts_at" validate:"required"`
	EndsAt          string   `json:"ends_at" validate:"required"`
}

type SubmitScoreRequest struct {
	Score int `json:"score" validate:"required"`
}
