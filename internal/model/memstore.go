package model

import (
	"strings"
	"sync"
	"time"
)

type MemoryStore struct {
	mtx    sync.RWMutex
	data   map[string]User
	nextID int
}

func NewStore() *MemoryStore {
	return &MemoryStore{
		data: make(map[string]User),
	}
}

func (s *MemoryStore) Create(u User) (User, error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	lowEmail := strings.ToLower(u.Email)

	if _, ok := s.data[lowEmail]; ok {
		return User{}, ErrMailTaken
	}

	s.nextID++
	u.ID = s.nextID

	u.CreatedAt = time.Now()

	s.data[lowEmail] = u

	return u, nil
}

func (s *MemoryStore) GetByMail(email string) (User, error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	user, ok := s.data[email]
	if !ok {
		return User{}, ErrNotFound
	}

	return user, nil
}
