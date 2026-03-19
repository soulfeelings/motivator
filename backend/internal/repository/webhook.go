package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hustlers/motivator-backend/internal/model"
)

type WebhookRepository interface {
	Create(ctx context.Context, w *model.Webhook) error
	GetByID(ctx context.Context, id string) (*model.Webhook, error)
	ListByCompany(ctx context.Context, companyID string) ([]model.Webhook, error)
	ListActiveByEvent(ctx context.Context, companyID, event string) ([]model.Webhook, error)
	Delete(ctx context.Context, id string) error
}

type webhookRepo struct {
	pool *pgxpool.Pool
}

func NewWebhookRepository(pool *pgxpool.Pool) WebhookRepository {
	return &webhookRepo{pool: pool}
}

func (r *webhookRepo) Create(ctx context.Context, w *model.Webhook) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO webhooks (company_id, name, url, platform, events) VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, is_active, created_at, updated_at`,
		w.CompanyID, w.Name, w.URL, w.Platform, w.Events,
	).Scan(&w.ID, &w.IsActive, &w.CreatedAt, &w.UpdatedAt)
}

func (r *webhookRepo) GetByID(ctx context.Context, id string) (*model.Webhook, error) {
	w := &model.Webhook{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, company_id, name, url, platform, events, is_active, created_at, updated_at FROM webhooks WHERE id = $1`, id,
	).Scan(&w.ID, &w.CompanyID, &w.Name, &w.URL, &w.Platform, &w.Events, &w.IsActive, &w.CreatedAt, &w.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, model.ErrNotFound
	}
	return w, err
}

func (r *webhookRepo) ListByCompany(ctx context.Context, companyID string) ([]model.Webhook, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, name, url, platform, events, is_active, created_at, updated_at
		 FROM webhooks WHERE company_id = $1 ORDER BY created_at ASC`, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var webhooks []model.Webhook
	for rows.Next() {
		var w model.Webhook
		if err := rows.Scan(&w.ID, &w.CompanyID, &w.Name, &w.URL, &w.Platform, &w.Events, &w.IsActive, &w.CreatedAt, &w.UpdatedAt); err != nil {
			return nil, err
		}
		webhooks = append(webhooks, w)
	}
	return webhooks, rows.Err()
}

func (r *webhookRepo) ListActiveByEvent(ctx context.Context, companyID, event string) ([]model.Webhook, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, name, url, platform, events, is_active, created_at, updated_at
		 FROM webhooks WHERE company_id = $1 AND is_active = true AND $2 = ANY(events)`, companyID, event)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var webhooks []model.Webhook
	for rows.Next() {
		var w model.Webhook
		if err := rows.Scan(&w.ID, &w.CompanyID, &w.Name, &w.URL, &w.Platform, &w.Events, &w.IsActive, &w.CreatedAt, &w.UpdatedAt); err != nil {
			return nil, err
		}
		webhooks = append(webhooks, w)
	}
	return webhooks, rows.Err()
}

func (r *webhookRepo) Delete(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM webhooks WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return model.ErrNotFound
	}
	return nil
}
