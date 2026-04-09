package server

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/pawelgrzybek/go-notes/gen/notes/v1"
	"github.com/pawelgrzybek/go-notes/internal/notes"
)

func (s *Server) GetNote(ctx context.Context, req *pb.GetNoteRequest) (*pb.GetNoteResponse, error) {
	noteID := req.GetId()
	note, err := s.svc.Get(ctx, int64(noteID))
	if err != nil {
		if errors.Is(err, notes.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "note with id %d not found", noteID)
		}
		return nil, status.Error(codes.Internal, "failed to get note")
	}

	return &pb.GetNoteResponse{
		Note: noteToProto(note),
	}, nil
}
