package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func test(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Task Strated...")

	for i := 0; i < 5; i++ {
		time.Sleep(1 * time.Second)
		fmt.Println("Working on %d.", i)
	}

	fmt.Fprintf(w, "Task Completed Sucessfully!!!!")
	fmt.Println("Task finished!!!")
}

func main() {

	server := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	http.HandleFunc("/", test)

	go func() {
		fmt.Println("Server is starting...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	// Channel to listen OS Signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit // Block until we receive our signal.
	fmt.Println("Server is shutting down...")

	//create context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Graceful Shutdown the server
	if err := server.Shutdown(ctx); err != nil {
		fmt.Println("Server forced to shutdown:", err)
	} else {
		fmt.Println("Server Gracefully stopped")
	}
}
