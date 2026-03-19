package service

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hustlers/motivator-backend/internal/model"
	"github.com/hustlers/motivator-backend/internal/repository"
)

type AchievementService struct {
	pool         *pgxpool.Pool
	achievements repository.AchievementRepository
	members      repository.MembershipRepository
	badges       repository.BadgeRepository
	notifier     Notifier
}

func NewAchievementService(pool *pgxpool.Pool, achievements repository.AchievementRepository, members repository.MembershipRepository, badges repository.BadgeRepository, notifier Notifier) *AchievementService {
	return &AchievementService{pool: pool, achievements: achievements, members: members, badges: badges, notifier: notifier}
}

func (s *AchievementService) Create(ctx context.Context, companyID string, req model.CreateAchievementRequest) (*model.Achievement, error) {
	a := &model.Achievement{
		CompanyID:   companyID,
		Name:        req.Name,
		Description: req.Description,
		Metric:      req.Metric,
		Operator:    req.Operator,
		Threshold:   req.Threshold,
		BadgeID:     req.BadgeID,
		XPReward:    req.XPReward,
		CoinReward:  req.CoinReward,
	}
	if err := s.achievements.Create(ctx, a); err != nil {
		return nil, err
	}
	return a, nil
}

func (s *AchievementService) ListByCompany(ctx context.Context, companyID string) ([]model.Achievement, error) {
	return s.achievements.ListByCompany(ctx, companyID)
}

func (s *AchievementService) Delete(ctx context.Context, id string) error {
	return s.achievements.Delete(ctx, id)
}

func (s *AchievementService) ListMemberAchievements(ctx context.Context, membershipID string) ([]model.MemberAchievement, error) {
	return s.achievements.ListMemberAchievements(ctx, membershipID)
}

// EvaluateMetric checks all active achievements for a given metric and awards completed ones.
// Returns the list of newly completed achievements.
func (s *AchievementService) EvaluateMetric(ctx context.Context, membershipID, companyID, metric string, value int) ([]model.MemberAchievement, error) {
	achievements, err := s.achievements.ListActiveByMetric(ctx, companyID, metric)
	if err != nil {
		return nil, err
	}

	var completed []model.MemberAchievement
	for _, a := range achievements {
		if !evaluateCondition(value, a.Operator, a.Threshold) {
			continue
		}

		done, err := s.achievements.HasCompleted(ctx, membershipID, a.ID)
		if err != nil {
			log.Printf("error checking completion for achievement=%s member=%s: %v", a.ID, membershipID, err)
			continue
		}
		if done {
			continue
		}

		tx, err := s.pool.Begin(ctx)
		if err != nil {
			log.Printf("error starting tx for achievement=%s: %v", a.ID, err)
			continue
		}

		ma, err := s.achievements.RecordCompletion(ctx, tx, membershipID, a.ID)
		if err != nil {
			tx.Rollback(ctx)
			log.Printf("error recording achievement=%s: %v", a.ID, err)
			continue
		}

		if a.CoinReward > 0 {
			if err := s.members.AwardCoins(ctx, tx, membershipID, a.CoinReward); err != nil {
				tx.Rollback(ctx)
				log.Printf("error awarding coins for achievement=%s: %v", a.ID, err)
				continue
			}
		}

		if a.BadgeID != nil {
			if _, err := s.badges.AwardBadge(ctx, tx, membershipID, *a.BadgeID); err != nil {
				log.Printf("error awarding badge for achievement=%s: %v (may already have it)", a.ID, err)
				// Non-fatal: badge may already be awarded
			}
		}

		if err := tx.Commit(ctx); err != nil {
			log.Printf("error committing achievement=%s: %v", a.ID, err)
			continue
		}

		if a.XPReward > 0 {
			if _, err := s.members.AwardXP(ctx, membershipID, a.XPReward); err != nil {
				log.Printf("error awarding XP for achievement=%s: %v", a.ID, err)
			}
		}

		NotifyAchievementCompleted(ctx, s.notifier, membershipID, a.Name, a.XPReward, a.CoinReward)

		ma.Achievement = &a
		completed = append(completed, *ma)
	}

	return completed, nil
}

func evaluateCondition(value int, op model.MetricOperator, threshold int) bool {
	switch op {
	case model.OpGTE:
		return value >= threshold
	case model.OpLTE:
		return value <= threshold
	case model.OpEQ:
		return value == threshold
	case model.OpGT:
		return value > threshold
	case model.OpLT:
		return value < threshold
	default:
		return false
	}
}
