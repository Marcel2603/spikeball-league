-- name: CreateLeague :one
INSERT INTO leagues (public_id, admin_id, name)
VALUES (?, ?, ?)
RETURNING *;

-- name: GetLeagueByPublicID :one
SELECT * FROM leagues
WHERE public_id = ?;

-- name: GetLeagueByAdminID :one
SELECT * FROM leagues
WHERE admin_id = ? LIMIT 1;

-- name: CheckLeagueExists :one
SELECT 1 FROM leagues
WHERE admin_id = ? OR public_id = ?
LIMIT 1;

-- name: DeleteLeague :exec
DELETE FROM leagues
WHERE admin_id = ?;

-- name: AddPlayer :one
INSERT INTO players (league_id, name)
VALUES (?, ?)
RETURNING *;

-- name: GetPlayersByLeagueID :many
SELECT * FROM players
WHERE league_id = ?
ORDER BY created_at ASC;

-- name: RemovePlayer :exec
DELETE FROM players
WHERE id = ? AND league_id = ?;

-- name: AddTeam :one
INSERT INTO teams (league_id, name, player1_id, player2_id)
VALUES (?, ?, ?, ?)
RETURNING *;

-- name: GetTeamsByLeagueID :many
SELECT * FROM teams
WHERE league_id = ?
ORDER BY created_at ASC;

-- name: RemoveTeam :exec
DELETE FROM teams
WHERE id = ? AND league_id = ?;

-- name: LogMatch :one
INSERT INTO matches (
    league_id, team1_id, team2_id, team1_score, team2_score
) VALUES (
    ?, ?, ?, ?, ?
)
RETURNING *;

-- name: GetMatchesByLeagueID :many
SELECT * FROM matches
WHERE league_id = ?
ORDER BY created_at DESC;

-- name: DeleteMatch :exec
DELETE FROM matches
WHERE id = ? AND league_id = ?;
