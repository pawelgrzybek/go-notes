package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pawelgrzybek/go-notes/internal/models"
)

func TestHandlerList(t *testing.T) {
	t.Run("empty list returns 200 with empty array", func(t *testing.T) {
		h, _ := setupTestServer(t)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
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

	t.Run("returns all notes", func(t *testing.T) {
		h, db := setupTestServer(t)
		insertNote(t, db, "first")
		insertNote(t, db, "second")

		req := httptest.NewRequest(http.MethodGet, "/", nil)
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
			t.Fatalf("expected 2 notes, got %d", len(notes))
		}
		if notes[0].Note != "first" || notes[1].Note != "second" {
			t.Fatalf("unexpected notes: %+v", notes)
		}
	})
}
