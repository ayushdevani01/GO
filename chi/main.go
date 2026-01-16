package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	repo := NewMockUserRepository()
	service := NewUserService(repo)
	handler := NewUserHandler(service)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	handler.RegisterUser(r)
	handler.GetAllUsers(r)

	fmt.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Println("Failed to start server:", err)
	}
}
