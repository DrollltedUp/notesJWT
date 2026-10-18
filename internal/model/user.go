package model

import (
	"errors"
	"time"
)

type User struct {
	ID           int
	Email        string
	PasswordHash string
	CreatedAt    time.Time
}

var (
	ErrMailTaken           = errors.New("email already exist")
	ErrNotFound            = errors.New("email not found")
	ErrInvalidOutput       = errors.New("invalid mail or password format")
	ErrInvalidCredentation = errors.New("invalid mail or password")
)

type Store interface {
	Create(u User) (User, error)
	GetByMail(email string) (User, error)
}
