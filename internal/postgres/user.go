package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/DrollltedUp/notesJWT/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ model.Store = (*UserStore)(nil)

type UserStore struct {
	pool *pgxpool.Pool
}

func NewUserStore(pool *pgxpool.Pool) *UserStore {
	return &UserStore{
		pool: pool,
	}
}

func (s *UserStore) Create(u model.User) (model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	err := s.pool.QueryRow(ctx, `INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id, created_at`, u.Email, u.PasswordHash).Scan(&u.ID, &u.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == uniqueViolation {
			return model.User{}, model.ErrMailTaken
		}
		return model.User{}, fmt.Errorf("create user: %w", err)
	}
	return u, nil
}

func (s *UserStore) GetByMail(email string) (model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	var u model.User
	err := s.pool.QueryRow(ctx, `SELECT id, email, password_hash, created_at FROM users WHERE email = $1`, email).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.User{}, model.ErrNotFound
	}
	if err != nil {
		return model.User{}, fmt.Errorf("get user by email: %w", err)
	}

	return u, nil
}
