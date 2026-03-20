package repository

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hustlers/motivator-backend/internal/model"
)

type IntegrationRepository interface {
	Create(ctx context.Context, i *model.Integration) error
	GetByID(ctx context.Context, id string) (*model.Integration, error)
	GetByWebhookSecret(ctx context.Context, secret string) (*model.Integration, error)
	ListByCompany(ctx context.Context, companyID string) ([]model.Integration, error)
	Delete(ctx context.Context, id string) error
	CreateMapping(ctx context.Context, m *model.IntegrationMapping) error
	ListMappings(ctx context.Context, integrationID string) ([]model.IntegrationMapping, error)
	GetMappingByEvent(ctx context.Context, integrationID, externalEvent string) (*model.IntegrationMapping, error)
	DeleteMapping(ctx context.Context, id string) error
	LogEvent(ctx context.Context, e *model.IntegrationEvent) error
	ListEvents(ctx context.Context, integrationID string, limit int) ([]model.IntegrationEvent, error)
}

type integrationRepo struct {
	pool *pgxpool.Pool
}

func NewIntegrationRepository(pool *pgxpool.Pool) IntegrationRepository {
	return &integrationRepo{pool: pool}
}

func (r *integrationRepo) Create(ctx context.Context, i *model.Integration) error {
	configJSON, _ := json.Marshal(i.Config)
	return r.pool.QueryRow(ctx,
		`INSERT INTO integrations (company_id, provider, name, config, webhook_secret)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, is_active, created_at, updated_at`,
		i.CompanyID, i.Provider, i.Name, configJSON, i.WebhookSecret,
	).Scan(&i.ID, &i.IsActive, &i.CreatedAt, &i.UpdatedAt)
}

func (r *integrationRepo) GetByID(ctx context.Context, id string) (*model.Integration, error) {
	i := &model.Integration{}
	var configJSON []byte
	err := r.pool.QueryRow(ctx,
		`SELECT id, company_id, provider, name, config, webhook_secret, is_active, created_at, updated_at
		 FROM integrations WHERE id = $1`, id,
	).Scan(&i.ID, &i.CompanyID, &i.Provider, &i.Name, &configJSON, &i.WebhookSecret, &i.IsActive, &i.CreatedAt, &i.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, model.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	json.Unmarshal(configJSON, &i.Config)
	return i, nil
}

func (r *integrationRepo) GetByWebhookSecret(ctx context.Context, secret string) (*model.Integration, error) {
	i := &model.Integration{}
	var configJSON []byte
	err := r.pool.QueryRow(ctx,
		`SELECT id, company_id, provider, name, config, webhook_secret, is_active, created_at, updated_at
		 FROM integrations WHERE webhook_secret = $1 AND is_active = true`, secret,
	).Scan(&i.ID, &i.CompanyID, &i.Provider, &i.Name, &configJSON, &i.WebhookSecret, &i.IsActive, &i.CreatedAt, &i.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, model.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	json.Unmarshal(configJSON, &i.Config)
	return i, nil
}

func (r *integrationRepo) ListByCompany(ctx context.Context, companyID string) ([]model.Integration, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, company_id, provider, name, config, webhook_secret, is_active, created_at, updated_at
		 FROM integrations WHERE company_id = $1 ORDER BY created_at ASC`, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var integrations []model.Integration
	for rows.Next() {
		var i model.Integration
		var configJSON []byte
		if err := rows.Scan(&i.ID, &i.CompanyID, &i.Provider, &i.Name, &configJSON, &i.WebhookSecret, &i.IsActive, &i.CreatedAt, &i.UpdatedAt); err != nil {
			return nil, err
		}
		json.Unmarshal(configJSON, &i.Config)
		integrations = append(integrations, i)
	}
	return integrations, rows.Err()
}

func (r *integrationRepo) Delete(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM integrations WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return model.ErrNotFound
	}
	return nil
}

func (r *integrationRepo) CreateMapping(ctx context.Context, m *model.IntegrationMapping) error {
	transformJSON, _ := json.Marshal(m.Transform)
	if m.UserField == "" {
		m.UserField = "email"
	}
	return r.pool.QueryRow(ctx,
		`INSERT INTO integration_mappings (integration_id, external_event, metric, user_field, transform)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, is_active, created_at`,
		m.IntegrationID, m.ExternalEvent, m.Metric, m.UserField, transformJSON,
	).Scan(&m.ID, &m.IsActive, &m.CreatedAt)
}

func (r *integrationRepo) ListMappings(ctx context.Context, integrationID string) ([]model.IntegrationMapping, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, integration_id, external_event, metric, user_field, transform, is_active, created_at
		 FROM integration_mappings WHERE integration_id = $1 ORDER BY created_at ASC`, integrationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var mappings []model.IntegrationMapping
	for rows.Next() {
		var m model.IntegrationMapping
		var transformJSON []byte
		if err := rows.Scan(&m.ID, &m.IntegrationID, &m.ExternalEvent, &m.Metric, &m.UserField, &transformJSON, &m.IsActive, &m.CreatedAt); err != nil {
			return nil, err
		}
		json.Unmarshal(transformJSON, &m.Transform)
		mappings = append(mappings, m)
	}
	return mappings, rows.Err()
}

func (r *integrationRepo) GetMappingByEvent(ctx context.Context, integrationID, externalEvent string) (*model.IntegrationMapping, error) {
	m := &model.IntegrationMapping{}
	var transformJSON []byte
	err := r.pool.QueryRow(ctx,
		`SELECT id, integration_id, external_event, metric, user_field, transform, is_active, created_at
		 FROM integration_mappings WHERE integration_id = $1 AND external_event = $2 AND is_active = true`,
		integrationID, externalEvent,
	).Scan(&m.ID, &m.IntegrationID, &m.ExternalEvent, &m.Metric, &m.UserField, &transformJSON, &m.IsActive, &m.CreatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	json.Unmarshal(transformJSON, &m.Transform)
	return m, nil
}

func (r *integrationRepo) DeleteMapping(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM integration_mappings WHERE id = $1`, id)
	return err
}

func (r *integrationRepo) LogEvent(ctx context.Context, e *model.IntegrationEvent) error {
	rawJSON, _ := json.Marshal(e.RawData)
	return r.pool.QueryRow(ctx,
		`INSERT INTO integration_events (integration_id, external_event, raw_data, metric, user_email, value, processed, error)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		 RETURNING id, received_at`,
		e.IntegrationID, e.ExternalEvent, rawJSON, e.Metric, e.UserEmail, e.Value, e.Processed, e.Error,
	).Scan(&e.ID, &e.ReceivedAt)
}

func (r *integrationRepo) ListEvents(ctx context.Context, integrationID string, limit int) ([]model.IntegrationEvent, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, integration_id, external_event, metric, user_email, value, processed, error, received_at
		 FROM integration_events WHERE integration_id = $1 ORDER BY received_at DESC LIMIT $2`, integrationID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []model.IntegrationEvent
	for rows.Next() {
		var e model.IntegrationEvent
		if err := rows.Scan(&e.ID, &e.IntegrationID, &e.ExternalEvent, &e.Metric, &e.UserEmail, &e.Value, &e.Processed, &e.Error, &e.ReceivedAt); err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	return events, rows.Err()
}
