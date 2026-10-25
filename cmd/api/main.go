package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/DrollltedUp/notesJWT/internal/auth"
	"github.com/DrollltedUp/notesJWT/internal/model"
	"github.com/DrollltedUp/notesJWT/internal/note"
	"github.com/DrollltedUp/notesJWT/internal/postgres"
)

func main() {
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Error("JWT_SECRET is not set")
		os.Exit(1)
	}
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Error("DATABASE_URL is not set")
		os.Exit(1)
	}

	pool, err := postgres.NewPool(context.Background(), dbURL)
	if err != nil {
		log.Error("cannot connect to database", "err", err)
		os.Exit(1)
	}

	defer pool.Close()

	svc := model.NewService(model.NewStore())
	tokens := auth.NewTokenManager(secret, 24*time.Hour)
	crud := note.NewServiceCRUD(note.NewNotesMemoryStore())
	srv := NewServerHTTP(svc, log, tokens, crud)

	httpServer := &http.Server{
		Addr:              ":9097",
		Handler:           srv.routes(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
	}

	log.Info("Server started", "addr", httpServer.Addr)
	if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error("server failed", "err: ", err)
		os.Exit(1)
	}
}
