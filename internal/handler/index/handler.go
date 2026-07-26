package index

import (
	"github.com/Marcel2603/spikeball-league/views"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	component := views.Index()
	err := component.Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
