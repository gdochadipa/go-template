package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/user/go-templates/template-sqlite/internal/config"
	"github.com/user/go-templates/template-sqlite/internal/user"
	"github.com/user/go-templates/template-sqlite/pkg/logger"
	"go.uber.org/zap"
	_ "modernc.org/sqlite"
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
	db, err := sql.Open(cfg.DB.Driver, cfg.DB.Source)
	if err != nil {
		logger.Fatal("cannot open db", zap.Error(err))
	}
	// Verify connection
	if err := db.Ping(); err != nil {
		logger.Fatal("cannot connect to db", zap.Error(err))
	}
	defer db.Close()

	// Initialize Architecture Layers (Feature-based)
	userRepo := user.NewSqliteRepository(db)
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
