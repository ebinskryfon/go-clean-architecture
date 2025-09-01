package server

import (
	"context"
	"go-clean-architecture/internal/adapter/controller"
	"go-clean-architecture/pkg/response"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// Server represents the HTTP server
type Server struct {
	router         *gin.Engine
	httpServer     *http.Server
	userController *controller.UserController
}

// NewServer creates a new HTTP server instance
func NewServer(userController *controller.UserController) *Server {
	// Set gin mode based on environment
	if os.Getenv("GIN_MODE") != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Add middlewares
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(corsMiddleware())
	router.Use(timeoutMiddleware(30 * time.Second))

	server := &Server{
		router:         router,
		userController: userController,
	}

	server.setupRoutes()
	return server
}

// setupRoutes configures all API routes
func (s *Server) setupRoutes() {
	// Health check endpoint
	s.router.GET("/health", s.healthCheck)

	// API v1 routes
	v1 := s.router.Group("/api/v1")
	{
		// User routes
		users := v1.Group("/users")
		{
			users.POST("", s.userController.CreateUser)
			users.GET("", s.userController.GetAllUsers)
			users.GET("/:id", s.userController.GetUser)
			users.PUT("/:id", s.userController.UpdateUser)
			users.DELETE("/:id", s.userController.DeleteUser)
			users.PUT("/:id/activate", s.userController.ActivateUser)
			users.PUT("/:id/deactivate", s.userController.DeactivateUser)
		}
	}

	// 404 handler
	s.router.NoRoute(func(c *gin.Context) {
		response.NotFound(c, "Endpoint not found")
	})
}

// Start starts the HTTP server
func (s *Server) Start(port string) error {
	if port == "" {
		port = getEnv("PORT", "8080")
	}

	s.httpServer = &http.Server{
		Addr:         ":" + port,
		Handler:      s.router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("Server starting on port %s", port)

	// Start server in a goroutine
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	return nil
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	log.Println("Shutting down server...")
	return s.httpServer.Shutdown(ctx)
}

// healthCheck handles health check requests
func (s *Server) healthCheck(c *gin.Context) {
	response.Success(c, "Server is healthy", gin.H{
		"status":    "ok",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"version":   "1.0.0",
	})
}

// corsMiddleware adds CORS headers
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Requested-With")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// timeoutMiddleware adds request timeout
func timeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create a context with timeout
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		// Replace the request context
		c.Request = c.Request.WithContext(ctx)

		// Channel to capture if handler completed
		done := make(chan struct{})

		go func() {
			defer close(done)
			c.Next()
		}()

		select {
		case <-done:
			// Handler completed successfully
			return
		case <-ctx.Done():
			// Timeout occurred
			if !c.Writer.Written() {
				response.InternalError(c, "Request timeout", "The request took too long to process")
			}
			c.Abort()
			return
		}
	}
}

// getEnv gets environment variable with fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
