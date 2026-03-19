package service

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hustlers/motivator-backend/internal/model"
	"github.com/hustlers/motivator-backend/internal/repository"
)

type BadgeService struct {
	pool    *pgxpool.Pool
	badges  repository.BadgeRepository
	members repository.MembershipRepository
}

func NewBadgeService(pool *pgxpool.Pool, badges repository.BadgeRepository, members repository.MembershipRepository) *BadgeService {
	return &BadgeService{pool: pool, badges: badges, members: members}
}

func (s *BadgeService) Create(ctx context.Context, companyID string, req model.CreateBadgeRequest) (*model.Badge, error) {
	badge := &model.Badge{
		CompanyID:   companyID,
		Name:        req.Name,
		Description: req.Description,
		IconURL:     req.IconURL,
		XPReward:    req.XPReward,
		CoinReward:  req.CoinReward,
	}
	if err := s.badges.Create(ctx, badge); err != nil {
		return nil, err
	}
	return badge, nil
}

func (s *BadgeService) ListByCompany(ctx context.Context, companyID string) ([]model.Badge, error) {
	return s.badges.ListByCompany(ctx, companyID)
}

func (s *BadgeService) Delete(ctx context.Context, id string) error {
	return s.badges.Delete(ctx, id)
}

func (s *BadgeService) AwardBadge(ctx context.Context, membershipID, badgeID string) (*model.MemberBadge, error) {
	badge, err := s.badges.GetByID(ctx, badgeID)
	if err != nil {
		return nil, err
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			log.Printf("rollback error: %v", err)
		}
	}()

	mb, err := s.badges.AwardBadge(ctx, tx, membershipID, badgeID)
	if err != nil {
		return nil, err
	}

	// Award XP and coins from badge
	if badge.XPReward > 0 || badge.CoinReward > 0 {
		if badge.CoinReward > 0 {
			if err := s.members.AwardCoins(ctx, tx, membershipID, badge.CoinReward); err != nil {
				return nil, err
			}
		}
		// XP awarded outside tx since AwardXP uses pool directly
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	// Award XP after commit (uses pool, updates level too)
	if badge.XPReward > 0 {
		if _, err := s.members.AwardXP(ctx, membershipID, badge.XPReward); err != nil {
			log.Printf("failed to award XP for badge: %v", err)
		}
	}

	mb.Badge = badge
	return mb, nil
}

func (s *BadgeService) ListMemberBadges(ctx context.Context, membershipID string) ([]model.MemberBadge, error) {
	return s.badges.ListMemberBadges(ctx, membershipID)
}
