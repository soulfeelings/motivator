package service

import (
	"context"
	"fmt"

	"github.com/hustlers/motivator-game/internal/model"
	"github.com/hustlers/motivator-game/internal/repository"
)

type GameService struct {
	bases   repository.BaseRepository
	army    repository.ArmyRepository
	battles repository.BattleRepository
}

func NewGameService(bases repository.BaseRepository, army repository.ArmyRepository, battles repository.BattleRepository) *GameService {
	return &GameService{bases: bases, army: army, battles: battles}
}

func (s *GameService) GetOrCreateBase(ctx context.Context, membershipID string) (*model.BaseOverview, error) {
	base, err := s.bases.GetOrCreate(ctx, membershipID)
	if err != nil {
		return nil, err
	}
	buildings, err := s.bases.ListBuildings(ctx, base.ID)
	if err != nil {
		return nil, err
	}
	army, err := s.army.GetArmy(ctx, base.ID)
	if err != nil {
		return nil, err
	}
	return &model.BaseOverview{Base: *base, Buildings: buildings, Army: army}, nil
}

func (s *GameService) GetBaseOverview(ctx context.Context, baseID string) (*model.BaseOverview, error) {
	base, err := s.bases.GetByID(ctx, baseID)
	if err != nil {
		return nil, err
	}
	buildings, err := s.bases.ListBuildings(ctx, base.ID)
	if err != nil {
		return nil, err
	}
	army, err := s.army.GetArmy(ctx, base.ID)
	if err != nil {
		return nil, err
	}
	return &model.BaseOverview{Base: *base, Buildings: buildings, Army: army}, nil
}

func (s *GameService) ListBases(ctx context.Context) ([]model.Base, error) {
	return s.bases.ListAll(ctx)
}

func (s *GameService) DepositCoins(ctx context.Context, baseID string, amount int) error {
	return s.bases.UpdateCoins(ctx, baseID, amount)
}

func (s *GameService) Build(ctx context.Context, baseID string, req model.BuildRequest) (*model.BaseBuilding, error) {
	types, err := s.bases.GetBuildingTypes(ctx)
	if err != nil {
		return nil, err
	}

	var bt *model.BuildingType
	for _, t := range types {
		if t.ID == req.BuildingID {
			bt = &t
			break
		}
	}
	if bt == nil {
		return nil, fmt.Errorf("unknown building type: %s", req.BuildingID)
	}

	base, err := s.bases.GetByID(ctx, baseID)
	if err != nil {
		return nil, err
	}
	if base.CoinsBalance < bt.Cost {
		return nil, model.ErrInsufficientFunds
	}

	if err := s.bases.UpdateCoins(ctx, baseID, -bt.Cost); err != nil {
		return nil, err
	}

	b := &model.BaseBuilding{
		BaseID:     baseID,
		BuildingID: req.BuildingID,
		GridX:      req.GridX,
		GridY:      req.GridY,
	}
	if err := s.bases.AddBuilding(ctx, b); err != nil {
		s.bases.UpdateCoins(ctx, baseID, bt.Cost) // refund
		return nil, err
	}
	return b, nil
}

func (s *GameService) HireUnits(ctx context.Context, baseID string, req model.HireRequest) error {
	ut, err := s.army.GetUnitType(ctx, req.UnitID)
	if err != nil {
		return fmt.Errorf("unknown unit type: %s", req.UnitID)
	}

	totalCost := ut.Cost * req.Count
	base, err := s.bases.GetByID(ctx, baseID)
	if err != nil {
		return err
	}
	if base.CoinsBalance < totalCost {
		return model.ErrInsufficientFunds
	}

	if err := s.bases.UpdateCoins(ctx, baseID, -totalCost); err != nil {
		return err
	}
	return s.army.HireUnits(ctx, baseID, req.UnitID, req.Count)
}

func (s *GameService) GetBuildingTypes(ctx context.Context) ([]model.BuildingType, error) {
	return s.bases.GetBuildingTypes(ctx)
}

func (s *GameService) GetUnitTypes(ctx context.Context) ([]model.UnitType, error) {
	return s.army.GetUnitTypes(ctx)
}

func (s *GameService) GetBattleHistory(ctx context.Context, baseID string) ([]model.Battle, error) {
	return s.battles.ListByBase(ctx, baseID)
}

func (s *GameService) GetBattle(ctx context.Context, id string) (*model.Battle, error) {
	return s.battles.GetByID(ctx, id)
}
