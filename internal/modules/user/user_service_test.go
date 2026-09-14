package user_test

import (
	"context"
	"errors"
	"testing"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/configs"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/modules/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockUserRepo struct {
	user.UserRepository
	GetUserByIdFunc func(ctx context.Context, id int) (*user.UserModel, error)
}

func (m *mockUserRepo) GetUserById(ctx context.Context, id int) (*user.UserModel, error) {
	if m.GetUserByIdFunc != nil {
		return m.GetUserByIdFunc(ctx, id)
	}
	return nil, nil
}

func TestUserService_GetUserById_TableDriven(t *testing.T) {
	cfg := &configs.Config{}

	tests := []struct {
		name        string
		userId      int
		mockGet     func(ctx context.Context, id int) (*user.UserModel, error)
		expectedErr error
		expectedRes *user.UserModel
	}{
		{
			name:   "Success - User found",
			userId: 1,
			mockGet: func(ctx context.Context, id int) (*user.UserModel, error) {
				return &user.UserModel{ID: 1, Username: "bayu", Email: "bayu@example.com"}, nil
			},
			expectedErr: nil,
			expectedRes: &user.UserModel{ID: 1, Username: "bayu", Email: "bayu@example.com"},
		},
		{
			name:   "Success - User not found (returns nil)",
			userId: 999,
			mockGet: func(ctx context.Context, id int) (*user.UserModel, error) {
				return nil, nil
			},
			expectedErr: nil,
			expectedRes: nil,
		},
		{
			name:   "Failed - Database query error",
			userId: 2,
			mockGet: func(ctx context.Context, id int) (*user.UserModel, error) {
				return nil, errors.New("db query failed")
			},
			expectedErr: errors.New("db query failed"),
			expectedRes: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockUserRepo{GetUserByIdFunc: tt.mockGet}
			svc := user.NewUserService(cfg, repo)

			res, err := svc.GetUserById(context.Background(), tt.userId)
			if tt.expectedErr != nil {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr.Error())
				assert.Nil(t, res)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedRes, res)
			}
		})
	}
}
