package note

import (
	"strings"
	"unicode/utf8"
)

type ServiceCRUD struct {
	store Store
}

func NewServiceCRUD(s Store) *ServiceCRUD {
	return &ServiceCRUD{
		store: s,
	}
}

func (s *ServiceCRUD) Create(userID int, title, body string) (Note, error) {
	if title == "" || body == "" {
		return Note{}, ErrInvalidInput
	}

	note := Note{
		OwnerID: userID,
		Title:   title,
		Body:    body,
	}

	created, err := s.store.Create(note)
	if err != nil {
		return Note{}, err
	}

	return created, err
}

func (s *ServiceCRUD) Get(userID, id int) (Note, error) {
	return s.owned(userID, id)
}

func (s *ServiceCRUD) List(userID int) ([]Note, error) {
	return s.store.ListByOwner(userID)
}

func (s *ServiceCRUD) Update(userID, id int, title, body string) (Note, error) {
	n, err := s.owned(userID, id)
	if err != nil {
		return Note{}, err
	}

	title, err = validate(title, body)
	if err != nil {
		return Note{}, err
	}

	n.Title, n.Body = title, body
	return s.store.Update(n)
}

func (s *ServiceCRUD) Delete(userID, id int) error {
	if _, err := s.owned(userID, id); err != nil {
		return err
	}

	return s.store.Delete(id)
}

func (s *ServiceCRUD) owned(userID, id int) (Note, error) {
	n, err := s.store.GetByID(id)
	if err != nil {
		return Note{}, err
	}

	if n.OwnerID != userID {
		return Note{}, ErrNotFound
	}

	return n, nil
}

func validate(title, body string) (string, error) {
	title = strings.TrimSpace(title)
	if title == "" || utf8.RuneCountInString(title) > 200 || utf8.RuneCountInString(body) > 10000 {
		return "", ErrInvalidInput
	}
	return title, nil
}
