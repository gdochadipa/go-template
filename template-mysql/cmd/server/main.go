package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/go-sql-driver/mysql"
	"github.com/user/go-templates/template-mysql/internal/config"
	"github.com/user/go-templates/template-mysql/internal/user"
	"github.com/user/go-templates/template-mysql/pkg/logger"
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
	db, err := sql.Open(cfg.DB.Driver, cfg.DB.Source)
	if err != nil {
		logger.Fatal("cannot open db", zap.Error(err))
	}
	// Verify connection
	if err := db.Ping(); err != nil {
		logger.Fatal("cannot connect to db", zap.Error(err))
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	defer db.Close()

	// Initialize Layers
	userRepo := user.NewMysqlRepository(db)
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
