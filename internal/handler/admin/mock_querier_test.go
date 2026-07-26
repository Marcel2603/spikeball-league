package admin

import (
	"context"
	"github.com/Marcel2603/spikeball-league/internal/db"
)

type MockQuerier struct {
	AddPlayerFunc            func(ctx context.Context, arg db.AddPlayerParams) (db.Player, error)
	AddTeamFunc              func(ctx context.Context, arg db.AddTeamParams) (db.Team, error)
	CheckLeagueExistsFunc    func(ctx context.Context, arg db.CheckLeagueExistsParams) (int64, error)
	CreateLeagueFunc         func(ctx context.Context, arg db.CreateLeagueParams) (db.League, error)
	DeleteLeagueFunc         func(ctx context.Context, adminID string) error
	DeleteMatchFunc          func(ctx context.Context, arg db.DeleteMatchParams) error
	GetLeagueByAdminIDFunc   func(ctx context.Context, adminID string) (db.League, error)
	GetLeagueByPublicIDFunc  func(ctx context.Context, publicID string) (db.League, error)
	GetMatchesByLeagueIDFunc func(ctx context.Context, leagueID int64) ([]db.Match, error)
	GetPlayersByLeagueIDFunc func(ctx context.Context, leagueID int64) ([]db.Player, error)
	GetTeamsByLeagueIDFunc   func(ctx context.Context, leagueID int64) ([]db.Team, error)
	LogMatchFunc             func(ctx context.Context, arg db.LogMatchParams) (db.Match, error)
	RemovePlayerFunc         func(ctx context.Context, arg db.RemovePlayerParams) error
	RemoveTeamFunc           func(ctx context.Context, arg db.RemoveTeamParams) error
}

func (m *MockQuerier) AddPlayer(ctx context.Context, arg db.AddPlayerParams) (db.Player, error) {
	return m.AddPlayerFunc(ctx, arg)
}
func (m *MockQuerier) AddTeam(ctx context.Context, arg db.AddTeamParams) (db.Team, error) {
	return m.AddTeamFunc(ctx, arg)
}
func (m *MockQuerier) CheckLeagueExists(ctx context.Context, arg db.CheckLeagueExistsParams) (int64, error) {
	return m.CheckLeagueExistsFunc(ctx, arg)
}
func (m *MockQuerier) CreateLeague(ctx context.Context, arg db.CreateLeagueParams) (db.League, error) {
	return m.CreateLeagueFunc(ctx, arg)
}
func (m *MockQuerier) DeleteLeague(ctx context.Context, adminID string) error {
	return m.DeleteLeagueFunc(ctx, adminID)
}
func (m *MockQuerier) DeleteMatch(ctx context.Context, arg db.DeleteMatchParams) error {
	return m.DeleteMatchFunc(ctx, arg)
}
func (m *MockQuerier) GetLeagueByAdminID(ctx context.Context, adminID string) (db.League, error) {
	return m.GetLeagueByAdminIDFunc(ctx, adminID)
}
func (m *MockQuerier) GetLeagueByPublicID(ctx context.Context, publicID string) (db.League, error) {
	return m.GetLeagueByPublicIDFunc(ctx, publicID)
}
func (m *MockQuerier) GetMatchesByLeagueID(ctx context.Context, leagueID int64) ([]db.Match, error) {
	return m.GetMatchesByLeagueIDFunc(ctx, leagueID)
}
func (m *MockQuerier) GetPlayersByLeagueID(ctx context.Context, leagueID int64) ([]db.Player, error) {
	return m.GetPlayersByLeagueIDFunc(ctx, leagueID)
}
func (m *MockQuerier) GetTeamsByLeagueID(ctx context.Context, leagueID int64) ([]db.Team, error) {
	return m.GetTeamsByLeagueIDFunc(ctx, leagueID)
}
func (m *MockQuerier) LogMatch(ctx context.Context, arg db.LogMatchParams) (db.Match, error) {
	return m.LogMatchFunc(ctx, arg)
}
func (m *MockQuerier) RemovePlayer(ctx context.Context, arg db.RemovePlayerParams) error {
	return m.RemovePlayerFunc(ctx, arg)
}
func (m *MockQuerier) RemoveTeam(ctx context.Context, arg db.RemoveTeamParams) error {
	return m.RemoveTeamFunc(ctx, arg)
}
