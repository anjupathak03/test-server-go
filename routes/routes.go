package routes

import (
	"test-server/handlers"
	"test-server/repository"

	"github.com/gorilla/mux"
)

func SetupRouter(repo *repository.TodoRepository) *mux.Router {
	router := mux.NewRouter()
	todoHandler := handlers.NewTodoHandler(repo)

	// Todo routes
	router.HandleFunc("/api/todos", todoHandler.CreateTodo).Methods("POST")
	router.HandleFunc("/api/todos", todoHandler.GetAllTodos).Methods("GET")
	router.HandleFunc("/api/todos/{id}", todoHandler.GetTodo).Methods("GET")
	router.HandleFunc("/api/todos/{id}", todoHandler.UpdateTodo).Methods("PUT")
	router.HandleFunc("/api/todos/{id}", todoHandler.DeleteTodo).Methods("DELETE")

	return router
}
