package repository

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hustlers/motivator-backend/internal/model"
)

type GamePlanRepository interface {
	Create(ctx context.Context, gp *model.GamePlan) error
	GetByID(ctx context.Context, id string) (*model.GamePlan, error)
	ListByCompany(ctx context.Context, companyID string) ([]model.GamePlan, error)
	Update(ctx context.Context, id string, req model.UpdateGamePlanRequest) (*model.GamePlan, error)
	UpdateFlowData(ctx context.Context, id string, flowData model.FlowData) error
	SetActive(ctx context.Context, id string, active bool) error
	Delete(ctx context.Context, id string) error
}

type gamePlanRepo struct {
	pool *pgxpool.Pool
}

func NewGamePlanRepository(pool *pgxpool.Pool) GamePlanRepository {
	return &gamePlanRepo{pool: pool}
}

func (r *gamePlanRepo) Create(ctx context.Context, gp *model.GamePlan) error {
	flowJSON, err := json.Marshal(gp.FlowData)
	if err != nil {
		return err
	}
	return r.pool.QueryRow(ctx,
		`INSERT INTO game_plans (company_id, name, description, flow_data)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, is_active, created_at, updated_at`,
		gp.CompanyID, gp.Name, gp.Description, flowJSON,
	).Scan(&gp.ID, &gp.IsActive, &gp.CreatedAt, &gp.UpdatedAt)
}

func (r *gamePlanRepo) GetByID(ctx context.Context, id string) (*model.GamePlan, error) {
	gp := &model.GamePlan{}
	var flowJSON []byte
	err := r.pool.QueryRow(ctx,
		`SELECT id, company_id, name, description, flow_data, is_active, created_at, updated_at
		 FROM game_plans WHERE id = $1`, id,
	).Scan(&gp.ID, &gp.CompanyID, &gp.Name, &gp.Description, &flowJSON, &gp.IsActive, &gp.CreatedAt, &gp.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, model.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	json.Unmarshal(flowJSON, &gp.FlowData)
	return gp, nil
}

func (r *gamePlanRepo) ListByCompany(ctx context.Context, companyID string) ([]model.GamePlan, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, name, description, flow_data, is_active, created_at, updated_at
		 FROM game_plans WHERE company_id = $1 ORDER BY created_at DESC`, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plans []model.GamePlan
	for rows.Next() {
		var gp model.GamePlan
		var flowJSON []byte
		if err := rows.Scan(&gp.ID, &gp.CompanyID, &gp.Name, &gp.Description, &flowJSON, &gp.IsActive, &gp.CreatedAt, &gp.UpdatedAt); err != nil {
			return nil, err
		}
		json.Unmarshal(flowJSON, &gp.FlowData)
		plans = append(plans, gp)
	}
	return plans, rows.Err()
}

func (r *gamePlanRepo) Update(ctx context.Context, id string, req model.UpdateGamePlanRequest) (*model.GamePlan, error) {
	gp := &model.GamePlan{}
	var flowJSON []byte

	if req.FlowData != nil {
		fj, _ := json.Marshal(req.FlowData)
		err := r.pool.QueryRow(ctx,
			`UPDATE game_plans SET
				name = COALESCE($2, name),
				description = COALESCE($3, description),
				flow_data = $4,
				updated_at = now()
			 WHERE id = $1
			 RETURNING id, company_id, name, description, flow_data, is_active, created_at, updated_at`,
			id, req.Name, req.Description, fj,
		).Scan(&gp.ID, &gp.CompanyID, &gp.Name, &gp.Description, &flowJSON, &gp.IsActive, &gp.CreatedAt, &gp.UpdatedAt)
		if err == pgx.ErrNoRows {
			return nil, model.ErrNotFound
		}
		if err != nil {
			return nil, err
		}
	} else {
		err := r.pool.QueryRow(ctx,
			`UPDATE game_plans SET
				name = COALESCE($2, name),
				description = COALESCE($3, description),
				updated_at = now()
			 WHERE id = $1
			 RETURNING id, company_id, name, description, flow_data, is_active, created_at, updated_at`,
			id, req.Name, req.Description,
		).Scan(&gp.ID, &gp.CompanyID, &gp.Name, &gp.Description, &flowJSON, &gp.IsActive, &gp.CreatedAt, &gp.UpdatedAt)
		if err == pgx.ErrNoRows {
			return nil, model.ErrNotFound
		}
		if err != nil {
			return nil, err
		}
	}
	json.Unmarshal(flowJSON, &gp.FlowData)
	return gp, nil
}

func (r *gamePlanRepo) UpdateFlowData(ctx context.Context, id string, flowData model.FlowData) error {
	fj, err := json.Marshal(flowData)
	if err != nil {
		return err
	}
	_, err = r.pool.Exec(ctx, `UPDATE game_plans SET flow_data = $2, updated_at = now() WHERE id = $1`, id, fj)
	return err
}

func (r *gamePlanRepo) SetActive(ctx context.Context, id string, active bool) error {
	tag, err := r.pool.Exec(ctx, `UPDATE game_plans SET is_active = $2, updated_at = now() WHERE id = $1`, id, active)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return model.ErrNotFound
	}
	return nil
}

func (r *gamePlanRepo) Delete(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM game_plans WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return model.ErrNotFound
	}
	return nil
}
