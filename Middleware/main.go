package main

import (
	"fmt"
	"net/http"
)

func mainHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the Middleware Example!")
}

func middleware2(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Middleware 2: Before Handler")
		next.ServeHTTP(w, r)
		fmt.Println("Middleware 2: After Handler")
	})
}

func middleware1(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Middleware 1: Before Handler")
		next.ServeHTTP(w, r)
		fmt.Println("Middleware 1: After Handler")
	})
}

func middleware0(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Middleware 0: Before Handler")
		next.ServeHTTP(w, r)
		fmt.Println("Middleware 0: After Handler")
	})
}
func main() {
	mainhandler := http.HandlerFunc(mainHandler)
	http.Handle("/", middleware0(middleware1(middleware2(mainhandler))))

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
