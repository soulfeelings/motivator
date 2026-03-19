package handler

import (
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/requestid"

	"github.com/hustlers/motivator-backend/internal/middleware"
	"github.com/hustlers/motivator-backend/internal/model"
	"github.com/hustlers/motivator-backend/internal/service"
	"github.com/hustlers/motivator-backend/pkg/response"
)

type NotificationHandler struct {
	svc *service.NotificationService
}

func NewNotificationHandler(svc *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{svc: svc}
}

// RegisterToken godoc
// @Summary Register device token
// @Description Register a push notification device token
// @Tags notifications
// @Accept json
// @Produce json
// @Param id path string true "Company ID"
// @Param body body model.RegisterTokenRequest true "Token data"
// @Success 201 {object} model.DeviceToken
// @Router /companies/{id}/notifications/register [post]
func (h *NotificationHandler) RegisterToken(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	m := middleware.GetMembership(c)

	var req model.RegisterTokenRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	if req.Token == "" {
		return response.BadRequest(c, "token is required")
	}
	if req.Platform == "" {
		req.Platform = "android"
	}

	log.Printf("trace=%s | registering device token for member=%s platform=%s", traceID, m.ID, req.Platform)

	dt, err := h.svc.RegisterToken(c.Context(), m.ID, req)
	if err != nil {
		log.Printf("trace=%s | error registering token: %v", traceID, err)
		return response.InternalError(c)
	}
	return response.Created(c, dt)
}

// UnregisterToken godoc
// @Summary Unregister device token
// @Description Remove a push notification device token
// @Tags notifications
// @Accept json
// @Produce json
// @Param id path string true "Company ID"
// @Param body body model.RegisterTokenRequest true "Token to remove"
// @Success 200 {object} response.Response
// @Router /companies/{id}/notifications/unregister [post]
func (h *NotificationHandler) UnregisterToken(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	m := middleware.GetMembership(c)

	var req model.RegisterTokenRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	log.Printf("trace=%s | unregistering device token for member=%s", traceID, m.ID)

	if err := h.svc.UnregisterToken(c.Context(), m.ID, req.Token); err != nil {
		log.Printf("trace=%s | error unregistering token: %v", traceID, err)
		return response.InternalError(c)
	}
	return response.Success(c, fiber.Map{"unregistered": true})
}
