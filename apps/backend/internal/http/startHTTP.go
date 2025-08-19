// Package http provides work with net
package http

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Kost0/L0/internal/cache"
	"github.com/Kost0/L0/internal/handlers"
	"github.com/Kost0/L0/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/swaggo/http-swagger"
)

// StartHTTPServer starts the server and handlers
// Accepts:
//   - ctx: context
//   - repo: repository
//   - cache: struct for work with cache
func StartHTTPServer(ctx context.Context, repo *repository.SQLOrderRepository, cache *cache.OrderCache) {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World")
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	h := handlers.Handler{
		Repo:  repo,
		Cache: cache,
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
