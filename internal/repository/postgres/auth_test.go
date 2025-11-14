package postgres

import (
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestSaveRefreshToken(t *testing.T) {
	manager, mock, cleanup := newTestManager(t)
	defer cleanup()

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO sessions (user_id, token, expires_at) VALUES ($1, $2, $3)`)).
		WithArgs("user", "token", sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	if err := manager.SaveRefreshToken("token", "user"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetRefreshTokenSuccess(t *testing.T) {
	manager, mock, cleanup := newTestManager(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"id", "token", "expires_at"}).
		AddRow("1", "refresh", time.Now())

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, token, expires_at FROM sessions WHERE user_id = $1`)).
		WithArgs("user").
		WillReturnRows(rows)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT is_admin FROM users WHERE user_id = $1`)).
		WithArgs("user").
		WillReturnRows(sqlmock.NewRows([]string{"is_admin"}).AddRow(true))

	token, role, err := manager.GetRefreshToken("user")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if role != "admin" || token.Token != "refresh" {
		t.Fatalf("unexpected data: %+v role=%s", token, role)
	}
}

func TestGetRefreshTokenRoleError(t *testing.T) {
	manager, mock, cleanup := newTestManager(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"id", "token", "expires_at"}).
		AddRow("1", "refresh", time.Now())

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, token, expires_at FROM sessions WHERE user_id = $1`)).
		WithArgs("user").
		WillReturnRows(rows)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT is_admin FROM users WHERE user_id = $1`)).
		WithArgs("user").
		WillReturnError(sql.ErrNoRows)

	_, _, err := manager.GetRefreshToken("user")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestDeleteRefreshToken(t *testing.T) {
	manager, mock, cleanup := newTestManager(t)
	defer cleanup()

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM sessions WHERE user_id = $1 and token = $2`)).
		WithArgs("user", "token").
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err := manager.DeleteRefreshToken("user", "token"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUpdateRefreshToken(t *testing.T) {
	manager, mock, cleanup := newTestManager(t)
	defer cleanup()

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE sessions SET expires_at = $1, token = $2 WHERE user_id = $3`)).
		WithArgs(sqlmock.AnyArg(), "token", "user").
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err := manager.UpdateRefreshToken("user", "token"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSetUserStatus(t *testing.T) {
	manager, mock, cleanup := newTestManager(t)
	defer cleanup()

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE users SET is_admin = $1 WHERE user_id = $2`)).
		WithArgs(true, "user").
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err := manager.SetUserStatus("user", true); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
