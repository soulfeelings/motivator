package model

import "time"

type QuestStatus string

const (
	QuestDraft     QuestStatus = "draft"
	QuestActive    QuestStatus = "active"
	QuestVoting    QuestStatus = "voting"
	QuestRevealed  QuestStatus = "revealed"
	QuestCompleted QuestStatus = "completed"
)

type Quest struct {
	ID          string      `json:"id"`
	CompanyID   string      `json:"company_id"`
	Name        string      `json:"name"`
	Description *string     `json:"description,omitempty"`
	Status      QuestStatus `json:"status"`
	XPReward    int         `json:"xp_reward"`
	CoinReward  int         `json:"coin_reward"`
	BonusXP     int         `json:"bonus_xp"`
	BonusCoins  int         `json:"bonus_coins"`
	Deadline    time.Time   `json:"deadline"`
	RevealAt    *time.Time  `json:"reveal_at,omitempty"`
	CreatedAt   time.Time   `json:"created_at"`
	CompletedAt *time.Time  `json:"completed_at,omitempty"`
	PairCount   int         `json:"pair_count,omitempty"`
	SentCount   int         `json:"sent_count,omitempty"`
}

type QuestPair struct {
	ID         string     `json:"id"`
	QuestID    string     `json:"quest_id"`
	SenderID   string     `json:"sender_id"`
	ReceiverID string     `json:"receiver_id"`
	Message    *string    `json:"message,omitempty"`
	SentAt     *time.Time `json:"sent_at,omitempty"`
	VoteCount  int        `json:"vote_count,omitempty"`
}

type QuestVote struct {
	ID      string `json:"id"`
	QuestID string `json:"quest_id"`
	PairID  string `json:"pair_id"`
	VoterID string `json:"voter_id"`
}

type CreateQuestRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	XPReward    int     `json:"xp_reward"`
	CoinReward  int     `json:"coin_reward"`
	BonusXP     int     `json:"bonus_xp"`
	BonusCoins  int     `json:"bonus_coins"`
	Deadline    string  `json:"deadline"`
}

type SendMessageRequest struct {
	Message string `json:"message" validate:"required"`
}

type VoteRequest struct {
	PairID string `json:"pair_id" validate:"required"`
}

// ReceivedMessage is what a receiver sees (anonymous before reveal).
type ReceivedMessage struct {
	PairID   string     `json:"pair_id"`
	Message  string     `json:"message"`
	SentAt   time.Time  `json:"sent_at"`
	SenderID *string    `json:"sender_id,omitempty"` // nil until revealed
}
