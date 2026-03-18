package middleware

import (
	"github.com/gofiber/fiber/v3"

	"github.com/hustlers/motivator-backend/internal/model"
	"github.com/hustlers/motivator-backend/internal/repository"
	"github.com/hustlers/motivator-backend/pkg/response"
)

type RBACMiddleware struct {
	members repository.MembershipRepository
}

func NewRBAC(members repository.MembershipRepository) *RBACMiddleware {
	return &RBACMiddleware{members: members}
}

// LoadMembership resolves the current user's membership for the company in :id param.
func (r *RBACMiddleware) LoadMembership() fiber.Handler {
	return func(c fiber.Ctx) error {
		userID := GetUserID(c)
		companyID := c.Params("id")
		if companyID == "" {
			return response.BadRequest(c, "missing company id")
		}

		membership, err := r.members.GetByUserAndCompany(c.Context(), userID, companyID)
		if err != nil {
			return response.Forbidden(c, "not a member of this company")
		}
		if !membership.IsActive {
			return response.Forbidden(c, "membership is deactivated")
		}

		c.Locals("membership", membership)
		return c.Next()
	}
}

// RequireRole ensures the current user has one of the allowed roles.
func RequireRole(allowed ...model.Role) fiber.Handler {
	roleSet := make(map[model.Role]bool, len(allowed))
	for _, r := range allowed {
		roleSet[r] = true
	}
	return func(c fiber.Ctx) error {
		m := GetMembership(c)
		if m == nil {
			return response.Forbidden(c, "no membership context")
		}
		if !roleSet[m.Role] {
			return response.Forbidden(c, "insufficient role")
		}
		return c.Next()
	}
}

func GetMembership(c fiber.Ctx) *model.Membership {
	m, _ := c.Locals("membership").(*model.Membership)
	return m
}
