package views

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/Marcel2603/spikeball-league/internal/db"
	"github.com/PuerkitoBio/goquery"
)

func TestLeagueDashboard(t *testing.T) {
	league := db.League{Name: "Test League", PublicID: "pub1"}
	teams := []db.Team{{ID: 1, Name: "Team A"}}
	matches := []db.Match{}
	standings := []TeamStanding{{Name: "Team A", GamesPlayed: 1, Wins: 1, Losses: 0, PointsScored: 21, PointsDiff: 6, WinRate: 100.0}}

	component := LeagueDashboard(league, teams, matches, standings)
	var buf bytes.Buffer
	err := component.Render(context.Background(), &buf)
	if err != nil {
		t.Fatalf("failed to render component: %v", err)
	}

	doc, err := goquery.NewDocumentFromReader(&buf)
	if err != nil {
		t.Fatalf("failed to parse html: %v", err)
	}

	if doc.Find("h2").First().Text() != "Test League" {
		t.Errorf("expected title not found, got %s", doc.Find("h2").First().Text())
	}
}

func TestLeagueDashboard_Empty(t *testing.T) {
	league := db.League{Name: "Empty League", PublicID: "pub2"}
	teams := []db.Team{}
	matches := []db.Match{}
	standings := []TeamStanding{}

	component := LeagueDashboard(league, teams, matches, standings)
	var buf bytes.Buffer
	err := component.Render(context.Background(), &buf)
	if err != nil {
		t.Fatalf("failed to render component: %v", err)
	}

	doc, err := goquery.NewDocumentFromReader(&buf)
	if err != nil {
		t.Fatalf("failed to parse html: %v", err)
	}

	if !strings.Contains(doc.Find("h2").First().Text(), "Empty League") {
		t.Errorf("expected title not found")
	}
}

func TestLeagueDashboard_FailingWriter(_ *testing.T) {
	league := db.League{}
	teams := []db.Team{{ID: 1}}
	matches := []db.Match{{ID: 1}}
	standings := []TeamStanding{{Name: "A"}}

	for i := 0; i < 15000; i++ {
		_ = LeagueDashboard(league, teams, matches, standings).Render(context.Background(), &failWriter{target: i})
		_ = LeagueDashboard(league, []db.Team{}, []db.Match{}, []TeamStanding{}).Render(context.Background(), &failWriter{target: i})
	}
}
