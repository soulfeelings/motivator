package service

import (
	"context"
	"fmt"
	"time"

	"github.com/hustlers/motivator-backend/internal/model"
	"github.com/hustlers/motivator-backend/internal/repository"
)

type TournamentService struct {
	tournaments repository.TournamentRepository
	members     repository.MembershipRepository
}

func NewTournamentService(tournaments repository.TournamentRepository, members repository.MembershipRepository) *TournamentService {
	return &TournamentService{tournaments: tournaments, members: members}
}

func (s *TournamentService) Create(ctx context.Context, companyID string, req model.CreateTournamentRequest) (*model.Tournament, error) {
	startsAt, err := time.Parse(time.RFC3339, req.StartsAt)
	if err != nil {
		return nil, fmt.Errorf("invalid starts_at format, use RFC3339")
	}
	endsAt, err := time.Parse(time.RFC3339, req.EndsAt)
	if err != nil {
		return nil, fmt.Errorf("invalid ends_at format, use RFC3339")
	}

	t := &model.Tournament{
		CompanyID:       companyID,
		Name:            req.Name,
		Description:     req.Description,
		Season:          req.Season,
		Metric:          req.Metric,
		PrizePool:       req.PrizePool,
		XPPrizes:        req.XPPrizes,
		CoinPrizes:      req.CoinPrizes,
		MaxParticipants: req.MaxParticipants,
		StartsAt:        startsAt,
		EndsAt:          endsAt,
	}
	if len(t.XPPrizes) == 0 {
		t.XPPrizes = []int{100, 50, 25}
	}
	if len(t.CoinPrizes) == 0 {
		t.CoinPrizes = []int{50, 25, 10}
	}
	if err := s.tournaments.Create(ctx, t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *TournamentService) GetByID(ctx context.Context, id string) (*model.Tournament, error) {
	return s.tournaments.GetByID(ctx, id)
}

func (s *TournamentService) ListByCompany(ctx context.Context, companyID string) ([]model.Tournament, error) {
	return s.tournaments.ListByCompany(ctx, companyID)
}

func (s *TournamentService) UpdateStatus(ctx context.Context, id string, status model.TournamentStatus) error {
	return s.tournaments.UpdateStatus(ctx, id, status)
}

func (s *TournamentService) Delete(ctx context.Context, id string) error {
	return s.tournaments.Delete(ctx, id)
}

func (s *TournamentService) Join(ctx context.Context, tournamentID, membershipID string) (*model.TournamentParticipant, error) {
	t, err := s.tournaments.GetByID(ctx, tournamentID)
	if err != nil {
		return nil, err
	}
	if t.Status != model.TournamentRegistration && t.Status != model.TournamentActive {
		return nil, fmt.Errorf("tournament is not accepting participants")
	}
	if t.MaxParticipants != nil && t.ParticipantCount >= *t.MaxParticipants {
		return nil, fmt.Errorf("tournament is full")
	}

	p := &model.TournamentParticipant{TournamentID: tournamentID, MembershipID: membershipID}
	if err := s.tournaments.Join(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *TournamentService) Leave(ctx context.Context, tournamentID, membershipID string) error {
	return s.tournaments.Leave(ctx, tournamentID, membershipID)
}

func (s *TournamentService) SubmitScore(ctx context.Context, tournamentID, membershipID string, score int) error {
	return s.tournaments.UpdateScore(ctx, tournamentID, membershipID, score)
}

func (s *TournamentService) GetStandings(ctx context.Context, tournamentID string) ([]model.TournamentParticipant, error) {
	return s.tournaments.GetStandings(ctx, tournamentID)
}

func (s *TournamentService) Complete(ctx context.Context, tournamentID string) error {
	if err := s.tournaments.UpdateRanks(ctx, tournamentID); err != nil {
		return err
	}

	t, err := s.tournaments.GetByID(ctx, tournamentID)
	if err != nil {
		return err
	}

	standings, err := s.tournaments.GetStandings(ctx, tournamentID)
	if err != nil {
		return err
	}

	// Award prizes to top 3
	for i, p := range standings {
		if i >= 3 {
			break
		}
		if i < len(t.XPPrizes) && t.XPPrizes[i] > 0 {
			s.members.AwardXP(ctx, p.MembershipID, t.XPPrizes[i])
		}
	}

	return s.tournaments.Complete(ctx, tournamentID)
}
