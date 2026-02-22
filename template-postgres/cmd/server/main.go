package main

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/user/go-templates/template-postgres/internal/config"
	"github.com/user/go-templates/template-postgres/internal/user"
	"github.com/user/go-templates/template-postgres/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	// Load Configuration
	cfg, err := config.LoadConfig("config")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Initialize Logger
	logger, err := logger.New(cfg.Log.Level)
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	// Connect to Database
	dbPool, err := pgxpool.New(context.Background(), cfg.DB.Source)
	if err != nil {
		logger.Fatal("cannot connect to db", zap.Error(err))
	}
	defer dbPool.Close()

	// Initialize Layers (Feature-based)
	userRepo := user.NewPostgresRepository(dbPool)
	userService := user.NewService(userRepo, logger)
	userHandler := user.NewHandler(userService)

	// Router Setup
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api/v1", func(r chi.Router) {
		userHandler.RegisterRoutes(r)
	})

	// Start Server
	logger.Info("server starting", zap.String("port", cfg.Server.Port))
	if err := http.ListenAndServe(":"+cfg.Server.Port, r); err != nil {
		logger.Fatal("server failed", zap.Error(err))
	}
}
