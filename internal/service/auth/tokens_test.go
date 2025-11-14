package auth

import (
	"database/sql"
	"io"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/sirupsen/logrus"

	"github.com/Hirogava/avito-pr/internal/config/logger"
	authModels "github.com/Hirogava/avito-pr/internal/models/auth"
	"github.com/Hirogava/avito-pr/internal/repository/postgres"
	"github.com/golang-jwt/jwt/v5"
)

const testSecret = "super-secret"

func TestMain(m *testing.M) {
	if err := os.Setenv("JWT_SECRET", testSecret); err != nil {
		panic(err)
	}
	secret = testSecret
	logger.Logger = logrus.New()
	logger.Logger.SetOutput(io.Discard)
	code := m.Run()
	os.Exit(code)
}

func newMockManager(t *testing.T) (*postgres.Manager, sqlmock.Sqlmock, func()) {
	t.Helper()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	manager := &postgres.Manager{Conn: db}

	return manager, mock, func() {
		_ = db.Close()
	}
}

func TestGenerateRefreshTokenSuccess(t *testing.T) {
	manager, mock, cleanup := newMockManager(t)
	defer cleanup()

	userID := "user-id"
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO sessions (user_id, token, expires_at) VALUES ($1, $2, $3)`)).
		WithArgs(userID, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	token, err := GenerateRefreshToken(manager, userID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token == "" {
		t.Fatal("expected non-empty token")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestGenerateRefreshTokenSaveError(t *testing.T) {
	manager, mock, cleanup := newMockManager(t)
	defer cleanup()

	userID := "user-id"
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO sessions (user_id, token, expires_at) VALUES ($1, $2, $3)`)).
		WithArgs(userID, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(sql.ErrConnDone)

	token, err := GenerateRefreshToken(manager, userID)
	if err == nil {
		t.Fatal("expected error")
	}
	if token != "" {
		t.Fatalf("expected empty token, got %s", token)
	}
}

func TestGenerateAccessToken(t *testing.T) {
	claims := jwt.MapClaims{
		"id":   "user",
		"role": "admin",
		"exp":  time.Now().Add(time.Minute).Unix(),
	}

	token, err := GenerateAccessToken(claims)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token == "" {
		t.Fatal("expected non-empty token")
	}

	parsed, err := ParseToken(token)
	if err != nil {
		t.Fatalf("failed to parse generated token: %v", err)
	}
	if !parsed.Valid {
		t.Fatal("expected token to be valid")
	}
}

func TestParseTokenInvalid(t *testing.T) {
	if _, err := ParseToken(""); err == nil {
		t.Fatal("expected error for empty token")
	}
}

func TestAddAccessTime(t *testing.T) {
	exp := AddAccessTime()
	target := time.Now().Add(15 * time.Minute)
	actual := time.Unix(exp, 0)

	diff := actual.Sub(target)
	if diff < -time.Minute || diff > time.Minute {
		t.Fatalf("expected exp to be ~15m ahead, diff=%v", diff)
	}
}

func expectSessionQuery(mock sqlmock.Sqlmock, userID string, token authModels.RefreshToken, isAdmin bool) {
	rows := sqlmock.NewRows([]string{"id", "token", "expires_at"}).
		AddRow(token.ID, token.Token, token.ExpiredAt)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, token, expires_at FROM sessions WHERE user_id = $1`)).
		WithArgs(userID).
		WillReturnRows(rows)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT is_admin FROM users WHERE user_id = $1`)).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"is_admin"}).AddRow(isAdmin))
}

func TestValidateRefreshTokenSuccess(t *testing.T) {
	manager, mock, cleanup := newMockManager(t)
	defer cleanup()

	userID := "user"
	token := authModels.RefreshToken{
		ID:        "1",
		Token:     "refresh-token",
		ExpiredAt: time.Now().Add(time.Hour),
	}

	expectSessionQuery(mock, userID, token, true)

	rt, role, err := ValidateRefreshToken(manager, userID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if role != "admin" {
		t.Fatalf("expected role admin, got %s", role)
	}
	if rt.Token != token.Token {
		t.Fatalf("unexpected token %s", rt.Token)
	}
}

func TestValidateRefreshTokenMissing(t *testing.T) {
	manager, mock, cleanup := newMockManager(t)
	defer cleanup()

	userID := "missing"
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, token, expires_at FROM sessions WHERE user_id = $1`)).
		WithArgs(userID).
		WillReturnError(sql.ErrNoRows)

	_, _, err := ValidateRefreshToken(manager, userID)
	if err != jwt.ErrTokenExpired {
		t.Fatalf("expected jwt.ErrTokenExpired, got %v", err)
	}
}

func TestValidateRefreshTokenExpired(t *testing.T) {
	manager, mock, cleanup := newMockManager(t)
	defer cleanup()

	userID := "user"
	token := authModels.RefreshToken{
		ID:        "2",
		Token:     "old-token",
		ExpiredAt: time.Now().Add(-time.Hour),
	}

	expectSessionQuery(mock, userID, token, false)
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM sessions WHERE user_id = $1 and token = $2`)).
		WithArgs(userID, token.Token).
		WillReturnResult(sqlmock.NewResult(0, 1))

	_, _, err := ValidateRefreshToken(manager, userID)
	if err != jwt.ErrTokenExpired {
		t.Fatalf("expected jwt.ErrTokenExpired, got %v", err)
	}
}

func TestValidateRefreshTokenDeleteError(t *testing.T) {
	manager, mock, cleanup := newMockManager(t)
	defer cleanup()

	userID := "user"
	token := authModels.RefreshToken{
		ID:        "3",
		Token:     "old-token",
		ExpiredAt: time.Now().Add(-time.Minute),
	}

	expectSessionQuery(mock, userID, token, false)
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM sessions WHERE user_id = $1 and token = $2`)).
		WithArgs(userID, token.Token).
		WillReturnError(sql.ErrConnDone)

	_, _, err := ValidateRefreshToken(manager, userID)
	if err == nil {
		t.Fatal("expected error")
	}
	if err == jwt.ErrTokenExpired {
		t.Fatal("expected db error but got token expired")
	}
}

func TestGetClaims(t *testing.T) {
	claims := jwt.MapClaims{
		"id":   "user",
		"role": "user",
		"exp":  time.Now().Add(time.Minute).Unix(),
	}

	token, err := GenerateAccessToken(claims)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	got, err := GetClaims(token)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["id"] != "user" {
		t.Fatalf("expected id=user, got %v", got["id"])
	}
}

func TestGetClaimsInvalid(t *testing.T) {
	if _, err := GetClaims("invalid"); err == nil {
		t.Fatal("expected error for invalid token")
	}
}
