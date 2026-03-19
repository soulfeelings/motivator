package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hustlers/motivator-backend/internal/model"
)

type BadgeRepository interface {
	Create(ctx context.Context, badge *model.Badge) error
	GetByID(ctx context.Context, id string) (*model.Badge, error)
	ListByCompany(ctx context.Context, companyID string) ([]model.Badge, error)
	Delete(ctx context.Context, id string) error
	AwardBadge(ctx context.Context, tx pgx.Tx, membershipID, badgeID string) (*model.MemberBadge, error)
	ListMemberBadges(ctx context.Context, membershipID string) ([]model.MemberBadge, error)
}

type badgeRepo struct {
	pool *pgxpool.Pool
}

func NewBadgeRepository(pool *pgxpool.Pool) BadgeRepository {
	return &badgeRepo{pool: pool}
}

func (r *badgeRepo) Create(ctx context.Context, badge *model.Badge) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO badges (company_id, name, description, icon_url, xp_reward, coin_reward)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 RETURNING id, created_at`,
		badge.CompanyID, badge.Name, badge.Description, badge.IconURL, badge.XPReward, badge.CoinReward,
	).Scan(&badge.ID, &badge.CreatedAt)
}

func (r *badgeRepo) GetByID(ctx context.Context, id string) (*model.Badge, error) {
	b := &model.Badge{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, company_id, name, description, icon_url, xp_reward, coin_reward, created_at
		 FROM badges WHERE id = $1`, id,
	).Scan(&b.ID, &b.CompanyID, &b.Name, &b.Description, &b.IconURL, &b.XPReward, &b.CoinReward, &b.CreatedAt)
	if err == pgx.ErrNoRows {
		return nil, model.ErrNotFound
	}
	return b, err
}

func (r *badgeRepo) ListByCompany(ctx context.Context, companyID string) ([]model.Badge, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, name, description, icon_url, xp_reward, coin_reward, created_at
		 FROM badges WHERE company_id = $1 ORDER BY created_at ASC`, companyID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var badges []model.Badge
	for rows.Next() {
		var b model.Badge
		if err := rows.Scan(&b.ID, &b.CompanyID, &b.Name, &b.Description, &b.IconURL, &b.XPReward, &b.CoinReward, &b.CreatedAt); err != nil {
			return nil, err
		}
		badges = append(badges, b)
	}
	return badges, rows.Err()
}

func (r *badgeRepo) Delete(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM badges WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return model.ErrNotFound
	}
	return nil
}

func (r *badgeRepo) AwardBadge(ctx context.Context, tx pgx.Tx, membershipID, badgeID string) (*model.MemberBadge, error) {
	mb := &model.MemberBadge{}
	err := tx.QueryRow(ctx,
		`INSERT INTO member_badges (membership_id, badge_id)
		 VALUES ($1, $2)
		 RETURNING id, membership_id, badge_id, awarded_at`,
		membershipID, badgeID,
	).Scan(&mb.ID, &mb.MembershipID, &mb.BadgeID, &mb.AwardedAt)
	return mb, err
}

func (r *badgeRepo) ListMemberBadges(ctx context.Context, membershipID string) ([]model.MemberBadge, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT mb.id, mb.membership_id, mb.badge_id, mb.awarded_at,
		        b.id, b.company_id, b.name, b.description, b.icon_url, b.xp_reward, b.coin_reward, b.created_at
		 FROM member_badges mb
		 JOIN badges b ON b.id = mb.badge_id
		 WHERE mb.membership_id = $1
		 ORDER BY mb.awarded_at DESC`, membershipID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.MemberBadge
	for rows.Next() {
		var mb model.MemberBadge
		b := &model.Badge{}
		if err := rows.Scan(&mb.ID, &mb.MembershipID, &mb.BadgeID, &mb.AwardedAt,
			&b.ID, &b.CompanyID, &b.Name, &b.Description, &b.IconURL, &b.XPReward, &b.CoinReward, &b.CreatedAt); err != nil {
			return nil, err
		}
		mb.Badge = b
		result = append(result, mb)
	}
	return result, rows.Err()
}
