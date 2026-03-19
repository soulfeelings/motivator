package main

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/requestid"

	"github.com/hustlers/motivator-game/internal/config"
	"github.com/hustlers/motivator-game/internal/handler"
	"github.com/hustlers/motivator-game/internal/middleware"
	"github.com/hustlers/motivator-game/internal/repository"
	"github.com/hustlers/motivator-game/internal/service"
)

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
	baseRepo := repository.NewBaseRepository(pool)
	armyRepo := repository.NewArmyRepository(pool)
	battleRepo := repository.NewBattleRepository(pool)

	// Services
	gameSvc := service.NewGameService(baseRepo, armyRepo, battleRepo)

	// Handler
	gameHandler := handler.NewGameHandler(gameSvc)

	// Middleware
	auth := middleware.NewAuth(cfg.SupabaseJWTSecret)

	// Fiber
	app := fiber.New(fiber.Config{
		AppName: "Motivator Game Server",
	})

	app.Use(requestid.New())
	app.Use(logger.New(logger.Config{
		Format:     "${time} | ${status} | ${latency} | ${method} ${path} | ${ip} | trace=${locals:requestid} | ${error}\n",
		TimeFormat: "2006-01-02 15:04:05",
	}))
	app.Use(cors.New())

	app.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok", "service": "game-server"})
	})

	handler.RegisterRoutes(app, gameHandler, auth)

	port := cfg.Port
	if port == "8080" {
		port = "8081"
	}
	log.Printf("game server starting on :%s", port)
	log.Fatal(app.Listen(":" + port))
}
