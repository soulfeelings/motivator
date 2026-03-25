package model

import (
	"encoding/json"
	"time"
)

type GameType string

const (
	GameTrivia         GameType = "trivia"
	GamePhotoChallenge GameType = "photo_challenge"
	GameTwoTruths      GameType = "two_truths"
)

type GameStatus string

const (
	GameDraft     GameStatus = "draft"
	GameActive    GameStatus = "active"
	GameVoting    GameStatus = "voting"
	GameCompleted GameStatus = "completed"
)

type SocialGame struct {
	ID          string          `json:"id"`
	CompanyID   string          `json:"company_id"`
	GameType    GameType        `json:"game_type"`
	Name        string          `json:"name"`
	Description *string         `json:"description,omitempty"`
	Status      GameStatus      `json:"status"`
	Config      json.RawMessage `json:"config"`
	XPReward    int             `json:"xp_reward"`
	CoinReward  int             `json:"coin_reward"`
	StartsAt    time.Time       `json:"starts_at"`
	EndsAt      time.Time       `json:"ends_at"`
	CreatedBy   string          `json:"created_by"`
	CreatedAt   time.Time       `json:"created_at"`
	CompletedAt *time.Time      `json:"completed_at,omitempty"`
}

type SocialGameQuestion struct {
	ID           string          `json:"id"`
	GameID       string          `json:"game_id"`
	Question     string          `json:"question"`
	Options      json.RawMessage `json:"options"`
	CorrectIndex int             `json:"correct_index"`
	SortOrder    int             `json:"sort_order"`
}

type SocialGameAnswer struct {
	ID            string    `json:"id"`
	GameID        string    `json:"game_id"`
	QuestionID    string    `json:"question_id"`
	MemberID      string    `json:"member_id"`
	SelectedIndex int       `json:"selected_index"`
	IsCorrect     bool      `json:"is_correct"`
	AnsweredAt    time.Time `json:"answered_at"`
}

type SocialGameSubmission struct {
	ID          string          `json:"id"`
	GameID      string          `json:"game_id"`
	MemberID    string          `json:"member_id"`
	Content     *string         `json:"content,omitempty"`
	Statements  json.RawMessage `json:"statements,omitempty"`
	LieIndex    *int            `json:"lie_index,omitempty"`
	SubmittedAt time.Time       `json:"submitted_at"`
	VoteCount   int             `json:"vote_count,omitempty"`
}

type SocialGameVote struct {
	ID           string    `json:"id"`
	GameID       string    `json:"game_id"`
	SubmissionID string    `json:"submission_id"`
	VoterID      string    `json:"voter_id"`
	VoteValue    *int      `json:"vote_value,omitempty"`
	VotedAt      time.Time `json:"voted_at"`
}

// Request types

type CreateSocialGameRequest struct {
	GameType    GameType        `json:"game_type"`
	Name        string          `json:"name"`
	Description *string         `json:"description,omitempty"`
	Config      json.RawMessage `json:"config,omitempty"`
	XPReward    int             `json:"xp_reward"`
	CoinReward  int             `json:"coin_reward"`
	StartsAt    string          `json:"starts_at"`
	EndsAt      string          `json:"ends_at"`
}

type AddQuestionRequest struct {
	Question     string          `json:"question"`
	Options      json.RawMessage `json:"options"`
	CorrectIndex int             `json:"correct_index"`
	SortOrder    int             `json:"sort_order"`
}

type SubmitAnswerRequest struct {
	QuestionID    string `json:"question_id" validate:"required"`
	SelectedIndex int    `json:"selected_index"`
}

type SubmitEntryRequest struct {
	Content    *string         `json:"content,omitempty"`
	Statements json.RawMessage `json:"statements,omitempty"`
	LieIndex   *int            `json:"lie_index,omitempty"`
}

type CastVoteRequest struct {
	SubmissionID string `json:"submission_id" validate:"required"`
	VoteValue    *int   `json:"vote_value,omitempty"`
}

// Results types

type SocialGameResults struct {
	Game              SocialGame              `json:"game"`
	ParticipantCount  int                     `json:"participant_count"`
	TotalMembers      int                     `json:"total_members"`
	ParticipationRate float64                 `json:"participation_rate"`
	Leaderboard       []LeaderboardEntry      `json:"leaderboard"`
	Winner            *LeaderboardEntry       `json:"winner,omitempty"`
}

type LeaderboardEntry struct {
	MemberID string `json:"member_id"`
	Score    int    `json:"score"`
}
