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

type AchievementHandler struct {
	svc *service.AchievementService
}

func NewAchievementHandler(svc *service.AchievementService) *AchievementHandler {
	return &AchievementHandler{svc: svc}
}

// Create godoc
// @Summary Create an achievement rule
// @Description Define a new achievement with metric conditions (admin+ only)
// @Tags achievements
// @Accept json
// @Produce json
// @Param id path string true "Company ID"
// @Param body body model.CreateAchievementRequest true "Achievement data"
// @Success 201 {object} model.Achievement
// @Failure 400 {object} response.Response
// @Router /companies/{id}/achievements [post]
func (h *AchievementHandler) Create(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	companyID := c.Params("id")

	var req model.CreateAchievementRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	if req.Name == "" || req.Metric == "" {
		return response.BadRequest(c, "name and metric are required")
	}
	if req.Operator == "" {
		req.Operator = model.OpGTE
	}

	log.Printf("trace=%s | creating achievement name=%s metric=%s company=%s", traceID, req.Name, req.Metric, companyID)

	achievement, err := h.svc.Create(c.Context(), companyID, req)
	if err != nil {
		log.Printf("trace=%s | error creating achievement: %v", traceID, err)
		return response.InternalError(c)
	}

	log.Printf("trace=%s | achievement created id=%s", traceID, achievement.ID)
	return response.Created(c, achievement)
}

// List godoc
// @Summary List achievements
// @Description List all achievement rules for a company
// @Tags achievements
// @Produce json
// @Param id path string true "Company ID"
// @Success 200 {object} response.Response
// @Router /companies/{id}/achievements [get]
func (h *AchievementHandler) List(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	companyID := c.Params("id")

	log.Printf("trace=%s | listing achievements for company=%s", traceID, companyID)

	achievements, err := h.svc.ListByCompany(c.Context(), companyID)
	if err != nil {
		log.Printf("trace=%s | error listing achievements: %v", traceID, err)
		return response.InternalError(c)
	}
	return response.Success(c, achievements)
}

// Delete godoc
// @Summary Delete an achievement
// @Description Delete an achievement rule (admin+ only)
// @Tags achievements
// @Param id path string true "Company ID"
// @Param achievementId path string true "Achievement ID"
// @Success 200 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /companies/{id}/achievements/{achievementId} [delete]
func (h *AchievementHandler) Delete(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	achievementID := c.Params("achievementId")

	log.Printf("trace=%s | deleting achievement id=%s", traceID, achievementID)

	if err := h.svc.Delete(c.Context(), achievementID); err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return response.NotFound(c, "achievement not found")
		}
		log.Printf("trace=%s | error deleting achievement: %v", traceID, err)
		return response.InternalError(c)
	}
	return response.Success(c, fiber.Map{"deleted": true})
}

// EvaluateMetric godoc
// @Summary Report a metric value
// @Description Report a metric value for a member, triggers achievement evaluation
// @Tags achievements
// @Accept json
// @Produce json
// @Param id path string true "Company ID"
// @Param memberId path string true "Membership ID"
// @Param body body model.EvaluateMetricRequest true "Metric data"
// @Success 200 {object} response.Response
// @Router /companies/{id}/members/{memberId}/metrics [post]
func (h *AchievementHandler) EvaluateMetric(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	companyID := c.Params("id")
	memberID := c.Params("memberId")

	var req model.EvaluateMetricRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	if req.Metric == "" {
		return response.BadRequest(c, "metric is required")
	}

	log.Printf("trace=%s | evaluating metric=%s value=%d for member=%s company=%s", traceID, req.Metric, req.Value, memberID, companyID)

	completed, err := h.svc.EvaluateMetric(c.Context(), memberID, companyID, req.Metric, req.Value)
	if err != nil {
		log.Printf("trace=%s | error evaluating metric: %v", traceID, err)
		return response.InternalError(c)
	}

	log.Printf("trace=%s | %d achievements completed", traceID, len(completed))
	return response.Success(c, fiber.Map{
		"completed_count": len(completed),
		"completed":       completed,
	})
}

// ListMemberAchievements godoc
// @Summary List member's completed achievements
// @Description Returns all achievements completed by a member
// @Tags achievements
// @Produce json
// @Param id path string true "Company ID"
// @Param memberId path string true "Membership ID"
// @Success 200 {object} response.Response
// @Router /companies/{id}/members/{memberId}/achievements [get]
func (h *AchievementHandler) ListMemberAchievements(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	memberID := c.Params("memberId")

	log.Printf("trace=%s | listing achievements for member=%s", traceID, memberID)

	achievements, err := h.svc.ListMemberAchievements(c.Context(), memberID)
	if err != nil {
		log.Printf("trace=%s | error listing member achievements: %v", traceID, err)
		return response.InternalError(c)
	}
	return response.Success(c, achievements)
}
