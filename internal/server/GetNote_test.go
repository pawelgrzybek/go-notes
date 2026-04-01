package server_test

import (
	"context"
	"testing"

	pb "github.com/pawelgrzybek/go-notes/gen/notes/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGetNote(t *testing.T) {
	t.Run("returns note by id", func(t *testing.T) {
		s := setupServer(t)

		created, _ := s.CreateNote(context.Background(), &pb.CreateNoteRequest{Note: new("my note")})

		resp, err := s.GetNote(context.Background(), &pb.GetNoteRequest{Id: created.Note.Id})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if resp.Note.GetNote() != "my note" {
			t.Errorf("expected 'my note', got %q", resp.Note.GetNote())
		}
	})

	t.Run("returns NotFound when note does not exist", func(t *testing.T) {
		s := setupServer(t)

		_, err := s.GetNote(context.Background(), &pb.GetNoteRequest{Id: new(int32(99))})
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
