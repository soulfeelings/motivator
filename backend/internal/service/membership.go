package service

import (
	"context"

	"github.com/hustlers/motivator-backend/internal/model"
	"github.com/hustlers/motivator-backend/internal/repository"
)

type MembershipService struct {
	members repository.MembershipRepository
	badges  repository.BadgeRepository
}

func NewMembershipService(members repository.MembershipRepository, badges repository.BadgeRepository) *MembershipService {
	return &MembershipService{members: members, badges: badges}
}

func (s *MembershipService) GetByID(ctx context.Context, id string) (*model.Membership, error) {
	return s.members.GetByID(ctx, id)
}

func (s *MembershipService) GetByUserAndCompany(ctx context.Context, userID, companyID string) (*model.Membership, error) {
	return s.members.GetByUserAndCompany(ctx, userID, companyID)
}

func (s *MembershipService) ListByCompany(ctx context.Context, companyID string, page model.PaginationRequest) ([]model.Membership, int, error) {
	return s.members.ListByCompany(ctx, companyID, page)
}

func (s *MembershipService) ListByUser(ctx context.Context, userID string) ([]model.Membership, error) {
	return s.members.ListByUser(ctx, userID)
}

func (s *MembershipService) Update(ctx context.Context, id string, req model.UpdateMembershipRequest) (*model.Membership, error) {
	return s.members.Update(ctx, id, req)
}

func (s *MembershipService) AwardXP(ctx context.Context, id string, amount int) (*model.Membership, error) {
	return s.members.AwardXP(ctx, id, amount)
}

func (s *MembershipService) GetProfile(ctx context.Context, id string) (*model.ProfileResponse, error) {
	m, err := s.members.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	badges, err := s.badges.ListMemberBadges(ctx, id)
	if err != nil {
		return nil, err
	}
	badgeList := make([]model.Badge, 0, len(badges))
	for _, mb := range badges {
		if mb.Badge != nil {
			badgeList = append(badgeList, *mb.Badge)
		}
	}
	return &model.ProfileResponse{Membership: *m, Badges: badgeList}, nil
}

func (s *MembershipService) Deactivate(ctx context.Context, id string) error {
	return s.members.Deactivate(ctx, id)
}
