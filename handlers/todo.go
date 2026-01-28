package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"test-server/models"

	"github.com/gorilla/mux"
)

type TodoRepository interface {
	Create(*models.CreateTodoRequest) (*models.Todo, error)
	GetAll() ([]models.Todo, error)
	GetByID(int) (*models.Todo, error)
	Update(int, *models.UpdateTodoRequest) (*models.Todo, error)
	Delete(int) error
}

type TodoHandler struct {
	repo TodoRepository
}

func NewTodoHandler(repo TodoRepository) *TodoHandler {
	return &TodoHandler{repo: repo}
}

func (h *TodoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	var req models.CreateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if req.Title == "" {
		respondWithError(w, http.StatusBadRequest, "Title is required")
		return
	}

	todo, err := h.repo.Create(&req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, todo)
}

func (h *TodoHandler) GetAllTodos(w http.ResponseWriter, r *http.Request) {
	todos, err := h.repo.GetAll()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, todos)
}

func (h *TodoHandler) GetTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid todo ID")
		return
	}

	todo, err := h.repo.GetByID(id)
	if err != nil {
		if err.Error() == "todo not found" {
			respondWithError(w, http.StatusNotFound, "Todo not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, todo)
}

func (h *TodoHandler) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid todo ID")
		return
	}

	var req models.UpdateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	todo, err := h.repo.Update(id, &req)
	if err != nil {
		if err.Error() == "todo not found" {
			respondWithError(w, http.StatusNotFound, "Todo not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, todo)
}

func (h *TodoHandler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid todo ID")
		return
	}

	err = h.repo.Delete(id)
	if err != nil {
		if err.Error() == "todo not found" {
			respondWithError(w, http.StatusNotFound, "Todo not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Todo deleted successfully"})
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
