package memberships_test

import (
	"context"
	"testing"
	"time"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/configs"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/model/memberships"
	service "github.com/ArielioBayu/go-simple-blog-v2/internal/service/memberships"
	"golang.org/x/crypto/bcrypt"
)

type mockMembershipsRepo struct {
	deleteExpiredCalled bool
	userIdPassed        int
	user                *memberships.UserModel
	existingToken       *memberships.RefreshTokenModel
	insertedToken       memberships.RefreshTokenModel
}

func (m *mockMembershipsRepo) GetUser(ctx context.Context, email, username string) (*memberships.UserModel, error) {
	return m.user, nil
}
func (m *mockMembershipsRepo) GetUserById(ctx context.Context, id int) (*memberships.UserModel, error) {
	return m.user, nil
}
func (m *mockMembershipsRepo) GetUserByEmail(ctx context.Context, email string) (*memberships.UserModel, error) {
	return m.user, nil
}
func (m *mockMembershipsRepo) GetRefreshToken(ctx context.Context, userId int, now time.Time) (*memberships.RefreshTokenModel, error) {
	return m.existingToken, nil
}
func (m *mockMembershipsRepo) GetIdRefreshToken(ctx context.Context, request memberships.RefreshTokenRequest) (*memberships.RefreshTokenModel, error) {
	return m.existingToken, nil
}
func (m *mockMembershipsRepo) CreateUser(ctx context.Context, model memberships.UserModel) error {
	return nil
}
func (m *mockMembershipsRepo) InsertRefreshToken(ctx context.Context, model memberships.RefreshTokenModel) error {
	m.insertedToken = model
	return nil
}
func (m *mockMembershipsRepo) DeleteExpiredRefreshTokens(ctx context.Context, userId int, now time.Time) error {
	m.deleteExpiredCalled = true
	m.userIdPassed = userId
	return nil
}

func TestSignIn_PurgesExpiredTokens(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("secret123"), bcrypt.DefaultCost)
	mockRepo := &mockMembershipsRepo{
		user: &memberships.UserModel{
			ID:       50,
			Email:    "test@example.com",
			Username: "testuser",
			Password: string(hashedPassword),
		},
		existingToken: nil,
	}

	cfg := &configs.Config{
		Service: configs.Service{
			SecretKey: "secret-key",
		},
	}

	srv := service.NewMembershipsService(cfg, mockRepo)

	token, refreshToken, err := srv.SignIn(context.Background(), memberships.SignInRequest{
		Email:    "test@example.com",
		Password: "secret123",
	})

	if err != nil {
		t.Fatalf("unexpected error during SignIn: %v", err)
	}

	if token == "" || refreshToken == "" {
		t.Errorf("expected token and refreshToken to be non-empty")
	}

	if !mockRepo.deleteExpiredCalled {
		t.Errorf("expected DeleteExpiredRefreshTokens to be called on SignIn")
	}

	if mockRepo.userIdPassed != 50 {
		t.Errorf("expected userId 50 passed to DeleteExpiredRefreshTokens, got %d", mockRepo.userIdPassed)
	}
}
