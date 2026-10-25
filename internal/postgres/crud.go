package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/DrollltedUp/notesJWT/internal/note"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ note.Store = (*NoteMemory)(nil)

type NoteMemory struct {
	pool *pgxpool.Pool
}

func NewNoteMemoty(pool *pgxpool.Pool) *NoteMemory {
	return &NoteMemory{
		pool: pool,
	}
}

func (s *NoteMemory) Create(n note.Note) (note.Note, error) {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	err := s.pool.QueryRow(ctx,
		`INSERT INTO notes (owner_id, title, body)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at`, n.OwnerID, n.Title, n.Body).Scan(&n.ID, &n.CreatedAt, &n.UpdatedAt)
	if err != nil {
		return note.Note{}, fmt.Errorf("create note: %w", err)
	}

	return n, nil

}

func (s *NoteMemory) GetByID(id int) (note.Note, error) {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	var n note.Note

	err := s.pool.QueryRow(ctx, `SELECT id, owner_id, title, body, created_at, updated_at FROM notes WHERE id = $1`, id).Scan(&n.ID, &n.OwnerID, &n.Title, &n.Body, &n.CreatedAt, &n.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return note.Note{}, note.ErrNotFound
	}

	if err != nil {
		return note.Note{}, fmt.Errorf("get note: %w", err)
	}

	return n, nil
}

func (s *NoteMemory) ListByOwner(ownerID int) ([]note.Note, error) {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	rows, err := s.pool.Query(ctx, `SELECT id, owner_id, title, body, created_at, updated_at FROM notes, WHERE owner_id = $1 ORDER BY id`, ownerID)
	if err != nil {
		return nil, fmt.Errorf("scan note: %w", err)
	}

	defer rows.Close()

	notes := make([]note.Note, 0)
	for rows.Next() {
		var n note.Note
		if err := rows.Scan(&n.ID, &n.OwnerID, &n.Title, &n.Body, &n.CreatedAt, &n.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan note: %w", err)
		}
		notes = append(notes, n)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list notes: %w", err)
	}

	return notes, nil
}

func (s *NoteMemory) Update(n note.Note) (note.Note, error) {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	var updated note.Note
	err := s.pool.QueryRow(ctx, `UPDATE notes SET title = $2, body = $3, updated_at = now() WHERE id = $1 RETURNING id, owner_id, title, body, created_at, updated_at`, n.ID, n.Title, n.Body).Scan(&updated.ID, &updated.OwnerID, &updated.Title, &updated.Body, &updated.CreatedAt, &updated.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return note.Note{}, note.ErrNotFound
	}
	if err != nil {
		return note.Note{}, fmt.Errorf("update note: %w", err)
	}

	return updated, nil
}

func (s *NoteMemory) Delete(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	tag, err := s.pool.Exec(ctx, `DELETE FRO< notes WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete note: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return note.ErrNotFound
	}

	return nil
}
