package service

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hustlers/motivator-backend/internal/model"
	"github.com/hustlers/motivator-backend/internal/repository"
)

type RewardService struct {
	pool    *pgxpool.Pool
	rewards repository.RewardRepository
	members repository.MembershipRepository
}

func NewRewardService(pool *pgxpool.Pool, rewards repository.RewardRepository, members repository.MembershipRepository) *RewardService {
	return &RewardService{pool: pool, rewards: rewards, members: members}
}

func (s *RewardService) Create(ctx context.Context, companyID string, req model.CreateRewardRequest) (*model.Reward, error) {
	rw := &model.Reward{
		CompanyID:   companyID,
		Name:        req.Name,
		Description: req.Description,
		CostCoins:   req.CostCoins,
		Stock:       req.Stock,
	}
	if err := s.rewards.Create(ctx, rw); err != nil {
		return nil, err
	}
	return rw, nil
}

func (s *RewardService) ListByCompany(ctx context.Context, companyID string) ([]model.Reward, error) {
	return s.rewards.ListByCompany(ctx, companyID)
}

func (s *RewardService) Delete(ctx context.Context, id string) error {
	return s.rewards.Delete(ctx, id)
}

func (s *RewardService) Redeem(ctx context.Context, membershipID, rewardID string) (*model.Redemption, error) {
	reward, err := s.rewards.GetByID(ctx, rewardID)
	if err != nil {
		return nil, err
	}
	if !reward.IsActive {
		return nil, fmt.Errorf("reward is not available")
	}
	if reward.Stock != nil && *reward.Stock <= 0 {
		return nil, fmt.Errorf("reward is out of stock")
	}

	member, err := s.members.GetByID(ctx, membershipID)
	if err != nil {
		return nil, err
	}
	if member.Coins < reward.CostCoins {
		return nil, fmt.Errorf("insufficient coins: have %d, need %d", member.Coins, reward.CostCoins)
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

	// Deduct coins
	if err := s.members.AwardCoins(ctx, tx, membershipID, -reward.CostCoins); err != nil {
		return nil, err
	}

	// Decrement stock
	if reward.Stock != nil {
		if err := s.rewards.DecrementStock(ctx, tx, rewardID); err != nil {
			return nil, err
		}
	}

	// Create redemption
	rd := &model.Redemption{
		MembershipID: membershipID,
		RewardID:     rewardID,
		CoinsSpent:   reward.CostCoins,
	}
	if err := s.rewards.CreateRedemption(ctx, tx, rd); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	rd.Reward = reward
	return rd, nil
}

func (s *RewardService) ListRedemptions(ctx context.Context, companyID string) ([]model.Redemption, error) {
	return s.rewards.ListRedemptions(ctx, companyID)
}

func (s *RewardService) UpdateRedemptionStatus(ctx context.Context, id string, status model.RedemptionStatus) error {
	return s.rewards.UpdateRedemptionStatus(ctx, id, status)
}
