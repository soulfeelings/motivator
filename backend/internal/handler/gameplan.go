package handler

import (
	"errors"
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/requestid"

	"github.com/hustlers/motivator-backend/internal/model"
	"github.com/hustlers/motivator-backend/internal/service"
	"github.com/hustlers/motivator-backend/pkg/response"
)

type GamePlanHandler struct {
	svc *service.GamePlanService
}

func NewGamePlanHandler(svc *service.GamePlanService) *GamePlanHandler {
	return &GamePlanHandler{svc: svc}
}

// Create godoc
// @Summary Create a game plan
// @Description Create a new game plan with visual flow data (admin+ only)
// @Tags game-plans
// @Accept json
// @Produce json
// @Param id path string true "Company ID"
// @Param body body model.CreateGamePlanRequest true "Game plan data"
// @Success 201 {object} model.GamePlan
// @Router /companies/{id}/game-plans [post]
func (h *GamePlanHandler) Create(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	companyID := c.Params("id")

	var req model.CreateGamePlanRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	if req.Name == "" {
		return response.BadRequest(c, "name is required")
	}

	log.Printf("trace=%s | creating game plan name=%s company=%s", traceID, req.Name, companyID)

	gp, err := h.svc.Create(c.Context(), companyID, req)
	if err != nil {
		log.Printf("trace=%s | error creating game plan: %v", traceID, err)
		return response.InternalError(c)
	}
	return response.Created(c, gp)
}

// GetByID godoc
// @Summary Get game plan
// @Description Get a game plan with flow data
// @Tags game-plans
// @Produce json
// @Param id path string true "Company ID"
// @Param planId path string true "Game Plan ID"
// @Success 200 {object} model.GamePlan
// @Router /companies/{id}/game-plans/{planId} [get]
func (h *GamePlanHandler) GetByID(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	planID := c.Params("planId")

	log.Printf("trace=%s | fetching game plan id=%s", traceID, planID)

	gp, err := h.svc.GetByID(c.Context(), planID)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return response.NotFound(c, "game plan not found")
		}
		return response.InternalError(c)
	}
	return response.Success(c, gp)
}

// List godoc
// @Summary List game plans
// @Description List all game plans for a company
// @Tags game-plans
// @Produce json
// @Param id path string true "Company ID"
// @Success 200 {object} response.Response
// @Router /companies/{id}/game-plans [get]
func (h *GamePlanHandler) List(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	companyID := c.Params("id")

	log.Printf("trace=%s | listing game plans company=%s", traceID, companyID)

	plans, err := h.svc.ListByCompany(c.Context(), companyID)
	if err != nil {
		log.Printf("trace=%s | error listing game plans: %v", traceID, err)
		return response.InternalError(c)
	}
	return response.Success(c, plans)
}

// Update godoc
// @Summary Update game plan
// @Description Update game plan name/description/flow (admin+ only)
// @Tags game-plans
// @Accept json
// @Produce json
// @Param id path string true "Company ID"
// @Param planId path string true "Game Plan ID"
// @Param body body model.UpdateGamePlanRequest true "Update data"
// @Success 200 {object} model.GamePlan
// @Router /companies/{id}/game-plans/{planId} [patch]
func (h *GamePlanHandler) Update(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	planID := c.Params("planId")

	var req model.UpdateGamePlanRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	log.Printf("trace=%s | updating game plan id=%s", traceID, planID)

	gp, err := h.svc.Update(c.Context(), planID, req)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return response.NotFound(c, "game plan not found")
		}
		return response.InternalError(c)
	}
	return response.Success(c, gp)
}

// SaveFlow godoc
// @Summary Save flow data
// @Description Save the visual flow data for a game plan
// @Tags game-plans
// @Accept json
// @Produce json
// @Param id path string true "Company ID"
// @Param planId path string true "Game Plan ID"
// @Param body body model.FlowData true "Flow data"
// @Success 200 {object} response.Response
// @Router /companies/{id}/game-plans/{planId}/flow [put]
func (h *GamePlanHandler) SaveFlow(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	planID := c.Params("planId")

	var flowData model.FlowData
	if err := c.Bind().JSON(&flowData); err != nil {
		return response.BadRequest(c, "invalid flow data")
	}

	log.Printf("trace=%s | saving flow for plan=%s nodes=%d edges=%d", traceID, planID, len(flowData.Nodes), len(flowData.Edges))

	if err := h.svc.SaveFlow(c.Context(), planID, flowData); err != nil {
		log.Printf("trace=%s | error saving flow: %v", traceID, err)
		return response.InternalError(c)
	}
	return response.Success(c, fiber.Map{"saved": true})
}

// Activate godoc
// @Summary Activate/deactivate game plan
// @Description Toggle game plan active state (admin+ only)
// @Tags game-plans
// @Produce json
// @Param id path string true "Company ID"
// @Param planId path string true "Game Plan ID"
// @Success 200 {object} response.Response
// @Router /companies/{id}/game-plans/{planId}/activate [post]
func (h *GamePlanHandler) Activate(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	planID := c.Params("planId")

	log.Printf("trace=%s | activating game plan id=%s", traceID, planID)

	if err := h.svc.SetActive(c.Context(), planID, true); err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return response.NotFound(c, "game plan not found")
		}
		return response.InternalError(c)
	}
	return response.Success(c, fiber.Map{"activated": true})
}

// Deactivate godoc
// @Summary Deactivate game plan
// @Tags game-plans
// @Param id path string true "Company ID"
// @Param planId path string true "Game Plan ID"
// @Success 200 {object} response.Response
// @Router /companies/{id}/game-plans/{planId}/deactivate [post]
func (h *GamePlanHandler) Deactivate(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	planID := c.Params("planId")

	log.Printf("trace=%s | deactivating game plan id=%s", traceID, planID)

	if err := h.svc.SetActive(c.Context(), planID, false); err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return response.NotFound(c, "game plan not found")
		}
		return response.InternalError(c)
	}
	return response.Success(c, fiber.Map{"deactivated": true})
}

// Delete godoc
// @Summary Delete game plan
// @Tags game-plans
// @Param id path string true "Company ID"
// @Param planId path string true "Game Plan ID"
// @Success 200 {object} response.Response
// @Router /companies/{id}/game-plans/{planId} [delete]
func (h *GamePlanHandler) Delete(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	planID := c.Params("planId")

	log.Printf("trace=%s | deleting game plan id=%s", traceID, planID)

	if err := h.svc.Delete(c.Context(), planID); err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return response.NotFound(c, "game plan not found")
		}
		return response.InternalError(c)
	}
	return response.Success(c, fiber.Map{"deleted": true})
}
