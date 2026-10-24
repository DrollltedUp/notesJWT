package note

import (
	"errors"
	"testing"
)

func TestOwnership(t *testing.T) {
	svc := NewServiceCRUD(NewNotesMemoryStore())

	n, err := svc.Create(1, "Title", "Body")
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if n.ID == 0 {
		t.Fatal("ID не выдан")
	}

	if _, err := svc.Get(2, n.ID); !errors.Is(err, ErrNotFound) {
		t.Errorf("Get чужой: err = %v, want ErrNotFound", err)
	}
	if _, err := svc.Update(2, n.ID, "взлом", "x"); !errors.Is(err, ErrNotFound) {
		t.Errorf("Update чужой: err = %v, want ErrNotFound", err)
	}
	if err := svc.Delete(2, n.ID); !errors.Is(err, ErrNotFound) {
		t.Errorf("Delete чужой: err = %v, want ErrNotFound", err)
	}
	if list, _ := svc.List(2); len(list) != 0 {
		t.Errorf("List(2) = %d заметок, want 0", len(list))
	}

	got, err := svc.Get(1, n.ID)
	if err != nil || got.Title != "Title" {
		t.Fatalf("после попыток взлома: %+v, err=%v", got, err)
	}

	upd, err := svc.Update(1, n.ID, "New", "b")
	if err != nil || upd.Title != "New" || upd.OwnerID != 1 {
		t.Errorf("Update своей: %+v, err=%v", upd, err)
	}
	if err := svc.Delete(1, n.ID); err != nil {
		t.Errorf("Delete своей: %v", err)
	}
	if _, err := svc.Get(1, n.ID); !errors.Is(err, ErrNotFound) {
		t.Errorf("Get после удаления: err = %v, want ErrNotFound", err)
	}
}
