package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hustlers/motivator-backend/internal/model"
)

type AnalyticsRepository interface {
	GetOverview(ctx context.Context, companyID string) (*model.AnalyticsOverview, error)
	GetTopPerformers(ctx context.Context, companyID string, limit int) ([]model.TopPerformer, error)
	GetAchievementStats(ctx context.Context, companyID string) ([]model.AchievementStat, error)
	GetChallengeStats(ctx context.Context, companyID string) (*model.ChallengeStat, error)
	GetRewardStats(ctx context.Context, companyID string) ([]model.RewardStat, error)
	GetXPDistribution(ctx context.Context, companyID string) ([]model.XPDistribution, error)
}

type analyticsRepo struct {
	pool *pgxpool.Pool
}

func NewAnalyticsRepository(pool *pgxpool.Pool) AnalyticsRepository {
	return &analyticsRepo{pool: pool}
}

func (r *analyticsRepo) GetOverview(ctx context.Context, companyID string) (*model.AnalyticsOverview, error) {
	o := &model.AnalyticsOverview{}

	r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM memberships WHERE company_id = $1 AND is_active = true`, companyID,
	).Scan(&o.TotalMembers)

	r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM memberships WHERE company_id = $1 AND is_active = true AND xp > 0`, companyID,
	).Scan(&o.ActiveMembers)

	r.pool.QueryRow(ctx,
		`SELECT COALESCE(SUM(xp), 0) FROM memberships WHERE company_id = $1`, companyID,
	).Scan(&o.TotalXPAwarded)

	r.pool.QueryRow(ctx,
		`SELECT COALESCE(SUM(coins), 0) FROM memberships WHERE company_id = $1`, companyID,
	).Scan(&o.TotalCoinsAwarded)

	r.pool.QueryRow(ctx,
		`SELECT COALESCE(SUM(rd.coins_spent), 0) FROM redemptions rd
		 JOIN rewards rw ON rw.id = rd.reward_id WHERE rw.company_id = $1`, companyID,
	).Scan(&o.TotalCoinsSpent)

	r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM member_badges mb
		 JOIN memberships m ON m.id = mb.membership_id WHERE m.company_id = $1`, companyID,
	).Scan(&o.TotalBadgesAwarded)

	r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM member_achievements ma
		 JOIN memberships m ON m.id = ma.membership_id WHERE m.company_id = $1`, companyID,
	).Scan(&o.TotalAchievements)

	r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM challenges WHERE company_id = $1`, companyID,
	).Scan(&o.TotalChallenges)

	r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM redemptions rd
		 JOIN rewards rw ON rw.id = rd.reward_id WHERE rw.company_id = $1`, companyID,
	).Scan(&o.TotalRedemptions)

	return o, nil
}

func (r *analyticsRepo) GetTopPerformers(ctx context.Context, companyID string, limit int) ([]model.TopPerformer, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT m.id, m.display_name, m.xp, m.level,
		        (SELECT COUNT(*) FROM member_badges WHERE membership_id = m.id),
		        (SELECT COUNT(*) FROM member_achievements WHERE membership_id = m.id)
		 FROM memberships m WHERE m.company_id = $1 AND m.is_active = true
		 ORDER BY m.xp DESC LIMIT $2`, companyID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var performers []model.TopPerformer
	for rows.Next() {
		var p model.TopPerformer
		if err := rows.Scan(&p.MembershipID, &p.DisplayName, &p.XP, &p.Level, &p.Badges, &p.Achievements); err != nil {
			return nil, err
		}
		performers = append(performers, p)
	}
	return performers, rows.Err()
}

func (r *analyticsRepo) GetAchievementStats(ctx context.Context, companyID string) ([]model.AchievementStat, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT a.id, a.name, a.metric,
		        (SELECT COUNT(*) FROM member_achievements WHERE achievement_id = a.id)
		 FROM achievements a WHERE a.company_id = $1 ORDER BY 4 DESC LIMIT 20`, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []model.AchievementStat
	for rows.Next() {
		var s model.AchievementStat
		if err := rows.Scan(&s.AchievementID, &s.Name, &s.Metric, &s.Completions); err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}
	return stats, rows.Err()
}

func (r *analyticsRepo) GetChallengeStats(ctx context.Context, companyID string) (*model.ChallengeStat, error) {
	s := &model.ChallengeStat{}
	r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM challenges WHERE company_id = $1`, companyID).Scan(&s.TotalChallenges)
	r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM challenges WHERE company_id = $1 AND status = 'completed'`, companyID).Scan(&s.Completed)
	r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM challenges WHERE company_id = $1 AND status = 'active'`, companyID).Scan(&s.Active)
	r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM challenges WHERE company_id = $1 AND status = 'pending'`, companyID).Scan(&s.Pending)
	r.pool.QueryRow(ctx, `SELECT COALESCE(AVG(xp_reward), 0) FROM challenges WHERE company_id = $1`, companyID).Scan(&s.AvgXPReward)
	return s, nil
}

func (r *analyticsRepo) GetRewardStats(ctx context.Context, companyID string) ([]model.RewardStat, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT rw.id, rw.name, rw.cost_coins,
		        COUNT(rd.id), COALESCE(SUM(rd.coins_spent), 0)
		 FROM rewards rw
		 LEFT JOIN redemptions rd ON rd.reward_id = rw.id
		 WHERE rw.company_id = $1
		 GROUP BY rw.id, rw.name, rw.cost_coins
		 ORDER BY COUNT(rd.id) DESC`, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []model.RewardStat
	for rows.Next() {
		var s model.RewardStat
		if err := rows.Scan(&s.RewardID, &s.Name, &s.CostCoins, &s.TotalRedeemed, &s.TotalCoinsSpent); err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}
	return stats, rows.Err()
}

func (r *analyticsRepo) GetXPDistribution(ctx context.Context, companyID string) ([]model.XPDistribution, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT level, COUNT(*) FROM memberships WHERE company_id = $1 AND is_active = true GROUP BY level ORDER BY level`, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dist []model.XPDistribution
	for rows.Next() {
		var d model.XPDistribution
		if err := rows.Scan(&d.Level, &d.Count); err != nil {
			return nil, err
		}
		dist = append(dist, d)
	}
	return dist, rows.Err()
}
