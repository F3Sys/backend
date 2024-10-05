package server

import (
	"net/http"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"backend/internal/database"
)

type Server struct {
	DB database.Service
}

func NewServer() *http.Server {
	NewServer := &Server{
		DB: database.New(),
	}

	// Declare Server config
	server := &http.Server{
		Addr:         "0.0.0.0:8080",
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
