package service

import (
	"context"
	"fmt"

	"github.com/hustlers/motivator-backend/internal/model"
	"github.com/hustlers/motivator-backend/internal/repository"
)

type TeamService struct {
	teams   repository.TeamRepository
	members repository.MembershipRepository
}

func NewTeamService(teams repository.TeamRepository, members repository.MembershipRepository) *TeamService {
	return &TeamService{teams: teams, members: members}
}

func (s *TeamService) Create(ctx context.Context, companyID string, req model.CreateTeamRequest) (*model.Team, error) {
	t := &model.Team{
		CompanyID:   companyID,
		Name:        req.Name,
		Description: req.Description,
		Color:       req.Color,
	}
	if err := s.teams.Create(ctx, t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *TeamService) GetByID(ctx context.Context, id string) (*model.Team, error) {
	return s.teams.GetByID(ctx, id)
}

func (s *TeamService) ListByCompany(ctx context.Context, companyID string) ([]model.Team, error) {
	return s.teams.ListByCompany(ctx, companyID)
}

func (s *TeamService) Delete(ctx context.Context, id string) error {
	return s.teams.Delete(ctx, id)
}

func (s *TeamService) AddMember(ctx context.Context, teamID, membershipID string) (*model.TeamMember, error) {
	tm := &model.TeamMember{TeamID: teamID, MembershipID: membershipID}
	if err := s.teams.AddMember(ctx, tm); err != nil {
		return nil, err
	}
	return tm, nil
}

func (s *TeamService) RemoveMember(ctx context.Context, teamID, membershipID string) error {
	return s.teams.RemoveMember(ctx, teamID, membershipID)
}

func (s *TeamService) ListMembers(ctx context.Context, teamID string) ([]model.TeamMember, error) {
	return s.teams.ListMembers(ctx, teamID)
}

func (s *TeamService) CreateBattle(ctx context.Context, companyID string, req model.CreateTeamBattleRequest) (*model.TeamBattle, error) {
	if req.TeamAID == req.TeamBID {
		return nil, fmt.Errorf("cannot battle the same team")
	}
	b := &model.TeamBattle{
		CompanyID:  companyID,
		TeamAID:    req.TeamAID,
		TeamBID:    req.TeamBID,
		Metric:     req.Metric,
		Target:     req.Target,
		XPReward:   req.XPReward,
		CoinReward: req.CoinReward,
	}
	if b.XPReward == 0 {
		b.XPReward = 100
	}
	if b.CoinReward == 0 {
		b.CoinReward = 50
	}
	if err := s.teams.CreateBattle(ctx, b); err != nil {
		return nil, err
	}
	return b, nil
}

func (s *TeamService) ListBattles(ctx context.Context, companyID string) ([]model.TeamBattle, error) {
	return s.teams.ListBattles(ctx, companyID)
}

func (s *TeamService) ReportScore(ctx context.Context, battleID string, req model.ReportTeamScoreRequest) (*model.TeamBattle, error) {
	b, err := s.teams.GetBattle(ctx, battleID)
	if err != nil {
		return nil, err
	}
	if b.Status == model.TeamBattleCompleted {
		return nil, fmt.Errorf("battle already completed")
	}

	isTeamA := req.TeamID == b.TeamAID
	if !isTeamA && req.TeamID != b.TeamBID {
		return nil, model.ErrForbidden
	}

	if err := s.teams.UpdateBattleScore(ctx, battleID, isTeamA, req.Score); err != nil {
		return nil, err
	}

	// Activate if pending
	if b.Status == model.TeamBattlePending {
		s.teams.UpdateBattleScore(ctx, battleID, isTeamA, req.Score)
	}

	// Refresh
	b, _ = s.teams.GetBattle(ctx, battleID)

	// Auto-complete if target reached
	if b.Target > 0 && (b.TeamAScore >= b.Target || b.TeamBScore >= b.Target) {
		winnerID := b.TeamAID
		if b.TeamBScore > b.TeamAScore {
			winnerID = b.TeamBID
		}
		s.teams.CompleteBattle(ctx, battleID, winnerID)
		b.Status = model.TeamBattleCompleted
		b.WinnerID = &winnerID

		// Award XP+coins to winning team members
		s.awardTeamReward(ctx, winnerID, b.XPReward, b.CoinReward)
	}

	return b, nil
}

func (s *TeamService) awardTeamReward(ctx context.Context, teamID string, xp, coins int) {
	members, err := s.teams.ListMembers(ctx, teamID)
	if err != nil {
		return
	}
	for _, m := range members {
		if xp > 0 {
			s.members.AwardXP(ctx, m.MembershipID, xp)
		}
	}
}
