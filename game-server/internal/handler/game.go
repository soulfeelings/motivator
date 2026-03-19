package handler

import (
	"errors"
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/requestid"

	"github.com/hustlers/motivator-game/internal/middleware"
	"github.com/hustlers/motivator-game/internal/model"
	"github.com/hustlers/motivator-game/internal/service"
	"github.com/hustlers/motivator-game/pkg/response"
)

type GameHandler struct {
	svc *service.GameService
}

func NewGameHandler(svc *service.GameService) *GameHandler {
	return &GameHandler{svc: svc}
}

func (h *GameHandler) GetMyBase(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	membershipID := c.Query("membership_id")
	if membershipID == "" {
		return response.BadRequest(c, "membership_id query param required")
	}

	log.Printf("trace=%s | get/create base for member=%s", traceID, membershipID)

	overview, err := h.svc.GetOrCreateBase(c.Context(), membershipID)
	if err != nil {
		log.Printf("trace=%s | error: %v", traceID, err)
		return response.InternalError(c)
	}
	return response.Success(c, overview)
}

func (h *GameHandler) GetBase(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	baseID := c.Params("baseId")

	log.Printf("trace=%s | get base=%s", traceID, baseID)

	overview, err := h.svc.GetBaseOverview(c.Context(), baseID)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return response.NotFound(c, "base not found")
		}
		return response.InternalError(c)
	}
	return response.Success(c, overview)
}

func (h *GameHandler) ListBases(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	log.Printf("trace=%s | listing all bases", traceID)

	bases, err := h.svc.ListBases(c.Context())
	if err != nil {
		return response.InternalError(c)
	}
	return response.Success(c, bases)
}

func (h *GameHandler) Build(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	baseID := c.Params("baseId")

	var req model.BuildRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	log.Printf("trace=%s | building %s at base=%s", traceID, req.BuildingID, baseID)

	building, err := h.svc.Build(c.Context(), baseID, req)
	if err != nil {
		if errors.Is(err, model.ErrInsufficientFunds) {
			return response.BadRequest(c, "not enough coins")
		}
		log.Printf("trace=%s | error building: %v", traceID, err)
		return response.BadRequest(c, err.Error())
	}
	return response.Created(c, building)
}

func (h *GameHandler) HireUnits(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	baseID := c.Params("baseId")

	var req model.HireRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	log.Printf("trace=%s | hiring %d %s at base=%s", traceID, req.Count, req.UnitID, baseID)

	if err := h.svc.HireUnits(c.Context(), baseID, req); err != nil {
		if errors.Is(err, model.ErrInsufficientFunds) {
			return response.BadRequest(c, "not enough coins")
		}
		return response.BadRequest(c, err.Error())
	}
	return response.Success(c, fiber.Map{"hired": true})
}

func (h *GameHandler) DepositCoins(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	baseID := c.Params("baseId")

	var req model.DepositCoinsRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	log.Printf("trace=%s | depositing %d coins to base=%s", traceID, req.Amount, baseID)

	if err := h.svc.DepositCoins(c.Context(), baseID, req.Amount); err != nil {
		return response.InternalError(c)
	}
	return response.Success(c, fiber.Map{"deposited": true})
}

func (h *GameHandler) Attack(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	baseID := c.Params("baseId")

	var req model.AttackRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	log.Printf("trace=%s | base=%s attacking base=%s", traceID, baseID, req.DefenderBaseID)

	battle, err := h.svc.Attack(c.Context(), baseID, req.DefenderBaseID)
	if err != nil {
		if errors.Is(err, model.ErrCannotAttackSelf) {
			return response.BadRequest(c, "cannot attack your own base")
		}
		log.Printf("trace=%s | error in battle: %v", traceID, err)
		return response.InternalError(c)
	}
	return response.Success(c, battle)
}

func (h *GameHandler) GetBattleHistory(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	baseID := c.Params("baseId")

	log.Printf("trace=%s | battle history for base=%s", traceID, baseID)

	battles, err := h.svc.GetBattleHistory(c.Context(), baseID)
	if err != nil {
		return response.InternalError(c)
	}
	return response.Success(c, battles)
}

func (h *GameHandler) GetBattle(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	battleID := c.Params("battleId")

	log.Printf("trace=%s | get battle=%s", traceID, battleID)

	battle, err := h.svc.GetBattle(c.Context(), battleID)
	if err != nil {
		return response.NotFound(c, "battle not found")
	}
	return response.Success(c, battle)
}

func (h *GameHandler) GetBuildingTypes(c fiber.Ctx) error {
	types, err := h.svc.GetBuildingTypes(c.Context())
	if err != nil {
		return response.InternalError(c)
	}
	return response.Success(c, types)
}

func (h *GameHandler) GetUnitTypes(c fiber.Ctx) error {
	types, err := h.svc.GetUnitTypes(c.Context())
	if err != nil {
		return response.InternalError(c)
	}
	return response.Success(c, types)
}

func RegisterRoutes(app *fiber.App, h *GameHandler, auth *middleware.AuthMiddleware) {
	api := app.Group("/api/v1/game", auth.Required())

	api.Get("/base", h.GetMyBase)
	api.Get("/bases", h.ListBases)
	api.Get("/bases/:baseId", h.GetBase)
	api.Post("/bases/:baseId/build", h.Build)
	api.Post("/bases/:baseId/hire", h.HireUnits)
	api.Post("/bases/:baseId/deposit", h.DepositCoins)
	api.Post("/bases/:baseId/attack", h.Attack)
	api.Get("/bases/:baseId/battles", h.GetBattleHistory)
	api.Get("/battles/:battleId", h.GetBattle)
	api.Get("/building-types", h.GetBuildingTypes)
	api.Get("/unit-types", h.GetUnitTypes)
}
