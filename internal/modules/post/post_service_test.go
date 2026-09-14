package post_test

import (
	"context"
	"errors"
	"testing"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/configs"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/constants"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/modules/activity"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/modules/comment"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/modules/post"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockPostRepo struct {
	post.PostRepository
	CreatePostFunc      func(ctx context.Context, model post.PostModel) error
	GetAllPostFunc      func(ctx context.Context, limit, offset, userID int) (post.GetAllPostResponse, error)
	GetPostByIdFunc     func(ctx context.Context, id, userID int) (*post.Data, error)
	CheckPostExistsFunc func(ctx context.Context, id int) (bool, error)
}

func (m *mockPostRepo) CreatePost(ctx context.Context, model post.PostModel) error {
	if m.CreatePostFunc != nil {
		return m.CreatePostFunc(ctx, model)
	}
	return nil
}

func (m *mockPostRepo) GetAllPost(ctx context.Context, limit, offset, userID int) (post.GetAllPostResponse, error) {
	if m.GetAllPostFunc != nil {
		return m.GetAllPostFunc(ctx, limit, offset, userID)
	}
	return post.GetAllPostResponse{}, nil
}

func (m *mockPostRepo) GetPostById(ctx context.Context, id, userID int) (*post.Data, error) {
	if m.GetPostByIdFunc != nil {
		return m.GetPostByIdFunc(ctx, id, userID)
	}
	return nil, nil
}

type mockCommentRepo struct {
	comment.CommentRepository
	GetCommentByIdFunc func(ctx context.Context, postId int) ([]comment.GetComment, error)
}

func (m *mockCommentRepo) GetCommentById(ctx context.Context, postId int) ([]comment.GetComment, error) {
	if m.GetCommentByIdFunc != nil {
		return m.GetCommentByIdFunc(ctx, postId)
	}
	return nil, nil
}

type mockActivityRepo struct {
	activity.ActivityRepository
	CountLikedByPostIDFunc func(ctx context.Context, postId int) (int, error)
}

func (m *mockActivityRepo) CountLikedByPostID(ctx context.Context, postId int) (int, error) {
	if m.CountLikedByPostIDFunc != nil {
		return m.CountLikedByPostIDFunc(ctx, postId)
	}
	return 0, nil
}

func TestPostService_CreatePost_TableDriven(t *testing.T) {
	cfg := &configs.Config{}

	tests := []struct {
		name        string
		userId      int
		input       post.PostRequest
		mockCreate  func(ctx context.Context, model post.PostModel) error
		expectedErr error
	}{
		{
			name:   "Success - Post created",
			userId: 1,
			input: post.PostRequest{
				PostTitle:    "Hello Golang",
				PostContent:  "Golang is powerful",
				PostHashtags: []string{"golang", "backend"},
			},
			mockCreate: func(ctx context.Context, model post.PostModel) error {
				assert.Equal(t, 1, model.UserId)
				assert.Equal(t, "Hello Golang", model.PostTitle)
				assert.Equal(t, "golang,backend", model.PostHashtags)
				return nil
			},
			expectedErr: nil,
		},
		{
			name:   "Failed - Database insert error",
			userId: 1,
			input: post.PostRequest{
				PostTitle:   "Error Post",
				PostContent: "Content",
			},
			mockCreate: func(ctx context.Context, model post.PostModel) error {
				return errors.New("insert post failed")
			},
			expectedErr: errors.New("insert post failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pRepo := &mockPostRepo{CreatePostFunc: tt.mockCreate}
			srv := post.NewPostService(cfg, pRepo, nil, nil)

			err := srv.CreatePost(context.Background(), tt.userId, tt.input)
			if tt.expectedErr != nil {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestPostService_GetPostById_TableDriven(t *testing.T) {
	cfg := &configs.Config{}

	mockPostData := &post.Data{
		ID:          1,
		UserId:      5,
		Username:    "bayu",
		PostTitle:   "First Post",
		PostContent: "First post body",
		IsLiked:     true,
	}

	tests := []struct {
		name          string
		postId        int
		userId        int
		mockGetPost   func(ctx context.Context, id, userID int) (*post.Data, error)
		mockGetCount  func(ctx context.Context, postId int) (int, error)
		mockGetComm   func(ctx context.Context, postId int) ([]comment.GetComment, error)
		expectedErr   error
		expectedLikes int
		expectedComms int
	}{
		{
			name:   "Success - Complete detail with comments and like count",
			postId: 1,
			userId: 5,
			mockGetPost: func(ctx context.Context, id, userID int) (*post.Data, error) {
				return mockPostData, nil
			},
			mockGetCount: func(ctx context.Context, postId int) (int, error) {
				return 15, nil
			},
			mockGetComm: func(ctx context.Context, postId int) ([]comment.GetComment, error) {
				return []comment.GetComment{
					{ID: 1, UserId: 2, Username: "ariel", CommentContent: "Great!"},
				}, nil
			},
			expectedErr:   nil,
			expectedLikes: 15,
			expectedComms: 1,
		},
		{
			name:   "Failed - Post Not Found",
			postId: 99,
			userId: 5,
			mockGetPost: func(ctx context.Context, id, userID int) (*post.Data, error) {
				return nil, constants.ErrPostNotFound
			},
			expectedErr: constants.ErrPostNotFound,
		},
		{
			name:   "Failed - Database error on CountLikedByPostID",
			postId: 1,
			userId: 5,
			mockGetPost: func(ctx context.Context, id, userID int) (*post.Data, error) {
				return mockPostData, nil
			},
			mockGetCount: func(ctx context.Context, postId int) (int, error) {
				return 0, errors.New("db count error")
			},
			expectedErr: errors.New("db count error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pRepo := &mockPostRepo{GetPostByIdFunc: tt.mockGetPost}
			cRepo := &mockCommentRepo{GetCommentByIdFunc: tt.mockGetComm}
			aRepo := &mockActivityRepo{CountLikedByPostIDFunc: tt.mockGetCount}

			srv := post.NewPostService(cfg, pRepo, cRepo, aRepo)
			res, err := srv.GetPostById(context.Background(), tt.postId, tt.userId)

			if tt.expectedErr != nil {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr.Error())
				assert.Nil(t, res)
			} else {
				require.NoError(t, err)
				require.NotNil(t, res)
				assert.Equal(t, tt.expectedLikes, res.LikedCount)
				assert.Len(t, res.Comments, tt.expectedComms)
			}
		})
	}
}
