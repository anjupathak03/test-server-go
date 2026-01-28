package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"test-server/models"
)

type TodoRepository struct {
	db *sql.DB
}

func NewTodoRepository(db *sql.DB) *TodoRepository {
	return &TodoRepository{db: db}
}

func (r *TodoRepository) Create(todo *models.CreateTodoRequest) (*models.Todo, error) {
	query := `INSERT INTO todos (title, description) VALUES (?, ?)`
	result, err := r.db.Exec(query, todo.Title, todo.Description)
	if err != nil {
		return nil, fmt.Errorf("failed to create todo: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	return r.GetByID(int(id))
}

func (r *TodoRepository) GetAll() ([]models.Todo, error) {
	query := `SELECT id, title, description, completed, created_at, updated_at FROM todos ORDER BY created_at DESC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query todos: %w", err)
	}
	defer rows.Close()

	var todos []models.Todo
	for rows.Next() {
		var todo models.Todo
		err := rows.Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Completed, &todo.CreatedAt, &todo.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan todo: %w", err)
		}
		todos = append(todos, todo)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating todos: %w", err)
	}

	return todos, nil
}

func (r *TodoRepository) GetByID(id int) (*models.Todo, error) {
	query := `SELECT id, title, description, completed, created_at, updated_at FROM todos WHERE id = ?`
	var todo models.Todo
	err := r.db.QueryRow(query, id).Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Completed, &todo.CreatedAt, &todo.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("todo not found")
		}
		return nil, fmt.Errorf("failed to get todo: %w", err)
	}

	return &todo, nil
}

func (r *TodoRepository) Update(id int, req *models.UpdateTodoRequest) (*models.Todo, error) {
	var setParts []string
	var args []interface{}

	if req.Title != nil {
		setParts = append(setParts, "title = ?")
		args = append(args, *req.Title)
	}
	if req.Description != nil {
		setParts = append(setParts, "description = ?")
		args = append(args, *req.Description)
	}
	if req.Completed != nil {
		setParts = append(setParts, "completed = ?")
		args = append(args, *req.Completed)
	}

	if len(setParts) == 0 {
		return r.GetByID(id)
	}

	args = append(args, id)
	query := fmt.Sprintf("UPDATE todos SET %s WHERE id = ?", strings.Join(setParts, ", "))

	_, err := r.db.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to update todo: %w", err)
	}

	return r.GetByID(id)
}

func (r *TodoRepository) Delete(id int) error {
	query := `DELETE FROM todos WHERE id = ?`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete todo: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("todo not found")
	}

	return nil
}
