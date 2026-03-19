package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hustlers/motivator-backend/internal/model"
)

type DeviceTokenRepository interface {
	Register(ctx context.Context, dt *model.DeviceToken) error
	Unregister(ctx context.Context, membershipID, token string) error
	ListByMembership(ctx context.Context, membershipID string) ([]model.DeviceToken, error)
	ListByMemberships(ctx context.Context, membershipIDs []string) ([]model.DeviceToken, error)
}

type deviceTokenRepo struct {
	pool *pgxpool.Pool
}

func NewDeviceTokenRepository(pool *pgxpool.Pool) DeviceTokenRepository {
	return &deviceTokenRepo{pool: pool}
}

func (r *deviceTokenRepo) Register(ctx context.Context, dt *model.DeviceToken) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO device_tokens (membership_id, token, platform)
		 VALUES ($1, $2, $3)
		 ON CONFLICT (membership_id, token) DO UPDATE SET platform = $3
		 RETURNING id, created_at`,
		dt.MembershipID, dt.Token, dt.Platform,
	).Scan(&dt.ID, &dt.CreatedAt)
}

func (r *deviceTokenRepo) Unregister(ctx context.Context, membershipID, token string) error {
	_, err := r.pool.Exec(ctx,
		`DELETE FROM device_tokens WHERE membership_id = $1 AND token = $2`,
		membershipID, token)
	return err
}

func (r *deviceTokenRepo) ListByMembership(ctx context.Context, membershipID string) ([]model.DeviceToken, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, membership_id, token, platform, created_at FROM device_tokens WHERE membership_id = $1`,
		membershipID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tokens []model.DeviceToken
	for rows.Next() {
		var dt model.DeviceToken
		if err := rows.Scan(&dt.ID, &dt.MembershipID, &dt.Token, &dt.Platform, &dt.CreatedAt); err != nil {
			return nil, err
		}
		tokens = append(tokens, dt)
	}
	return tokens, rows.Err()
}

func (r *deviceTokenRepo) ListByMemberships(ctx context.Context, membershipIDs []string) ([]model.DeviceToken, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, membership_id, token, platform, created_at FROM device_tokens WHERE membership_id = ANY($1)`,
		membershipIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tokens []model.DeviceToken
	for rows.Next() {
		var dt model.DeviceToken
		if err := rows.Scan(&dt.ID, &dt.MembershipID, &dt.Token, &dt.Platform, &dt.CreatedAt); err != nil {
			return nil, err
		}
		tokens = append(tokens, dt)
	}
	return tokens, rows.Err()
}
