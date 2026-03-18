package service

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hustlers/motivator-backend/internal/model"
	"github.com/hustlers/motivator-backend/internal/repository"
)

type CompanyService struct {
	pool       *pgxpool.Pool
	companies  repository.CompanyRepository
	members    repository.MembershipRepository
}

func NewCompanyService(pool *pgxpool.Pool, companies repository.CompanyRepository, members repository.MembershipRepository) *CompanyService {
	return &CompanyService{pool: pool, companies: companies, members: members}
}

func (s *CompanyService) Create(ctx context.Context, userID string, req model.CreateCompanyRequest) (*model.Company, *model.Membership, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			log.Printf("rollback error: %v", err)
		}
	}()

	company := &model.Company{
		Name: req.Name,
		Slug: req.Slug,
	}
	if err := s.companies.Create(ctx, tx, company); err != nil {
		return nil, nil, err
	}

	membership := &model.Membership{
		UserID:    userID,
		CompanyID: company.ID,
		Role:      model.RoleOwner,
	}
	if err := s.members.Create(ctx, tx, membership); err != nil {
		return nil, nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, nil, err
	}

	return company, membership, nil
}

func (s *CompanyService) GetByID(ctx context.Context, id string) (*model.Company, error) {
	return s.companies.GetByID(ctx, id)
}

func (s *CompanyService) Update(ctx context.Context, id string, req model.UpdateCompanyRequest) (*model.Company, error) {
	return s.companies.Update(ctx, id, req)
}

func (s *CompanyService) Delete(ctx context.Context, id string) error {
	return s.companies.Delete(ctx, id)
}
