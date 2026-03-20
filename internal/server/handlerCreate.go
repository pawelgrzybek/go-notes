package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/pawelgrzybek/go-notes/internal/store"
)

func handlerCreate(s store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Note string `json:"note"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Note == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		note, err := s.Create(body.Note)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(note); err != nil {
			log.Printf("encode response: %v", err)
		}
	}
}
