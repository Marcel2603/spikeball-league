package views

import (
	"bytes"
	"context"
	"testing"

	"github.com/Marcel2603/spikeball-league/internal/db"
	"github.com/PuerkitoBio/goquery"
)

func TestMatchForm(t *testing.T) {
	teams := []db.Team{{ID: 1, Name: "Team A"}, {ID: 2, Name: "Team B"}}
	component := MatchForm(teams, "pub1")
	var buf bytes.Buffer
	err := component.Render(context.Background(), &buf)
	if err != nil {
		t.Fatalf("failed to render component: %v", err)
	}

	doc, err := goquery.NewDocumentFromReader(&buf)
	if err != nil {
		t.Fatalf("failed to parse html: %v", err)
	}

	if doc.Find("form[hx-post='/league/pub1/matches']").Length() == 0 {
		t.Errorf("expected form not found")
	}
}

func TestMatchForm_FailingWriter(_ *testing.T) {
	teams := []db.Team{{ID: 1, Name: "Team A"}, {ID: 2, Name: "Team B"}}
	for i := 0; i < 15000; i++ {
		_ = MatchForm(teams, "pub1").Render(context.Background(), &failWriter{target: i})
	}
}
