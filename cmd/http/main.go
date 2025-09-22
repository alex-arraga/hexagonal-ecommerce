package main

import (
	"context"
	"go-ecommerce/internal/adapters/api/http/handlers"
	"go-ecommerce/internal/adapters/api/http/routes"
	"go-ecommerce/internal/adapters/config"
	"go-ecommerce/internal/adapters/logger"
	"go-ecommerce/internal/adapters/security"
	"go-ecommerce/internal/adapters/storage/cache/redis"
	"go-ecommerce/internal/adapters/storage/database/postgres"
	"go-ecommerce/internal/core/services"
	"net/http"
	"time"

	"go-ecommerce/internal/adapters/storage/database/postgres/repository"
	"log/slog"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func main() {
	config, err := config.New()
	if err != nil {
		slog.Error("Error loading environment variables", "error", err)
		os.Exit(1)
	}

	// Set logger
	logger.Set(config.App)
	slog.Info("Starting application", "app", config.App.Name, "env", config.App.Env)

	// Init database
	ctx := context.Background()
	db, err := postgres.New(ctx, config.DB)
	if err != nil {
		slog.Error("Error initializing database connection", "error", err)
		os.Exit(1)
	}
	defer postgres.Close()

	// execute database migrations
	err = postgres.Migrate(db)
	if err != nil {
		slog.Error("Error executing database migrations", "error", err)
		os.Exit(1)
	}
	slog.Info("Successfully migrated the database")

	// Init cache service
	cache, err := redis.New(ctx, config.Redis)
	if err != nil {
		slog.Error("Error initializing cache connection", "error", err)
		os.Exit(1)
	}
	defer cache.Close()
	slog.Info("Successfully connected to the cache server")

	hasher := &security.Hasher{}

	// dependency injection
	userRepo := repository.NewUserRepo(db)
	userSrv := services.NewUserService(userRepo, cache, hasher)
	userHandler := handlers.NewUserHandler(userSrv)

	// router and load routes
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: config.HTTP.AllowedOrigins,
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		MaxAge:         300,
	}))

	routes.LoadUserRoutes(router, userHandler)

	// Configurar servidor HTTP
	s := &http.Server{
		Handler:      router,
		Addr:         ":" + config.HTTP.Port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	slog.Info("Server listening", "port:", config.HTTP.Port)

	err = s.ListenAndServe()
	if err != nil {
		slog.Error("Server failed", "error", err)
		os.Exit(1)
	}

}
