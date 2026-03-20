package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/pawelgrzybek/go-notes/internal/store"
)

func handlerList(s store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		list, err := s.List()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(list); err != nil {
			log.Printf("encode response: %v", err)
		}
	}
}
