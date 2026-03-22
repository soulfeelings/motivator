package main

import (
	"context"
	"log"

	"github.com/gofiber/contrib/v3/swaggo"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/requestid"

	"github.com/hustlers/motivator-backend/internal/config"
	"github.com/hustlers/motivator-backend/internal/handler"
	"github.com/hustlers/motivator-backend/internal/middleware"
	"github.com/hustlers/motivator-backend/internal/repository"
	"github.com/hustlers/motivator-backend/internal/service"

	_ "github.com/hustlers/motivator-backend/docs"
)

// @title Motivator API
// @version 1.0
// @description Motivator backend API
// @host localhost:8080
// @BasePath /api/v1
func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	ctx := context.Background()
	pool, err := config.NewDatabasePool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	// Repositories
	companyRepo := repository.NewCompanyRepository(pool)
	memberRepo := repository.NewMembershipRepository(pool)
	inviteRepo := repository.NewInviteRepository(pool)
	badgeRepo := repository.NewBadgeRepository(pool)
	achievementRepo := repository.NewAchievementRepository(pool)
	challengeRepo := repository.NewChallengeRepository(pool)
	rewardRepo := repository.NewRewardRepository(pool)
	gamePlanRepo := repository.NewGamePlanRepository(pool)
	deviceTokenRepo := repository.NewDeviceTokenRepository(pool)
	teamRepo := repository.NewTeamRepository(pool)
	tournamentRepo := repository.NewTournamentRepository(pool)
	webhookRepo := repository.NewWebhookRepository(pool)
	integrationRepo := repository.NewIntegrationRepository(pool)

	// Services
	companySvc := service.NewCompanyService(pool, companyRepo, memberRepo)
	memberSvc := service.NewMembershipService(memberRepo, badgeRepo)
	inviteSvc := service.NewInviteService(pool, inviteRepo, memberRepo)
	badgeSvc := service.NewBadgeService(pool, badgeRepo, memberRepo)
	notifSvc := service.NewNotificationService(deviceTokenRepo)
	achievementSvc := service.NewAchievementService(pool, achievementRepo, memberRepo, badgeRepo, notifSvc)
	challengeSvc := service.NewChallengeService(challengeRepo, memberRepo, notifSvc)
	rewardSvc := service.NewRewardService(pool, rewardRepo, memberRepo)
	gamePlanSvc := service.NewGamePlanService(gamePlanRepo)
	teamSvc := service.NewTeamService(teamRepo, memberRepo)
	tournamentSvc := service.NewTournamentService(tournamentRepo, memberRepo)
	webhookSvc := service.NewWebhookService(webhookRepo)
	integrationSvc := service.NewIntegrationService(integrationRepo, memberRepo, achievementRepo)
	analyticsRepo := repository.NewAnalyticsRepository(pool)
	analyticsSvc := service.NewAnalyticsService(analyticsRepo)
	questRepo := repository.NewQuestRepository(pool)
	questSvc := service.NewQuestService(questRepo, memberRepo)

	// Handlers
	handlers := handler.Handlers{
		Company:     handler.NewCompanyHandler(companySvc),
		Membership:  handler.NewMembershipHandler(memberSvc),
		Invite:      handler.NewInviteHandler(inviteSvc),
		Badge:       handler.NewBadgeHandler(badgeSvc),
		Achievement: handler.NewAchievementHandler(achievementSvc),
		Leaderboard: handler.NewLeaderboardHandler(memberSvc),
		Challenge:   handler.NewChallengeHandler(challengeSvc),
		Reward:      handler.NewRewardHandler(rewardSvc),
		GamePlan:     handler.NewGamePlanHandler(gamePlanSvc),
		Notification: handler.NewNotificationHandler(notifSvc),
		Team:         handler.NewTeamHandler(teamSvc),
		Tournament:   handler.NewTournamentHandler(tournamentSvc),
		Webhook:      handler.NewWebhookHandler(webhookSvc),
		Integration:  handler.NewIntegrationHandler(integrationSvc),
		Analytics:    handler.NewAnalyticsHandler(analyticsSvc),
		Quest:        handler.NewQuestHandler(questSvc),
	}

	// Middleware
	auth := middleware.NewAuth(cfg.SupabaseJWTSecret)
	rbac := middleware.NewRBAC(memberRepo)

	// Fiber app
	app := fiber.New(fiber.Config{
		AppName: "Motivator API",
	})

	app.Use(cors.New())
	app.Use(requestid.New())
	app.Use(logger.New(logger.Config{
		Format:     "${time} | ${status} | ${latency} | ${method} ${path} | ${ip} | trace=${locals:requestid} | ${error}\n",
		TimeFormat: "2006-01-02 15:04:05",
	}))

	// Swagger
	app.Get("/swagger/*", swaggo.HandlerDefault)

	// Health check (no auth)
	app.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// API routes
	handler.RegisterRoutes(app, handlers, auth, rbac)

	log.Printf("server starting on :%s", cfg.Port)
	log.Fatal(app.Listen(":" + cfg.Port))
}
