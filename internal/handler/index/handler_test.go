package index

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Marcel2603/spikeball-league/cmd/config"
)

func TestHandler(t *testing.T) {
	config.Configuration = config.Config{}

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Handler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}
