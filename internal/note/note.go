package note

import (
	"errors"
	"time"
)

type Note struct {
	ID        int
	OwnerID   int
	Title     string
	Body      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

var (
	ErrNotFound     = errors.New("note not found")
	ErrInvalidInput = errors.New("invalid note")
)

type Store interface {
	Create(n Note) (Note, error)
	GetByID(id int) (Note, error)
	ListByOwner(ownerID int) ([]Note, error)
	Update(n Note) (Note, error)
	Delete(id int) error
}
