package service

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hustlers/motivator-backend/internal/model"
	"github.com/hustlers/motivator-backend/internal/repository"
)

type MembershipService struct {
	pool    *pgxpool.Pool
	members repository.MembershipRepository
	badges  repository.BadgeRepository
}

func NewMembershipService(pool *pgxpool.Pool, members repository.MembershipRepository, badges repository.BadgeRepository) *MembershipService {
	return &MembershipService{pool: pool, members: members, badges: badges}
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

func (s *MembershipService) AddToCompany(ctx context.Context, userID, companyID string, role model.Role, displayName *string) (*model.Membership, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			log.Printf("rollback error: %v", err)
		}
	}()

	m := &model.Membership{
		UserID:      userID,
		CompanyID:   companyID,
		Role:        role,
		DisplayName: displayName,
	}
	if err := s.members.Create(ctx, tx, m); err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *MembershipService) Deactivate(ctx context.Context, id string) error {
	return s.members.Deactivate(ctx, id)
}
