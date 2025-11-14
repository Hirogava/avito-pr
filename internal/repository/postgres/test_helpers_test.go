package postgres

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func newTestManager(t *testing.T) (*Manager, sqlmock.Sqlmock, func()) {
	t.Helper()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("create sqlmock: %v", err)
	}

	manager := &Manager{Conn: db}
	return manager, mock, func() {
		_ = db.Close()
	}
}
