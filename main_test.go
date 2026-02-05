package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"test-server/database"
	"test-server/models"
	"test-server/repository"
	"test-server/routes"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

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

func keployAgentBaseURL() string {
	if uri := strings.TrimRight(os.Getenv("KEPLOY_AGENT_URI"), "/"); uri != "" {
		return uri
	}
	if port := os.Getenv("KEPLOY_AGENT_PORT"); port != "" {
		return "http://localhost:" + port + "/agent"
	}
	return "http://localhost:6789/agent"
}

func TestIntegrationCreateTodo(t *testing.T) {
	startKeploySession(t, "TestIntegrationCreateTodo")
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
	startKeploySession(t, "TestIntegrationGetAllTodos")
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
	startKeploySession(t, "TestIntegrationGetTodoByID")
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
	// startKeploySession(t, "TestExternalHTTPSCall")

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
	// startKeploySession(t, "TestMySQLHealth")

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

func TestMongoHealth(t *testing.T) {
	// startKeploySession(t, "TestMongoHealth")

	uri := getEnvOrDefault("MONGO_URI", "")
	if uri == "" {
		host := getEnvOrDefault("MONGO_HOST", "localhost")
		port := getEnvOrDefault("MONGO_PORT", "27017")
		uri = "mongodb://" + net.JoinHostPort(host, port)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		t.Fatalf("Failed to create MongoDB client: %v", err)
	}
	defer func() {
		_ = client.Disconnect(context.Background())
	}()

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		t.Fatalf("Failed to ping MongoDB: %v", err)
	}
}
// func TestIntegrationUpdateTodo(t *testing.T) {
// 	setupTestDB(t)
// 	defer teardownTestDB(t)

// 	repo := repository.NewTodoRepository(database.DB)
// 	router := routes.SetupRouter(repo)

// 	// Create a todo
// 	created, _ := repo.Create(&models.CreateTodoRequest{
// 		Title:       "Original Title",
// 		Description: "Original Description",
// 	})

// 	// Update the todo
// 	newTitle := "Updated Title"
// 	completed := true
// 	updateReq := models.UpdateTodoRequest{
// 		Title:     &newTitle,
// 		Completed: &completed,
// 	}
// 	body, _ := json.Marshal(updateReq)

// 	req := httptest.NewRequest("PUT", fmt.Sprintf("/api/todos/%d", created.ID), bytes.NewBuffer(body))
// 	req.Header.Set("Content-Type", "application/json")
// 	w := httptest.NewRecorder()

// 	router.ServeHTTP(w, req)

// 	assert.Equal(t, http.StatusOK, w.Code)

// 	var todo models.Todo
// 	json.NewDecoder(w.Body).Decode(&todo)

// 	assert.Equal(t, newTitle, todo.Title)
// 	assert.True(t, todo.Completed)
// }

// func TestIntegrationDeleteTodo(t *testing.T) {
// 	setupTestDB(t)
// 	defer teardownTestDB(t)

// 	repo := repository.NewTodoRepository(database.DB)
// 	router := routes.SetupRouter(repo)

// 	// Create a todo
// 	created, _ := repo.Create(&models.CreateTodoRequest{
// 		Title:       "To Be Deleted",
// 		Description: "This will be deleted",
// 	})

// 	// Delete the todo
// 	req := httptest.NewRequest("DELETE", fmt.Sprintf("/api/todos/%d", created.ID), nil)
// 	w := httptest.NewRecorder()

// 	router.ServeHTTP(w, req)

// 	assert.Equal(t, http.StatusOK, w.Code)

// 	// Verify it's deleted
// 	_, err := repo.GetByID(created.ID)
// 	assert.Error(t, err)
// }

// func TestIntegrationFullWorkflow(t *testing.T) {
// 	setupTestDB(t)
// 	defer teardownTestDB(t)

// 	repo := repository.NewTodoRepository(database.DB)
// 	router := routes.SetupRouter(repo)

// 	// 1. Create a todo
// 	createReq := models.CreateTodoRequest{
// 		Title:       "Workflow Todo",
// 		Description: "Testing full workflow",
// 	}
// 	body, _ := json.Marshal(createReq)

// 	req := httptest.NewRequest("POST", "/api/todos", bytes.NewBuffer(body))
// 	req.Header.Set("Content-Type", "application/json")
// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, req)

// 	var created models.Todo
// 	json.NewDecoder(w.Body).Decode(&created)

// 	// 2. Get the todo
// 	req = httptest.NewRequest("GET", fmt.Sprintf("/api/todos/%d", created.ID), nil)
// 	w = httptest.NewRecorder()
// 	router.ServeHTTP(w, req)
// 	assert.Equal(t, http.StatusOK, w.Code)

// 	// 3. Update the todo
// 	completed := true
// 	updateReq := models.UpdateTodoRequest{Completed: &completed}
// 	body, _ = json.Marshal(updateReq)
// 	req = httptest.NewRequest("PUT", fmt.Sprintf("/api/todos/%d", created.ID), bytes.NewBuffer(body))
// 	req.Header.Set("Content-Type", "application/json")
// 	w = httptest.NewRecorder()
// 	router.ServeHTTP(w, req)

// 	var updated models.Todo
// 	json.NewDecoder(w.Body).Decode(&updated)
// 	assert.True(t, updated.Completed)

// 	// 4. Get all todos
// 	req = httptest.NewRequest("GET", "/api/todos", nil)
// 	w = httptest.NewRecorder()
// 	router.ServeHTTP(w, req)
// 	assert.Equal(t, http.StatusOK, w.Code)

// 	// 5. Delete the todo
// 	req = httptest.NewRequest("DELETE", fmt.Sprintf("/api/todos/%d", created.ID), nil)
// 	w = httptest.NewRecorder()
// 	router.ServeHTTP(w, req)
// 	assert.Equal(t, http.StatusOK, w.Code)

// 	// Wait a bit for the delete to propagate
// 	time.Sleep(100 * time.Millisecond)
// }

func startKeploySession(t *testing.T, sessionName string) {
	client := &http.Client{}
	url := keployAgentBaseURL() + "/hooks/start-session"
	payload := map[string]string{"name": sessionName}
	jsonPayload, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		// Just log, don't fail if agent is not running (e.g. running tests without keploy)
		t.Logf("Failed to call agent start-session: %v", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Logf("Agent returned non-200: %d", resp.StatusCode)
	}
}
