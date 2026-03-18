package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hustlers/motivator-backend/internal/model"
)

type CompanyRepository interface {
	Create(ctx context.Context, tx pgx.Tx, company *model.Company) error
	GetByID(ctx context.Context, id string) (*model.Company, error)
	GetBySlug(ctx context.Context, slug string) (*model.Company, error)
	Update(ctx context.Context, id string, req model.UpdateCompanyRequest) (*model.Company, error)
	Delete(ctx context.Context, id string) error
}

type companyRepo struct {
	pool *pgxpool.Pool
}

func NewCompanyRepository(pool *pgxpool.Pool) CompanyRepository {
	return &companyRepo{pool: pool}
}

func (r *companyRepo) Create(ctx context.Context, tx pgx.Tx, company *model.Company) error {
	return tx.QueryRow(ctx,
		`INSERT INTO companies (name, slug, logo_url) VALUES ($1, $2, $3)
		 RETURNING id, created_at, updated_at`,
		company.Name, company.Slug, company.LogoURL,
	).Scan(&company.ID, &company.CreatedAt, &company.UpdatedAt)
}

func (r *companyRepo) GetByID(ctx context.Context, id string) (*model.Company, error) {
	c := &model.Company{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, name, slug, logo_url, created_at, updated_at FROM companies WHERE id = $1`, id,
	).Scan(&c.ID, &c.Name, &c.Slug, &c.LogoURL, &c.CreatedAt, &c.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, model.ErrNotFound
	}
	return c, err
}

func (r *companyRepo) GetBySlug(ctx context.Context, slug string) (*model.Company, error) {
	c := &model.Company{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, name, slug, logo_url, created_at, updated_at FROM companies WHERE slug = $1`, slug,
	).Scan(&c.ID, &c.Name, &c.Slug, &c.LogoURL, &c.CreatedAt, &c.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, model.ErrNotFound
	}
	return c, err
}

func (r *companyRepo) Update(ctx context.Context, id string, req model.UpdateCompanyRequest) (*model.Company, error) {
	c := &model.Company{}
	err := r.pool.QueryRow(ctx,
		`UPDATE companies SET
			name     = COALESCE($2, name),
			slug     = COALESCE($3, slug),
			logo_url = COALESCE($4, logo_url),
			updated_at = now()
		 WHERE id = $1
		 RETURNING id, name, slug, logo_url, created_at, updated_at`,
		id, req.Name, req.Slug, req.LogoURL,
	).Scan(&c.ID, &c.Name, &c.Slug, &c.LogoURL, &c.CreatedAt, &c.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, model.ErrNotFound
	}
	return c, err
}

func (r *companyRepo) Delete(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM companies WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return model.ErrNotFound
	}
	return nil
}
