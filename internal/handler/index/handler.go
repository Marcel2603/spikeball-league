package index

import (
	"net/http"

	custommw "github.com/Marcel2603/spikeball-league/internal/middleware"
	"github.com/Marcel2603/spikeball-league/views"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	theme := custommw.ThemeFromContext(r.Context())
	component := views.Index(theme)
	err := component.Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
