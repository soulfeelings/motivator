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

type IntegrationHandler struct {
	svc *service.IntegrationService
}

func NewIntegrationHandler(svc *service.IntegrationService) *IntegrationHandler {
	return &IntegrationHandler{svc: svc}
}

func (h *IntegrationHandler) Create(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	companyID := c.Params("id")
	var req model.CreateIntegrationRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	if req.Provider == "" || req.Name == "" {
		return response.BadRequest(c, "provider and name are required")
	}
	log.Printf("trace=%s | creating integration provider=%s name=%s company=%s", traceID, req.Provider, req.Name, companyID)
	i, err := h.svc.Create(c.Context(), companyID, req)
	if err != nil {
		return response.InternalError(c)
	}
	return response.Created(c, i)
}

func (h *IntegrationHandler) List(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	companyID := c.Params("id")
	log.Printf("trace=%s | listing integrations company=%s", traceID, companyID)
	is, err := h.svc.ListByCompany(c.Context(), companyID)
	if err != nil {
		return response.InternalError(c)
	}
	return response.Success(c, is)
}

func (h *IntegrationHandler) Delete(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	integrationID := c.Params("integrationId")
	log.Printf("trace=%s | deleting integration=%s", traceID, integrationID)
	if err := h.svc.Delete(c.Context(), integrationID); err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return response.NotFound(c, "integration not found")
		}
		return response.InternalError(c)
	}
	return response.Success(c, fiber.Map{"deleted": true})
}

func (h *IntegrationHandler) CreateMapping(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	integrationID := c.Params("integrationId")
	var req model.CreateMappingRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	log.Printf("trace=%s | creating mapping event=%s metric=%s integration=%s", traceID, req.ExternalEvent, req.Metric, integrationID)
	m, err := h.svc.CreateMapping(c.Context(), integrationID, req)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}
	return response.Created(c, m)
}

func (h *IntegrationHandler) ListMappings(c fiber.Ctx) error {
	integrationID := c.Params("integrationId")
	ms, err := h.svc.ListMappings(c.Context(), integrationID)
	if err != nil {
		return response.InternalError(c)
	}
	return response.Success(c, ms)
}

func (h *IntegrationHandler) DeleteMapping(c fiber.Ctx) error {
	mappingID := c.Params("mappingId")
	if err := h.svc.DeleteMapping(c.Context(), mappingID); err != nil {
		return response.InternalError(c)
	}
	return response.Success(c, fiber.Map{"deleted": true})
}

func (h *IntegrationHandler) ListEvents(c fiber.Ctx) error {
	integrationID := c.Params("integrationId")
	events, err := h.svc.ListEvents(c.Context(), integrationID)
	if err != nil {
		return response.InternalError(c)
	}
	return response.Success(c, events)
}

// InboundWebhook is the public endpoint that external services (Jira, GitHub, etc.) call.
// No auth required — validated by webhook secret in URL.
func (h *IntegrationHandler) InboundWebhook(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	secret := c.Params("secret")

	var body map[string]any
	if err := c.Bind().JSON(&body); err != nil {
		return response.BadRequest(c, "invalid JSON body")
	}

	log.Printf("trace=%s | inbound webhook secret=%s", traceID, secret[:8])

	if err := h.svc.ProcessWebhook(c.Context(), secret, body); err != nil {
		log.Printf("trace=%s | webhook processing error: %v", traceID, err)
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, fiber.Map{"received": true})
}
