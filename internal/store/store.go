package store

import (
	"database/sql"
	"errors"

	"github.com/pawelgrzybek/go-notes/internal/models"
)

type Store interface {
	List() ([]models.Note, error)
	Get(noteID int) (models.Note, error)
	Create(note string) (models.Note, error)
	DeleteOne(noteID int) (models.Note, error)
	DeleteAll() ([]models.Note, error)
}

type notesStore struct {
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &notesStore{
		db: db,
	}
}

func (s *notesStore) List() ([]models.Note, error) {
	rows, err := s.db.Query("SELECT id, note, created_at FROM notes")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	notes := []models.Note{}
	for rows.Next() {
		var n models.Note
		if err := rows.Scan(&n.ID, &n.Note, &n.CreatedAt); err != nil {
			return nil, err
		}
		notes = append(notes, n)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return notes, nil

}

func (s *notesStore) Get(noteID int) (models.Note, error) {
	var note models.Note
	err := s.db.QueryRow("SELECT id, note, created_at FROM notes WHERE id = ?", noteID).Scan(&note.ID, &note.Note, &note.CreatedAt)
	if err != nil {
		return models.Note{}, err
	}

	return note, nil
}

func (s *notesStore) Create(note string) (models.Note, error) {
	result, err := s.db.Exec("INSERT INTO notes (note) VALUES (?)", note)
	if err != nil {
		return models.Note{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return models.Note{}, err

	}

	var n models.Note
	err = s.db.QueryRow("SELECT id, note, created_at FROM notes WHERE id = ?", id).Scan(&n.ID, &n.Note, &n.CreatedAt)
	if err != nil {
		return models.Note{}, err
	}

	return n, nil

}

func (s *notesStore) DeleteOne(noteID int) (models.Note, error) {
	var note models.Note
	err := s.db.QueryRow("SELECT id, note, created_at FROM notes WHERE id = ?", noteID).Scan(&note.ID, &note.Note, &note.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return models.Note{}, sql.ErrNoRows
	}
	if err != nil {
		return models.Note{}, err
	}

	_, err = s.db.Exec("DELETE FROM notes WHERE id = ?", noteID)
	if err != nil {
		return models.Note{}, err
	}

	return note, nil
}

func (s *notesStore) DeleteAll() ([]models.Note, error) {
	rows, err := s.db.Query("SELECT id, note, created_at FROM notes")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	notes := []models.Note{}
	for rows.Next() {
		var n models.Note
		if err := rows.Scan(&n.ID, &n.Note, &n.CreatedAt); err != nil {
			return nil, err
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
