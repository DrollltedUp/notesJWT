package model

import (
	"errors"
	"math/rand"
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
	Create(email, password string) (User, error)
	GetByMail(email string) (User, error)
}

func Create(email, password string) (User, error) {
	user := User{
		ID:           rand.Intn(1000000),
		Email:        email,
		PasswordHash: password,
		CreatedAt:    time.Now(),
	}

	if user.Email == email {
		return User{}, ErrMailTaken
	}

	return user, nil
}

func (u User) GetByMail(email string) (User, error) {
	if u.Email != email {
		return User{}, ErrNotFound
	}

	return u, nil
}
