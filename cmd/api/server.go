package main

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/DrollltedUp/notesJWT/internal/auth"
	"github.com/DrollltedUp/notesJWT/internal/model"
)

type errorDTO struct {
	Message string    `json:"message"`
	Time    time.Time `json:"time"`
}

type ServerHTTP struct {
	svc    *model.Service
	log    *slog.Logger
	tokens *auth.TokenManager
}

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type registerResponse struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
}

type ctxKey struct{}

func NewServerHTTP(svc *model.Service, log *slog.Logger, tokens *auth.TokenManager) *ServerHTTP {
	return &ServerHTTP{
		svc:    svc,
		log:    log,
		tokens: tokens,
	}
}

func (s *ServerHTTP) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /register", s.handleRegister)
	mux.HandleFunc("POST /login", s.handleLogin)
	mux.Handle("GET /me", s.authMiddleware(http.HandlerFunc(s.handleMe)))

	return mux
}

func (s *ServerHTTP) handleRegister(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var request registerRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "bad request")
		return
	}

	user, err := s.svc.Register(request.Email, request.Password)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrInvalidOutput):
			writeError(w, http.StatusBadRequest, "not validate information: email or password")
		case errors.Is(err, model.ErrMailTaken):
			writeError(w, http.StatusConflict, "this email was taken")
		default:
			writeError(w, http.StatusInternalServerError, "internal error")
		}
		return
	}
	writeJSON(w, http.StatusCreated, registerResponse{
		ID:    user.ID,
		Email: user.Email,
	})

}

func (s *ServerHTTP) handleLogin(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var loginRequest loginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		writeError(w, http.StatusBadRequest, "bad request")
		return
	}

	user, err := s.svc.Login(loginRequest.Email, loginRequest.Password)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrInvalidCredentation):
			writeError(w, http.StatusUnauthorized, "invalid email or password")
		default:
			writeError(w, http.StatusInternalServerError, "internal error")
		}
		return
	}

	token, err := s.tokens.Issue(user.ID)
	if err != nil {
		s.log.Error("issue token failed", "err", err)
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	writeJSON(w, http.StatusOK, loginResponse{
		Token: token,
	})

}

func (s *ServerHTTP) handleMe(w http.ResponseWriter, r *http.Request) {
	id, ok := userIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	writeJSON(w, http.StatusOK, map[string]int{"user_id": id})
}

func (s *ServerHTTP) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		token, ok := strings.CutPrefix(header, "Bearer ")
		if token == "" || !ok {
			writeError(w, http.StatusUnauthorized, "invalid token")
			return
		}

		id, err := s.tokens.Parse(token)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "invalid token")
			return
		}
		ctx := context.WithValue(r.Context(), ctxKey{}, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func writeJSON(w http.ResponseWriter, statusCode int, v any) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, statusCode int, msg string) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(errorDTO{Message: msg, Time: time.Now()})
}

func userIDFromContext(ctx context.Context) (int, bool) {
	id, ok := ctx.Value(ctxKey{}).(int)
	return id, ok
}
