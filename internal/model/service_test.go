package model

import (
	"errors"
	"strings"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestIsValidPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		want     bool
	}{
		{"пустой", "", false},
		{"7 символов", "1234567", false},
		{"ровно 8 символов", "12345678", true},
		{"кириллица 8 букв", "парольчик", true}, // 9 букв — считаются символы, а не байты
		{"72 байта", strings.Repeat("a", 72), true},
		{"73 байта", strings.Repeat("a", 73), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isValidatePassword(tt.password)
			if got != tt.want {
				t.Errorf("isValidPassword(%q) = %v, want %v", tt.password, got, tt.want)
			}
		})
	}

}

func TestIsValidEmail(t *testing.T) {
	tests := []struct {
		name  string
		email string
		want  bool
	}{
		{"нормальный", "a@b.com", true},
		{"без @", "abc", false},
		{"пустой", "", false},
		{"с именем", "Имя <a@b.com>", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isValidateEmail(tt.email)
			if got != tt.want {
				t.Errorf("isValidEmail (%q) = %v. want %v", tt.email, got, tt.want)
			}
		})
	}
}

func TestRegister_OK(t *testing.T) {
	svc := NewService(NewStore())
	u, err := svc.Register("A@B.com", "supersecret")
	if err != nil {
		t.Fatalf("unexpected error: %v", err) // Fatalf: с ошибкой дальше смотреть нечего
	}
	if u.ID == 0 {
		t.Errorf("ID не проставлен")
	}
	if u.Email != "a@b.com" {
		t.Errorf("Email = %q, want %q", u.Email, "a@b.com")
	}
	if u.PasswordHash == "supersecret" {
		t.Errorf("пароль сохранён открытым текстом!")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte("supersecret")); err != nil {
		t.Errorf("хэш не соответствует паролю: %v", err)
	}
}

func TestRegister_InvalidInput(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		password string
	}{
		{"плохой email", "abc", "supersecret"},
		{"короткий пароль", "a@b.com", "123"},
		{"длинный пароль", "a@b.com", strings.Repeat("a", 73)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewService(NewStore())
			_, err := svc.Register(tt.email, tt.password)

			if !errors.Is(err, ErrInvalidOutput) {
				t.Errorf("Register(%q, %q) err = %v, want ErrInvalidOutput", tt.email, tt.password, err)
			}
		})
	}
}

func TestRegister_DuplicateEmail(t *testing.T) {
	svc := NewService(NewStore())

	if _, err := svc.Register("a@b.com", "supersecret"); err != nil {
		t.Fatalf("первая регистрация: unexpected error: %v", err)
	}

	_, err := svc.Register("a@b.com", "supersecret")
	if !errors.Is(err, ErrMailTaken) {
		t.Errorf("вторая регистрация: err = %v, want ErrMailTaken", err)
	}
}

func TestRegister_EmailCaseInsensitive(t *testing.T) {
	svc := NewService(NewStore())

	if _, err := svc.Register("a@b.com", "supersecret"); err != nil {
		t.Fatalf("первая регистрация: unexpected error: %v", err)
	}

	_, err := svc.Register("  A@B.COM ", "supersecret")
	if !errors.Is(err, ErrMailTaken) {
		t.Errorf("регистрация A@B.COM: err = %v, want ErrMailTaken", err)
	}
}

func TestMemoryStore_CreateAndGet(t *testing.T) {
	st := NewStore()

	u1, err := st.Create(User{Email: "first@b.com", PasswordHash: "h1"})
	if err != nil {
		t.Fatalf("Create #1: unexpected error: %v", err)
	}
	if u1.ID != 1 {
		t.Errorf("u1.ID = %d, want 1", u1.ID)
	}

	u2, err := st.Create(User{Email: "second@b.com", PasswordHash: "h2"})
	if err != nil {
		t.Fatalf("Create #2: unexpected error: %v", err)
	}
	if u2.ID != 2 {
		t.Errorf("u2.ID = %d, want 2", u2.ID)
	}

	_, err = st.Create(User{Email: "first@b.com", PasswordHash: "h3"})
	if !errors.Is(err, ErrMailTaken) {
		t.Errorf("повтор email: err = %v, want ErrMailTaken", err)
	}

	got, err := st.GetByMail("first@b.com")
	if err != nil {
		t.Fatalf("GetByMail(first@b.com): unexpected error: %v", err)
	}
	if got.Email != "first@b.com" {
		t.Errorf("got.Email = %q, want %q", got.Email, "first@b.com")
	}
	if got.ID != u1.ID {
		t.Errorf("got.ID = %d, want %d", got.ID, u1.ID)
	}

	_, err = st.GetByMail("nobody@b.com")
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("GetByMail(nobody): err = %v, want ErrNotFound", err)
	}
}
