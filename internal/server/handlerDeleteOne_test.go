package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pawelgrzybek/go-notes/internal/models"
)

func TestHandlerDeleteOne(t *testing.T) {
	t.Run("delete existing note returns 200", func(t *testing.T) {
		h, db := setupTestServer(t)
		insertNote(t, db, "to delete")

		req := httptest.NewRequest(http.MethodDelete, "/1", nil)
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", rec.Code)
		}

		var note models.Note
		if err := json.NewDecoder(rec.Body).Decode(&note); err != nil {
			t.Fatal(err)
		}
		if note.Note != "to delete" {
			t.Fatalf("expected note 'to delete', got '%s'", note.Note)
		}
	})

	t.Run("delete non-existent note returns 404", func(t *testing.T) {
		h, _ := setupTestServer(t)

		req := httptest.NewRequest(http.MethodDelete, "/999", nil)
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Fatalf("expected status 404, got %d", rec.Code)
		}
	})

	t.Run("invalid ID returns 400", func(t *testing.T) {
		h, _ := setupTestServer(t)

		req := httptest.NewRequest(http.MethodDelete, "/abc", nil)
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected status 400, got %d", rec.Code)
		}
	})
}
