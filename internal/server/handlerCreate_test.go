package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/pawelgrzybek/go-notes/internal/models"
)

func TestHandlerCreate(t *testing.T) {
	t.Run("valid body returns 201 with created note", func(t *testing.T) {
		h, _ := setupTestServer(t)

		body := strings.NewReader(`{"note":"my note"}`)
		req := httptest.NewRequest(http.MethodPost, "/", body)
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)

		if rec.Code != http.StatusCreated {
			t.Fatalf("expected status 201, got %d", rec.Code)
		}

		var note models.Note
		if err := json.NewDecoder(rec.Body).Decode(&note); err != nil {
			t.Fatal(err)
		}
		if note.Note != "my note" {
			t.Fatalf("expected note 'my note', got '%s'", note.Note)
		}
		if note.ID == 0 {
			t.Fatal("expected non-zero ID")
		}
	})

	t.Run("empty note returns 400", func(t *testing.T) {
		h, _ := setupTestServer(t)

		body := strings.NewReader(`{"note":""}`)
		req := httptest.NewRequest(http.MethodPost, "/", body)
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected status 400, got %d", rec.Code)
		}
	})

	t.Run("invalid JSON returns 400", func(t *testing.T) {
		h, _ := setupTestServer(t)

		body := strings.NewReader(`not json`)
		req := httptest.NewRequest(http.MethodPost, "/", body)
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected status 400, got %d", rec.Code)
		}
	})
}
