package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pawelgrzybek/go-notes/internal/models"
)

func TestHandlerGet(t *testing.T) {
	t.Run("valid ID returns note", func(t *testing.T) {
		h, db := setupTestServer(t)
		insertNote(t, db, "hello")

		req := httptest.NewRequest(http.MethodGet, "/1", nil)
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", rec.Code)
		}

		var note models.Note
		if err := json.NewDecoder(rec.Body).Decode(&note); err != nil {
			t.Fatal(err)
		}
		if note.Note != "hello" {
			t.Fatalf("expected note 'hello', got '%s'", note.Note)
		}
	})

	t.Run("invalid ID returns 400", func(t *testing.T) {
		h, _ := setupTestServer(t)

		req := httptest.NewRequest(http.MethodGet, "/abc", nil)
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected status 400, got %d", rec.Code)
		}
	})

	t.Run("non-existent ID returns 404", func(t *testing.T) {
		h, _ := setupTestServer(t)

		req := httptest.NewRequest(http.MethodGet, "/999", nil)
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Fatalf("expected status 404, got %d", rec.Code)
		}
	})
}
