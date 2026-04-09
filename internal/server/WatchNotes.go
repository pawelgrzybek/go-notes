package server

import (
	pb "github.com/pawelgrzybek/go-notes/gen/notes/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) WatchNotes(req *pb.WatchNotesRequest, stream pb.NoteService_WatchNotesServer) error {
	ch := s.svc.Watch(stream.Context())

	if err := s.sendSnapshot(stream); err != nil {
		return err
	}

	for {
		select {
		case <-stream.Context().Done():
			return stream.Context().Err()
		case _, ok := <-ch:
			if !ok {
				return nil
			}
			if err := s.sendSnapshot(stream); err != nil {
				return err
			}
		}
	}
}

func (s *Server) sendSnapshot(stream pb.NoteService_WatchNotesServer) error {
	result, err := s.svc.List(stream.Context())
	if err != nil {
		return status.Errorf(codes.Internal, "failed to list notes: %v", err)
	}

	protoNotes := make([]*pb.Note, len(result))
	for i, n := range result {
		protoNotes[i] = noteToProto(n)
	}

	return stream.Send(&pb.WatchNotesResponse{Notes: protoNotes})
}
