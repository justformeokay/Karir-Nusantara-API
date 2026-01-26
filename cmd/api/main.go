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

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"github.com/karirnusantara/api/internal/config"
	"github.com/karirnusantara/api/internal/database"
	"github.com/karirnusantara/api/internal/middleware"
	"github.com/karirnusantara/api/internal/modules/admin"
	"github.com/karirnusantara/api/internal/modules/applications"
	"github.com/karirnusantara/api/internal/modules/auth"
	"github.com/karirnusantara/api/internal/modules/chat"
	"github.com/karirnusantara/api/internal/modules/company"
	"github.com/karirnusantara/api/internal/modules/cvs"
	"github.com/karirnusantara/api/internal/modules/dashboard"
	"github.com/karirnusantara/api/internal/modules/jobs"
	"github.com/karirnusantara/api/internal/modules/passwordreset"
	"github.com/karirnusantara/api/internal/modules/policies"
	"github.com/karirnusantara/api/internal/modules/profile"
	"github.com/karirnusantara/api/internal/modules/quota"
	"github.com/karirnusantara/api/internal/modules/recommendations"
	"github.com/karirnusantara/api/internal/modules/tickets"
	"github.com/karirnusantara/api/internal/modules/wishlist"
	"github.com/karirnusantara/api/internal/shared/email"
	"github.com/karirnusantara/api/internal/shared/invoice"
	"github.com/karirnusantara/api/internal/shared/response"
	"github.com/karirnusantara/api/internal/shared/validator"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database connection
	db, err := database.NewMySQL(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("Connected to database successfully")

	// Initialize validator
	v := validator.New()

	// Initialize email service first (needed by auth service)
	emailConfig := email.LoadConfigFromEnv()
	emailService := email.NewService(emailConfig)

	// Initialize middleware - need authService for auth middleware
	// Create auth service first for middleware initialization
	authRepo := auth.NewRepository(db)
	authService := auth.NewServiceWithEmail(authRepo, &cfg.JWT, emailService)

	authMiddleware := middleware.NewAuthMiddleware(authService)

	// Initialize other repositories
	jobsRepo := jobs.NewRepository(db)
	cvsRepo := cvs.NewRepository(db)
	applicationsRepo := applications.NewRepository(db)
	wishlistRepo := wishlist.NewRepository(db)
	quotaRepo := quota.NewRepository(db)
	dashboardRepo := dashboard.NewRepository(db)
	companyRepo := company.NewRepository(db)
	chatRepo := chat.NewRepository(db)
	profileRepo := profile.NewRepository(db)
	ticketsRepo := tickets.NewRepository(db)
	passwordResetRepo := passwordreset.NewRepository(db)

	// Initialize other services
	quotaService := quota.NewService(quotaRepo)
	jobsService := jobs.NewServiceWithEmail(jobsRepo, companyRepo, quotaService, emailService)
	cvsService := cvs.NewService(cvsRepo)
	applicationsService := applications.NewService(applicationsRepo, cvsService, jobsService)
	wishlistService := wishlist.NewService(wishlistRepo)
	dashboardService := dashboard.NewService(dashboardRepo)
	companyService := company.NewService(companyRepo)
	chatService := chat.NewService(chatRepo)
	profileService := profile.NewService(profileRepo)
	passwordResetService := passwordreset.NewService(passwordResetRepo, emailService)
	ticketsService := tickets.NewService(ticketsRepo)

	// Initialize invoice service
	invoiceService := invoice.NewService("./docs/invoices")

	// Initialize handlers
	authHandler := auth.NewHandler(authService, v, emailService)
	jobsHandler := jobs.NewHandler(jobsService, v)
	cvsHandler := cvs.NewHandler(cvsService, v)
	applicationsHandler := applications.NewHandler(applicationsService, v)
	wishlistHandler := wishlist.NewHandler(wishlistService, v)
	quotaHandler := quota.NewHandler(quotaService, v)
	dashboardHandler := dashboard.NewHandler(dashboardService)
	chatHandler := chat.NewHandler(chatService, v, "./docs")
	profileHandler := profile.NewHandler(profileService, v, "./docs")
	passwordResetHandler := passwordreset.NewHandler(passwordResetService)
	ticketsHandler := tickets.NewHandler(ticketsService, v)

	// Initialize recommendations module
	recommendationsService := recommendations.NewService()
	recommendationsHandler := recommendations.NewHandler(recommendationsService, jobsService, cvsService, profileService)

	// Initialize company file service
	companyFileService := company.NewFileService("./docs/companies")
	companyHandler := company.NewHandler(companyService, companyFileService)

	// Setup router
	r := chi.NewRouter()

	// Global middleware
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(middleware.RequestLogger)
	r.Use(middleware.NewCORS(cfg.CORS.AllowedOrigins))
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.Timeout(60 * time.Second))

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		response.JSON(w, http.StatusOK, map[string]string{
			"status":  "healthy",
			"service": "karir-nusantara-api",
			"version": "1.0.0",
		})
	})

	// Static files for company documents
	fs := http.FileServer(http.Dir("./docs"))
	r.Handle("/docs/*", http.StripPrefix("/docs/", fs))

	// API routes
	r.Route("/api/v1", func(r chi.Router) {
		// Register module routes with middleware functions
		auth.RegisterRoutes(r, authHandler, authMiddleware.Authenticate)
		jobs.RegisterRoutes(r, jobsHandler, authMiddleware.Authenticate, authMiddleware.RequireCompany, authMiddleware.RequireJobSeeker)
		cvs.RegisterRoutes(r, cvsHandler, authMiddleware.Authenticate, authMiddleware.RequireJobSeeker)
		profile.RegisterRoutes(r, profileHandler, authMiddleware.Authenticate, authMiddleware.RequireJobSeeker)
		applications.RegisterRoutes(r, applicationsHandler, authMiddleware.Authenticate, authMiddleware.RequireJobSeeker, authMiddleware.RequireCompany)
		wishlist.RegisterRoutes(r, wishlistHandler, authMiddleware.Authenticate, authMiddleware.RequireJobSeeker)
		quota.RegisterRoutes(r, quotaHandler, authMiddleware.Authenticate, authMiddleware.RequireCompany)
		dashboard.RegisterRoutes(r, dashboardHandler, authMiddleware.Authenticate, authMiddleware.RequireCompany)
		company.RegisterRoutes(r, companyHandler, authMiddleware.Authenticate)
		chat.RegisterRoutes(r, chatHandler, authMiddleware)
		policies.RegisterRoutes(r)
		recommendations.RegisterRoutes(r, recommendationsHandler, authMiddleware.Authenticate)
		passwordreset.RegisterRoutes(r, passwordResetHandler)
		tickets.RegisterRoutes(r, ticketsHandler, authMiddleware)

		// Admin module routes
		adminModule := admin.NewModuleWithQuota(db, cfg, authMiddleware, quotaService, emailService, invoiceService)
		adminModule.RegisterRoutes(r)
	})

	// 404 handler
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		response.JSON(w, http.StatusNotFound, map[string]string{
			"error":   "not_found",
			"message": "The requested resource was not found",
		})
	})

	// 405 handler
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		response.JSON(w, http.StatusMethodNotAllowed, map[string]string{
			"error":   "method_not_allowed",
			"message": "The requested method is not allowed for this resource",
		})
	})

	// Create server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.App.Port),
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Server starting on port %s", cfg.App.Port)
		log.Printf("Environment: %s", cfg.App.Env)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped gracefully")
}
