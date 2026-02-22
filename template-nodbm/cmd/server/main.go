package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/user/go-templates/template-nodbm/internal/config"
	"github.com/user/go-templates/template-nodbm/internal/user"
	"github.com/user/go-templates/template-nodbm/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	// Load Configuration
	cfg, err := config.LoadConfig("config")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Initialize Logger
	log, err := logger.New(cfg.Log.Level)
	if err != nil {
		panic(err)
	}
	defer log.Sync()

	// Initialize Architecture Layers (Feature-based)
	userRepo := user.NewMemoryRepository()
	userService := user.NewService(userRepo, log)
	userHandler := user.NewHandler(userService)

	// Router Setup
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api/v1", func(r chi.Router) {
		userHandler.RegisterRoutes(r)
	})

	// Start Server
	log.Info("server starting", zap.String("port", cfg.Server.Port))
	if err := http.ListenAndServe(":"+cfg.Server.Port, r); err != nil {
		log.Fatal("server failed", zap.Error(err))
	}
}
