package service

import (
	"context"

	"github.com/hustlers/motivator-backend/internal/model"
	"github.com/hustlers/motivator-backend/internal/repository"
)

type GamePlanService struct {
	plans repository.GamePlanRepository
}

func NewGamePlanService(plans repository.GamePlanRepository) *GamePlanService {
	return &GamePlanService{plans: plans}
}

func (s *GamePlanService) Create(ctx context.Context, companyID string, req model.CreateGamePlanRequest) (*model.GamePlan, error) {
	gp := &model.GamePlan{
		CompanyID:   companyID,
		Name:        req.Name,
		Description: req.Description,
		FlowData:    req.FlowData,
	}
	if err := s.plans.Create(ctx, gp); err != nil {
		return nil, err
	}
	return gp, nil
}

func (s *GamePlanService) GetByID(ctx context.Context, id string) (*model.GamePlan, error) {
	return s.plans.GetByID(ctx, id)
}

func (s *GamePlanService) ListByCompany(ctx context.Context, companyID string) ([]model.GamePlan, error) {
	return s.plans.ListByCompany(ctx, companyID)
}

func (s *GamePlanService) Update(ctx context.Context, id string, req model.UpdateGamePlanRequest) (*model.GamePlan, error) {
	return s.plans.Update(ctx, id, req)
}

func (s *GamePlanService) SaveFlow(ctx context.Context, id string, flowData model.FlowData) error {
	return s.plans.UpdateFlowData(ctx, id, flowData)
}

func (s *GamePlanService) SetActive(ctx context.Context, id string, active bool) error {
	return s.plans.SetActive(ctx, id, active)
}

func (s *GamePlanService) Delete(ctx context.Context, id string) error {
	return s.plans.Delete(ctx, id)
}
