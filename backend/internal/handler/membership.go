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

type MembershipHandler struct {
	svc           *service.MembershipService
	supabaseAdmin *service.SupabaseAdmin
}

func NewMembershipHandler(svc *service.MembershipService, supabaseAdmin *service.SupabaseAdmin) *MembershipHandler {
	return &MembershipHandler{svc: svc, supabaseAdmin: supabaseAdmin}
}

// List godoc
// @Summary List company members
// @Description Returns paginated list of active members
// @Tags members
// @Produce json
// @Param id path string true "Company ID"
// @Param page query int false "Page number"
// @Param per_page query int false "Items per page"
// @Success 200 {object} response.PaginatedResponse
// @Router /companies/{id}/members [get]
func (h *MembershipHandler) List(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	companyID := c.Params("id")

	var page model.PaginationRequest
	if err := c.Bind().Query(&page); err != nil {
		return response.BadRequest(c, "invalid pagination params")
	}

	log.Printf("trace=%s | listing members for company=%s", traceID, companyID)

	members, total, err := h.svc.ListByCompany(c.Context(), companyID, page)
	if err != nil {
		log.Printf("trace=%s | error listing members: %v", traceID, err)
		return response.InternalError(c)
	}

	page.Normalize()
	totalPages := (total + page.PerPage - 1) / page.PerPage
	return response.Paginated(c, members, response.Meta{
		Page: page.Page, PerPage: page.PerPage, Total: total, TotalPages: totalPages,
	})
}

// GetByID godoc
// @Summary Get member by ID
// @Description Returns a single membership
// @Tags members
// @Produce json
// @Param id path string true "Company ID"
// @Param memberId path string true "Membership ID"
// @Success 200 {object} model.Membership
// @Failure 404 {object} response.Response
// @Router /companies/{id}/members/{memberId} [get]
func (h *MembershipHandler) GetByID(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	memberID := c.Params("memberId")

	log.Printf("trace=%s | fetching member id=%s", traceID, memberID)

	member, err := h.svc.GetByID(c.Context(), memberID)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return response.NotFound(c, "member not found")
		}
		log.Printf("trace=%s | error fetching member: %v", traceID, err)
		return response.InternalError(c)
	}
	return response.Success(c, member)
}

// Update godoc
// @Summary Update member
// @Description Update member role, display name, or job title (admin+ only)
// @Tags members
// @Accept json
// @Produce json
// @Param id path string true "Company ID"
// @Param memberId path string true "Membership ID"
// @Param body body model.UpdateMembershipRequest true "Update data"
// @Success 200 {object} model.Membership
// @Router /companies/{id}/members/{memberId} [patch]
func (h *MembershipHandler) Update(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	memberID := c.Params("memberId")

	var req model.UpdateMembershipRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	log.Printf("trace=%s | updating member id=%s", traceID, memberID)

	member, err := h.svc.Update(c.Context(), memberID, req)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return response.NotFound(c, "member not found")
		}
		log.Printf("trace=%s | error updating member: %v", traceID, err)
		return response.InternalError(c)
	}
	return response.Success(c, member)
}

// Deactivate godoc
// @Summary Remove member
// @Description Deactivate a member (admin+ only)
// @Tags members
// @Param id path string true "Company ID"
// @Param memberId path string true "Membership ID"
// @Success 200 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /companies/{id}/members/{memberId} [delete]
func (h *MembershipHandler) Deactivate(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	memberID := c.Params("memberId")

	log.Printf("trace=%s | deactivating member id=%s", traceID, memberID)

	if err := h.svc.Deactivate(c.Context(), memberID); err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return response.NotFound(c, "member not found")
		}
		log.Printf("trace=%s | error deactivating member: %v", traceID, err)
		return response.InternalError(c)
	}
	return response.Success(c, fiber.Map{"deactivated": true})
}

