package main

import (
	"context"
	"log"

	"github.com/gofiber/contrib/v3/swaggo"
	"github.com/gofiber/fiber/v3"
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

	// Services
	companySvc := service.NewCompanyService(pool, companyRepo, memberRepo)
	memberSvc := service.NewMembershipService(memberRepo)
	inviteSvc := service.NewInviteService(pool, inviteRepo, memberRepo)

	// Handlers
	handlers := handler.Handlers{
		Company:    handler.NewCompanyHandler(companySvc),
		Membership: handler.NewMembershipHandler(memberSvc),
		Invite:     handler.NewInviteHandler(inviteSvc),
	}

	// Middleware
	auth := middleware.NewAuth(cfg.SupabaseJWTSecret)
	rbac := middleware.NewRBAC(memberRepo)

	// Fiber app
	app := fiber.New(fiber.Config{
		AppName: "Motivator API",
	})

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
