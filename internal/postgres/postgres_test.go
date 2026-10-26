package postgres

// import (
// 	"context"
// 	"errors"
// 	"os"
// 	"testing"

// 	"github.com/jackc/pgx/v5/pgxpool"

// 	"github.com/DrollltedUp/notesJWT/internal/model"
// 	"github.com/DrollltedUp/notesJWT/internal/note"
// )

// func testPool(t *testing.T) *pgxpool.Pool {
// 	t.Helper()

// 	url := os.Getenv("TEST_DATABASE_URL")
// 	if url == "" {
// 		t.Skip("TEST_DATABASE_URL не задан: тест на настоящей базе пропущен")
// 	}

// 	pool, err := NewPool(context.Background(), url)
// 	if err != nil {
// 		t.Fatalf("подключение к тестовой базе: %v", err)
// 	}
// 	t.Cleanup(pool.Close)

// 	_, err = pool.Exec(context.Background(), `TRUNCATE users, notes RESTART IDENTITY CASCADE`)
// 	if err != nil {
// 		t.Fatalf("очистка тестовой базы: %v", err)
// 	}

// 	return pool
// }

// func TestUserStore(t *testing.T) {
// 	users := NewUserStore(testPool(t))

// 	u, err := users.Create(model.User{Email: "a@b.com", PasswordHash: "hash"})
// 	if err != nil {
// 		t.Fatalf("Create: %v", err)
// 	}
// 	if u.ID == 0 || u.CreatedAt.IsZero() {
// 		t.Errorf("база не выдала ID или время: %+v", u)
// 	}

// 	if _, err := users.Create(model.User{Email: "a@b.com", PasswordHash: "other"}); !errors.Is(err, model.ErrMailTaken) {
// 		t.Errorf("повторный email: err = %v, want ErrMailTaken", err)
// 	}

// 	got, err := users.GetByMail("a@b.com")
// 	if err != nil || got.ID != u.ID || got.PasswordHash != "hash" {
// 		t.Errorf("GetByMail: %+v, err=%v", got, err)
// 	}

// 	if _, err := users.GetByMail("nobody@b.com"); !errors.Is(err, model.ErrNotFound) {
// 		t.Errorf("неизвестный email: err = %v, want ErrNotFound", err)
// 	}
// }

// func TestNoteStore_Ownership(t *testing.T) {
// 	pool := testPool(t)
// 	users := NewUserStore(pool)
// 	svc := note.NewServiceCRUD(NewNoteStore(pool))

// 	a, _ := users.Create(model.User{Email: "a@b.com", PasswordHash: "h"})
// 	b, _ := users.Create(model.User{Email: "c@d.com", PasswordHash: "h"})

// 	n, err := svc.Create(a.ID, "Title", "Body")
// 	if err != nil {
// 		t.Fatalf("Create: %v", err)
// 	}
// 	if n.ID == 0 {
// 		t.Fatal("ID не выдан")
// 	}

// 	if _, err := svc.Get(b.ID, n.ID); !errors.Is(err, note.ErrNotFound) {
// 		t.Errorf("Get чужой: err = %v, want ErrNotFound", err)
// 	}
// 	if _, err := svc.Update(b.ID, n.ID, "взлом", "x"); !errors.Is(err, note.ErrNotFound) {
// 		t.Errorf("Update чужой: err = %v, want ErrNotFound", err)
// 	}
// 	if err := svc.Delete(b.ID, n.ID); !errors.Is(err, note.ErrNotFound) {
// 		t.Errorf("Delete чужой: err = %v, want ErrNotFound", err)
// 	}
// 	if list, _ := svc.List(b.ID); len(list) != 0 {
// 		t.Errorf("List(B) = %d заметок, want 0", len(list))
// 	}

// 	got, err := svc.Get(a.ID, n.ID)
// 	if err != nil || got.Title != "Title" {
// 		t.Fatalf("после попыток взлома: %+v, err=%v", got, err)
// 	}

// 	upd, err := svc.Update(a.ID, n.ID, "New", "b")
// 	if err != nil || upd.Title != "New" || upd.OwnerID != a.ID {
// 		t.Errorf("Update своей: %+v, err=%v", upd, err)
// 	}
// 	if !upd.CreatedAt.Equal(n.CreatedAt) {
// 		t.Errorf("CreatedAt изменился: %v → %v", n.CreatedAt, upd.CreatedAt)
// 	}
// 	if !upd.UpdatedAt.After(n.UpdatedAt) {
// 		t.Errorf("UpdatedAt не обновился")
// 	}

// 	if err := svc.Delete(a.ID, n.ID); err != nil {
// 		t.Errorf("Delete своей: %v", err)
// 	}
// 	if _, err := svc.Get(a.ID, n.ID); !errors.Is(err, note.ErrNotFound) {
// 		t.Errorf("Get после удаления: err = %v, want ErrNotFound", err)
// 	}
// }

// // Злой «заголовок» — это просто текст, а не команда.
// func TestNoteStore_SQLInjectionIsHarmless(t *testing.T) {
// 	pool := testPool(t)
// 	u, _ := NewUserStore(pool).Create(model.User{Email: "a@b.com", PasswordHash: "h"})
// 	notes := NewNoteStore(pool)

// 	evil := `x'); DROP TABLE notes; --`
// 	n, err := notes.Create(note.Note{OwnerID: u.ID, Title: evil})
// 	if err != nil {
// 		t.Fatalf("Create: %v", err)
// 	}
// 	got, err := notes.GetByID(n.ID)
// 	if err != nil || got.Title != evil {
// 		t.Errorf("заголовок исказился: %q, err=%v", got.Title, err)
// 	}
// }
