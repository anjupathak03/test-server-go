package repository

import (
	"testing"
	"time"

	"test-server/models"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestCreateTodo(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	repo := NewTodoRepository(db)

	req := &models.CreateTodoRequest{
		Title:       "Test Todo",
		Description: "Test Description",
	}

	now := time.Now()

	mock.ExpectExec("INSERT INTO todos").
		WithArgs(req.Title, req.Description).
		WillReturnResult(sqlmock.NewResult(1, 1))

	rows := sqlmock.NewRows([]string{"id", "title", "description", "completed", "created_at", "updated_at"}).
		AddRow(1, "Test Todo", "Test Description", false, now, now)
	
	mock.ExpectQuery("SELECT (.+) FROM todos WHERE id = ?").
		WithArgs(1).
		WillReturnRows(rows)

	todo, err := repo.Create(req)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if todo.Title != req.Title {
		t.Errorf("Expected title %s, got %s", req.Title, todo.Title)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestGetAllTodos(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	repo := NewTodoRepository(db)

	now := time.Now()

	rows := sqlmock.NewRows([]string{"id", "title", "description", "completed", "created_at", "updated_at"}).
		AddRow(1, "Todo 1", "Description 1", false, now, now).
		AddRow(2, "Todo 2", "Description 2", true, now, now)

	mock.ExpectQuery("SELECT (.+) FROM todos ORDER BY created_at DESC").
		WillReturnRows(rows)

	todos, err := repo.GetAll()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(todos) != 2 {
		t.Errorf("Expected 2 todos, got %d", len(todos))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestGetTodoByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	repo := NewTodoRepository(db)

	now := time.Now()

	mock.ExpectQuery("SELECT (.+) FROM todos WHERE id = ?").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "completed", "created_at", "updated_at"}).
			AddRow(1, "Test Todo", "Test Description", false, now, now))

	todo, err := repo.GetByID(1)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if todo.ID != 1 {
		t.Errorf("Expected ID 1, got %d", todo.ID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestUpdateTodo(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	repo := NewTodoRepository(db)

	title := "Updated Title"
	completed := true
	req := &models.UpdateTodoRequest{
		Title:     &title,
		Completed: &completed,
	}

	now := time.Now()

	mock.ExpectExec("UPDATE todos SET (.+) WHERE id = ?").
		WithArgs(title, completed, 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectQuery("SELECT (.+) FROM todos WHERE id = ?").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "completed", "created_at", "updated_at"}).
			AddRow(1, title, "Test Description", completed, now, now))

	todo, err := repo.Update(1, req)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if todo.Title != title {
		t.Errorf("Expected title %s, got %s", title, todo.Title)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestDeleteTodo(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	repo := NewTodoRepository(db)

	mock.ExpectExec("DELETE FROM todos WHERE id = ?").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.Delete(1)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}
