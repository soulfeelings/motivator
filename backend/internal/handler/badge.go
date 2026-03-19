package handler

import (
	"errors"
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/requestid"

	"github.com/hustlers/motivator-backend/internal/middleware"
	"github.com/hustlers/motivator-backend/internal/model"
	"github.com/hustlers/motivator-backend/internal/service"
	"github.com/hustlers/motivator-backend/pkg/response"
)

type BadgeHandler struct {
	svc *service.BadgeService
}

func NewBadgeHandler(svc *service.BadgeService) *BadgeHandler {
	return &BadgeHandler{svc: svc}
}

// Create godoc
// @Summary Create a badge
// @Description Create a new badge for the company (admin+ only)
// @Tags badges
// @Accept json
// @Produce json
// @Param id path string true "Company ID"
// @Param body body model.CreateBadgeRequest true "Badge data"
// @Success 201 {object} model.Badge
// @Failure 400 {object} response.Response
// @Router /companies/{id}/badges [post]
func (h *BadgeHandler) Create(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	companyID := c.Params("id")

	var req model.CreateBadgeRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	if req.Name == "" {
		return response.BadRequest(c, "name is required")
	}

	log.Printf("trace=%s | creating badge name=%s company=%s", traceID, req.Name, companyID)

	badge, err := h.svc.Create(c.Context(), companyID, req)
	if err != nil {
		log.Printf("trace=%s | error creating badge: %v", traceID, err)
		return response.InternalError(c)
	}

	log.Printf("trace=%s | badge created id=%s", traceID, badge.ID)
	return response.Created(c, badge)
}

// List godoc
// @Summary List badges
// @Description List all badges for a company
// @Tags badges
// @Produce json
// @Param id path string true "Company ID"
// @Success 200 {object} response.Response
// @Router /companies/{id}/badges [get]
func (h *BadgeHandler) List(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	companyID := c.Params("id")

	log.Printf("trace=%s | listing badges for company=%s", traceID, companyID)

	badges, err := h.svc.ListByCompany(c.Context(), companyID)
	if err != nil {
		log.Printf("trace=%s | error listing badges: %v", traceID, err)
		return response.InternalError(c)
	}
	return response.Success(c, badges)
}

// Delete godoc
// @Summary Delete a badge
// @Description Delete a badge (admin+ only)
// @Tags badges
// @Param id path string true "Company ID"
// @Param badgeId path string true "Badge ID"
// @Success 200 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /companies/{id}/badges/{badgeId} [delete]
func (h *BadgeHandler) Delete(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	badgeID := c.Params("badgeId")

	log.Printf("trace=%s | deleting badge id=%s", traceID, badgeID)

	if err := h.svc.Delete(c.Context(), badgeID); err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return response.NotFound(c, "badge not found")
		}
		log.Printf("trace=%s | error deleting badge: %v", traceID, err)
		return response.InternalError(c)
	}
	return response.Success(c, fiber.Map{"deleted": true})
}

// Award godoc
// @Summary Award badge to member
// @Description Award a badge to a member and grant XP/coins (admin+ only)
// @Tags badges
// @Accept json
// @Produce json
// @Param id path string true "Company ID"
// @Param memberId path string true "Membership ID"
// @Param body body model.AwardBadgeRequest true "Badge to award"
// @Success 200 {object} model.MemberBadge
// @Router /companies/{id}/members/{memberId}/badges [post]
func (h *BadgeHandler) Award(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	memberID := c.Params("memberId")
	_ = middleware.GetUserID(c)

	var req model.AwardBadgeRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	if req.BadgeID == "" {
		return response.BadRequest(c, "badge_id is required")
	}

	log.Printf("trace=%s | awarding badge=%s to member=%s", traceID, req.BadgeID, memberID)

	mb, err := h.svc.AwardBadge(c.Context(), memberID, req.BadgeID)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return response.NotFound(c, "badge or member not found")
		}
		log.Printf("trace=%s | error awarding badge: %v", traceID, err)
		return response.InternalError(c)
	}

	log.Printf("trace=%s | badge awarded id=%s", traceID, mb.ID)
	return response.Success(c, mb)
}
