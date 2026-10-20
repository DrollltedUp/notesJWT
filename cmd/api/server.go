package main

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/DrollltedUp/notesJWT/internal/model"
)

type errorDTO struct {
	Message string    `json:"message"`
	Time    time.Time `json:"time"`
}

type ServerHTTP struct {
	svc *model.Service
	log *slog.Logger
}

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type registerResponse struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

func NewServerHTTP(svc *model.Service, log *slog.Logger) *ServerHTTP {
	return &ServerHTTP{
		svc: svc,
		log: log,
	}
}

func (s *ServerHTTP) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /register", s.handleRegister)

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
