-- name: ListNotes :many
SELECT id, note, created_at FROM notes;

-- name: GetNote :one
SELECT id, note, created_at FROM notes WHERE id = ?;

-- name: CreateNote :one
INSERT INTO notes (note) VALUES (?) RETURNING id, note, created_at;

-- name: DeleteNoteByID :one
DELETE FROM notes WHERE id = ? RETURNING id, note, created_at;

-- name: DeleteAllNotes :many
DELETE FROM notes RETURNING id, note, created_at;
