package admin

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Marcel2603/spikeball-league/internal/db"
	"github.com/Marcel2603/spikeball-league/views"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	q    db.Querier
	host string
}

func NewHandler(q db.Querier, host string) *Handler {
	return &Handler{q: q, host: host}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Get("/admin/{adminID}", h.AdminDashboard)
	r.Delete("/admin/{adminID}", h.DeleteLeague)
	r.Post("/admin/{adminID}/players", h.AddPlayer)
	r.Delete("/admin/{adminID}/players/{playerID}", h.RemovePlayer)
	r.Post("/admin/{adminID}/teams", h.AddTeam)
	r.Delete("/admin/{adminID}/teams/{teamID}", h.RemoveTeam)
	r.Delete("/admin/{adminID}/matches/{matchID}", h.RemoveMatch)
}

func (h *Handler) AdminDashboard(w http.ResponseWriter, r *http.Request) {
	adminID := chi.URLParam(r, "adminID")
	league, err := h.q.GetLeagueByAdminID(r.Context(), adminID)
	if err != nil {
		http.Error(w, "League not found", http.StatusNotFound)
		return
	}

	players, err := h.q.GetPlayersByLeagueID(r.Context(), league.ID)
	if err != nil {
		players = []db.Player{}
	}

	teams, err := h.q.GetTeamsByLeagueID(r.Context(), league.ID)
	if err != nil {
		teams = []db.Team{}
	}

	matches, err := h.q.GetMatchesByLeagueID(r.Context(), league.ID)
	if err != nil {
		matches = []db.Match{}
	}

	component := views.AdminDashboard(league, players, teams, matches, h.host)
	component.Render(r.Context(), w)
}

func (h *Handler) DeleteLeague(w http.ResponseWriter, r *http.Request) {
	adminID := chi.URLParam(r, "adminID")
	err := h.q.DeleteLeague(r.Context(), adminID)
	if err != nil {
		slog.Error("Failed to delete league", "err", err)
		http.Error(w, "Failed to delete league", http.StatusInternalServerError)
		return
	}

	// HTMX redirect to home
	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) AddPlayer(w http.ResponseWriter, r *http.Request) {
	adminID := chi.URLParam(r, "adminID")
	league, err := h.q.GetLeagueByAdminID(r.Context(), adminID)
	if err != nil {
		http.Error(w, "League not found", http.StatusNotFound)
		return
	}

	name := r.FormValue("player_name")
	if name != "" {
		_, err = h.q.AddPlayer(r.Context(), db.AddPlayerParams{
			LeagueID: league.ID,
			Name:     name,
		})
		if err != nil {
			slog.Error("Failed to add player", "err", err)
		}
	}

	w.Header().Set("HX-Redirect", "/admin/"+adminID)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) RemovePlayer(w http.ResponseWriter, r *http.Request) {
	adminID := chi.URLParam(r, "adminID")
	playerIDStr := chi.URLParam(r, "playerID")
	playerID, err := strconv.ParseInt(playerIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid player ID", http.StatusBadRequest)
		return
	}

	league, err := h.q.GetLeagueByAdminID(r.Context(), adminID)
	if err != nil {
		http.Error(w, "League not found", http.StatusNotFound)
		return
	}

	h.q.RemovePlayer(r.Context(), db.RemovePlayerParams{
		ID:       playerID,
		LeagueID: league.ID,
	})

	w.Header().Set("HX-Redirect", "/admin/"+adminID)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) AddTeam(w http.ResponseWriter, r *http.Request) {
	adminID := chi.URLParam(r, "adminID")
	league, err := h.q.GetLeagueByAdminID(r.Context(), adminID)
	if err != nil {
		http.Error(w, "League not found", http.StatusNotFound)
		return
	}

	name := r.FormValue("team_name")
	p1, _ := strconv.ParseInt(r.FormValue("player1_id"), 10, 64)
	p2, _ := strconv.ParseInt(r.FormValue("player2_id"), 10, 64)

	if name != "" && p1 != 0 && p2 != 0 && p1 != p2 {
		_, err = h.q.AddTeam(r.Context(), db.AddTeamParams{
			LeagueID:  league.ID,
			Name:      name,
			Player1ID: p1,
			Player2ID: p2,
		})
		if err != nil {
			slog.Error("Failed to add team", "err", err)
		}
	}

	w.Header().Set("HX-Redirect", "/admin/"+adminID)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) RemoveTeam(w http.ResponseWriter, r *http.Request) {
	adminID := chi.URLParam(r, "adminID")
	teamIDStr := chi.URLParam(r, "teamID")
	teamID, err := strconv.ParseInt(teamIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid team ID", http.StatusBadRequest)
		return
	}

	league, err := h.q.GetLeagueByAdminID(r.Context(), adminID)
	if err != nil {
		http.Error(w, "League not found", http.StatusNotFound)
		return
	}

	h.q.RemoveTeam(r.Context(), db.RemoveTeamParams{
		ID:       teamID,
		LeagueID: league.ID,
	})

	w.Header().Set("HX-Redirect", "/admin/"+adminID)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) RemoveMatch(w http.ResponseWriter, r *http.Request) {
	adminID := chi.URLParam(r, "adminID")
	matchIDStr := chi.URLParam(r, "matchID")
	matchID, err := strconv.ParseInt(matchIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid match ID", http.StatusBadRequest)
		return
	}

	league, err := h.q.GetLeagueByAdminID(r.Context(), adminID)
	if err != nil {
		http.Error(w, "League not found", http.StatusNotFound)
		return
	}

	h.q.DeleteMatch(r.Context(), db.DeleteMatchParams{
		ID:       matchID,
		LeagueID: league.ID,
	})

	w.Header().Set("HX-Redirect", "/admin/"+adminID)
	w.WriteHeader(http.StatusOK)
}
