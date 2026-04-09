package server

import (
	"context"
	"errors"

	pb "github.com/pawelgrzybek/go-notes/gen/notes/v1"
	"github.com/pawelgrzybek/go-notes/internal/notes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) DeleteNote(ctx context.Context, req *pb.DeleteNoteRequest) (*pb.DeleteNoteResponse, error) {
	note, err := s.svc.Delete(ctx, int64(req.GetId()))
	if err != nil {
		if errors.Is(err, notes.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "note with id %d not found", req.GetId())
		}
		return nil, status.Errorf(codes.Internal, "failed to delete note: %v", err)
	}

	return &pb.DeleteNoteResponse{
		Note: noteToProto(note),
	}, nil
}
