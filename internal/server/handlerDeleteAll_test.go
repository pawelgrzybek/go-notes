package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pawelgrzybek/go-notes/internal/models"
)

func TestHandlerDeleteAll(t *testing.T) {
	t.Run("delete all returns deleted notes", func(t *testing.T) {
		h, db := setupTestServer(t)
		insertNote(t, db, "one")
		insertNote(t, db, "two")

		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", rec.Code)
		}

		var notes []models.Note
		if err := json.NewDecoder(rec.Body).Decode(&notes); err != nil {
			t.Fatal(err)
		}
		if len(notes) != 2 {
			t.Fatalf("expected 2 deleted notes, got %d", len(notes))
		}
	})

	t.Run("delete all when empty returns empty array", func(t *testing.T) {
		h, _ := setupTestServer(t)

		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", rec.Code)
		}

		var notes []models.Note
		if err := json.NewDecoder(rec.Body).Decode(&notes); err != nil {
			t.Fatal(err)
		}
		if len(notes) != 0 {
			t.Fatalf("expected empty list, got %d notes", len(notes))
		}
	})
}
