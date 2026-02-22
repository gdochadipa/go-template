package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/user/go-templates/template-http-proto/internal/config"
	"github.com/user/go-templates/template-http-proto/internal/user"
	"github.com/user/go-templates/template-http-proto/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.LoadConfig("config")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	logger, err := logger.New(cfg.Log.Level)
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	userSvc := user.NewService(logger)
	userHandler := user.NewHandler(userSvc)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api/v1", func(r chi.Router) {
		userHandler.RegisterRoutes(r)
	})

	logger.Info("HTTP Proto server starting", zap.String("port", cfg.Server.Port))
	if err := http.ListenAndServe(":"+cfg.Server.Port, r); err != nil {
		logger.Fatal("server failed", zap.Error(err))
	}
}
