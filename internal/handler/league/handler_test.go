package league

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/Marcel2603/spikeball-league/internal/db"
	"github.com/go-chi/chi/v5"
)

func TestIndex(t *testing.T) {
	mock := &MockQuerier{}
	h := NewHandler(mock, "http://localhost")

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	h.Index(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("got status %d; want %d", w.Code, http.StatusOK)
	}
}

func TestCreateLeague(t *testing.T) {
	mock := &MockQuerier{
		CreateLeagueFunc: func(_ context.Context, arg db.CreateLeagueParams) (db.League, error) {
			if arg.Name == "FailLeague" {
				return db.League{}, fmt.Errorf("db error")
			}
			return db.League{
				ID:       1,
				PublicID: "pub-1",
				AdminID:  "adm-1",
				Name:     arg.Name,
			}, nil
		},
	}
	h := NewHandler(mock, "http://localhost")

	t.Run("Valid creation", func(t *testing.T) {
		form := url.Values{}
		form.Add("name", "Test League")
		req := httptest.NewRequest(http.MethodPost, "/leagues", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()

		h.CreateLeague(w, req)

		if w.Code != http.StatusSeeOther {
			t.Errorf("got status %d; want %d", w.Code, http.StatusSeeOther)
		}

		if !strings.HasPrefix(w.Header().Get("Location"), "/admin/") {
			t.Errorf("got redirect %q; want prefix /admin/", w.Header().Get("Location"))
		}
	})

	t.Run("DB error", func(t *testing.T) {
		form := url.Values{}
		form.Add("name", "FailLeague")
		req := httptest.NewRequest(http.MethodPost, "/leagues", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()

		h.CreateLeague(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("got status %d; want %d", w.Code, http.StatusInternalServerError)
		}
	})
}

func TestCheckLeagueExists(t *testing.T) {
	mock := &MockQuerier{
		CheckLeagueExistsFunc: func(_ context.Context, arg db.CheckLeagueExistsParams) (int64, error) {
			if arg.AdminID == "exists" || arg.PublicID == "exists" {
				return 1, nil
			}
			return 0, sql.ErrNoRows
		},
	}
	h := NewHandler(mock, "http://localhost")

	r := chi.NewRouter()
	h.RegisterRoutes(r)

	t.Run("Exists", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/leagues/exists/exists", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("got status %d; want %d", w.Code, http.StatusOK)
		}
	})

	t.Run("Not found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/leagues/missing/exists", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("got status %d; want %d", w.Code, http.StatusNotFound)
		}
	})
}

func TestLogMatch(t *testing.T) {
	mock := &MockQuerier{
		GetLeagueByPublicIDFunc: func(_ context.Context, publicID string) (db.League, error) {
			if publicID == "notfound" {
				return db.League{}, sql.ErrNoRows
			}
			return db.League{ID: 1, PublicID: publicID}, nil
		},
		LogMatchFunc: func(_ context.Context, _ db.LogMatchParams) (db.Match, error) {
			return db.Match{ID: 1}, nil
		},
	}
	h := NewHandler(mock, "http://localhost")
	r := chi.NewRouter()
	h.RegisterRoutes(r)

	t.Run("Valid match", func(t *testing.T) {
		form := url.Values{}
		form.Add("team1_id", "1")
		form.Add("team2_id", "2")
		form.Add("team1_score", "21")
		form.Add("team2_score", "15")
		req := httptest.NewRequest(http.MethodPost, "/league/pub-1/matches", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("got status %d; want %d", w.Code, http.StatusOK)
		}
	})

	t.Run("Same teams", func(t *testing.T) {
		form := url.Values{}
		form.Add("team1_id", "1")
		form.Add("team2_id", "1")
		form.Add("team1_score", "21")
		form.Add("team2_score", "15")
		req := httptest.NewRequest(http.MethodPost, "/league/pub-1/matches", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("got status %d; want %d", w.Code, http.StatusBadRequest)
		}
	})
}

func TestLeagueDashboard(t *testing.T) {
	mock := &MockQuerier{
		GetLeagueByPublicIDFunc: func(_ context.Context, publicID string) (db.League, error) {
			if publicID == "notfound" {
				return db.League{}, sql.ErrNoRows
			}
			return db.League{ID: 1, Name: "Test League", PublicID: publicID}, nil
		},
		GetTeamsByLeagueIDFunc: func(_ context.Context, _ int64) ([]db.Team, error) {
			return []db.Team{}, nil
		},
		GetMatchesByLeagueIDFunc: func(_ context.Context, _ int64) ([]db.Match, error) {
			return []db.Match{}, nil
		},
	}
	h := NewHandler(mock, "http://localhost")
	r := chi.NewRouter()
	h.RegisterRoutes(r)

	t.Run("Found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/league/pub-1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("got status %d; want %d", w.Code, http.StatusOK)
		}
	})

	t.Run("Not found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/league/notfound", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("got status %d; want %d", w.Code, http.StatusNotFound)
		}
	})
}
