package server_test

import (
	"context"
	"testing"

	pb "github.com/pawelgrzybek/go-notes/gen/notes/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestDeleteNote(t *testing.T) {
	t.Run("deletes note and returns it", func(t *testing.T) {
		s := setupServer(t)

		created, _ := s.CreateNote(context.Background(), &pb.CreateNoteRequest{Note: new("to delete")})

		resp, err := s.DeleteNote(context.Background(), &pb.DeleteNoteRequest{Id: created.Note.Id})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if resp.Note.GetNote() != "to delete" {
			t.Errorf("expected 'to delete', got %q", resp.Note.GetNote())
		}

		// verify it's actually gone
		listResp, _ := s.ListNotes(context.Background(), &pb.ListNotesRequest{})
		if len(listResp.Notes) != 0 {
			t.Errorf("expected 0 notes after delete, got %d", len(listResp.Notes))
		}
	})

	t.Run("returns NotFound when note does not exist", func(t *testing.T) {
		s := setupServer(t)

		_, err := s.DeleteNote(context.Background(), &pb.DeleteNoteRequest{Id: new(int32(99))})
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		st, ok := status.FromError(err)
		if !ok {
			t.Fatalf("expected gRPC status error, got %v", err)
		}

		if st.Code() != codes.NotFound {
			t.Errorf("expected code NotFound, got %v", st.Code())
		}
	})
}
