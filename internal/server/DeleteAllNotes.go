package server

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/pawelgrzybek/go-notes/gen/notes/v1"
)

func (s *Server) DeleteAllNotes(ctx context.Context, req *pb.DeleteAllNotesRequest) (*pb.DeleteAllNotesResponse, error) {
	result, err := s.svc.DeleteAll(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete all notes: %v", err)
	}

	protoNotes := make([]*pb.Note, len(result))
	for i, n := range result {
		protoNotes[i] = noteToProto(n)
	}

	return &pb.DeleteAllNotesResponse{
		Notes: protoNotes,
	}, nil
}
