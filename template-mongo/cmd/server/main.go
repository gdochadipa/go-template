package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/user/go-templates/template-mongo/internal/config"
	"github.com/user/go-templates/template-mongo/internal/user"
	"github.com/user/go-templates/template-mongo/pkg/logger"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongoDriver.Connect(ctx, options.Client().ApplyURI(cfg.DB.URI))
	if err != nil {
		logger.Fatal("cannot connect to mongo", zap.Error(err))
	}
	defer client.Disconnect(context.Background())

	db := client.Database(cfg.DB.Database)

	// Initialize Architecture Layers (Feature-based)
	userRepo := user.NewMongoRepository(db)
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
