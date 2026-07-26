package admin

import (
	"bytes"
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

func TestAdminDashboard(t *testing.T) {
	mock := &MockQuerier{
		GetLeagueByAdminIDFunc: func(_ context.Context, adminID string) (db.League, error) {
			if adminID == "notfound" {
				return db.League{}, sql.ErrNoRows
			}
			return db.League{ID: 1, Name: "Test League", AdminID: adminID}, nil
		},
		GetPlayersByLeagueIDFunc: func(_ context.Context, _ int64) ([]db.Player, error) {
			return []db.Player{{ID: 1, Name: "Alice"}}, nil
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
		req := httptest.NewRequest(http.MethodGet, "/admin/adm-1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("got status %d; want %d", w.Code, http.StatusOK)
		}
		if !bytes.Contains(w.Body.Bytes(), []byte("Test League")) {
			t.Errorf("expected body to contain 'Test League'")
		}
	})

	t.Run("Not found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/admin/notfound", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("got status %d; want %d", w.Code, http.StatusNotFound)
		}
	})
}

func TestDeleteLeague(t *testing.T) {
	mock := &MockQuerier{
		DeleteLeagueFunc: func(_ context.Context, adminID string) error {
			if adminID == "fail" {
				return fmt.Errorf("db error")
			}
			return nil
		},
	}
	h := NewHandler(mock, "http://localhost")
	r := chi.NewRouter()
	h.RegisterRoutes(r)

	t.Run("Success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/admin/adm-1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("got status %d; want %d", w.Code, http.StatusOK)
		}
		if w.Header().Get("HX-Redirect") != "/" {
			t.Errorf("got redirect %q; want /", w.Header().Get("HX-Redirect"))
		}
	})

	t.Run("Failure", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/admin/fail", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("got status %d; want %d", w.Code, http.StatusInternalServerError)
		}
	})
}

func TestAddPlayer(t *testing.T) {
	mock := &MockQuerier{
		GetLeagueByAdminIDFunc: func(_ context.Context, adminID string) (db.League, error) {
			if adminID == "notfound" {
				return db.League{}, sql.ErrNoRows
			}
			return db.League{ID: 1, AdminID: adminID}, nil
		},
		AddPlayerFunc: func(_ context.Context, arg db.AddPlayerParams) (db.Player, error) {
			if arg.Name == "FailPlayer" {
				return db.Player{}, fmt.Errorf("db error")
			}
			return db.Player{ID: 1, Name: arg.Name}, nil
		},
	}
	h := NewHandler(mock, "http://localhost")
	r := chi.NewRouter()
	h.RegisterRoutes(r)

	t.Run("Valid player", func(t *testing.T) {
		form := url.Values{}
		form.Add("player_name", "Alice")
		req := httptest.NewRequest(http.MethodPost, "/admin/adm-1/players", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("got status %d; want %d", w.Code, http.StatusOK)
		}
	})

	t.Run("League not found", func(t *testing.T) {
		form := url.Values{}
		form.Add("player_name", "Alice")
		req := httptest.NewRequest(http.MethodPost, "/admin/notfound/players", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("got status %d; want %d", w.Code, http.StatusNotFound)
		}
	})
}

func TestRemovePlayer(t *testing.T) {
	mock := &MockQuerier{
		GetLeagueByAdminIDFunc: func(_ context.Context, adminID string) (db.League, error) {
			if adminID == "notfound" {
				return db.League{}, sql.ErrNoRows
			}
			return db.League{ID: 1, AdminID: adminID}, nil
		},
		RemovePlayerFunc: func(_ context.Context, _ db.RemovePlayerParams) error {
			return nil
		},
	}
	h := NewHandler(mock, "http://localhost")
	r := chi.NewRouter()
	h.RegisterRoutes(r)

	t.Run("Valid remove", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/admin/adm-1/players/1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("got status %d; want %d", w.Code, http.StatusOK)
		}
	})

	t.Run("League not found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/admin/notfound/players/1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("got status %d; want %d", w.Code, http.StatusNotFound)
		}
	})
}

func TestAddTeam(t *testing.T) {
	mock := &MockQuerier{
		GetLeagueByAdminIDFunc: func(_ context.Context, adminID string) (db.League, error) {
			if adminID == "notfound" {
				return db.League{}, sql.ErrNoRows
			}
			return db.League{ID: 1, AdminID: adminID}, nil
		},
		AddTeamFunc: func(_ context.Context, arg db.AddTeamParams) (db.Team, error) {
			if arg.Name == "FailTeam" {
				return db.Team{}, fmt.Errorf("db error")
			}
			return db.Team{ID: 1, Name: arg.Name}, nil
		},
	}
	h := NewHandler(mock, "http://localhost")
	r := chi.NewRouter()
	h.RegisterRoutes(r)

	t.Run("Valid team", func(t *testing.T) {
		form := url.Values{}
		form.Add("team_name", "Team Alpha")
		form.Add("player1_id", "1")
		form.Add("player2_id", "2")
		req := httptest.NewRequest(http.MethodPost, "/admin/adm-1/teams", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("got status %d; want %d", w.Code, http.StatusOK)
		}
	})
}

func TestRemoveTeam(t *testing.T) {
	mock := &MockQuerier{
		GetLeagueByAdminIDFunc: func(_ context.Context, adminID string) (db.League, error) {
			if adminID == "notfound" {
				return db.League{}, sql.ErrNoRows
			}
			return db.League{ID: 1, AdminID: adminID}, nil
		},
		RemoveTeamFunc: func(_ context.Context, _ db.RemoveTeamParams) error {
			return nil
		},
	}
	h := NewHandler(mock, "http://localhost")
	r := chi.NewRouter()
	h.RegisterRoutes(r)

	t.Run("Valid remove", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/admin/adm-1/teams/1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("got status %d; want %d", w.Code, http.StatusOK)
		}
	})
}

func TestRemoveMatch(t *testing.T) {
	mock := &MockQuerier{
		GetLeagueByAdminIDFunc: func(_ context.Context, adminID string) (db.League, error) {
			if adminID == "notfound" {
				return db.League{}, sql.ErrNoRows
			}
			return db.League{ID: 1, AdminID: adminID}, nil
		},
		DeleteMatchFunc: func(_ context.Context, _ db.DeleteMatchParams) error {
			return nil
		},
	}
	h := NewHandler(mock, "http://localhost")
	r := chi.NewRouter()
	h.RegisterRoutes(r)

	t.Run("Valid remove", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/admin/adm-1/matches/1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("got status %d; want %d", w.Code, http.StatusOK)
		}
	})
}
