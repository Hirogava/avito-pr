package prs

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupRequest(t *testing.T, method, path string, body []byte) (*gin.Context, *httptest.ResponseRecorder) {
	t.Helper()

	gin.SetMode(gin.TestMode)

	req, err := http.NewRequest(method, path, bytes.NewReader(body))
	if err != nil {
		t.Fatalf("new request error: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	return c, w
}

func TestCreatePRForbiddenWithoutRole(t *testing.T) {
	c, w := setupRequest(t, http.MethodPost, "/pullRequest/create", nil)

	CreatePR(c, nil)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d", w.Code)
	}
}

func TestCreatePRBadRequest(t *testing.T) {
	c, w := setupRequest(t, http.MethodPost, "/pullRequest/create", []byte(`{"pull_request_id": 1}`))
	c.Set("role", "admin")

	CreatePR(c, nil)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestMergePRForbiddenWithoutRole(t *testing.T) {
	c, w := setupRequest(t, http.MethodPost, "/pullRequest/merge", nil)

	MergePR(c, nil)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d", w.Code)
	}
}

func TestMergePRBadRequest(t *testing.T) {
	c, w := setupRequest(t, http.MethodPost, "/pullRequest/merge", []byte(`{"pull_request_id": 1}`))
	c.Set("role", "admin")

	MergePR(c, nil)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestReassignPRForbiddenWithoutRole(t *testing.T) {
	c, w := setupRequest(t, http.MethodPost, "/pullRequest/reassign", nil)

	ReassignAuthor(c, nil)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d", w.Code)
	}
}

func TestReassignPRBadRequest(t *testing.T) {
	c, w := setupRequest(t, http.MethodPost, "/pullRequest/reassign", []byte(`{"pull_request_id": "pr"}`))
	c.Set("role", "admin")

	ReassignAuthor(c, nil)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}
