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

type CompanyHandler struct {
	svc *service.CompanyService
}

func NewCompanyHandler(svc *service.CompanyService) *CompanyHandler {
	return &CompanyHandler{svc: svc}
}

// Create godoc
// @Summary Create a company
// @Description Creates a new company and assigns the caller as owner
// @Tags companies
// @Accept json
// @Produce json
// @Param body body model.CreateCompanyRequest true "Company data"
// @Success 201 {object} model.Company
// @Failure 400 {object} response.Response
// @Router /companies [post]
func (h *CompanyHandler) Create(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	userID := middleware.GetUserID(c)

	var req model.CreateCompanyRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	if req.Name == "" || req.Slug == "" {
		return response.BadRequest(c, "name and slug are required")
	}

	log.Printf("trace=%s | creating company slug=%s for user=%s", traceID, req.Slug, userID)

	company, _, err := h.svc.Create(c.Context(), userID, req)
	if err != nil {
		log.Printf("trace=%s | error creating company: %v", traceID, err)
		return response.InternalError(c)
	}

	log.Printf("trace=%s | company created id=%s", traceID, company.ID)
	return response.Created(c, company)
}

// GetByID godoc
// @Summary Get company by ID
// @Description Returns company details
// @Tags companies
// @Produce json
// @Param id path string true "Company ID"
// @Success 200 {object} model.Company
// @Failure 404 {object} response.Response
// @Router /companies/{id} [get]
func (h *CompanyHandler) GetByID(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	id := c.Params("id")

	log.Printf("trace=%s | fetching company id=%s", traceID, id)

	company, err := h.svc.GetByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return response.NotFound(c, "company not found")
		}
		log.Printf("trace=%s | error fetching company: %v", traceID, err)
		return response.InternalError(c)
	}
	return response.Success(c, company)
}

// Update godoc
// @Summary Update company
// @Description Updates company details (owner/admin only)
// @Tags companies
// @Accept json
// @Produce json
// @Param id path string true "Company ID"
// @Param body body model.UpdateCompanyRequest true "Update data"
// @Success 200 {object} model.Company
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /companies/{id} [patch]
func (h *CompanyHandler) Update(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	id := c.Params("id")

	var req model.UpdateCompanyRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	log.Printf("trace=%s | updating company id=%s", traceID, id)

	company, err := h.svc.Update(c.Context(), id, req)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return response.NotFound(c, "company not found")
		}
		log.Printf("trace=%s | error updating company: %v", traceID, err)
		return response.InternalError(c)
	}
	return response.Success(c, company)
}

// Delete godoc
// @Summary Delete company
// @Description Deletes a company (owner only)
// @Tags companies
// @Param id path string true "Company ID"
// @Success 200 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /companies/{id} [delete]
func (h *CompanyHandler) Delete(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	id := c.Params("id")

	log.Printf("trace=%s | deleting company id=%s", traceID, id)

	if err := h.svc.Delete(c.Context(), id); err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return response.NotFound(c, "company not found")
		}
		log.Printf("trace=%s | error deleting company: %v", traceID, err)
		return response.InternalError(c)
	}

	log.Printf("trace=%s | company deleted id=%s", traceID, id)
	return response.Success(c, fiber.Map{"deleted": true})
}
