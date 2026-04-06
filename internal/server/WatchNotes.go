package server

import (
	pb "github.com/pawelgrzybek/go-notes/gen/notes/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) WatchNotes(req *pb.WatchNotesRequest, stream pb.NoteService_WatchNotesServer) error {
	rows, err := s.q.ListNotes(stream.Context())
	if err != nil {
		return status.Errorf(codes.Internal, "failed to list notes: %v", err)
	}

	notes := make([]*pb.Note, len(rows))
	for i, r := range rows {
		notes[i] = noteToProto(r)
	}
	if err := stream.Send(&pb.WatchNotesResponse{Notes: notes}); err != nil {
		return err
	}

	ch := s.Notifier.Subscribe()
	defer s.Notifier.Unsubscribe(ch)

	for {
		select {
		case <-stream.Context().Done():
			return stream.Context().Err()
		case _, ok := <-ch:
			if !ok {
				return nil
			}
			rows, err := s.q.ListNotes(stream.Context())
			if err != nil {
				return status.Errorf(codes.Internal, "failed to list notes: %v", err)
			}

			notes := make([]*pb.Note, len(rows))
			for i, r := range rows {
				notes[i] = noteToProto(r)
			}
			if err := stream.Send(&pb.WatchNotesResponse{Notes: notes}); err != nil {
				return err
			}
		}
	}
}
