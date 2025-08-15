package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"healthcare-api/internal/config"
	"healthcare-api/internal/database"
	"healthcare-api/internal/handlers"
	"healthcare-api/internal/middleware"
	"healthcare-api/internal/repository"
	"healthcare-api/internal/service"
	"healthcare-api/internal/worker"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Setup logger
	logger := logrus.New()
	logger.SetLevel(logrus.Level(cfg.LogLevel))
	logger.SetFormatter(&logrus.JSONFormatter{})

	// Initialize database
	db, err := database.NewConnection(cfg.Database)
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := database.RunMigrations(cfg.Database.URL); err != nil {
		logger.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize repositories
	patientRepo := repository.NewPatientRepository(db)
	observationRepo := repository.NewObservationRepository(db)

	// Initialize services
	patientService := service.NewPatientService(patientRepo, logger)
	observationService := service.NewObservationService(observationRepo, logger)

	// Initialize worker pool
	workerPool := worker.NewWorkerPool(10, 1000, logger)
	
	// Register job handlers
	patientIndexHandler := worker.NewPatientIndexHandler(patientService, logger)
	observationProcessHandler := worker.NewObservationProcessHandler(observationService, logger)
	auditLogHandler := worker.NewAuditLogHandler(logger)
	
	workerPool.RegisterHandler(patientIndexHandler)
	workerPool.RegisterHandler(observationProcessHandler)
	workerPool.RegisterHandler(auditLogHandler)
	
	// Start worker pool
	workerPool.Start()
	defer workerPool.Stop()

	// Initialize handlers
	patientHandler := handlers.NewPatientHandler(patientService, logger)
	observationHandler := handlers.NewObservationHandler(observationService, logger)

	// Setup router
	router := setupRouter(cfg, patientHandler, observationHandler, logger)

	// Setup server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.Server.IdleTimeout) * time.Second,
	}

	// Start server in goroutine
	go func() {
		logger.Infof("Starting Healthcare API server on port %d", cfg.Server.Port)
		logger.Info("API Documentation: https://github.com/your-org/healthcare-api/blob/main/docs/API.md")
		logger.Info("Health Check: http://localhost:%d/health", cfg.Server.Port)
		
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down Healthcare API server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatalf("Server forced to shutdown: %v", err)
	}

	logger.Info("Healthcare API server exited")
}

func setupRouter(cfg *config.Config, patientHandler *handlers.PatientHandler, observationHandler *handlers.ObservationHandler, logger *logrus.Logger) *gin.Engine {
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(cfg.JWT.Secret, logger)
	rateLimiter := middleware.NewRateLimiter(100.0, 20) // 100 req/min, burst 20
	validationMiddleware := middleware.NewValidationMiddleware()

	// Global middleware
	router.Use(middleware.Logger(logger))
	router.Use(middleware.Recovery(logger))
	router.Use(middleware.CORS())
	router.Use(rateLimiter.RateLimit())
	router.Use(middleware.Security())

	// Health check endpoint (no auth required)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"timestamp": time.Now().UTC(),
			"version":   "1.0.0",
			"service":   "healthcare-api",
		})
	})

	// API documentation endpoint
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service":       "Healthcare API",
			"version":       "1.0.0",
			"documentation": "https://github.com/your-org/healthcare-api/blob/main/docs/API.md",
			"fhir_version":  "R4",
			"endpoints": gin.H{
				"health":       "/health",
				"patients":     "/api/v1/patients",
				"observations": "/api/v1/observations",
			},
		})
	})

	// API v1 routes with authentication
	v1 := router.Group("/api/v1")
	v1.Use(authMiddleware.RequireAuth())
	{
		// Patient routes
		patients := v1.Group("/patients")
		patients.Use(authMiddleware.RequireScope("patient:read"))
		{
			patients.POST("", 
				authMiddleware.RequireScope("patient:write"),
				validationMiddleware.ValidatePatientCreate(),
				patientHandler.CreatePatient)
			patients.GET("/:id", patientHandler.GetPatient)
			patients.PUT("/:id", 
				authMiddleware.RequireScope("patient:write"),
				validationMiddleware.ValidatePatientUpdate(),
				patientHandler.UpdatePatient)
			patients.DELETE("/:id", 
				authMiddleware.RequireScope("patient:delete"),
				patientHandler.DeletePatient)
			patients.GET("", patientHandler.ListPatients)
		}

		// Observation routes
		observations := v1.Group("/observations")
		observations.Use(authMiddleware.RequireScope("observation:read"))
		{
			observations.POST("", 
				authMiddleware.RequireScope("observation:write"),
				validationMiddleware.ValidateObservationCreate(),
				observationHandler.CreateObservation)
			observations.GET("/:id", observationHandler.GetObservation)
			observations.PUT("/:id", 
				authMiddleware.RequireScope("observation:write"),
				validationMiddleware.ValidateObservationUpdate(),
				observationHandler.UpdateObservation)
			observations.DELETE("/:id", 
				authMiddleware.RequireScope("observation:delete"),
				observationHandler.DeleteObservation)
			observations.GET("", observationHandler.ListObservations)
		}
	}

	return router
}
