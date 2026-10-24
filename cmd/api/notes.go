package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/DrollltedUp/notesJWT/internal/note"
)

type noteRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type noteResponse struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func toNoteResponse(n note.Note) noteResponse {
	return noteResponse{
		ID:        n.ID,
		Title:     n.Title,
		Body:      n.Body,
		CreatedAt: n.CreatedAt,
		UpdatedAt: n.UpdatedAt,
	}
}

// Create
func (s *ServerHTTP) CreateHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := s.currentUserID(w, r)
	if !ok {
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	var req noteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "bad request")
		return
	}

	n, err := s.crud.Create(userID, req.Title, req.Body)
	if err != nil {
		s.noteError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, toNoteResponse(n))

}

// Get All Notes
func (s *ServerHTTP) GetAllHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := s.currentUserID(w, r)
	if !ok {
		return
	}
	notes, err := s.crud.List(userID)
	if err != nil {
		s.noteError(w, err)
		return
	}

	resp := make([]noteResponse, 0, len(notes))
	for _, n := range resp {
		resp = append(resp, n)
	}

	writeJSON(w, http.StatusOK, resp)

}

// Get Notes by ID
func (s *ServerHTTP) GetNotesByIDHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := s.currentUserID(w, r)
	if !ok {
		return
	}

	id, ok := parseNoteID(w, r)
	if !ok {
		return
	}

	note, err := s.crud.Get(userID, id)
	if err != nil {
		s.noteError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toNoteResponse(note))

}

// Update notes
func (s *ServerHTTP) UpdateNoteHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := s.currentUserID(w, r)
	if !ok {
		return
	}

	id, ok := parseNoteID(w, r)
	if !ok {
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	var req noteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "bad request")
		return
	}

	n, err := s.crud.Update(userID, id, req.Title, req.Body)
	if err != nil {
		s.noteError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toNoteResponse(n))
}

// Delete note
func (s *ServerHTTP) DeleteNoteHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := s.currentUserID(w, r)
	if !ok {
		return
	}

	id, ok := parseNoteID(w, r)
	if !ok {
		return
	}

	if err := s.crud.Delete(userID, id); err != nil {
		s.noteError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *ServerHTTP) currentUserID(w http.ResponseWriter, r *http.Request) (int, bool) {
	id, ok := userIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusInternalServerError, "internal error")
	}

	return id, ok
}

func parseNoteID(w http.ResponseWriter, r *http.Request) (int, bool) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid note id")
		return 0, false
	}

	return id, true
}

func (s *ServerHTTP) noteError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, note.ErrNotFound):
		writeError(w, http.StatusNotFound, "note not found")
	case errors.Is(err, note.ErrInvalidInput):
		writeError(w, http.StatusBadRequest, "invalid title or body")
	default:
		s.log.Error("note operation failed", "err", err)
		writeError(w, http.StatusInternalServerError, "internal error")
	}
}
