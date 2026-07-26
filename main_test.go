package main

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Marcel2603/spikeball-league/cmd/config"
	"github.com/Marcel2603/spikeball-league/internal/db"
	staticfiles "github.com/Marcel2603/spikeball-league/internal/handler/static-files"
)

func TestServerStarts(t *testing.T) {
	c := config.Config{}
	c.Server.Host = "localhost"

	mockLogger := slog.New(slog.NewTextHandler(io.Discard, nil))

	staticfiles.NewHandler(staticFiles)

	dbConn, err := db.InitDB(":memory:")
	if err != nil {
		t.Fatalf("Failed to init memory db: %v", err)
	}
	defer dbConn.Close()
	queries := db.New(dbConn)

	app, err := setupApp(c, mockLogger, queries)
	if err != nil {
		t.Fatalf("Failed to setup app: %v", err)
	}

	ts := httptest.NewServer(app)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d", resp.StatusCode)
	}
}
