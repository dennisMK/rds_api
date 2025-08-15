package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"healthcare-api/internal/handlers"
	"healthcare-api/internal/middleware"
)

// SetupRoutes configures all API routes with appropriate middleware
func SetupRoutes(
	patientHandler *handlers.PatientHandler,
	observationHandler *handlers.ObservationHandler,
	authMiddleware *middleware.AuthMiddleware,
	rateLimitMiddleware *middleware.RateLimitMiddleware,
	securityMiddleware *middleware.SecurityMiddleware,
	loggingMiddleware *middleware.LoggingMiddleware,
	validationMiddleware *middleware.ValidationMiddleware,
	auditMiddleware *middleware.AuditMiddleware,
) *mux.Router {
	router := mux.NewRouter()

	// Apply global middleware
	router.Use(securityMiddleware.SecurityHeaders)
	router.Use(loggingMiddleware.LogRequests)
	router.Use(rateLimitMiddleware.RateLimit)

	// Health check endpoints (no auth required)
	router.HandleFunc("/health", healthCheck).Methods("GET")
	router.HandleFunc("/health/ready", readinessCheck).Methods("GET")
	router.HandleFunc("/health/live", livenessCheck).Methods("GET")

	// API v1 routes
	api := router.PathPrefix("/api/v1").Subrouter()
	
	// Apply authentication to all API routes
	api.Use(authMiddleware.Authenticate)
	api.Use(auditMiddleware.LogAuditTrail)

	// Patient routes
	patientRoutes := api.PathPrefix("/patients").Subrouter()
	patientRoutes.Use(validationMiddleware.ValidatePatient)
	patientRoutes.HandleFunc("", patientHandler.CreatePatient).Methods("POST")
	patientRoutes.HandleFunc("", patientHandler.ListPatients).Methods("GET")
	patientRoutes.HandleFunc("/{id}", patientHandler.GetPatient).Methods("GET")
	patientRoutes.HandleFunc("/{id}", patientHandler.UpdatePatient).Methods("PUT")
	patientRoutes.HandleFunc("/{id}", patientHandler.DeletePatient).Methods("DELETE")

	// Observation routes
	observationRoutes := api.PathPrefix("/observations").Subrouter()
	observationRoutes.Use(validationMiddleware.ValidateObservation)
	observationRoutes.HandleFunc("", observationHandler.CreateObservation).Methods("POST")
	observationRoutes.HandleFunc("", observationHandler.ListObservations).Methods("GET")
	observationRoutes.HandleFunc("/{id}", observationHandler.GetObservation).Methods("GET")
	observationRoutes.HandleFunc("/{id}", observationHandler.UpdateObservation).Methods("PUT")
	observationRoutes.HandleFunc("/{id}", observationHandler.DeleteObservation).Methods("DELETE")

	// Metrics endpoint (protected)
	router.HandleFunc("/metrics", metricsHandler).Methods("GET")

	return router
}

// healthCheck provides basic health status
func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"healthy","timestamp":"` + 
		time.Now().UTC().Format(time.RFC3339) + `","version":"1.0.0"}`))
}

// readinessCheck verifies all dependencies are ready
func readinessCheck(w http.ResponseWriter, r *http.Request) {
	// TODO: Check database connectivity, external services, etc.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ready","timestamp":"` + 
		time.Now().UTC().Format(time.RFC3339) + `"}`))
}

// livenessCheck verifies the application is alive
func livenessCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"alive","timestamp":"` + 
		time.Now().UTC().Format(time.RFC3339) + `"}`))
}

// metricsHandler exposes Prometheus metrics
func metricsHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement Prometheus metrics endpoint
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("# Healthcare API Metrics\n# TODO: Implement metrics collection\n"))
}
