package note

import (
	"sort"
	"sync"
	"time"
)

type NotesMemoryStore struct {
	mtx    sync.RWMutex
	notes  map[int]Note
	nextID int
}

func NewNotesMemoryStore() *NotesMemoryStore {
	return &NotesMemoryStore{
		notes: make(map[int]Note),
	}
}

func (n *NotesMemoryStore) Create(note Note) (Note, error) {
	n.mtx.Lock()
	defer n.mtx.Unlock()

	n.nextID++
	note.ID = n.nextID

	now := time.Now()
	note.CreatedAt = now
	note.UpdatedAt = now

	n.notes[note.ID] = note
	return note, nil
}

func (n *NotesMemoryStore) GetByID(id int) (Note, error) {
	n.mtx.RLock()
	defer n.mtx.RUnlock()

	notes, ok := n.notes[id]
	if !ok {
		return Note{}, ErrNotFound
	}

	return notes, nil
}

func (n *NotesMemoryStore) ListByOwner(ownerID int) ([]Note, error) {
	n.mtx.RLock()
	defer n.mtx.RUnlock()

	var notes = make([]Note, 0)

	for _, s := range n.notes {
		if s.OwnerID == ownerID {
			notes = append(notes, s)
		}
	}
	sort.Slice(notes, func(i, j int) bool {
		return notes[i].ID < notes[j].ID
	})
	return notes, nil
}

func (n *NotesMemoryStore) Update(note Note) (Note, error) {
	n.mtx.Lock()
	defer n.mtx.Unlock()

	old, ok := n.notes[note.ID]
	if !ok {
		return Note{}, ErrNotFound
	}
	old.Title = note.Title
	old.Body = note.Body
	old.UpdatedAt = time.Now()
	n.notes[note.ID] = old
	return old, nil
}

func (n *NotesMemoryStore) Delete(id int) error {
	n.mtx.Lock()
	defer n.mtx.Unlock()

	if _, ok := n.notes[id]; !ok {
		return ErrNotFound
	}

	delete(n.notes, id)

	return nil

}
