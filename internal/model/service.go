package model

import (
	"errors"
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

	lowerEmail := normalizeEmail(email)

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

func (s *Service) Login(email, password string) (User, error) {
	lowerEmail := normalizeEmail(email)

	user, err := s.store.GetByMail(lowerEmail)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return User{}, ErrInvalidCredentation
		}
		return User{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return User{}, ErrInvalidCredentation
	}

	return user, nil
}

func isValidateEmail(email string) bool {
	addr, err := mail.ParseAddress(email)
	return err == nil && addr.Address == email
}

func isValidatePassword(password string) bool {
	return utf8.RuneCountInString(password) >= 8 && len(password) <= 72
}

func normalizeEmail(email string) string {
	normalizeEmail := strings.ToLower(strings.TrimSpace(email))
	return normalizeEmail
}

// СRUD Service
func (s *Service) Create(userID, title, body string) (Note, error) {

}
