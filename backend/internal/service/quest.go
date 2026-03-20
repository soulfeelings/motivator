package service

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/hustlers/motivator-backend/internal/model"
	"github.com/hustlers/motivator-backend/internal/repository"
)

type QuestService struct {
	quests  repository.QuestRepository
	members repository.MembershipRepository
}

func NewQuestService(quests repository.QuestRepository, members repository.MembershipRepository) *QuestService {
	return &QuestService{quests: quests, members: members}
}

func (s *QuestService) Create(ctx context.Context, companyID string, req model.CreateQuestRequest) (*model.Quest, error) {
	deadline := time.Now().Add(3 * 24 * time.Hour)
	if req.Deadline != "" {
		if d, err := time.Parse(time.RFC3339, req.Deadline); err == nil {
			deadline = d
		}
	}
	q := &model.Quest{
		CompanyID:  companyID,
		Name:       req.Name,
		XPReward:   req.XPReward,
		CoinReward: req.CoinReward,
		BonusXP:    req.BonusXP,
		BonusCoins: req.BonusCoins,
		Deadline:   deadline,
	}
	if q.Name == "" {
		q.Name = "Secret Motivator"
	}
	if q.XPReward == 0 {
		q.XPReward = 25
	}
	if q.CoinReward == 0 {
		q.CoinReward = 10
	}
	if q.BonusXP == 0 {
		q.BonusXP = 50
	}
	if q.BonusCoins == 0 {
		q.BonusCoins = 25
	}
	if err := s.quests.Create(ctx, q); err != nil {
		return nil, err
	}
	return q, nil
}

func (s *QuestService) GetByID(ctx context.Context, id string) (*model.Quest, error) {
	return s.quests.GetByID(ctx, id)
}

func (s *QuestService) ListByCompany(ctx context.Context, companyID string) ([]model.Quest, error) {
	return s.quests.ListByCompany(ctx, companyID)
}

func (s *QuestService) Delete(ctx context.Context, id string) error {
	return s.quests.Delete(ctx, id)
}

// Start activates the quest and creates random pairings.
func (s *QuestService) Start(ctx context.Context, questID, companyID string) error {
	q, err := s.quests.GetByID(ctx, questID)
	if err != nil {
		return err
	}
	if q.Status != model.QuestDraft {
		return fmt.Errorf("quest is not in draft status")
	}

	// Get all active members
	page := model.PaginationRequest{Page: 1, PerPage: 100}
	allMembers, _, err := s.members.ListByCompany(ctx, companyID, page)
	if err != nil {
		return err
	}
	if len(allMembers) < 2 {
		return fmt.Errorf("need at least 2 members to start a quest")
	}

	// Shuffle and create circular pairs: A→B, B→C, C→A
	shuffled := make([]model.Membership, len(allMembers))
	copy(shuffled, allMembers)
	rand.Shuffle(len(shuffled), func(i, j int) { shuffled[i], shuffled[j] = shuffled[j], shuffled[i] })

	for i := range shuffled {
		receiverIdx := (i + 1) % len(shuffled)
		pair := &model.QuestPair{
			QuestID:    questID,
			SenderID:   shuffled[i].ID,
			ReceiverID: shuffled[receiverIdx].ID,
		}
		if err := s.quests.CreatePair(ctx, pair); err != nil {
			return err
		}
	}

	return s.quests.UpdateStatus(ctx, questID, model.QuestActive)
}

// GetMyTarget returns who the current member needs to send a message to.
func (s *QuestService) GetMyTarget(ctx context.Context, questID, membershipID string) (*model.QuestPair, error) {
	return s.quests.GetPairBySender(ctx, questID, membershipID)
}

// SendMessage sends the anonymous message.
func (s *QuestService) SendMessage(ctx context.Context, questID, membershipID, message string) error {
	pair, err := s.quests.GetPairBySender(ctx, questID, membershipID)
	if err != nil {
		return err
	}
	if pair.SentAt != nil {
		return fmt.Errorf("message already sent")
	}
	return s.quests.SendMessage(ctx, pair.ID, message)
}

// GetReceivedMessages returns messages received by a member (anonymous until revealed).
func (s *QuestService) GetReceivedMessages(ctx context.Context, questID, membershipID string) ([]model.ReceivedMessage, error) {
	q, err := s.quests.GetByID(ctx, questID)
	if err != nil {
		return nil, err
	}
	revealed := q.Status == model.QuestRevealed || q.Status == model.QuestCompleted
	return s.quests.GetReceivedMessages(ctx, questID, membershipID, revealed)
}

// StartVoting transitions to voting phase.
func (s *QuestService) StartVoting(ctx context.Context, questID string) error {
	return s.quests.UpdateStatus(ctx, questID, model.QuestVoting)
}

// Vote for the best message.
func (s *QuestService) Vote(ctx context.Context, questID, membershipID, pairID string) error {
	v := &model.QuestVote{QuestID: questID, PairID: pairID, VoterID: membershipID}
	return s.quests.Vote(ctx, v)
}

// Reveal shows who sent what.
func (s *QuestService) Reveal(ctx context.Context, questID string) error {
	return s.quests.UpdateStatus(ctx, questID, model.QuestRevealed)
}

// Complete awards everyone who participated + bonus to winner.
func (s *QuestService) Complete(ctx context.Context, questID string) error {
	q, err := s.quests.GetByID(ctx, questID)
	if err != nil {
		return err
	}

	pairs, err := s.quests.ListPairs(ctx, questID)
	if err != nil {
		return err
	}

	// Award participation XP/coins to everyone who sent a message
	for _, p := range pairs {
		if p.SentAt != nil {
			s.members.AwardXP(ctx, p.SenderID, q.XPReward)
		}
	}

	// Award bonus to the winner (most votes)
	winner, err := s.quests.GetWinnerPair(ctx, questID)
	if err == nil && winner != nil {
		s.members.AwardXP(ctx, winner.SenderID, q.BonusXP)
	}

	return s.quests.Complete(ctx, questID)
}

// ListPairs returns all pairs (admin view).
func (s *QuestService) ListPairs(ctx context.Context, questID string) ([]model.QuestPair, error) {
	return s.quests.ListPairs(ctx, questID)
}
