package memberships_test

import (
	"context"
	"testing"
	"time"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/model/memberships"
	repo "github.com/ArielioBayu/go-simple-blog-v2/internal/repository/memberships"
	"github.com/ArielioBayu/go-simple-blog-v2/pkg/internalsql"
)

func TestDeleteExpiredRefreshTokens_Repository(t *testing.T) {
	dsn := "root:secret@tcp(localhost:3308)/db-simple-blog?parseTime=true&loc=Asia%2FJakarta"
	db, err := internalsql.Connect(dsn)
	if err != nil {
		t.Skipf("Skipping repository test if db is not accessible: %v", err)
		return
	}
	defer db.Close()

	r := repo.NewMembershipsRepository(db)
	ctx := context.Background()

	// Insert an expired token for test user 50
	expiredToken := "test-expired-token-" + time.Now().Format("150405.000000")
	err = r.InsertRefreshToken(ctx, memberships.RefreshTokenModel{
		UserId:       50,
		RefreshToken: expiredToken,
		ExpiredAt:    time.Now().Add(-1 * time.Hour), // Expired 1 hour ago
		CreatedAt:    time.Now().Add(-2 * time.Hour),
		UpdatedAt:    time.Now().Add(-2 * time.Hour),
		CreatedBy:    "test",
		UpdatedBy:    "test",
	})
	if err != nil {
		t.Fatalf("failed to insert test expired token: %v", err)
	}

	// Purge expired tokens
	err = r.DeleteExpiredRefreshTokens(ctx, 50, time.Now())
	if err != nil {
		t.Fatalf("failed to delete expired refresh tokens: %v", err)
	}

	// Verify that the expired token is no longer in the DB
	found, err := r.GetIdRefreshToken(ctx, memberships.RefreshTokenRequest{Token: expiredToken})
	if err != nil {
		t.Fatalf("unexpected error searching for deleted token: %v", err)
	}
	if found != nil {
		t.Errorf("expected token to be purged, but it was found: %+v", found)
	}
}
