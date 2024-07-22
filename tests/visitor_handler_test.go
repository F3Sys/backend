package tests

import (
	"backend/internal/database"
	"backend/internal/server"
	"encoding/json"
	"github.com/gofrs/uuid/v5"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
)

func TestVisitorHandler(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/visitor", nil)
	req.Header.Set(echo.HeaderXForwardedFor, "127.0.0.1")
	resp := httptest.NewRecorder()
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	s := &server.Server{
		Port: port,
		DB:   database.New(),
	}
	// Assertions
	if err := s.VisitorHandler(e.NewContext(req, resp)); err != nil {
		t.Errorf("handler() error = %v", err)
		return
	}
	if resp.Code != http.StatusOK {
		t.Errorf("handler() wrong status code = %v", resp.Code)
		return
	}
	var actual1 map[string]string
	// Decode the response body into the actual map
	if err := json.NewDecoder(resp.Body).Decode(&actual1); err != nil {
		t.Errorf("handler() error decoding response body: %v", err)
		return
	}
	// Compare the decoded response with the expected value
	t.Logf("actual1: %v", actual1)
	uuid1, err := uuid.FromString(actual1["uuid"])
	if err != nil {
		t.Errorf("handler() error parsing uuid: %v", err)
		return
	}

	// Call it again and check if it's the same uuid
	if err := s.VisitorHandler(e.NewContext(req, resp)); err != nil {
		t.Errorf("handler() error = %v", err)
		return
	}
	if resp.Code != http.StatusOK {
		t.Errorf("handler() wrong status code = %v", resp.Code)
		return
	}
	var actual2 map[string]string
	// Decode the response body into the actual map
	if err := json.NewDecoder(resp.Body).Decode(&actual2); err != nil {
		t.Errorf("handler() error decoding response body: %v", err)
		return
	}

	t.Logf("actual2: %v", actual2)
	uuid2, err := uuid.FromString(actual2["uuid"])
	if err != nil {
		t.Errorf("handler() error parsing uuid: %v", err)
		return
	}

	// Compare it to check if uuid is identical
	if uuid1.String() != uuid2.String() {
		t.Errorf("handler() wrong uuid1 = %v, uuid2 = %v", uuid1.String(), uuid2.String())
		return
	}
}
