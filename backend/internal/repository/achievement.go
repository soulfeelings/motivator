package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hustlers/motivator-backend/internal/model"
)

type AchievementRepository interface {
	Create(ctx context.Context, a *model.Achievement) error
	GetByID(ctx context.Context, id string) (*model.Achievement, error)
	ListByCompany(ctx context.Context, companyID string) ([]model.Achievement, error)
	ListActiveByMetric(ctx context.Context, companyID, metric string) ([]model.Achievement, error)
	Delete(ctx context.Context, id string) error
	HasCompleted(ctx context.Context, membershipID, achievementID string) (bool, error)
	RecordCompletion(ctx context.Context, tx pgx.Tx, membershipID, achievementID string) (*model.MemberAchievement, error)
	ListMemberAchievements(ctx context.Context, membershipID string) ([]model.MemberAchievement, error)
}

type achievementRepo struct {
	pool *pgxpool.Pool
}

func NewAchievementRepository(pool *pgxpool.Pool) AchievementRepository {
	return &achievementRepo{pool: pool}
}

func (r *achievementRepo) Create(ctx context.Context, a *model.Achievement) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO achievements (company_id, name, description, metric, operator, threshold, badge_id, xp_reward, coin_reward)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		 RETURNING id, is_active, created_at`,
		a.CompanyID, a.Name, a.Description, a.Metric, a.Operator, a.Threshold, a.BadgeID, a.XPReward, a.CoinReward,
	).Scan(&a.ID, &a.IsActive, &a.CreatedAt)
}

func (r *achievementRepo) GetByID(ctx context.Context, id string) (*model.Achievement, error) {
	a := &model.Achievement{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, company_id, name, description, metric, operator, threshold, badge_id, xp_reward, coin_reward, is_active, created_at
		 FROM achievements WHERE id = $1`, id,
	).Scan(&a.ID, &a.CompanyID, &a.Name, &a.Description, &a.Metric, &a.Operator, &a.Threshold, &a.BadgeID, &a.XPReward, &a.CoinReward, &a.IsActive, &a.CreatedAt)
	if err == pgx.ErrNoRows {
		return nil, model.ErrNotFound
	}
	return a, err
}

func (r *achievementRepo) ListByCompany(ctx context.Context, companyID string) ([]model.Achievement, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, name, description, metric, operator, threshold, badge_id, xp_reward, coin_reward, is_active, created_at
		 FROM achievements WHERE company_id = $1 ORDER BY created_at ASC`, companyID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanAchievements(rows)
}

func (r *achievementRepo) ListActiveByMetric(ctx context.Context, companyID, metric string) ([]model.Achievement, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, name, description, metric, operator, threshold, badge_id, xp_reward, coin_reward, is_active, created_at
		 FROM achievements WHERE company_id = $1 AND metric = $2 AND is_active = true`, companyID, metric,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanAchievements(rows)
}

func (r *achievementRepo) Delete(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM achievements WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return model.ErrNotFound
	}
	return nil
}

func (r *achievementRepo) HasCompleted(ctx context.Context, membershipID, achievementID string) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM member_achievements WHERE membership_id = $1 AND achievement_id = $2)`,
		membershipID, achievementID,
	).Scan(&exists)
	return exists, err
}

func (r *achievementRepo) RecordCompletion(ctx context.Context, tx pgx.Tx, membershipID, achievementID string) (*model.MemberAchievement, error) {
	ma := &model.MemberAchievement{}
	err := tx.QueryRow(ctx,
		`INSERT INTO member_achievements (membership_id, achievement_id)
		 VALUES ($1, $2)
		 RETURNING id, membership_id, achievement_id, completed_at`,
		membershipID, achievementID,
	).Scan(&ma.ID, &ma.MembershipID, &ma.AchievementID, &ma.CompletedAt)
	return ma, err
}

func (r *achievementRepo) ListMemberAchievements(ctx context.Context, membershipID string) ([]model.MemberAchievement, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT ma.id, ma.membership_id, ma.achievement_id, ma.completed_at,
		        a.id, a.company_id, a.name, a.description, a.metric, a.operator, a.threshold, a.badge_id, a.xp_reward, a.coin_reward, a.is_active, a.created_at
		 FROM member_achievements ma
		 JOIN achievements a ON a.id = ma.achievement_id
		 WHERE ma.membership_id = $1
		 ORDER BY ma.completed_at DESC`, membershipID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.MemberAchievement
	for rows.Next() {
		var ma model.MemberAchievement
		a := &model.Achievement{}
		if err := rows.Scan(&ma.ID, &ma.MembershipID, &ma.AchievementID, &ma.CompletedAt,
			&a.ID, &a.CompanyID, &a.Name, &a.Description, &a.Metric, &a.Operator, &a.Threshold, &a.BadgeID, &a.XPReward, &a.CoinReward, &a.IsActive, &a.CreatedAt); err != nil {
			return nil, err
		}
		ma.Achievement = a
		result = append(result, ma)
	}
	return result, rows.Err()
}

func scanAchievements(rows pgx.Rows) ([]model.Achievement, error) {
	var achievements []model.Achievement
	for rows.Next() {
		var a model.Achievement
		if err := rows.Scan(&a.ID, &a.CompanyID, &a.Name, &a.Description, &a.Metric, &a.Operator, &a.Threshold, &a.BadgeID, &a.XPReward, &a.CoinReward, &a.IsActive, &a.CreatedAt); err != nil {
			return nil, err
		}
		achievements = append(achievements, a)
	}
	return achievements, rows.Err()
}
