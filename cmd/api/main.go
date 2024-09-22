package main

import (
	"backend/internal/server"
	"log/slog"
	"os"
)

func main() {
	newServer := server.NewServer()

	err := newServer.ListenAndServe()
	if err != nil {
		slog.Default().Error("cannot start newServer", "erro", err)
		os.Exit(1)
	}
}
