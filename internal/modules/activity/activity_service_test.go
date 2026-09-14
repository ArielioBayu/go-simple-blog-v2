package activity_test

import (
	"context"
	"errors"
	"testing"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/configs"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/modules/activity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockActivityRepo struct {
	activity.ActivityRepository
	UpsertActivitiesFunc   func(ctx context.Context, model activity.ActivityModel) error
	CountLikedByPostIDFunc func(ctx context.Context, postId int) (int, error)
}

func (m *mockActivityRepo) UpsertActivities(ctx context.Context, model activity.ActivityModel) error {
	if m.UpsertActivitiesFunc != nil {
		return m.UpsertActivitiesFunc(ctx, model)
	}
	return nil
}

func (m *mockActivityRepo) CountLikedByPostID(ctx context.Context, postId int) (int, error) {
	if m.CountLikedByPostIDFunc != nil {
		return m.CountLikedByPostIDFunc(ctx, postId)
	}
	return 0, nil
}

func TestActivityService_InsertUpdateActivities_TableDriven(t *testing.T) {
	cfg := &configs.Config{}

	tests := []struct {
		name        string
		postId      int
		userId      int
		input       activity.ActivityRequest
		mockUpsert  func(ctx context.Context, model activity.ActivityModel) error
		expectedErr error
	}{
		{
			name:   "Success - Like post",
			postId: 10,
			userId: 3,
			input:  activity.ActivityRequest{IsLiked: true},
			mockUpsert: func(ctx context.Context, model activity.ActivityModel) error {
				assert.Equal(t, 10, model.PostId)
				assert.Equal(t, 3, model.UserId)
				assert.True(t, model.IsLiked)
				return nil
			},
			expectedErr: nil,
		},
		{
			name:   "Success - Unlike post",
			postId: 10,
			userId: 3,
			input:  activity.ActivityRequest{IsLiked: false},
			mockUpsert: func(ctx context.Context, model activity.ActivityModel) error {
				assert.Equal(t, 10, model.PostId)
				assert.Equal(t, 3, model.UserId)
				assert.False(t, model.IsLiked)
				return nil
			},
			expectedErr: nil,
		},
		{
			name:   "Failed - Database Upsert Error",
			postId: 10,
			userId: 3,
			input:  activity.ActivityRequest{IsLiked: true},
			mockUpsert: func(ctx context.Context, model activity.ActivityModel) error {
				return errors.New("deadlock / connection error")
			},
			expectedErr: errors.New("deadlock / connection error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockActivityRepo{UpsertActivitiesFunc: tt.mockUpsert}
			srv := activity.NewActivityService(cfg, repo)

			err := srv.InsertUpdateActivities(context.Background(), tt.postId, tt.userId, tt.input)
			if tt.expectedErr != nil {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
