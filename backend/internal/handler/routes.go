package handler

import (
	"github.com/gofiber/fiber/v3"

	"github.com/hustlers/motivator-backend/internal/middleware"
	"github.com/hustlers/motivator-backend/internal/model"
)

type Handlers struct {
	Company     *CompanyHandler
	Membership  *MembershipHandler
	Invite      *InviteHandler
	Badge       *BadgeHandler
	Achievement *AchievementHandler
	Leaderboard *LeaderboardHandler
}

func RegisterRoutes(app *fiber.App, h Handlers, auth *middleware.AuthMiddleware, rbac *middleware.RBACMiddleware) {
	api := app.Group("/api/v1")

	// Protected routes
	protected := api.Group("", auth.Required())

	// GET /me
	protected.Get("/me", h.Membership.Me)

	// Companies
	protected.Post("/companies", h.Company.Create)

	// Company-scoped routes (require membership)
	company := protected.Group("/companies/:id", rbac.LoadMembership())
	company.Get("", h.Company.GetByID)
	company.Patch("", middleware.RequireRole(model.RoleOwner, model.RoleAdmin), h.Company.Update)
	company.Delete("", middleware.RequireRole(model.RoleOwner), h.Company.Delete)

	// Members
	company.Get("/members", h.Membership.List)
	company.Get("/members/:memberId", h.Membership.GetByID)
	company.Get("/members/:memberId/profile", h.Membership.GetProfile)
	company.Post("/members/:memberId/xp", middleware.RequireRole(model.RoleOwner, model.RoleAdmin), h.Membership.AwardXP)
	company.Patch("/members/:memberId", middleware.RequireRole(model.RoleOwner, model.RoleAdmin), h.Membership.Update)
	company.Delete("/members/:memberId", middleware.RequireRole(model.RoleOwner, model.RoleAdmin), h.Membership.Deactivate)

	// Badges
	company.Get("/badges", h.Badge.List)
	company.Post("/badges", middleware.RequireRole(model.RoleOwner, model.RoleAdmin), h.Badge.Create)
	company.Delete("/badges/:badgeId", middleware.RequireRole(model.RoleOwner, model.RoleAdmin), h.Badge.Delete)
	company.Post("/members/:memberId/badges", middleware.RequireRole(model.RoleOwner, model.RoleAdmin), h.Badge.Award)

	// Achievements
	company.Get("/achievements", h.Achievement.List)
	company.Post("/achievements", middleware.RequireRole(model.RoleOwner, model.RoleAdmin), h.Achievement.Create)
	company.Delete("/achievements/:achievementId", middleware.RequireRole(model.RoleOwner, model.RoleAdmin), h.Achievement.Delete)
	company.Get("/members/:memberId/achievements", h.Achievement.ListMemberAchievements)
	company.Post("/members/:memberId/metrics", h.Achievement.EvaluateMetric)

	// Leaderboard
	company.Get("/leaderboard", h.Leaderboard.Get)

	// Invites (admin+)
	company.Post("/invites", middleware.RequireRole(model.RoleOwner, model.RoleAdmin), h.Invite.Create)
	company.Get("/invites", middleware.RequireRole(model.RoleOwner, model.RoleAdmin), h.Invite.List)
	company.Delete("/invites/:inviteId", middleware.RequireRole(model.RoleOwner, model.RoleAdmin), h.Invite.Revoke)

	// Accept invite (auth required, but no company membership needed)
	protected.Post("/invites/:token/accept", h.Invite.Accept)
}
