package postgres

import (
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	dbErrors "github.com/Hirogava/avito-pr/internal/errors/db"
	"github.com/Hirogava/avito-pr/internal/models/reqres"
)

func TestCreateTeamSuccess(t *testing.T) {
	manager, mock, cleanup := newTestManager(t)
	defer cleanup()

	req := reqres.TeamAddRequest{
		TeamName: "backend",
		Members: []reqres.TeamMemberResponse{
			{UserID: "u1", Username: "alice", IsActive: true},
			{UserID: "u2", Username: "bob", IsActive: false},
		},
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT EXISTS (SELECT 1 FROM teams WHERE team_name = $1)`)).
		WithArgs(req.TeamName).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO teams (team_name, created_at) VALUES ($1, $2)`)).
		WithArgs(req.TeamName, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	insertUser := regexp.QuoteMeta(`
			INSERT INTO users (user_id, username, team_name, is_active, created_at)
			VALUES ($1, $2, $3, $4, NOW())
			ON CONFLICT (user_id) DO UPDATE
			SET username = EXCLUDED.username,
				team_name = EXCLUDED.team_name,
				is_active = EXCLUDED.is_active,
				updated_at = NOW();
		`)
	mock.ExpectExec(insertUser).
		WithArgs("u1", "alice", req.TeamName, req.Members[0].IsActive).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(insertUser).
		WithArgs("u2", "bob", req.TeamName, req.Members[1].IsActive).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	team, err := manager.CreateTeam(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if team.TeamName != req.TeamName || len(team.Members) != len(req.Members) {
		t.Fatalf("unexpected team response: %#v", team)
	}
}

func TestCreateTeamAlreadyExists(t *testing.T) {
	manager, mock, cleanup := newTestManager(t)
	defer cleanup()

	req := reqres.TeamAddRequest{TeamName: "backend"}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT EXISTS (SELECT 1 FROM teams WHERE team_name = $1)`)).
		WithArgs(req.TeamName).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))
	mock.ExpectRollback()

	_, err := manager.CreateTeam(req)
	if err != dbErrors.ErrorTeamAlreadyExists {
		t.Fatalf("expected ErrorTeamAlreadyExists, got %v", err)
	}
}

func TestGetTeamSuccess(t *testing.T) {
	manager, mock, cleanup := newTestManager(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"user_id", "username", "is_active"}).
		AddRow("u1", "alice", true).
		AddRow("u2", "bob", false)

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT user_id, username, is_active
		FROM users
		WHERE team_name = $1
		ORDER BY username;
	`)).WithArgs("backend").
		WillReturnRows(rows)

	team, err := manager.GetTeam("backend")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(team.Members) != 2 {
		t.Fatalf("expected 2 members, got %d", len(team.Members))
	}
}

func TestGetTeamNotFound(t *testing.T) {
	manager, mock, cleanup := newTestManager(t)
	defer cleanup()

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT user_id, username, is_active
		FROM users
		WHERE team_name = $1
		ORDER BY username;
	`)).WithArgs("missing").
		WillReturnError(sql.ErrNoRows)

	_, err := manager.GetTeam("missing")
	if err != dbErrors.ErrorTeamNotFound {
		t.Fatalf("expected ErrorTeamNotFound, got %v", err)
	}
}
