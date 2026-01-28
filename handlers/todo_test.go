package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"test-server/models"

	"github.com/gorilla/mux"
)

// MockTodoRepository mocks the TodoRepository for testing
type MockTodoRepository struct {
	CreateFunc  func(*models.CreateTodoRequest) (*models.Todo, error)
	GetAllFunc  func() ([]models.Todo, error)
	GetByIDFunc func(int) (*models.Todo, error)
	UpdateFunc  func(int, *models.UpdateTodoRequest) (*models.Todo, error)
	DeleteFunc  func(int) error
}

func (m *MockTodoRepository) Create(req *models.CreateTodoRequest) (*models.Todo, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(req)
	}
	return nil, nil
}

func (m *MockTodoRepository) GetAll() ([]models.Todo, error) {
	if m.GetAllFunc != nil {
		return m.GetAllFunc()
	}
	return nil, nil
}

func (m *MockTodoRepository) GetByID(id int) (*models.Todo, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(id)
	}
	return nil, nil
}

func (m *MockTodoRepository) Update(id int, req *models.UpdateTodoRequest) (*models.Todo, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(id, req)
	}
	return nil, nil
}

func (m *MockTodoRepository) Delete(id int) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return nil
}

func TestCreateTodo(t *testing.T) {
	mockRepo := &MockTodoRepository{
		CreateFunc: func(req *models.CreateTodoRequest) (*models.Todo, error) {
			return &models.Todo{
				ID:          1,
				Title:       req.Title,
				Description: req.Description,
				Completed:   false,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}, nil
		},
	}

	handler := &TodoHandler{repo: mockRepo}

	reqBody := models.CreateTodoRequest{
		Title:       "Test Todo",
		Description: "Test Description",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/api/todos", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.CreateTodo(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
	}

	var todo models.Todo
	json.NewDecoder(w.Body).Decode(&todo)

	if todo.Title != reqBody.Title {
		t.Errorf("Expected title %s, got %s", reqBody.Title, todo.Title)
	}
}

func TestGetAllTodos(t *testing.T) {
	mockRepo := &MockTodoRepository{
		GetAllFunc: func() ([]models.Todo, error) {
			return []models.Todo{
				{ID: 1, Title: "Todo 1", Description: "Desc 1", Completed: false},
				{ID: 2, Title: "Todo 2", Description: "Desc 2", Completed: true},
			}, nil
		},
	}

	handler := &TodoHandler{repo: mockRepo}

	req := httptest.NewRequest("GET", "/api/todos", nil)
	w := httptest.NewRecorder()

	handler.GetAllTodos(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var todos []models.Todo
	json.NewDecoder(w.Body).Decode(&todos)

	if len(todos) != 2 {
		t.Errorf("Expected 2 todos, got %d", len(todos))
	}
}

func TestGetTodo(t *testing.T) {
	mockRepo := &MockTodoRepository{
		GetByIDFunc: func(id int) (*models.Todo, error) {
			return &models.Todo{
				ID:          id,
				Title:       "Test Todo",
				Description: "Test Description",
				Completed:   false,
			}, nil
		},
	}

	handler := &TodoHandler{repo: mockRepo}

	req := httptest.NewRequest("GET", "/api/todos/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	handler.GetTodo(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var todo models.Todo
	json.NewDecoder(w.Body).Decode(&todo)

	if todo.ID != 1 {
		t.Errorf("Expected ID 1, got %d", todo.ID)
	}
}

func TestUpdateTodo(t *testing.T) {
	title := "Updated Title"
	mockRepo := &MockTodoRepository{
		UpdateFunc: func(id int, req *models.UpdateTodoRequest) (*models.Todo, error) {
			return &models.Todo{
				ID:          id,
				Title:       *req.Title,
				Description: "Test Description",
				Completed:   false,
			}, nil
		},
	}

	handler := &TodoHandler{repo: mockRepo}

	reqBody := models.UpdateTodoRequest{
		Title: &title,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("PUT", "/api/todos/1", bytes.NewBuffer(body))
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	handler.UpdateTodo(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var todo models.Todo
	json.NewDecoder(w.Body).Decode(&todo)

	if todo.Title != title {
		t.Errorf("Expected title %s, got %s", title, todo.Title)
	}
}

func TestDeleteTodo(t *testing.T) {
	mockRepo := &MockTodoRepository{
		DeleteFunc: func(id int) error {
			return nil
		},
	}

	handler := &TodoHandler{repo: mockRepo}

	req := httptest.NewRequest("DELETE", "/api/todos/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	handler.DeleteTodo(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

func TestGetTodoNotFound(t *testing.T) {
	mockRepo := &MockTodoRepository{
		GetByIDFunc: func(id int) (*models.Todo, error) {
			return nil, errors.New("todo not found")
		},
	}

	handler := &TodoHandler{repo: mockRepo}

	req := httptest.NewRequest("GET", "/api/todos/999", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "999"})
	w := httptest.NewRecorder()

	handler.GetTodo(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
	}
}
