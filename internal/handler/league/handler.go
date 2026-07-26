package league

import (
	"log/slog"
	"net/http"
	"sort"
	"strconv"

	"github.com/Marcel2603/spikeball-league/internal/db"
	"github.com/Marcel2603/spikeball-league/views"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Handler struct {
	q    db.Querier
	host string
}

func NewHandler(q db.Querier, host string) *Handler {
	return &Handler{q: q, host: host}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Get("/", h.Index)
	r.Post("/leagues", h.CreateLeague)
	r.Get("/leagues/{id}/exists", h.CheckLeagueExists)

	r.Get("/league/{publicID}", h.LeagueDashboard)
	r.Post("/league/{publicID}/matches", h.LogMatch)
}

func (*Handler) Index(w http.ResponseWriter, r *http.Request) {
	component := views.Index()
	component.Render(r.Context(), w)
}

func (h *Handler) CreateLeague(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	if name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	publicID := uuid.New().String()
	adminID := uuid.New().String()

	_, err = h.q.CreateLeague(r.Context(), db.CreateLeagueParams{
		PublicID: publicID,
		AdminID:  adminID,
		Name:     name,
	})
	if err != nil {
		slog.Error("Failed to create league", "err", err)
		http.Error(w, "Failed to create league", http.StatusInternalServerError)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/admin/"+adminID)
		w.WriteHeader(http.StatusOK)
	} else {
		http.Redirect(w, r, "/admin/"+adminID, http.StatusFound)
	}
}

func (h *Handler) CheckLeagueExists(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	_, err := h.q.CheckLeagueExists(r.Context(), db.CheckLeagueExistsParams{
		AdminID:  id,
		PublicID: id,
	})
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) LeagueDashboard(w http.ResponseWriter, r *http.Request) {
	publicID := chi.URLParam(r, "publicID")
	league, err := h.q.GetLeagueByPublicID(r.Context(), publicID)
	if err != nil {
		http.Error(w, "League not found", http.StatusNotFound)
		return
	}

	teams, err := h.q.GetTeamsByLeagueID(r.Context(), league.ID)
	if err != nil {
		teams = []db.Team{}
	}

	matches, err := h.q.GetMatchesByLeagueID(r.Context(), league.ID)
	if err != nil {
		matches = []db.Match{}
	}

	standings := CalculateStandings(teams, matches)
	component := views.LeagueDashboard(league, teams, matches, standings)
	component.Render(r.Context(), w)
}

func (h *Handler) LogMatch(w http.ResponseWriter, r *http.Request) {
	publicID := chi.URLParam(r, "publicID")
	league, err := h.q.GetLeagueByPublicID(r.Context(), publicID)
	if err != nil {
		http.Error(w, "League not found", http.StatusNotFound)
		return
	}

	t1, _ := strconv.ParseInt(r.FormValue("team1_id"), 10, 64)
	t2, _ := strconv.ParseInt(r.FormValue("team2_id"), 10, 64)
	t1s, _ := strconv.ParseInt(r.FormValue("team1_score"), 10, 64)
	t2s, _ := strconv.ParseInt(r.FormValue("team2_score"), 10, 64)

	// Basic validation
	if t1 == 0 || t2 == 0 {
		http.Error(w, "All teams must be selected", http.StatusBadRequest)
		return
	}
	if t1 == t2 {
		http.Error(w, "Teams must be distinct", http.StatusBadRequest)
		return
	}

	_, err = h.q.LogMatch(r.Context(), db.LogMatchParams{
		LeagueID:   league.ID,
		Team1ID:    t1,
		Team2ID:    t2,
		Team1Score: t1s,
		Team2Score: t2s,
	})
	if err != nil {
		slog.Error("Failed to log match", "err", err)
	}

	w.Header().Set("HX-Redirect", "/league/"+publicID)
	w.WriteHeader(http.StatusOK)
}

func CalculateStandings(teams []db.Team, matches []db.Match) []views.TeamStanding {
	stats := make(map[int64]*views.TeamStanding)
	for _, t := range teams {
		stats[t.ID] = &views.TeamStanding{
			ID:   t.ID,
			Name: t.Name,
		}
	}

	for _, m := range matches {
		t1Win := m.Team1Score > m.Team2Score
		t2Win := m.Team2Score > m.Team1Score

		t1Diff := int(m.Team1Score - m.Team2Score)
		t2Diff := int(m.Team2Score - m.Team1Score)

		updateStats(stats, m.Team1ID, t1Win, t1Diff)
		updateStats(stats, m.Team2ID, t2Win, t2Diff)
	}

	result := make([]views.TeamStanding, 0, len(stats))
	for _, s := range stats {
		if s.GamesPlayed > 0 {
			s.WinRate = float64(s.Wins) / float64(s.GamesPlayed) * 100
		}
		result = append(result, *s)
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].Wins != result[j].Wins {
			return result[i].Wins > result[j].Wins
		}
		return result[i].PointsDiff > result[j].PointsDiff
	})

	return result
}

func updateStats(stats map[int64]*views.TeamStanding, tID int64, won bool, diff int) {
	s, ok := stats[tID]
	if !ok {
		return
	}
	s.GamesPlayed++
	s.PointsDiff += diff
	if won {
		s.Wins++
	} else {
		s.Losses++
	}
}
