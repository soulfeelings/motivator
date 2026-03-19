package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hustlers/motivator-game/internal/model"
)

type ArmyRepository interface {
	GetArmy(ctx context.Context, baseID string) ([]model.ArmyUnit, error)
	HireUnits(ctx context.Context, baseID, unitID string, count int) error
	RemoveUnits(ctx context.Context, baseID, unitID string, count int) error
	GetUnitTypes(ctx context.Context) ([]model.UnitType, error)
	GetUnitType(ctx context.Context, id string) (*model.UnitType, error)
}

type armyRepo struct {
	pool *pgxpool.Pool
}

func NewArmyRepository(pool *pgxpool.Pool) ArmyRepository {
	return &armyRepo{pool: pool}
}

func (r *armyRepo) GetArmy(ctx context.Context, baseID string) ([]model.ArmyUnit, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, base_id, unit_id, count FROM army_units WHERE base_id = $1 AND count > 0`, baseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var units []model.ArmyUnit
	for rows.Next() {
		var u model.ArmyUnit
		if err := rows.Scan(&u.ID, &u.BaseID, &u.UnitID, &u.Count); err != nil {
			return nil, err
		}
		units = append(units, u)
	}
	return units, rows.Err()
}

func (r *armyRepo) HireUnits(ctx context.Context, baseID, unitID string, count int) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO army_units (base_id, unit_id, count) VALUES ($1, $2, $3)
		 ON CONFLICT (base_id, unit_id) DO UPDATE SET count = army_units.count + $3`,
		baseID, unitID, count)
	return err
}

func (r *armyRepo) RemoveUnits(ctx context.Context, baseID, unitID string, count int) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE army_units SET count = GREATEST(count - $3, 0) WHERE base_id = $1 AND unit_id = $2`,
		baseID, unitID, count)
	return err
}

func (r *armyRepo) GetUnitTypes(ctx context.Context) ([]model.UnitType, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, name, description, cost, hp, attack, defense, speed, category FROM unit_types ORDER BY cost ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var types []model.UnitType
	for rows.Next() {
		var ut model.UnitType
		if err := rows.Scan(&ut.ID, &ut.Name, &ut.Description, &ut.Cost, &ut.HP, &ut.Attack, &ut.Defense, &ut.Speed, &ut.Category); err != nil {
			return nil, err
		}
		types = append(types, ut)
	}
	return types, rows.Err()
}

func (r *armyRepo) GetUnitType(ctx context.Context, id string) (*model.UnitType, error) {
	ut := &model.UnitType{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, name, description, cost, hp, attack, defense, speed, category FROM unit_types WHERE id = $1`, id,
	).Scan(&ut.ID, &ut.Name, &ut.Description, &ut.Cost, &ut.HP, &ut.Attack, &ut.Defense, &ut.Speed, &ut.Category)
	return ut, err
}
