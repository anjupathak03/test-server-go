package main

import (
	"log"
	"net/http"
	"os"

	"test-server/database"
	"test-server/repository"
	"test-server/routes"
)

func main() {
	// Database configuration
	dbConfig := database.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "3306"),
		User:     getEnv("DB_USER", "root"),
		Password: getEnv("DB_PASSWORD", "password"),
		DBName:   getEnv("DB_NAME", "todo_db"),
	}

	// Initialize database
	if err := database.InitDB(dbConfig); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.CloseDB()

	// Create tables
	if err := database.CreateTodoTable(); err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}

	// Initialize repository
	todoRepo := repository.NewTodoRepository(database.DB)

	// Setup routes
	router := routes.SetupRouter(todoRepo)

	// Start server
	port := getEnv("PORT", "8080")
	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
