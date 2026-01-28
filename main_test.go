package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

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
