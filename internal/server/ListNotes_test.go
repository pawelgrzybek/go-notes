package server_test

import (
	"context"
	"testing"

	pb "github.com/pawelgrzybek/go-notes/gen/notes/v1"
)

func TestListNotes(t *testing.T) {
	t.Run("returns empty list when no notes", func(t *testing.T) {
		s := setupServer(t)

		resp, err := s.ListNotes(context.Background(), &pb.ListNotesRequest{})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(resp.Notes) != 0 {
			t.Fatalf("expected 0 notes, got %d", len(resp.Notes))
		}
	})

	t.Run("returns notes from store", func(t *testing.T) {
		s := setupServer(t)

		s.CreateNote(context.Background(), &pb.CreateNoteRequest{Note: new("first note")})
		s.CreateNote(context.Background(), &pb.CreateNoteRequest{Note: new("second note")})

		resp, err := s.ListNotes(context.Background(), &pb.ListNotesRequest{})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(resp.Notes) != 2 {
			t.Fatalf("expected 2 notes, got %d", len(resp.Notes))
		}

		if resp.Notes[0].GetNote() != "first note" {
			t.Errorf("expected 'first note', got %q", resp.Notes[0].GetNote())
		}

		if resp.Notes[1].GetNote() != "second note" {
			t.Errorf("expected 'second note', got %q", resp.Notes[1].GetNote())
		}
	})
}
