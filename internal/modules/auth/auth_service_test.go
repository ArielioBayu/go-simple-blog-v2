package auth_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/configs"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/constants"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/modules/auth"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/modules/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

type mockUserRepo struct {
	user.UserRepository
	GetUserFunc        func(ctx context.Context, email, username string) (*user.UserModel, error)
	GetUserByIdFunc    func(ctx context.Context, id int) (*user.UserModel, error)
	GetUserByEmailFunc func(ctx context.Context, email string) (*user.UserModel, error)
	CreateUserFunc     func(ctx context.Context, model user.UserModel) error
}

func (m *mockUserRepo) GetUser(ctx context.Context, email, username string) (*user.UserModel, error) {
	if m.GetUserFunc != nil {
		return m.GetUserFunc(ctx, email, username)
	}
	return nil, nil
}

func (m *mockUserRepo) GetUserById(ctx context.Context, id int) (*user.UserModel, error) {
	if m.GetUserByIdFunc != nil {
		return m.GetUserByIdFunc(ctx, id)
	}
	return nil, nil
}

func (m *mockUserRepo) GetUserByEmail(ctx context.Context, email string) (*user.UserModel, error) {
	if m.GetUserByEmailFunc != nil {
		return m.GetUserByEmailFunc(ctx, email)
	}
	return nil, nil
}

func (m *mockUserRepo) CreateUser(ctx context.Context, model user.UserModel) error {
	if m.CreateUserFunc != nil {
		return m.CreateUserFunc(ctx, model)
	}
	return nil
}

type mockAuthRepo struct {
	auth.AuthRepository
	InsertRefreshTokenFunc         func(ctx context.Context, model auth.RefreshTokenModel) error
	GetRefreshTokenFunc            func(ctx context.Context, userId int, now time.Time) (*auth.RefreshTokenModel, error)
	GetIdRefreshTokenFunc          func(ctx context.Context, request auth.RefreshTokenRequest) (*auth.RefreshTokenModel, error)
	DeleteExpiredRefreshTokensFunc func(ctx context.Context, userId int, now time.Time) error
}

func (m *mockAuthRepo) InsertRefreshToken(ctx context.Context, model auth.RefreshTokenModel) error {
	if m.InsertRefreshTokenFunc != nil {
		return m.InsertRefreshTokenFunc(ctx, model)
	}
	return nil
}

func (m *mockAuthRepo) GetRefreshToken(ctx context.Context, userId int, now time.Time) (*auth.RefreshTokenModel, error) {
	if m.GetRefreshTokenFunc != nil {
		return m.GetRefreshTokenFunc(ctx, userId, now)
	}
	return nil, nil
}

func (m *mockAuthRepo) GetIdRefreshToken(ctx context.Context, request auth.RefreshTokenRequest) (*auth.RefreshTokenModel, error) {
	if m.GetIdRefreshTokenFunc != nil {
		return m.GetIdRefreshTokenFunc(ctx, request)
	}
	return nil, nil
}

func (m *mockAuthRepo) DeleteExpiredRefreshTokens(ctx context.Context, userId int, now time.Time) error {
	if m.DeleteExpiredRefreshTokensFunc != nil {
		return m.DeleteExpiredRefreshTokensFunc(ctx, userId, now)
	}
	return nil
}

