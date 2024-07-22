package server

import (
	"fmt"
	_ "github.com/joho/godotenv/autoload"
	"net/http"
	"os"
	"strconv"
	"time"

	"backend/internal/database"
)

type Server struct {
	Port int

	DB database.Service
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	NewServer := &Server{
		Port: port,

		DB: database.New(),
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.Port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
