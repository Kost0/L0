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
	DB *sql.DB
}

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

	data, ok := cache.Get(orderID)
	if ok {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	data, err := repository.SelectOrder(h.DB, orderID)
	log.Printf("Selected order %+v", data)
	if err != nil {
		log.Println(err)
		return
	}

	if err = json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
