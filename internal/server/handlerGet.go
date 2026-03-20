package server

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/pawelgrzybek/go-notes/internal/store"
)

func handlerGet(s store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		noteID, err := strconv.Atoi(r.PathValue("noteID"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		note, err := s.Get(noteID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(note); err != nil {
			log.Printf("encode response: %v", err)
		}
	}
}
