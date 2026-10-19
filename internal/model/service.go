package model

import (
	"net/mail"
	"strings"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	store Store
}

func NewService(s Store) *Service {
	return &Service{
		store: s,
	}

}

func (s *Service) Register(email, password string) (User, error) {

	lowerEmail := strings.ToLower(strings.TrimSpace(email))

	if !isValidateEmail(lowerEmail) || !isValidatePassword(password) {
		return User{}, ErrInvalidOutput
	}

	// Хэширование пароля
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}

	user := User{
		Email:        lowerEmail,
		PasswordHash: string(hash),
	}

	created, err := s.store.Create(user)
	if err != nil {
		return User{}, err
	}

	return created, nil

}

func isValidateEmail(email string) bool {
	addr, err := mail.ParseAddress(email)
	return err == nil && addr.Address == email
}

func isValidatePassword(password string) bool {
	return utf8.RuneCountInString(password) >= 8 && len(password) <= 72
}
