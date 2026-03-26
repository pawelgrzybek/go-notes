package server

import (
	"context"
	"database/sql"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/pawelgrzybek/go-notes/gen/notes/v1"
)

func (s *Server) GetNote(ctx context.Context, req *pb.GetNoteRequest) (*pb.GetNoteResponse, error) {
	noteID := req.GetId()

	note, err := s.Store.Get(noteID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "note with id %d not found", noteID)
		}

		return nil, status.Error(codes.Internal, "failed to get note")
	}

	return &pb.GetNoteResponse{
		Note: note,
	}, nil

}
