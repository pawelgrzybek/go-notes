package store

import (
	"database/sql"
	"errors"
	"time"

	pb "github.com/pawelgrzybek/go-notes/gen/notes/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Store interface {
	List() ([]*pb.Note, error)
	Get(noteID int32) (*pb.Note, error)
	Create(note string) (*pb.Note, error)
	DeleteOne(noteID int32) (*pb.Note, error)
	DeleteAll() ([]*pb.Note, error)
}

type notesStore struct {
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &notesStore{
		db: db,
	}
}

func (s *notesStore) List() ([]*pb.Note, error) {
	rows, err := s.db.Query("SELECT id, note, created_at FROM notes")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	notes := []*pb.Note{}
	for rows.Next() {
		var id int32
		var note string
		var createdAt time.Time

		if err := rows.Scan(&id, &note, &createdAt); err != nil {
			return nil, err
		}

		n := &pb.Note{
			Id:        &id,
			Note:      &note,
			CreatedAt: timestamppb.New(createdAt),
		}

		notes = append(notes, n)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return notes, nil
}

func (s *notesStore) Get(noteID int32) (*pb.Note, error) {
	var id int32
	var note string
	var createdAt time.Time
	err := s.db.QueryRow("SELECT id, note, created_at FROM notes WHERE id = ?", noteID).Scan(&id, &note, &createdAt)
	if err != nil {
		return &pb.Note{}, err
	}

	return &pb.Note{
		Id:        &id,
		Note:      &note,
		CreatedAt: timestamppb.New(createdAt),
	}, nil
}

func (s *notesStore) Create(n string) (*pb.Note, error) {
	result, err := s.db.Exec("INSERT INTO notes (note) VALUES (?)", n)
	if err != nil {
		return &pb.Note{}, err
	}

	idLast, err := result.LastInsertId()
	if err != nil {
		return &pb.Note{}, err

	}

	var id int32
	var note string
	var createdAt time.Time
	err = s.db.QueryRow("SELECT id, note, created_at FROM notes WHERE id = ?", idLast).Scan(&id, &note, &createdAt)
	if err != nil {
		return &pb.Note{}, err
	}

	return &pb.Note{
		Id:        &id,
		Note:      &note,
		CreatedAt: timestamppb.New(createdAt),
	}, nil

}

func (s *notesStore) DeleteOne(noteID int32) (*pb.Note, error) {
	var id int32
	var note string
	var createdAt time.Time

	err := s.db.QueryRow("SELECT id, note, created_at FROM notes WHERE id = ?", noteID).Scan(&id, &note, &createdAt)
	if errors.Is(err, sql.ErrNoRows) {
		return &pb.Note{}, sql.ErrNoRows
	}
	if err != nil {
		return &pb.Note{}, err
	}

	_, err = s.db.Exec("DELETE FROM notes WHERE id = ?", noteID)
	if err != nil {
		return &pb.Note{}, err
	}

	return &pb.Note{
		Id:        &id,
		Note:      &note,
		CreatedAt: timestamppb.New(createdAt),
	}, nil
}

func (s *notesStore) DeleteAll() ([]*pb.Note, error) {
	rows, err := s.db.Query("SELECT id, note, created_at FROM notes")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	notes := []*pb.Note{}
	for rows.Next() {
		var id int32
		var note string
		var createdAt time.Time
		if err := rows.Scan(&id, &note, &createdAt); err != nil {
			return nil, err
		}
		n := &pb.Note{
			Id:        &id,
			Note:      &note,
			CreatedAt: timestamppb.New(createdAt),
		}
		notes = append(notes, n)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	_, err = s.db.Exec("DELETE FROM notes")
	if err != nil {
		return nil, err
	}

	return notes, nil
}
