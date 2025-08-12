package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/Kost0/L0/internal/cache"
	"github.com/Kost0/L0/internal/repository"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	DB    *sql.DB
	Cache *cache.OrderCache
}

// GetOrderByID godoc
// @Summary Receive an order by ID
// @Description Gets information about an order by its ID
// @Produce json
// @Param orderID path string true "Order ID"
// @Success {object} CombinedData "OK
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

	data, ok := h.Cache.Get(orderID)
	if ok {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	data, err := repository.SelectOrder(h.DB, orderID)
	log.Printf("Selected order %+v", data)
	if err != nil {
		log.Println(err)
		return
	}

	h.Cache.Set(orderID, data)
	log.Printf("Order %s cached", orderID)

	if err = json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
