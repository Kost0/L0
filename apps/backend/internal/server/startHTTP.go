package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

func StartHTTPServer(ctx context.Context) {
	router := http.NewServeMux()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		<-ctx.Done()
		log.Println("Shutting down server...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Fatalf("Error shutting down server: %v\n", err)
		}
	}()

	log.Println("Starting HTTP server...")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Error starting HTTP server: %v\n", err)
	}
}
