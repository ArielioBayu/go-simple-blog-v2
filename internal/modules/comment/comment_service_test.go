package comment_test

import (
	"context"
	"errors"
	"testing"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/configs"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/constants"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/modules/comment"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockCommentRepo struct {
	comment.CommentRepository
	CreateCommentFunc  func(ctx context.Context, model comment.CommentModel) error
	GetCommentByIdFunc func(ctx context.Context, postId int) ([]comment.GetComment, error)
}

func (m *mockCommentRepo) CreateComment(ctx context.Context, model comment.CommentModel) error {
	if m.CreateCommentFunc != nil {
		return m.CreateCommentFunc(ctx, model)
	}
	return nil
}

func (m *mockCommentRepo) GetCommentById(ctx context.Context, postId int) ([]comment.GetComment, error) {
	if m.GetCommentByIdFunc != nil {
		return m.GetCommentByIdFunc(ctx, postId)
	}
	return nil, nil
}

type mockPostChecker struct {
	CheckPostExistsFunc func(ctx context.Context, postId int) (bool, error)
}

func (m *mockPostChecker) CheckPostExists(ctx context.Context, postId int) (bool, error) {
	if m.CheckPostExistsFunc != nil {
		return m.CheckPostExistsFunc(ctx, postId)
	}
	return true, nil
}

func TestCommentService_CreateComment_TableDriven(t *testing.T) {
	cfg := &configs.Config{}

	tests := []struct {
		name          string
		postId        int
		userId        int
		input         comment.CommentRequest
		mockCheckPost func(ctx context.Context, postId int) (bool, error)
		mockCreate    func(ctx context.Context, model comment.CommentModel) error
		expectedErr   error
	}{
		{
			name:   "Success - Valid comment created",
			postId: 1,
			userId: 2,
			input:  comment.CommentRequest{CommentContent: "Awesome blog post!"},
			mockCheckPost: func(ctx context.Context, postId int) (bool, error) {
				return true, nil
			},
			mockCreate: func(ctx context.Context, model comment.CommentModel) error {
				assert.Equal(t, 1, model.PostId)
				assert.Equal(t, 2, model.UserId)
				assert.Equal(t, "Awesome blog post!", model.CommentContent)
				return nil
			},
			expectedErr: nil,
		},
		{
			name:   "Failed - Post Not Found",
			postId: 999,
			userId: 2,
			input:  comment.CommentRequest{CommentContent: "Awesome!"},
			mockCheckPost: func(ctx context.Context, postId int) (bool, error) {
				return false, nil
			},
			expectedErr: constants.ErrPostNotFound,
		},
		{
			name:   "Failed - Error checking post existence",
			postId: 1,
			userId: 2,
			input:  comment.CommentRequest{CommentContent: "Awesome!"},
			mockCheckPost: func(ctx context.Context, postId int) (bool, error) {
				return false, errors.New("database connection down")
			},
			expectedErr: errors.New("database connection down"),
		},
		{
			name:   "Failed - Error on database insert comment",
			postId: 1,
			userId: 2,
			input:  comment.CommentRequest{CommentContent: "Awesome!"},
			mockCheckPost: func(ctx context.Context, postId int) (bool, error) {
				return true, nil
			},
			mockCreate: func(ctx context.Context, model comment.CommentModel) error {
				return errors.New("insert failed")
			},
			expectedErr: errors.New("insert failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockCommentRepo{CreateCommentFunc: tt.mockCreate}
			checker := &mockPostChecker{CheckPostExistsFunc: tt.mockCheckPost}
			srv := comment.NewCommentService(cfg, repo, checker)

			err := srv.CreateComment(context.Background(), tt.postId, tt.userId, tt.input)
			if tt.expectedErr != nil {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
