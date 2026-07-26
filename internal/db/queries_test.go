package db

import (
	"context"
	"testing"
)

func setupTestDB(t *testing.T) (*Queries, func()) {
	t.Helper()
	db, err := InitDB(":memory:")
	if err != nil {
		t.Fatalf("failed to init in-memory db: %v", err)
	}

	q := New(db)
	return q, func() {
		db.Close()
	}
}

func TestLeagueOperations(t *testing.T) {
	q, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Create a league
	league, err := q.CreateLeague(ctx, CreateLeagueParams{
		PublicID: "pub-1",
		AdminID:  "adm-1",
		Name:     "Test League",
	})
	if err != nil {
		t.Fatalf("failed to create league: %v", err)
	}
	if league.Name != "Test League" {
		t.Errorf("expected Test League, got %s", league.Name)
	}

	// Get league by Admin ID
	l1, err := q.GetLeagueByAdminID(ctx, "adm-1")
	if err != nil {
		t.Fatalf("failed to get league by admin id: %v", err)
	}
	if l1.Name != "Test League" {
		t.Errorf("expected Test League, got %s", l1.Name)
	}

	// Get league by Public ID
	l2, err := q.GetLeagueByPublicID(ctx, "pub-1")
	if err != nil {
		t.Fatalf("failed to get league by public id: %v", err)
	}
	if l2.Name != "Test League" {
		t.Errorf("expected Test League, got %s", l2.Name)
	}

	// Check exists
	exists, err := q.CheckLeagueExists(ctx, CheckLeagueExistsParams{
		AdminID:  "adm-1",
		PublicID: "pub-1",
	})
	if err != nil || exists != 1 {
		t.Errorf("expected league to exist, got error: %v, value: %v", err, exists)
	}

	// Delete league
	err = q.DeleteLeague(ctx, "adm-1")
	if err != nil {
		t.Fatalf("failed to delete league: %v", err)
	}

	// Check exists after delete
	_, err = q.CheckLeagueExists(ctx, CheckLeagueExistsParams{
		AdminID:  "adm-1",
		PublicID: "pub-1",
	})
	if err == nil {
		t.Error("expected error getting deleted league, got nil")
	}
}

func TestPlayerOperations(t *testing.T) {
	q, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()
	league, _ := q.CreateLeague(ctx, CreateLeagueParams{
		PublicID: "pub-1",
		AdminID:  "adm-1",
		Name:     "Test League",
	})

	// Add players
	p1, err := q.AddPlayer(ctx, AddPlayerParams{
		LeagueID: league.ID,
		Name:     "Alice",
	})
	if err != nil {
		t.Fatalf("failed to add player: %v", err)
	}

	_, err = q.AddPlayer(ctx, AddPlayerParams{
		LeagueID: league.ID,
		Name:     "Bob",
	})
	if err != nil {
		t.Fatalf("failed to add player: %v", err)
	}

	// Get players
	players, err := q.GetPlayersByLeagueID(ctx, league.ID)
	if err != nil {
		t.Fatalf("failed to get players: %v", err)
	}
	if len(players) != 2 {
		t.Errorf("expected 2 players, got %d", len(players))
	}

	// Remove player
	err = q.RemovePlayer(ctx, RemovePlayerParams{
		ID:       p1.ID,
		LeagueID: league.ID,
	})
	if err != nil {
		t.Fatalf("failed to remove player: %v", err)
	}

	players, _ = q.GetPlayersByLeagueID(ctx, league.ID)
	if len(players) != 1 || players[0].Name != "Bob" {
		t.Errorf("expected 1 player (Bob), got %d players", len(players))
	}
}

func TestTeamAndMatchOperations(t *testing.T) {
	q, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()
	league, _ := q.CreateLeague(ctx, CreateLeagueParams{
		PublicID: "pub-1",
		AdminID:  "adm-1",
		Name:     "Test League",
	})

	// Add players
	p1, _ := q.AddPlayer(ctx, AddPlayerParams{LeagueID: league.ID, Name: "A"})
	p2, _ := q.AddPlayer(ctx, AddPlayerParams{LeagueID: league.ID, Name: "B"})
	p3, _ := q.AddPlayer(ctx, AddPlayerParams{LeagueID: league.ID, Name: "C"})
	p4, _ := q.AddPlayer(ctx, AddPlayerParams{LeagueID: league.ID, Name: "D"})

	// Add Teams
	t1, err := q.AddTeam(ctx, AddTeamParams{
		LeagueID:  league.ID,
		Name:      "Team 1",
		Player1ID: p1.ID,
		Player2ID: p2.ID,
	})
	if err != nil {
		t.Fatalf("failed to add team: %v", err)
	}

	t2, err := q.AddTeam(ctx, AddTeamParams{
		LeagueID:  league.ID,
		Name:      "Team 2",
		Player1ID: p3.ID,
		Player2ID: p4.ID,
	})
	if err != nil {
		t.Fatalf("failed to add team: %v", err)
	}

	teams, err := q.GetTeamsByLeagueID(ctx, league.ID)
	if err != nil || len(teams) != 2 {
		t.Fatalf("failed to get teams, expected 2 got %d (err: %v)", len(teams), err)
	}

	// Log match
	m1, err := q.LogMatch(ctx, LogMatchParams{
		LeagueID:   league.ID,
		Team1ID:    t1.ID,
		Team2ID:    t2.ID,
		Team1Score: 21,
		Team2Score: 19,
	})
	if err != nil {
		t.Fatalf("failed to log match: %v", err)
	}

	matches, err := q.GetMatchesByLeagueID(ctx, league.ID)
	if err != nil || len(matches) != 1 {
		t.Fatalf("failed to get matches: %v", err)
	}

	// Delete match
	err = q.DeleteMatch(ctx, DeleteMatchParams{
		ID:       m1.ID,
		LeagueID: league.ID,
	})
	if err != nil {
		t.Fatalf("failed to delete match: %v", err)
	}

	matches, _ = q.GetMatchesByLeagueID(ctx, league.ID)
	if len(matches) != 0 {
		t.Errorf("expected 0 matches, got %d", len(matches))
	}

	// Remove team
	err = q.RemoveTeam(ctx, RemoveTeamParams{
		ID:       t1.ID,
		LeagueID: league.ID,
	})
	if err != nil {
		t.Fatalf("failed to remove team: %v", err)
	}
	teams, _ = q.GetTeamsByLeagueID(ctx, league.ID)
	if len(teams) != 1 {
		t.Errorf("expected 1 team, got %d", len(teams))
	}
}
