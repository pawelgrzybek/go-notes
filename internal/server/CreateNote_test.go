package server_test

import (
	"context"
	"testing"

	pb "github.com/pawelgrzybek/go-notes/gen/notes/v1"
)

func TestCreateNote(t *testing.T) {
	t.Run("creates note and returns it", func(t *testing.T) {
		s := setupServer(t)

		resp, err := s.CreateNote(context.Background(), &pb.CreateNoteRequest{Note: new("new note")})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if resp.Note.GetNote() != "new note" {
			t.Errorf("expected 'new note', got %q", resp.Note.GetNote())
		}

		if resp.Note.GetId() == 0 {
			t.Error("expected non-zero id")
		}

		if resp.Note.GetCreatedAt() == nil {
			t.Error("expected created_at to be set")
		}
	})
}
