package main

import (
	"fmt"
	"net/http"
	"sync"
)

type Store struct {
	mu   sync.RWMutex
	data map[string]string
}

var myStrore = Store{
	data: make(map[string]string),
}

func setHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	key := query.Get("key")
	value := query.Get("value")

	if key == "" || value == "" {
		http.Error(w, "Missing key or value", http.StatusBadRequest)
		return
	}

	myStrore.mu.Lock()
	defer myStrore.mu.Unlock()
	myStrore.data[key] = value

	fmt.Fprintf(w, "Saved %s = %s\n", key, value)
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	key := query.Get("key")
	if key == "" {
		http.Error(w, "Missing key", http.StatusBadRequest)
		return
	}

	myStrore.mu.RLock()
	defer myStrore.mu.RUnlock()

	value, ok := myStrore.data[key]
	if !ok {
		http.Error(w, "Key not found", http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, "Value for '%s' = %s\n", key, value)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	key := query.Get("key")

	if key == "" {
		http.Error(w, "Missing key", http.StatusBadRequest)
		return
	}

	myStrore.mu.Lock()
	defer myStrore.mu.Unlock()
	_, ok := myStrore.data[key]
	if !ok {
		http.Error(w, "Key not found", http.StatusNotFound)
		return
	}
	delete(myStrore.data, key)
	fmt.Fprintf(w, "Deleted key '%s'\n", key)
}
func main() {
	http.HandleFunc("/set", setHandler)
	http.HandleFunc("/get", getHandler)
	http.HandleFunc("/delete", deleteHandler)

	fmt.Println("Starting server at :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
