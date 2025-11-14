package users

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"

	"github.com/Hirogava/avito-pr/internal/repository/postgres"
)

func setupUsersContext(t *testing.T, method, path, body string) (*gin.Context, *httptest.ResponseRecorder) {
	t.Helper()

	gin.SetMode(gin.TestMode)
	req, err := http.NewRequest(method, path, strings.NewReader(body))
	if err != nil {
		t.Fatalf("new request error: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	return c, w
}

func TestSetIsActiveForbidden(t *testing.T) {
	c, w := setupUsersContext(t, http.MethodPost, "/users/setIsActive", "")

	SetIsActive(c, nil)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", w.Code)
	}
}

func TestSetIsActiveBadRequest(t *testing.T) {
	c, w := setupUsersContext(t, http.MethodPost, "/users/setIsActive", `{"user_id":1}`)
	c.Set("role", "admin")

	SetIsActive(c, nil)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestGetReviewBadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	req, err := http.NewRequest(http.MethodGet, "/users/getReview", nil)
	if err != nil {
		t.Fatalf("new request error: %v", err)
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	GetReview(c, nil)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestGerUsersInternalError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("create sqlmock: %v", err)
	}
	defer db.Close() //nolint:errcheck

	manager := &postgres.Manager{Conn: db}
	mock.ExpectQuery(`SELECT username, team_name, user_id, is_active FROM users`).WillReturnError(assertAnError{})

	gin.SetMode(gin.TestMode)
	req, _ := http.NewRequest(http.MethodGet, "/users", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	GerUsers(c, manager)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

type assertAnError struct{}

func (assertAnError) Error() string { return "error" }
