package server_test

import (
	"context"
	"testing"

	pb "github.com/pawelgrzybek/go-notes/gen/notes/v1"
)

func TestDeleteAllNotes(t *testing.T) {
	t.Run("deletes all notes and returns them", func(t *testing.T) {
		s := setupServer(t)

		s.CreateNote(context.Background(), &pb.CreateNoteRequest{Note: new("first")})
		s.CreateNote(context.Background(), &pb.CreateNoteRequest{Note: new("second")})

		resp, err := s.DeleteAllNotes(context.Background(), &pb.DeleteAllNotesRequest{})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(resp.Notes) != 2 {
			t.Fatalf("expected 2 deleted notes, got %d", len(resp.Notes))
		}

		if resp.Notes[0].GetNote() != "first" {
			t.Errorf("expected first note to be 'first', got %q", resp.Notes[0].GetNote())
		}

		if resp.Notes[1].GetNote() != "second" {
			t.Errorf("expected second note to be 'second', got %q", resp.Notes[1].GetNote())
		}

		// verify they're actually gone
		listResp, _ := s.ListNotes(context.Background(), &pb.ListNotesRequest{})
		if len(listResp.Notes) != 0 {
			t.Errorf("expected 0 notes after delete all, got %d", len(listResp.Notes))
		}
	})

	t.Run("returns empty list when no notes to delete", func(t *testing.T) {
		s := setupServer(t)

		resp, err := s.DeleteAllNotes(context.Background(), &pb.DeleteAllNotesRequest{})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(resp.Notes) != 0 {
			t.Fatalf("expected 0 notes, got %d", len(resp.Notes))
		}
	})
}