// AddMember godoc
// @Summary Add a new member
// @Description Creates a Supabase Auth user and adds them as a company member in one step
// @Tags members
// @Accept json
// @Produce json
// @Param id path string true "Company ID"
// @Param body body model.AddMemberRequest true "Member data"
// @Success 201 {object} model.AddMemberResponse
// @Router /companies/{id}/members/add [post]
func (h *MembershipHandler) AddMember(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	companyID := c.Params("id")

	var req model.AddMemberRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	if req.Email == "" || req.Password == "" {
		return response.BadRequest(c, "email and password are required")
	}
	if req.Role == "" {
		req.Role = model.RoleEmployee
	}

	log.Printf("trace=%s | adding member email=%s role=%s company=%s", traceID, req.Email, req.Role, companyID)

	// Create user in Supabase Auth
	user, err := h.supabaseAdmin.CreateUser(c.Context(), req.Email, req.Password)
	if err != nil {
		log.Printf("trace=%s | error creating supabase user: %v", traceID, err)
		return response.BadRequest(c, err.Error())
	}

	// Create membership
	membership, err := h.svc.AddToCompany(c.Context(), user.ID, companyID, req.Role, req.DisplayName)
	if err != nil {
		log.Printf("trace=%s | error creating membership: %v", traceID, err)
		return response.InternalError(c)
	}

	log.Printf("trace=%s | member added user=%s membership=%s", traceID, user.ID, membership.ID)
	return response.Created(c, model.AddMemberResponse{
		UserID:     user.ID,
		Email:      user.Email,
		Membership: *membership,
	})
}

// GetProfile godoc
// @Summary Get member profile
// @Description Returns member profile with XP, level, coins, and badges
// @Tags members
// @Produce json
// @Param id path string true "Company ID"
// @Param memberId path string true "Membership ID"
// @Success 200 {object} model.ProfileResponse
// @Router /companies/{id}/members/{memberId}/profile [get]
func (h *MembershipHandler) GetProfile(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	memberID := c.Params("memberId")

	log.Printf("trace=%s | fetching profile for member=%s", traceID, memberID)

	profile, err := h.svc.GetProfile(c.Context(), memberID)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return response.NotFound(c, "member not found")
		}
		log.Printf("trace=%s | error fetching profile: %v", traceID, err)
		return response.InternalError(c)
	}
	return response.Success(c, profile)
}

// AwardXP godoc
// @Summary Award XP to member
// @Description Award XP points to a member (admin+ only)
// @Tags members
// @Accept json
// @Produce json
// @Param id path string true "Company ID"
// @Param memberId path string true "Membership ID"
// @Param body body model.AwardXPRequest true "XP amount"
// @Success 200 {object} model.Membership
// @Router /companies/{id}/members/{memberId}/xp [post]
func (h *MembershipHandler) AwardXP(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	memberID := c.Params("memberId")

	var req model.AwardXPRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	if req.Amount <= 0 {
		return response.BadRequest(c, "amount must be positive")
	}

	log.Printf("trace=%s | awarding %d XP to member=%s reason=%s", traceID, req.Amount, memberID, req.Reason)

	member, err := h.svc.AwardXP(c.Context(), memberID, req.Amount)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return response.NotFound(c, "member not found")
		}
		log.Printf("trace=%s | error awarding XP: %v", traceID, err)
		return response.InternalError(c)
	}
	return response.Success(c, member)
}

// Me godoc
// @Summary Get current user profile
// @Description Returns current user's memberships across all companies
// @Tags user
// @Produce json
// @Success 200 {object} response.Response
// @Router /me [get]
func (h *MembershipHandler) Me(c fiber.Ctx) error {
	traceID := requestid.FromContext(c)
	userID := middleware.GetUserID(c)

	log.Printf("trace=%s | fetching profile for user=%s", traceID, userID)

	memberships, err := h.svc.ListByUser(c.Context(), userID)
	if err != nil {
		log.Printf("trace=%s | error fetching user memberships: %v", traceID, err)
		return response.InternalError(c)
	}

	return response.Success(c, fiber.Map{
		"user_id":     userID,
		"email":       middleware.GetEmail(c),
		"memberships": memberships,
	})
}
