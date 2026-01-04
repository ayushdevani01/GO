package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type contextKey string

const userIDKey contextKey = "userID"
const requestIDKey contextKey = "requestID"

func contextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := "user_test"
		requestID := "req_test"
		ctx := context.WithValue(r.Context(), userIDKey, userID)
		ctx = context.WithValue(ctx, requestIDKey, requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userIDKey)
	requestID := r.Context().Value(requestIDKey)
	fmt.Fprintf(w, "User ID: %s, Request ID: %s", userID, requestID)
}

func slowOperation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	select {
	case <-time.After(5 * time.Second):
		fmt.Fprintf(w, "Operation completed")
	case <-ctx.Done():
		err := ctx.Err()
		fmt.Fprintf(w, "Operation cancelled: %v", err)
	}
}

func main() {
	mainhandler := http.HandlerFunc(mainHandler)
	http.Handle("/", contextMiddleware(mainhandler))
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
