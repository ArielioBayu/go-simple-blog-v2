package internalsql_test

import (
	"testing"

	"github.com/ArielioBayu/go-simple-blog-v2/pkg/internalsql"
)

func TestConnect_Success(t *testing.T) {
	// DataSource pointing to running docker MySQL
	dsn := "root:secret@tcp(localhost:3308)/db-simple-blog?parseTime=true&loc=Asia%2FJakarta"
	db, err := internalsql.Connect(dsn)
	if err != nil {
		t.Skipf("Skipping integration test if db is not accessible: %v", err)
		return
	}
	defer db.Close()

	stats := db.Stats()
	if stats.MaxOpenConnections != 25 {
		t.Errorf("expected MaxOpenConnections 25, got %d", stats.MaxOpenConnections)
	}
}

func TestConnect_InvalidDSN(t *testing.T) {
	invalidDSN := "root:wrongpassword@tcp(localhost:3308)/db-simple-blog"
	_, err := internalsql.Connect(invalidDSN)
	if err == nil {
		t.Error("expected error with invalid credentials, got nil")
	}
}
