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
	Challenge   *ChallengeHandler
	Reward      *RewardHandler
	GamePlan     *GamePlanHandler
	Notification *NotificationHandler
	Team         *TeamHandler
	Tournament   *TournamentHandler
	Webhook      *WebhookHandler
	Integration  *IntegrationHandler
	Analytics    *AnalyticsHandler
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

	// Challenges
	company.Get("/challenges", h.Challenge.List)
	company.Post("/challenges", h.Challenge.Create)
	company.Post("/challenges/:challengeId/accept", h.Challenge.Accept)
	company.Post("/challenges/:challengeId/decline", h.Challenge.Decline)
	company.Post("/challenges/:challengeId/score", h.Challenge.ReportScore)

	// Rewards
	company.Get("/rewards", h.Reward.List)
	company.Post("/rewards", middleware.RequireRole(model.RoleOwner, model.RoleAdmin), h.Reward.Create)
	company.Delete("/rewards/:rewardId", middleware.RequireRole(model.RoleOwner, model.RoleAdmin), h.Reward.Delete)
	company.Post("/rewards/redeem", h.Reward.Redeem)
	company.Get("/rewards/redemptions", middleware.RequireRole(model.RoleOwner, model.RoleAdmin), h.Reward.ListRedemptions)
	company.Post("/rewards/redemptions/:redemptionId/fulfill", middleware.RequireRole(model.RoleOwner, model.RoleAdmin), h.Reward.FulfillRedemption)

	// Leaderboard
	company.Get("/leaderboard", h.Leaderboard.Get)

	// Game Plans
	company.Get("/game-plans", h.GamePlan.List)
	company.Post("/game-plans", middleware.RequireRole(model.RoleOwner, model.RoleAdmin), h.GamePlan.Create)
	company.Get("/game-plans/:planId", h.GamePlan.GetByID)
	company.Patch("/game-plans/:planId", middleware.RequireRole(model.RoleOwner, model.RoleAdmin), h.GamePlan.Update)
	company.Put("/game-plans/:planId/flow", middleware.RequireRole(model.RoleOwner, model.RoleAdmin), h.GamePlan.SaveFlow)
	company.Post("/game-plans/:planId/activate", middleware.RequireRole(model.RoleOwner, model.RoleAdmin), h.GamePlan.Activate)
	company.Post("/game-plans/:planId/deactivate", middleware.RequireRole(model.RoleOwner, model.RoleAdmin), h.GamePlan.Deactivate)
	company.Delete("/game-plans/:planId", middleware.RequireRole(model.RoleOwner, model.RoleAdmin), h.GamePlan.Delete)

	// Notifications
	company.Post("/notifications/register", h.Notification.RegisterToken)
	company.Post("/notifications/unregister", h.Notification.UnregisterToken)

	// Teams
	company.Get("/teams", h.Team.List)
	company.Post("/teams", middleware.RequireRole(model.RoleOwner, model.RoleAdmin), h.Team.Create)
	company.Get("/teams/:teamId", h.Team.GetByID)
	company.Delete("/teams/:teamId", middleware.RequireRole(model.RoleOwner, model.RoleAdmin), h.Team.Delete)
	company.Get("/teams/:teamId/members", h.Team.ListMembers)
	company.Post("/teams/:teamId/members", middleware.RequireRole(model.RoleOwner, model.RoleAdmin), h.Team.AddMember)
	company.Delete("/teams/:teamId/members/:membershipId", middleware.RequireRole(model.RoleOwner, model.RoleAdmin), h.Team.RemoveMember)

	// Team Battles
	company.Get("/team-battles", h.Team.ListBattles)
	company.Post("/team-battles", middleware.RequireRole(model.RoleOwner, model.RoleAdmin), h.Team.CreateBattle)
	company.Post("/team-battles/:battleId/score", h.Team.ReportScore)

	// Tournaments
	company.Get("/tournaments", h.Tournament.List)
	company.Post("/tournaments", middleware.RequireRole(model.RoleOwner, model.RoleAdmin), h.Tournament.Create)
	company.Get("/tournaments/:tournamentId", h.Tournament.GetByID)
	company.Patch("/tournaments/:tournamentId/status", middleware.RequireRole(model.RoleOwner, model.RoleAdmin), h.Tournament.UpdateStatus)
	company.Delete("/tournaments/:tournamentId", middleware.RequireRole(model.RoleOwner, model.RoleAdmin), h.Tournament.Delete)
	company.Post("/tournaments/:tournamentId/join", h.Tournament.Join)
	company.Post("/tournaments/:tournamentId/leave", h.Tournament.Leave)
	company.Post("/tournaments/:tournamentId/score", h.Tournament.SubmitScore)
	company.Get("/tournaments/:tournamentId/standings", h.Tournament.GetStandings)
	company.Post("/tournaments/:tournamentId/complete", middleware.RequireRole(model.RoleOwner, model.RoleAdmin), h.Tournament.Complete)

	// Webhooks (Slack/Teams)
	company.Get("/webhooks", middleware.RequireRole(model.RoleOwner, model.RoleAdmin), h.Webhook.List)
	company.Post("/webhooks", middleware.RequireRole(model.RoleOwner, model.RoleAdmin), h.Webhook.Create)
	company.Delete("/webhooks/:webhookId", middleware.RequireRole(model.RoleOwner, model.RoleAdmin), h.Webhook.Delete)

	// Invites (admin+)
	company.Post("/invites", middleware.RequireRole(model.RoleOwner, model.RoleAdmin), h.Invite.Create)
	company.Get("/invites", middleware.RequireRole(model.RoleOwner, model.RoleAdmin), h.Invite.List)
	company.Delete("/invites/:inviteId", middleware.RequireRole(model.RoleOwner, model.RoleAdmin), h.Invite.Revoke)

	// Integrations
	company.Get("/integrations", middleware.RequireRole(model.RoleOwner, model.RoleAdmin), h.Integration.List)
	company.Post("/integrations", middleware.RequireRole(model.RoleOwner, model.RoleAdmin), h.Integration.Create)
	company.Delete("/integrations/:integrationId", middleware.RequireRole(model.RoleOwner, model.RoleAdmin), h.Integration.Delete)
	company.Get("/integrations/:integrationId/mappings", middleware.RequireRole(model.RoleOwner, model.RoleAdmin), h.Integration.ListMappings)
	company.Post("/integrations/:integrationId/mappings", middleware.RequireRole(model.RoleOwner, model.RoleAdmin), h.Integration.CreateMapping)
	company.Delete("/integrations/:integrationId/mappings/:mappingId", middleware.RequireRole(model.RoleOwner, model.RoleAdmin), h.Integration.DeleteMapping)
	company.Get("/integrations/:integrationId/events", middleware.RequireRole(model.RoleOwner, model.RoleAdmin), h.Integration.ListEvents)

	// Analytics
	company.Get("/analytics", middleware.RequireRole(model.RoleOwner, model.RoleAdmin, model.RoleManager), h.Analytics.GetDashboard)

	// Accept invite (auth required, but no company membership needed)
	protected.Post("/invites/:token/accept", h.Invite.Accept)

	// Public inbound webhook (no auth — validated by secret)
	api.Post("/webhooks/inbound/:secret", h.Integration.InboundWebhook)
}
