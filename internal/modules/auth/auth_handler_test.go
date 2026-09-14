package auth_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/configs"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/constants"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/modules/auth"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

type mockAuthService struct {
	auth.AuthService
	SignUpFunc               func(ctx context.Context, req auth.SignUpRequest) error
	SignInFunc               func(ctx context.Context, req auth.SignInRequest) (string, string, error)
	GetIdRefreshTokenFunc    func(ctx context.Context, req auth.RefreshTokenRequest) (*auth.RefreshTokenModel, error)
	ValidateRefreshTokenFunc func(ctx context.Context, userId int, req auth.RefreshTokenRequest) (string, error)
}

func (m *mockAuthService) SignUp(ctx context.Context, req auth.SignUpRequest) error {
	if m.SignUpFunc != nil {
		return m.SignUpFunc(ctx, req)
	}
	return nil
}

func (m *mockAuthService) SignIn(ctx context.Context, req auth.SignInRequest) (string, string, error) {
	if m.SignInFunc != nil {
		return m.SignInFunc(ctx, req)
	}
	return "mock-token", "mock-refresh-token", nil
}

func (m *mockAuthService) GetIdRefreshToken(ctx context.Context, req auth.RefreshTokenRequest) (*auth.RefreshTokenModel, error) {
	if m.GetIdRefreshTokenFunc != nil {
		return m.GetIdRefreshTokenFunc(ctx, req)
	}
	return &auth.RefreshTokenModel{UserId: 1}, nil
}

func (m *mockAuthService) ValidateRefreshToken(ctx context.Context, userId int, req auth.RefreshTokenRequest) (string, error) {
	if m.ValidateRefreshTokenFunc != nil {
		return m.ValidateRefreshTokenFunc(ctx, userId, req)
	}
	return "mock-new-access-token", nil
}

func TestAuthHandler_SignUp_TableDriven(t *testing.T) {
	tests := []struct {
		name           string
		body           any
		mockSignUp     func(ctx context.Context, req auth.SignUpRequest) error
		expectedStatus int
	}{
		{
			name: "Success - Created 201",
			body: auth.SignUpRequest{Email: "test@mail.com", Username: "test", Password: "pass"},
			mockSignUp: func(ctx context.Context, req auth.SignUpRequest) error {
				return nil
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "Failed - Bad Request 400 Invalid JSON",
			body:           "invalid-json-string",
			mockSignUp:     nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Failed - Conflict 409 User Exists",
			body: auth.SignUpRequest{Email: "exist@mail.com", Username: "exist", Password: "pass"},
			mockSignUp: func(ctx context.Context, req auth.SignUpRequest) error {
				return constants.ErrUsernameOrEmailAlreadyExists
			},
			expectedStatus: http.StatusConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			svc := &mockAuthService{SignUpFunc: tt.mockSignUp}
			handler := auth.NewAuthHandler(svc, &configs.Config{})
			r.POST("/memberships/sign-up", handler.SignUp)

			var reqBody []byte
			if str, ok := tt.body.(string); ok {
				reqBody = []byte(str)
			} else {
				reqBody, _ = json.Marshal(tt.body)
			}

			req := httptest.NewRequest(http.MethodPost, "/memberships/sign-up", bytes.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)
			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestAuthHandler_SignIn_TableDriven(t *testing.T) {
	tests := []struct {
		name           string
		body           any
		mockSignIn     func(ctx context.Context, req auth.SignInRequest) (string, string, error)
		expectedStatus int
	}{
		{
			name: "Success - OK 200 with cookies and tokens",
			body: auth.SignInRequest{Email: "user@mail.com", Password: "pass"},
			mockSignIn: func(ctx context.Context, req auth.SignInRequest) (string, string, error) {
				return "jwt-access-token", "jwt-refresh-token", nil
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Failed - Not Found 404 User Not Found",
			body: auth.SignInRequest{Email: "notfound@mail.com", Password: "pass"},
			mockSignIn: func(ctx context.Context, req auth.SignInRequest) (string, string, error) {
				return "", "", constants.ErrDataNotFound
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "Failed - Bad Request 400 Invalid Password",
			body: auth.SignInRequest{Email: "user@mail.com", Password: "wrong"},
			mockSignIn: func(ctx context.Context, req auth.SignInRequest) (string, string, error) {
				return "", "", constants.ErrInvalidPassword
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			svc := &mockAuthService{SignInFunc: tt.mockSignIn}
			handler := auth.NewAuthHandler(svc, &configs.Config{})
			r.POST("/memberships/sign-in", handler.SignIn)

			reqBody, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/memberships/sign-in", bytes.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)
			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
