package auth

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/Hirogava/avito-pr/internal/config/logger"
)

func setupAuthContext(t *testing.T, method, path, body string) (*gin.Context, *httptest.ResponseRecorder) {
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

func TestAdminBadRequest(t *testing.T) {
	logger.Logger = logrus.New()
	logger.Logger.SetOutput(io.Discard)

	c, w := setupAuthContext(t, http.MethodPost, "/auth/admin", `{"user_id":1}`)

	Admin(c, nil)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestRefreshTokenBadRequest(t *testing.T) {
	logger.Logger = logrus.New()
	logger.Logger.SetOutput(io.Discard)

	c, w := setupAuthContext(t, http.MethodPost, "/auth/refresh", `{"access_token":1}`)
	c.Set("userID", "user")

	RefreshToken(c, nil)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}
