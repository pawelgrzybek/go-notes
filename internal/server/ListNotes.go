package server

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/pawelgrzybek/go-notes/gen/notes/v1"
)

func (s *Server) ListNotes(ctx context.Context, req *pb.ListNotesRequest) (*pb.ListNotesResponse, error) {
	list, err := s.Store.List()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list notes: %v", err)
	}

	return &pb.ListNotesResponse{
		Notes: list,
	}, nil

}
