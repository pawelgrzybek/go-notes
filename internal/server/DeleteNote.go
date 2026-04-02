package server

import (
	"context"
	"database/sql"
	"errors"

	pb "github.com/pawelgrzybek/go-notes/gen/notes/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) DeleteNote(ctx context.Context, req *pb.DeleteNoteRequest) (*pb.DeleteNoteResponse, error) {
	note, err := s.Store.DeleteOne(req.GetId())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "note with id %d not found", *req.Id)
		}
		return nil, status.Errorf(codes.Internal, "failed to delete note: %v", err)
	}

	if s.Notifier != nil {
		s.Notifier.Notify()
	}

	return &pb.DeleteNoteResponse{
		Note: note,
	}, nil
}
