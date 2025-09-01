package main

import (
	"context"
	"go-clean-architecture/internal/adapter/controller"
	"go-clean-architecture/internal/adapter/repository"
	"go-clean-architecture/internal/infrastructure/database"
	"go-clean-architecture/internal/infrastructure/server"
	"go-clean-architecture/internal/usecase"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize database connection
	dbConfig := database.NewConfig()
	db, err := database.Connect(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run database migrations
	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("Failed to run database migrations: %v", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)

	// Initialize use cases
	userUseCase := usecase.NewUserUseCase(userRepo)

	// Initialize controllers
	userController := controller.NewUserController(userUseCase)

	// Initialize HTTP server
	httpServer := server.NewServer(userController)

	// Start HTTP server
	if err := httpServer.Start(os.Getenv("PORT")); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	log.Println("Server started successfully")

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Received shutdown signal, initiating graceful shutdown...")

	// Create shutdown context with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	// Shutdown HTTP server
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	// Close database connection
	sqlDB, err := db.DB()
	if err == nil {
		if err := sqlDB.Close(); err != nil {
			log.Printf("Database close error: %v", err)
		} else {
			log.Println("Database connection closed")
		}
	}

	log.Println("Server shutdown complete")
}
