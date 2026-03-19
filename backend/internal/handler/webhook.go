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

type WebhookHandler struct {
	svc *service.WebhookService
}

func NewWebhookHandler(svc *service.WebhookService) *WebhookHandler {
	return &WebhookHandler{svc: svc}
}

func (h *WebhookHandler) Create(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	companyID := c.Params("id")
	var req model.CreateWebhookRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	if req.Name == "" || req.URL == "" || req.Platform == "" {
		return response.BadRequest(c, "name, url, and platform are required")
	}
	log.Printf("trace=%s | creating webhook name=%s platform=%s company=%s", traceID, req.Name, req.Platform, companyID)
	w, err := h.svc.Create(c.Context(), companyID, req)
	if err != nil {
		return response.InternalError(c)
	}
	return response.Created(c, w)
}

func (h *WebhookHandler) List(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	companyID := c.Params("id")
	log.Printf("trace=%s | listing webhooks company=%s", traceID, companyID)
	ws, err := h.svc.ListByCompany(c.Context(), companyID)
	if err != nil {
		return response.InternalError(c)
	}
	return response.Success(c, ws)
}

func (h *WebhookHandler) Delete(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	webhookID := c.Params("webhookId")
	log.Printf("trace=%s | deleting webhook=%s", traceID, webhookID)
	if err := h.svc.Delete(c.Context(), webhookID); err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return response.NotFound(c, "webhook not found")
		}
		return response.InternalError(c)
	}
	return response.Success(c, fiber.Map{"deleted": true})
}
