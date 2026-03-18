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

type InviteHandler struct {
	svc *service.InviteService
}

func NewInviteHandler(svc *service.InviteService) *InviteHandler {
	return &InviteHandler{svc: svc}
}

// Create godoc
// @Summary Create invite
// @Description Send an invite to join the company (admin+ only)
// @Tags invites
// @Accept json
// @Produce json
// @Param id path string true "Company ID"
// @Param body body model.CreateInviteRequest true "Invite data"
// @Success 201 {object} model.Invite
// @Failure 400 {object} response.Response
// @Router /companies/{id}/invites [post]
func (h *InviteHandler) Create(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	companyID := c.Params("id")
	userID := middleware.GetUserID(c)

	var req model.CreateInviteRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	if req.Email == "" {
		return response.BadRequest(c, "email is required")
	}
	if req.Role == "" {
		req.Role = model.RoleEmployee
	}

	log.Printf("trace=%s | creating invite email=%s company=%s by=%s", traceID, req.Email, companyID, userID)

	invite, err := h.svc.Create(c.Context(), companyID, userID, req)
	if err != nil {
		log.Printf("trace=%s | error creating invite: %v", traceID, err)
		return response.InternalError(c)
	}

	log.Printf("trace=%s | invite created id=%s", traceID, invite.ID)
	return response.Created(c, invite)
}

// List godoc
// @Summary List invites
// @Description List all invites for a company (admin+ only)
// @Tags invites
// @Produce json
// @Param id path string true "Company ID"
// @Param page query int false "Page number"
// @Param per_page query int false "Items per page"
// @Success 200 {object} response.PaginatedResponse
// @Router /companies/{id}/invites [get]
func (h *InviteHandler) List(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	companyID := c.Params("id")

	var page model.PaginationRequest
	if err := c.Bind().Query(&page); err != nil {
		return response.BadRequest(c, "invalid pagination params")
	}

	log.Printf("trace=%s | listing invites for company=%s", traceID, companyID)

	invites, total, err := h.svc.ListByCompany(c.Context(), companyID, page)
	if err != nil {
		log.Printf("trace=%s | error listing invites: %v", traceID, err)
		return response.InternalError(c)
	}

	page.Normalize()
	totalPages := (total + page.PerPage - 1) / page.PerPage
	return response.Paginated(c, invites, response.Meta{
		Page: page.Page, PerPage: page.PerPage, Total: total, TotalPages: totalPages,
	})
}

// Revoke godoc
// @Summary Revoke invite
// @Description Revoke a pending invite (admin+ only)
// @Tags invites
// @Param id path string true "Company ID"
// @Param inviteId path string true "Invite ID"
// @Success 200 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /companies/{id}/invites/{inviteId} [delete]
func (h *InviteHandler) Revoke(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	inviteID := c.Params("inviteId")

	log.Printf("trace=%s | revoking invite id=%s", traceID, inviteID)

	if err := h.svc.Revoke(c.Context(), inviteID); err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return response.NotFound(c, "invite not found or already processed")
		}
		log.Printf("trace=%s | error revoking invite: %v", traceID, err)
		return response.InternalError(c)
	}
	return response.Success(c, fiber.Map{"revoked": true})
}

// Accept godoc
// @Summary Accept invite
// @Description Accept an invite using the token (creates membership)
// @Tags invites
// @Produce json
// @Param token path string true "Invite token"
// @Success 200 {object} model.Membership
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /invites/{token}/accept [post]
func (h *InviteHandler) Accept(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	token := c.Params("token")
	userID := middleware.GetUserID(c)
	email := middleware.GetEmail(c)

	log.Printf("trace=%s | accepting invite token=%s user=%s", traceID, token, userID)

	membership, err := h.svc.Accept(c.Context(), token, userID, email)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrNotFound):
			return response.NotFound(c, "invite not found or already used")
		case errors.Is(err, model.ErrInviteExpired):
			return response.BadRequest(c, "invite has expired")
		case errors.Is(err, model.ErrEmailMismatch):
			return response.Forbidden(c, "email does not match invite")
		default:
			log.Printf("trace=%s | error accepting invite: %v", traceID, err)
			return response.InternalError(c)
		}
	}

	log.Printf("trace=%s | invite accepted, membership id=%s", traceID, membership.ID)
	return response.Success(c, membership)
}
