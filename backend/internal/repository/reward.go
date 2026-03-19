package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hustlers/motivator-backend/internal/model"
)

type RewardRepository interface {
	Create(ctx context.Context, r *model.Reward) error
	GetByID(ctx context.Context, id string) (*model.Reward, error)
	ListByCompany(ctx context.Context, companyID string) ([]model.Reward, error)
	Delete(ctx context.Context, id string) error
	DecrementStock(ctx context.Context, tx pgx.Tx, id string) error
	CreateRedemption(ctx context.Context, tx pgx.Tx, rd *model.Redemption) error
	ListRedemptions(ctx context.Context, companyID string) ([]model.Redemption, error)
	ListMemberRedemptions(ctx context.Context, membershipID string) ([]model.Redemption, error)
	UpdateRedemptionStatus(ctx context.Context, id string, status model.RedemptionStatus) error
}

type rewardRepo struct {
	pool *pgxpool.Pool
}

func NewRewardRepository(pool *pgxpool.Pool) RewardRepository {
	return &rewardRepo{pool: pool}
}

func (r *rewardRepo) Create(ctx context.Context, rw *model.Reward) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO rewards (company_id, name, description, cost_coins, stock)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, is_active, created_at, updated_at`,
		rw.CompanyID, rw.Name, rw.Description, rw.CostCoins, rw.Stock,
	).Scan(&rw.ID, &rw.IsActive, &rw.CreatedAt, &rw.UpdatedAt)
}

func (r *rewardRepo) GetByID(ctx context.Context, id string) (*model.Reward, error) {
	rw := &model.Reward{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, company_id, name, description, cost_coins, stock, is_active, created_at, updated_at
		 FROM rewards WHERE id = $1`, id,
	).Scan(&rw.ID, &rw.CompanyID, &rw.Name, &rw.Description, &rw.CostCoins, &rw.Stock, &rw.IsActive, &rw.CreatedAt, &rw.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, model.ErrNotFound
	}
	return rw, err
}

func (r *rewardRepo) ListByCompany(ctx context.Context, companyID string) ([]model.Reward, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, name, description, cost_coins, stock, is_active, created_at, updated_at
		 FROM rewards WHERE company_id = $1 AND is_active = true ORDER BY cost_coins ASC`, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rewards []model.Reward
	for rows.Next() {
		var rw model.Reward
		if err := rows.Scan(&rw.ID, &rw.CompanyID, &rw.Name, &rw.Description, &rw.CostCoins, &rw.Stock, &rw.IsActive, &rw.CreatedAt, &rw.UpdatedAt); err != nil {
			return nil, err
		}
		rewards = append(rewards, rw)
	}
	return rewards, rows.Err()
}

func (r *rewardRepo) Delete(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, `UPDATE rewards SET is_active = false, updated_at = now() WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return model.ErrNotFound
	}
	return nil
}

func (r *rewardRepo) DecrementStock(ctx context.Context, tx pgx.Tx, id string) error {
	_, err := tx.Exec(ctx, `UPDATE rewards SET stock = stock - 1, updated_at = now() WHERE id = $1 AND (stock IS NULL OR stock > 0)`, id)
	return err
}

func (r *rewardRepo) CreateRedemption(ctx context.Context, tx pgx.Tx, rd *model.Redemption) error {
	return tx.QueryRow(ctx,
		`INSERT INTO redemptions (membership_id, reward_id, coins_spent)
		 VALUES ($1, $2, $3)
		 RETURNING id, status, redeemed_at`,
		rd.MembershipID, rd.RewardID, rd.CoinsSpent,
	).Scan(&rd.ID, &rd.Status, &rd.RedeemedAt)
}

func (r *rewardRepo) ListRedemptions(ctx context.Context, companyID string) ([]model.Redemption, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT rd.id, rd.membership_id, rd.reward_id, rd.coins_spent, rd.status, rd.redeemed_at, rd.fulfilled_at,
		        rw.id, rw.company_id, rw.name, rw.description, rw.cost_coins, rw.stock, rw.is_active, rw.created_at, rw.updated_at
		 FROM redemptions rd
		 JOIN rewards rw ON rw.id = rd.reward_id
		 WHERE rw.company_id = $1
		 ORDER BY rd.redeemed_at DESC`, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanRedemptions(rows)
}

func (r *rewardRepo) ListMemberRedemptions(ctx context.Context, membershipID string) ([]model.Redemption, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT rd.id, rd.membership_id, rd.reward_id, rd.coins_spent, rd.status, rd.redeemed_at, rd.fulfilled_at,
		        rw.id, rw.company_id, rw.name, rw.description, rw.cost_coins, rw.stock, rw.is_active, rw.created_at, rw.updated_at
		 FROM redemptions rd
		 JOIN rewards rw ON rw.id = rd.reward_id
		 WHERE rd.membership_id = $1
		 ORDER BY rd.redeemed_at DESC`, membershipID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanRedemptions(rows)
}

func (r *rewardRepo) UpdateRedemptionStatus(ctx context.Context, id string, status model.RedemptionStatus) error {
	q := `UPDATE redemptions SET status = $2 WHERE id = $1`
	if status == model.RedemptionFulfilled {
		q = `UPDATE redemptions SET status = $2, fulfilled_at = now() WHERE id = $1`
	}
	tag, err := r.pool.Exec(ctx, q, id, status)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return model.ErrNotFound
	}
	return nil
}

func scanRedemptions(rows pgx.Rows) ([]model.Redemption, error) {
	var result []model.Redemption
	for rows.Next() {
		var rd model.Redemption
		rw := &model.Reward{}
		if err := rows.Scan(&rd.ID, &rd.MembershipID, &rd.RewardID, &rd.CoinsSpent, &rd.Status, &rd.RedeemedAt, &rd.FulfilledAt,
			&rw.ID, &rw.CompanyID, &rw.Name, &rw.Description, &rw.CostCoins, &rw.Stock, &rw.IsActive, &rw.CreatedAt, &rw.UpdatedAt); err != nil {
			return nil, err
		}
		rd.Reward = rw
		result = append(result, rd)
	}
	return result, rows.Err()
}
