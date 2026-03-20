package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/pawelgrzybek/go-notes/internal/store"
)

func handlerDeleteAll(s store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		notes, err := s.DeleteAll()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(notes); err != nil {
			log.Printf("encode response: %v", err)
		}
	}
}
