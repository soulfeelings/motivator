package service

import (
	"context"
	"fmt"

	"github.com/hustlers/motivator-backend/internal/model"
	"github.com/hustlers/motivator-backend/internal/repository"
)

type ChallengeService struct {
	challenges repository.ChallengeRepository
	members    repository.MembershipRepository
}

func NewChallengeService(challenges repository.ChallengeRepository, members repository.MembershipRepository) *ChallengeService {
	return &ChallengeService{challenges: challenges, members: members}
}

func (s *ChallengeService) Create(ctx context.Context, companyID, challengerID string, req model.CreateChallengeRequest) (*model.Challenge, error) {
	if challengerID == req.OpponentID {
		return nil, fmt.Errorf("cannot challenge yourself")
	}

	c := &model.Challenge{
		CompanyID:    companyID,
		ChallengerID: challengerID,
		OpponentID:   req.OpponentID,
		Metric:       req.Metric,
		Target:       req.Target,
		Wager:        req.Wager,
		XPReward:     req.XPReward,
	}
	if c.XPReward == 0 {
		c.XPReward = 50
	}
	if err := s.challenges.Create(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *ChallengeService) Accept(ctx context.Context, id, memberID string) error {
	c, err := s.challenges.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if c.OpponentID != memberID {
		return model.ErrForbidden
	}
	if c.Status != model.ChallengePending {
		return fmt.Errorf("challenge is not pending")
	}
	return s.challenges.UpdateStatus(ctx, id, model.ChallengeActive)
}

func (s *ChallengeService) Decline(ctx context.Context, id, memberID string) error {
	c, err := s.challenges.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if c.OpponentID != memberID {
		return model.ErrForbidden
	}
	if c.Status != model.ChallengePending {
		return fmt.Errorf("challenge is not pending")
	}
	return s.challenges.UpdateStatus(ctx, id, model.ChallengeDeclined)
}

func (s *ChallengeService) ReportScore(ctx context.Context, id, memberID string, score int) (*model.Challenge, error) {
	c, err := s.challenges.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if c.Status != model.ChallengeActive {
		return nil, fmt.Errorf("challenge is not active")
	}

	isChallenger := c.ChallengerID == memberID
	isOpponent := c.OpponentID == memberID
	if !isChallenger && !isOpponent {
		return nil, model.ErrForbidden
	}

	if err := s.challenges.UpdateScore(ctx, id, isChallenger, score); err != nil {
		return nil, err
	}

	// Refresh and check if both hit target
	c, _ = s.challenges.GetByID(ctx, id)
	if isChallenger {
		c.ChallengerScore = score
	} else {
		c.OpponentScore = score
	}

	// Auto-complete if either reaches target
	if c.ChallengerScore >= c.Target || c.OpponentScore >= c.Target {
		winnerID := c.ChallengerID
		if c.OpponentScore > c.ChallengerScore {
			winnerID = c.OpponentID
		}
		s.challenges.Complete(ctx, id, winnerID)
		c.Status = model.ChallengeCompleted
		c.WinnerID = &winnerID

		// Award XP to winner
		if c.XPReward > 0 {
			s.members.AwardXP(ctx, winnerID, c.XPReward)
		}
	}

	return c, nil
}

func (s *ChallengeService) ListByCompany(ctx context.Context, companyID string, status *model.ChallengeStatus) ([]model.Challenge, error) {
	return s.challenges.ListByCompany(ctx, companyID, status)
}

func (s *ChallengeService) ListByMember(ctx context.Context, membershipID string) ([]model.Challenge, error) {
	return s.challenges.ListByMember(ctx, membershipID)
}
