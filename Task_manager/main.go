package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

type Task struct {
	mu   sync.RWMutex
	ID   int    `json:"id"`
	Name string `json:"name"`
	Done bool   `json:"done"`
}

type TaskStore struct {
	mu     sync.RWMutex
	tasks  map[int]*Task
	nextID int
}

func NewTaskStore() *TaskStore {
	return &TaskStore{
		tasks:  make(map[int]*Task),
		nextID: 1,
	}
}

// ----------------------------------------------------------------
// Repository Layer
func (ts *TaskStore) AddTask(name string) *Task {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	task := &Task{
		ID:   ts.nextID,
		Name: name,
		Done: false,
	}
	ts.tasks[ts.nextID] = task
	ts.nextID++
	return task
}

func (ts *TaskStore) GetTask(id int) (*Task, bool) {
	ts.mu.RLock()
	defer ts.mu.RUnlock()
	task, exists := ts.tasks[id]
	return task, exists
}

func (ts *TaskStore) UpdateTask(id int, name string, done bool) bool {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	task, exists := ts.tasks[id]
	if !exists {
		return false
	}
	task.mu.Lock()
	defer task.mu.Unlock()
	task.Name = name
	task.Done = done
	return true
}

func (ts *TaskStore) DeleteTask(id int) bool {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	_, exists := ts.tasks[id]
	if !exists {
		return false
	}
	delete(ts.tasks, id)
	return true
}

func (ts *TaskStore) ListTasks() []*Task {
	ts.mu.RLock()
	defer ts.mu.RUnlock()
	tasks := make([]*Task, 0, len(ts.tasks))
	for _, task := range ts.tasks {
		tasks = append(tasks, task)
	}
	return tasks
}

type App struct {
	TaskStore *TaskStore
}

func NewApp() *App {
	return &App{
		TaskStore: NewTaskStore(),
	}
}

// ----------------------------------------------------------------

// ----------------------------------------------------------------
// handler Layer
func (app *App) handleTasks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		app.listTasks(w, r)
	case "POST":
		app.createTask(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (app *App) handleTaskByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/tasks/")
	if id == "" {
		http.Error(w, "Task ID is required", http.StatusBadRequest)
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid Task ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case "GET":
		app.getTask(w, r, idInt)
	case "PUT":
		app.updateTask(w, r, idInt)
	case "DELETE":
		app.deleteTask(w, r, idInt)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// ----------------------------------------------------------------

// ----------------------------------------------------------------
//service Layer

func (app *App) listTasks(w http.ResponseWriter, r *http.Request) {
	tasks := app.TaskStore.ListTasks()
	json.NewEncoder(w).Encode(tasks)
}

func (app *App) createTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	createdTask := app.TaskStore.AddTask(task.Name)
	json.NewEncoder(w).Encode(createdTask)
}

func (app *App) getTask(w http.ResponseWriter, r *http.Request, id int) {
	task, exists := app.TaskStore.GetTask(id)
	if !exists {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(task)
}
func (app *App) updateTask(w http.ResponseWriter, r *http.Request, id int) {
	var task Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	updated := app.TaskStore.UpdateTask(id, task.Name, task.Done)
	if !updated {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (app *App) deleteTask(w http.ResponseWriter, r *http.Request, id int) {
	deleted := app.TaskStore.DeleteTask(id)
	if !deleted {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ----------------------------------------------------------------

// ----------------------------------------------------------------
// Middleware Layer

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received %s request for %s\n", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		log.Printf("Completed %s request for %s\n", r.Method, r.URL.Path)
	})
}

// ----------------------------------------------------------------
func main() {

	// ----------------------------------------------------------------
	// Routing Layer
	app := NewApp()
	http.HandleFunc("/tasks", app.handleTasks)
	http.HandleFunc("/tasks/", app.handleTaskByID)
	// ----------------------------------------------------------------

	handler := loggingMiddleware(http.DefaultServeMux)
	server := &http.Server{
		Addr:         ":8080",
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	go func() {
		log.Println("Server is starting on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on :8080: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v\n", err)
	} else {
		log.Println("Server gracefully stopped")
	}
}
