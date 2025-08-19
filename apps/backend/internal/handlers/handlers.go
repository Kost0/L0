// Package handlers provides request processing
package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/Kost0/L0/internal/cache"
	"github.com/Kost0/L0/internal/repository"
	"github.com/go-chi/chi/v5"
)

// Handler contains tools for working with data
type Handler struct {
	Repo  repository.OrderRepository
	Cache cache.Cache
}

// GetOrderByID godoc
// @Summary Receive an order by ID
// @Description Gets information about an order by its ID
// @Produce json
// @Param orderID path string true "Order ID"
// @Success 200 {object} models.CombinedData "OK"
// @Failure 404 {string} string "There is no such order"
// @Failure 500 {string} string "Internal server error"
// @Router /orders/{orderID} [get]
func (h *Handler) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5000")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	orderID := chi.URLParam(r, "orderID")

	w.Header().Set("Content-Type", "application/json")

	start := time.Now()

	data, ok := h.Cache.Get(orderID)
	if ok {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Printf("The data was retrieved from the cache in %d milliseconds", time.Since(start))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()

	data, err := h.Repo.SelectWithRetry(ctx, orderID)
	log.Printf("Selected order %+v", data)
	log.Printf("The data was retrieved from the database in %d milliseconds", time.Since(start))
	if err != nil {
		log.Println(err)
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	h.Cache.Set(orderID, data)
	log.Printf("Order %s cached", orderID)

	if err = json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
