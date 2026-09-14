package user_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/configs"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/modules/user"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockUserService struct {
	user.UserService
	GetUserByIdFunc func(ctx context.Context, id int) (*user.UserModel, error)
}

func (m *mockUserService) GetUserById(ctx context.Context, id int) (*user.UserModel, error) {
	if m.GetUserByIdFunc != nil {
		return m.GetUserByIdFunc(ctx, id)
	}
	return nil, nil
}

func TestUserHandler_GetUser_TableDriven(t *testing.T) {
	tests := []struct {
		name           string
		contextUserID  int
		mockGet        func(ctx context.Context, id int) (*user.UserModel, error)
		expectedStatus int
	}{
		{
			name:          "Success - 200 OK",
			contextUserID: 5,
			mockGet: func(ctx context.Context, id int) (*user.UserModel, error) {
				return &user.UserModel{ID: 5, Username: "bayu", Email: "bayu@mail.com", CreatedAt: time.Now()}, nil
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:          "Failed - 401 Unauthorized (No User ID in Context)",
			contextUserID: 0,
			mockGet:       nil,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:          "Failed - 404 User Not Found",
			contextUserID: 99,
			mockGet: func(ctx context.Context, id int) (*user.UserModel, error) {
				return nil, nil
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:          "Failed - 500 Internal Server Error",
			contextUserID: 5,
			mockGet: func(ctx context.Context, id int) (*user.UserModel, error) {
				return nil, errors.New("db error")
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			svc := &mockUserService{GetUserByIdFunc: tt.mockGet}
			h := user.NewUserHandler(svc, &configs.Config{})

			r.GET("/memberships/get-user", func(c *gin.Context) {
				if tt.contextUserID > 0 {
					c.Set("id", tt.contextUserID)
				}
				h.GetUser(c)
			})

			req := httptest.NewRequest(http.MethodGet, "/memberships/get-user", nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)
			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