func TestAuthService_SignUp_TableDriven(t *testing.T) {
	cfg := &configs.Config{Service: configs.Service{SecretKey: "secret-key"}}

	tests := []struct {
		name        string
		input       auth.SignUpRequest
		mockGetUser func(ctx context.Context, email, username string) (*user.UserModel, error)
		mockCreate  func(ctx context.Context, model user.UserModel) error
		expectedErr error
	}{
		{
			name:  "Success - New user registered",
			input: auth.SignUpRequest{Email: "new@example.com", Username: "newuser", Password: "password123"},
			mockGetUser: func(ctx context.Context, email, username string) (*user.UserModel, error) {
				return nil, nil
			},
			mockCreate: func(ctx context.Context, model user.UserModel) error {
				return nil
			},
			expectedErr: nil,
		},
		{
			name:  "Failed - Username or Email already exists",
			input: auth.SignUpRequest{Email: "exist@example.com", Username: "existuser", Password: "password123"},
			mockGetUser: func(ctx context.Context, email, username string) (*user.UserModel, error) {
				return &user.UserModel{ID: 1, Email: email, Username: username}, nil
			},
			expectedErr: constants.ErrUsernameOrEmailAlreadyExists,
		},
		{
			name:  "Failed - Database query error on check user",
			input: auth.SignUpRequest{Email: "err@example.com", Username: "erruser", Password: "password123"},
			mockGetUser: func(ctx context.Context, email, username string) (*user.UserModel, error) {
				return nil, errors.New("db connection failure")
			},
			expectedErr: errors.New("db connection failure"),
		},
		{
			name:  "Failed - Database insert error on create user",
			input: auth.SignUpRequest{Email: "insert_err@example.com", Username: "erruser", Password: "password123"},
			mockGetUser: func(ctx context.Context, email, username string) (*user.UserModel, error) {
				return nil, nil
			},
			mockCreate: func(ctx context.Context, model user.UserModel) error {
				return errors.New("insert failed")
			},
			expectedErr: errors.New("insert failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uRepo := &mockUserRepo{
				GetUserFunc:    tt.mockGetUser,
				CreateUserFunc: tt.mockCreate,
			}
			aRepo := &mockAuthRepo{}
			svc := auth.NewAuthService(cfg, aRepo, uRepo)

			err := svc.SignUp(context.Background(), tt.input)
			if tt.expectedErr != nil {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestAuthService_SignIn_TableDriven(t *testing.T) {
	cfg := &configs.Config{Service: configs.Service{SecretKey: "secret-key-12345"}}
	hashedPass, err := bcrypt.GenerateFromPassword([]byte("validpassword"), bcrypt.DefaultCost)
	require.NoError(t, err)

	validUser := &user.UserModel{
		ID:       10,
		Email:    "valid@example.com",
		Username: "validuser",
		Password: string(hashedPass),
	}

	tests := []struct {
		name          string
		input         auth.SignInRequest
		mockGetUser   func(ctx context.Context, email string) (*user.UserModel, error)
		mockGetToken  func(ctx context.Context, userId int, now time.Time) (*auth.RefreshTokenModel, error)
		mockInsert    func(ctx context.Context, model auth.RefreshTokenModel) error
		expectedErr   error
		expectToken   bool
	}{
		{
			name:  "Success - Valid login with new refresh token generation",
			input: auth.SignInRequest{Email: "valid@example.com", Password: "validpassword"},
			mockGetUser: func(ctx context.Context, email string) (*user.UserModel, error) {
				return validUser, nil
			},
			mockGetToken: func(ctx context.Context, userId int, now time.Time) (*auth.RefreshTokenModel, error) {
				return nil, nil // Token belum ada, buat baru
			},
			mockInsert: func(ctx context.Context, model auth.RefreshTokenModel) error {
				return nil
			},
			expectedErr: nil,
			expectToken: true,
		},
		{
			name:  "Success - Valid login reusing active refresh token",
			input: auth.SignInRequest{Email: "valid@example.com", Password: "validpassword"},
			mockGetUser: func(ctx context.Context, email string) (*user.UserModel, error) {
				return validUser, nil
			},
			mockGetToken: func(ctx context.Context, userId int, now time.Time) (*auth.RefreshTokenModel, error) {
				return &auth.RefreshTokenModel{RefreshToken: "existing-refresh-token"}, nil
			},
			expectedErr: nil,
			expectToken: true,
		},
		{
			name:  "Failed - User not found",
			input: auth.SignInRequest{Email: "notfound@example.com", Password: "validpassword"},
			mockGetUser: func(ctx context.Context, email string) (*user.UserModel, error) {
				return nil, nil
			},
			expectedErr: constants.ErrDataNotFound,
			expectToken: false,
		},
		{
			name:  "Failed - Invalid password",
			input: auth.SignInRequest{Email: "valid@example.com", Password: "wrongpassword"},
			mockGetUser: func(ctx context.Context, email string) (*user.UserModel, error) {
				return validUser, nil
			},
			expectedErr: constants.ErrInvalidPassword,
			expectToken: false,
		},
		{
			name:  "Failed - Database error on GetUserByEmail",
			input: auth.SignInRequest{Email: "err@example.com", Password: "validpassword"},
			mockGetUser: func(ctx context.Context, email string) (*user.UserModel, error) {
				return nil, errors.New("db error")
			},
			expectedErr: errors.New("db error"),
			expectToken: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uRepo := &mockUserRepo{GetUserByEmailFunc: tt.mockGetUser}
			aRepo := &mockAuthRepo{
				GetRefreshTokenFunc:    tt.mockGetToken,
				InsertRefreshTokenFunc: tt.mockInsert,
				DeleteExpiredRefreshTokensFunc: func(ctx context.Context, userId int, now time.Time) error {
					return nil
				},
			}
			svc := auth.NewAuthService(cfg, aRepo, uRepo)

			token, refreshToken, err := svc.SignIn(context.Background(), tt.input)
			if tt.expectedErr != nil {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr.Error())
				assert.Empty(t, token)
				assert.Empty(t, refreshToken)
			} else {
				require.NoError(t, err)
				assert.NotEmpty(t, token)
				assert.NotEmpty(t, refreshToken)
			}
		})
	}
}

func TestAuthService_ValidateRefreshToken_TableDriven(t *testing.T) {
	cfg := &configs.Config{Service: configs.Service{SecretKey: "secret-key-12345"}}

	tests := []struct {
		name         string
		userId       int
		input        auth.RefreshTokenRequest
		mockGetToken func(ctx context.Context, userId int, now time.Time) (*auth.RefreshTokenModel, error)
		mockGetUser  func(ctx context.Context, id int) (*user.UserModel, error)
		expectedErr  error
	}{
		{
			name:   "Success - Valid token",
			userId: 1,
			input:  auth.RefreshTokenRequest{Token: "valid-token"},
			mockGetToken: func(ctx context.Context, userId int, now time.Time) (*auth.RefreshTokenModel, error) {
				return &auth.RefreshTokenModel{RefreshToken: "valid-token"}, nil
			},
			mockGetUser: func(ctx context.Context, id int) (*user.UserModel, error) {
				return &user.UserModel{ID: 1, Username: "bayu"}, nil
			},
			expectedErr: nil,
		},
		{
			name:   "Failed - Token expired / not found in DB",
			userId: 1,
			input:  auth.RefreshTokenRequest{Token: "expired-token"},
			mockGetToken: func(ctx context.Context, userId int, now time.Time) (*auth.RefreshTokenModel, error) {
				return nil, nil
			},
			expectedErr: constants.ErrTokenExpired,
		},
		{
			name:   "Failed - Token mismatch (Invalid Token)",
			userId: 1,
			input:  auth.RefreshTokenRequest{Token: "mismatch-token"},
			mockGetToken: func(ctx context.Context, userId int, now time.Time) (*auth.RefreshTokenModel, error) {
				return &auth.RefreshTokenModel{RefreshToken: "stored-token"}, nil
			},
			expectedErr: constants.ErrInvalidToken,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uRepo := &mockUserRepo{GetUserByIdFunc: tt.mockGetUser}
			aRepo := &mockAuthRepo{GetRefreshTokenFunc: tt.mockGetToken}
			svc := auth.NewAuthService(cfg, aRepo, uRepo)

			token, err := svc.ValidateRefreshToken(context.Background(), tt.userId, tt.input)
			if tt.expectedErr != nil {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr.Error())
				assert.Empty(t, token)
			} else {
				require.NoError(t, err)
				assert.NotEmpty(t, token)
			}
		})
	}
}
