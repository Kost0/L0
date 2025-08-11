package http

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Kost0/L0/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func StartHTTPServer(ctx context.Context, db *sql.DB) {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World")
	})

	h := handlers.Handler{
		DB: db,
	}

	r.Get("/orders/{orderID}", h.GetOrderByID)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
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
