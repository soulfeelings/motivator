package service

import (
	"context"

	"github.com/hustlers/motivator-backend/internal/model"
	"github.com/hustlers/motivator-backend/internal/repository"
)

type AnalyticsService struct {
	analytics repository.AnalyticsRepository
}

func NewAnalyticsService(analytics repository.AnalyticsRepository) *AnalyticsService {
	return &AnalyticsService{analytics: analytics}
}

func (s *AnalyticsService) GetDashboard(ctx context.Context, companyID string) (*model.AnalyticsDashboard, error) {
	overview, err := s.analytics.GetOverview(ctx, companyID)
	if err != nil {
		return nil, err
	}
	topPerformers, err := s.analytics.GetTopPerformers(ctx, companyID, 10)
	if err != nil {
		return nil, err
	}
	achievementStats, err := s.analytics.GetAchievementStats(ctx, companyID)
	if err != nil {
		return nil, err
	}
	challengeStats, err := s.analytics.GetChallengeStats(ctx, companyID)
	if err != nil {
		return nil, err
	}
	rewardStats, err := s.analytics.GetRewardStats(ctx, companyID)
	if err != nil {
		return nil, err
	}
	xpDist, err := s.analytics.GetXPDistribution(ctx, companyID)
	if err != nil {
		return nil, err
	}

	return &model.AnalyticsDashboard{
		Overview:         *overview,
		TopPerformers:    topPerformers,
		AchievementStats: achievementStats,
		ChallengeStats:   *challengeStats,
		RewardStats:      rewardStats,
		XPDistribution:   xpDist,
	}, nil
}
