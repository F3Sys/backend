package main

import (
	"backend/internal/server"
	"log/slog"
	"os"
)

func main() {
	slog.SetDefault(slog.Default())

	newServer := server.NewServer()

	err := newServer.ListenAndServe()
	if err != nil {
		slog.Error("cannot start newServer", "error", err)
		os.Exit(1)
	}
}
