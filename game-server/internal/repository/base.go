package repository

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hustlers/motivator-game/internal/model"
)

type BaseRepository interface {
	GetOrCreate(ctx context.Context, membershipID string) (*model.Base, error)
	GetByID(ctx context.Context, id string) (*model.Base, error)
	ListAll(ctx context.Context) ([]model.Base, error)
	UpdateCoins(ctx context.Context, id string, delta int) error
	AddBuilding(ctx context.Context, b *model.BaseBuilding) error
	ListBuildings(ctx context.Context, baseID string) ([]model.BaseBuilding, error)
	GetBuildingTypes(ctx context.Context) ([]model.BuildingType, error)
}

type baseRepo struct {
	pool *pgxpool.Pool
}

func NewBaseRepository(pool *pgxpool.Pool) BaseRepository {
	return &baseRepo{pool: pool}
}

func (r *baseRepo) GetOrCreate(ctx context.Context, membershipID string) (*model.Base, error) {
	b := &model.Base{}
	var layoutJSON []byte
	err := r.pool.QueryRow(ctx,
		`INSERT INTO bases (membership_id) VALUES ($1)
		 ON CONFLICT (membership_id) DO UPDATE SET updated_at = now()
		 RETURNING id, membership_id, name, level, layout, coins_balance, created_at, updated_at`,
		membershipID,
	).Scan(&b.ID, &b.MembershipID, &b.Name, &b.Level, &layoutJSON, &b.CoinsBalance, &b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(layoutJSON, &b.Layout)
	return b, nil
}

func (r *baseRepo) GetByID(ctx context.Context, id string) (*model.Base, error) {
	b := &model.Base{}
	var layoutJSON []byte
	err := r.pool.QueryRow(ctx,
		`SELECT id, membership_id, name, level, layout, coins_balance, created_at, updated_at FROM bases WHERE id = $1`, id,
	).Scan(&b.ID, &b.MembershipID, &b.Name, &b.Level, &layoutJSON, &b.CoinsBalance, &b.CreatedAt, &b.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, model.ErrNotFound
	}
	json.Unmarshal(layoutJSON, &b.Layout)
	return b, err
}

func (r *baseRepo) ListAll(ctx context.Context) ([]model.Base, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, membership_id, name, level, layout, coins_balance, created_at, updated_at FROM bases ORDER BY level DESC, coins_balance DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bases []model.Base
	for rows.Next() {
		var b model.Base
		var layoutJSON []byte
		if err := rows.Scan(&b.ID, &b.MembershipID, &b.Name, &b.Level, &layoutJSON, &b.CoinsBalance, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, err
		}
		json.Unmarshal(layoutJSON, &b.Layout)
		bases = append(bases, b)
	}
	return bases, rows.Err()
}

func (r *baseRepo) UpdateCoins(ctx context.Context, id string, delta int) error {
	_, err := r.pool.Exec(ctx, `UPDATE bases SET coins_balance = coins_balance + $2, updated_at = now() WHERE id = $1`, id, delta)
	return err
}

func (r *baseRepo) AddBuilding(ctx context.Context, b *model.BaseBuilding) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO base_buildings (base_id, building_id, grid_x, grid_y, hp)
		 VALUES ($1, $2, $3, $4, (SELECT hp FROM building_types WHERE id = $2))
		 RETURNING id, level, hp, built_at`,
		b.BaseID, b.BuildingID, b.GridX, b.GridY,
	).Scan(&b.ID, &b.Level, &b.HP, &b.BuiltAt)
}

func (r *baseRepo) ListBuildings(ctx context.Context, baseID string) ([]model.BaseBuilding, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, base_id, building_id, grid_x, grid_y, level, hp, built_at FROM base_buildings WHERE base_id = $1`, baseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var buildings []model.BaseBuilding
	for rows.Next() {
		var b model.BaseBuilding
		if err := rows.Scan(&b.ID, &b.BaseID, &b.BuildingID, &b.GridX, &b.GridY, &b.Level, &b.HP, &b.BuiltAt); err != nil {
			return nil, err
		}
		buildings = append(buildings, b)
	}
	return buildings, rows.Err()
}

func (r *baseRepo) GetBuildingTypes(ctx context.Context) ([]model.BuildingType, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, name, description, cost, build_time, hp, category, unlocks FROM building_types ORDER BY cost ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var types []model.BuildingType
	for rows.Next() {
		var bt model.BuildingType
		if err := rows.Scan(&bt.ID, &bt.Name, &bt.Description, &bt.Cost, &bt.BuildTime, &bt.HP, &bt.Category, &bt.Unlocks); err != nil {
			return nil, err
		}
		types = append(types, bt)
	}
	return types, rows.Err()
}
