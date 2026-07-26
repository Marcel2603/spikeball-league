package views

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/Marcel2603/spikeball-league/internal/db"
	"github.com/PuerkitoBio/goquery"
)

func TestAdminDashboard(t *testing.T) {
	league := db.League{Name: "Test League", AdminID: "admin1", PublicID: "pub1"}
	players := []db.Player{{ID: 1, Name: "Alice"}, {ID: 2, Name: "Bob"}}
	teams := []db.Team{{ID: 1, Name: "Team A", Player1ID: 1, Player2ID: 2}}
	var matches []db.Match

	component := AdminDashboard(league, players, teams, matches, "http://localhost")
	var buf bytes.Buffer
	err := component.Render(context.Background(), &buf)
	if err != nil {
		t.Fatalf("failed to render component: %v", err)
	}

	doc, err := goquery.NewDocumentFromReader(&buf)
	if err != nil {
		t.Fatalf("failed to parse html: %v", err)
	}

	if !strings.Contains(doc.Find("h2").First().Text(), "Test League") {
		t.Errorf("expected title not found, got %s", doc.Find("h2").First().Text())
	}
}

func TestAdminDashboard_Empty(t *testing.T) {
	league := db.League{Name: "Empty League", AdminID: "admin2", PublicID: "pub2"}
	var players []db.Player
	var teams []db.Team
	var matches []db.Match

	component := AdminDashboard(league, players, teams, matches, "http://localhost")
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

func TestAdminDashboard_Unknowns(t *testing.T) {
	league := db.League{Name: "League", AdminID: "admin2", PublicID: "pub2"}
	var players []db.Player
	teams := []db.Team{{ID: 1, Name: "Team A", Player1ID: 999, Player2ID: 999}}
	matches := []db.Match{{ID: 1, Team1ID: 999, Team2ID: 999}}

	component := AdminDashboard(league, players, teams, matches, "http://localhost")
	var buf bytes.Buffer
	err := component.Render(context.Background(), &buf)
	if err != nil {
		t.Fatalf("failed to render component: %v", err)
	}

	doc, err := goquery.NewDocumentFromReader(&buf)
	if err != nil {
		t.Fatalf("failed to parse html: %v", err)
	}

	if !strings.Contains(doc.Text(), "Unknown") {
		t.Errorf("expected Unknown player to be rendered")
	}
	if !strings.Contains(doc.Text(), "Unknown Team") {
		t.Errorf("expected Unknown Team to be rendered")
	}
}

func TestAdminDashboard_FailingWriter(_ *testing.T) {
	league := db.League{}
	players := []db.Player{{ID: 1}}
	teams := []db.Team{{ID: 1}}
	matches := []db.Match{{ID: 1}}

	for i := 0; i < 15000; i++ {
		_ = AdminDashboard(league, players, teams, matches, "").Render(context.Background(), &failWriter{target: i})
		_ = AdminDashboard(league, []db.Player{}, []db.Team{}, []db.Match{}, "").Render(context.Background(), &failWriter{target: i})
	}
}
