package main

import "net/http"

type noteRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

// Create
func (s *ServerHTTP) CreateHandler(w http.ResponseWriter, r *http.Request) {

}

// Get All Notes
func (s *ServerHTTP) GetAllHandler(w http.ResponseWriter, r *http.Request) {

}

// Get Notes by ID
func (s *ServerHTTP) GetNotesByIDHandler(w http.ResponseWriter, r *http.Request) {

}

// Update notes
func (s *ServerHTTP) UpdateNoteHandler(w http.ResponseWriter, r *http.Request) {

}
