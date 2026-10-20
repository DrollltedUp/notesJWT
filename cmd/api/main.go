package main

import (
	"errors"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/DrollltedUp/notesJWT/internal/model"
)

func main() {
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	svc := model.NewService(model.NewStore())
	srv := NewServerHTTP(svc, log)

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
