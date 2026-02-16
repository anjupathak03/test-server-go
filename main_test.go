package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"test-server/database"
	"test-server/models"
	"test-server/repository"
	"test-server/routes"

	"github.com/stretchr/testify/assert"
)

func setupTestDB(t *testing.T) {
	dbConfig := database.Config{
		Host:     getEnvOrDefault("DB_HOST", "localhost"),
		Port:     getEnvOrDefault("DB_PORT", "3306"),
		User:     getEnvOrDefault("DB_USER", "root"),
		Password: getEnvOrDefault("DB_PASSWORD", "password"),
		DBName:   getEnvOrDefault("DB_NAME", "todo_test_db"),
	}

	if err := database.InitDB(dbConfig); err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}

	if err := database.CreateTodoTable(); err != nil {
		t.Fatalf("Failed to create tables: %v", err)
	}

	// Clean up existing data
	database.DB.Exec("DELETE FROM todos")
}

func teardownTestDB(t *testing.T) {
	if database.DB != nil {
		database.DB.Exec("DELETE FROM todos")
		database.CloseDB()
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func TestIntegrationCreateTodo(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	repo := repository.NewTodoRepository(database.DB)
	router := routes.SetupRouter(repo)

	reqBody := models.CreateTodoRequest{
		Title:       "Integration Test Todo",
		Description: "Testing with real database",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/api/todos", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var todo models.Todo
	json.NewDecoder(w.Body).Decode(&todo)

	assert.Equal(t, reqBody.Title, todo.Title)
	assert.Equal(t, reqBody.Description, todo.Description)
	assert.False(t, todo.Completed)
	assert.NotZero(t, todo.ID)
}

func TestIntegrationGetAllTodos(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	repo := repository.NewTodoRepository(database.DB)
	router := routes.SetupRouter(repo)

	// Create test todos
	repo.Create(&models.CreateTodoRequest{Title: "Todo 1", Description: "Desc 1"})
	repo.Create(&models.CreateTodoRequest{Title: "Todo 2", Description: "Desc 2"})

	req := httptest.NewRequest("GET", "/api/todos", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var todos []models.Todo
	json.NewDecoder(w.Body).Decode(&todos)

	assert.GreaterOrEqual(t, len(todos), 2)
}

func TestIntegrationGetTodoByID(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	repo := repository.NewTodoRepository(database.DB)
	router := routes.SetupRouter(repo)

	// Create a todo
	created, _ := repo.Create(&models.CreateTodoRequest{
		Title:       "Test Todo",
		Description: "Test Description",
	})

	req := httptest.NewRequest("GET", fmt.Sprintf("/api/todos/%d", created.ID), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var todo models.Todo
	json.NewDecoder(w.Body).Decode(&todo)

	assert.Equal(t, created.ID, todo.ID)
	assert.Equal(t, created.Title, todo.Title)
}

func TestExternalHTTPSCall(t *testing.T) {
	url := getEnvOrDefault("EXTERNAL_API_URL", "https://postman-echo.com/get?foo=bar")
	client := &http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("User-Agent", "keploy-test-server/1.0")

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to call external API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		t.Fatalf("Unexpected status code: %d", resp.StatusCode)
	}

	if _, err := io.ReadAll(resp.Body); err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}
}

func TestMySQLHealth(t *testing.T) {
	dbConfig := database.Config{
		Host:     getEnvOrDefault("DB_HOST", "localhost"),
		Port:     getEnvOrDefault("DB_PORT", "3306"),
		User:     getEnvOrDefault("DB_USER", "root"),
		Password: getEnvOrDefault("DB_PASSWORD", "password"),
		DBName:   getEnvOrDefault("DB_NAME", "todo_test_db"),
	}

	if err := database.InitDB(dbConfig); err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.CloseDB()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := database.DB.PingContext(ctx); err != nil {
		t.Fatalf("Failed to ping MySQL: %v", err)
	}
}
