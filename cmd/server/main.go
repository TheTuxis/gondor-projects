package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/TheTuxis/gondor-projects/internal/config"
	"github.com/TheTuxis/gondor-projects/internal/handler"
	"github.com/TheTuxis/gondor-projects/internal/middleware"
	"github.com/TheTuxis/gondor-projects/internal/model"
	jwtpkg "github.com/TheTuxis/gondor-projects/internal/pkg/jwt"
	"github.com/TheTuxis/gondor-projects/internal/repository"
	"github.com/TheTuxis/gondor-projects/internal/service"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Init logger
	var logger *zap.Logger
	var err error
	if cfg.Environment == "production" {
		logger, err = zap.NewProduction()
	} else {
		logger, err = zap.NewDevelopment()
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to init logger: %v\n", err)
		os.Exit(1)
	}
	defer func() { _ = logger.Sync() }()

	// Connect to PostgreSQL
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		logger.Fatal("failed to connect to database", zap.Error(err))
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.Fatal("failed to get underlying sql.DB", zap.Error(err))
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	// Auto-migrate models
	if err := db.AutoMigrate(
		&model.Project{},
		&model.Task{},
		&model.Phase{},
		&model.ProjectMember{},
		&model.Deliverable{},
	); err != nil {
		logger.Fatal("failed to auto-migrate", zap.Error(err))
	}
	logger.Info("database migration completed")

	// Init Redis client
	var redisClient *redis.Client
	if cfg.RedisURL != "" {
		opts, err := redis.ParseURL("redis://" + cfg.RedisURL)
		if err != nil {
			// Fallback: treat as host:port
			opts = &redis.Options{Addr: cfg.RedisURL}
		}
		redisClient = redis.NewClient(opts)
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		if err := redisClient.Ping(ctx).Err(); err != nil {
			logger.Warn("redis connection failed, continuing without redis", zap.Error(err))
			redisClient = nil
		} else {
			logger.Info("connected to Redis")
		}
	}

	// Init JWT manager (validate-only — tokens are issued by gondor-users-security)
	jwtManager := jwtpkg.NewManager(cfg.JWTSecret)

	// Init repositories
	projectRepo := repository.NewProjectRepository(db)
	taskRepo := repository.NewTaskRepository(db)
	phaseRepo := repository.NewPhaseRepository(db)
	memberRepo := repository.NewMemberRepository(db)
	deliverableRepo := repository.NewDeliverableRepository(db)

	// Init services
	projectService := service.NewProjectService(projectRepo, logger)
	taskService := service.NewTaskService(taskRepo, logger)
	phaseService := service.NewPhaseService(phaseRepo, logger)
	memberService := service.NewMemberService(memberRepo, logger)
	deliverableService := service.NewDeliverableService(deliverableRepo, logger)

	// Init handlers
	healthHandler := handler.NewHealthHandler(db, redisClient)
	projectHandler := handler.NewProjectHandler(projectService)
	taskHandler := handler.NewTaskHandler(taskService)
	phaseHandler := handler.NewPhaseHandler(phaseService)
	memberHandler := handler.NewMemberHandler(memberService)
	deliverableHandler := handler.NewDeliverableHandler(deliverableService)

	// Setup Gin
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.LoggingMiddleware(logger))
	router.Use(middleware.AuthMiddleware(jwtManager))

	// Health & metrics (no auth required — handled by skip list)
	router.GET("/health", healthHandler.Health)
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Project routes
	v1 := router.Group("/v1")
	{
		v1.GET("/projects", projectHandler.List)
		v1.POST("/projects", projectHandler.Create)
		v1.GET("/projects/:id", projectHandler.GetByID)
		v1.PUT("/projects/:id", projectHandler.Update)
		v1.DELETE("/projects/:id", projectHandler.Delete)

		// Task routes
		v1.GET("/projects/:id/tasks", taskHandler.List)
		v1.POST("/projects/:id/tasks", taskHandler.Create)
		v1.GET("/projects/:id/tasks/:task_id", taskHandler.GetByID)
		v1.PUT("/projects/:id/tasks/:task_id", taskHandler.Update)
		v1.DELETE("/projects/:id/tasks/:task_id", taskHandler.Delete)

		// Phase routes
		v1.GET("/projects/:id/phases", phaseHandler.List)
		v1.POST("/projects/:id/phases", phaseHandler.Create)
		v1.PUT("/projects/:id/phases/:phase_id", phaseHandler.Update)
		v1.DELETE("/projects/:id/phases/:phase_id", phaseHandler.Delete)

		// Member routes
		v1.GET("/projects/:id/members", memberHandler.List)
		v1.POST("/projects/:id/members", memberHandler.Create)
		v1.DELETE("/projects/:id/members/:member_id", memberHandler.Delete)

		// Deliverable routes
		v1.GET("/projects/:id/deliverables", deliverableHandler.List)
		v1.POST("/projects/:id/deliverables", deliverableHandler.Create)
		v1.PUT("/projects/:id/deliverables/:deliverable_id", deliverableHandler.Update)
		v1.DELETE("/projects/:id/deliverables/:deliverable_id", deliverableHandler.Delete)
	}

	// Start server
	addr := fmt.Sprintf(":%d", cfg.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		logger.Info("starting server", zap.String("addr", addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("server failed", zap.Error(err))
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("server forced to shutdown", zap.Error(err))
	}

	if redisClient != nil {
		_ = redisClient.Close()
	}
	_ = sqlDB.Close()

	logger.Info("server stopped")
}
