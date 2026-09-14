package post_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/configs"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/constants"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/modules/post"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockPostService struct {
	post.PostService
	CreatePostFunc  func(ctx context.Context, userId int, request post.PostRequest) error
	GetAllPostFunc  func(ctx context.Context, pageSize, pageIndex, userID int) (post.GetAllPostResponse, error)
	GetPostByIdFunc func(ctx context.Context, id, userID int) (*post.GetPostResponse, error)
}

func (m *mockPostService) CreatePost(ctx context.Context, userId int, request post.PostRequest) error {
	if m.CreatePostFunc != nil {
		return m.CreatePostFunc(ctx, userId, request)
	}
	return nil
}

func (m *mockPostService) GetAllPost(ctx context.Context, pageSize, pageIndex, userID int) (post.GetAllPostResponse, error) {
	if m.GetAllPostFunc != nil {
		return m.GetAllPostFunc(ctx, pageSize, pageIndex, userID)
	}
	return post.GetAllPostResponse{}, nil
}

func (m *mockPostService) GetPostById(ctx context.Context, id, userID int) (*post.GetPostResponse, error) {
	if m.GetPostByIdFunc != nil {
		return m.GetPostByIdFunc(ctx, id, userID)
	}
	return nil, nil
}

func TestPostHandler_CreatePost_TableDriven(t *testing.T) {
	tests := []struct {
		name           string
		body           any
		mockCreate     func(ctx context.Context, userId int, request post.PostRequest) error
		expectedStatus int
	}{
		{
			name: "Success - 201 Created",
			body: post.PostRequest{
				PostTitle:    "Title",
				PostContent:  "Content",
				PostHashtags: []string{"tech"},
			},
			mockCreate: func(ctx context.Context, userId int, request post.PostRequest) error {
				return nil
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "Failed - 400 Bad Request Invalid JSON",
			body:           "invalid-json",
			mockCreate:     nil,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			svc := &mockPostService{CreatePostFunc: tt.mockCreate}
			h := post.NewPostHandler(svc, &configs.Config{})

			r.POST("/posts/create-post", func(c *gin.Context) {
				c.Set("id", 1)
				h.CreatePost(c)
			})

			var reqBody []byte
			if str, ok := tt.body.(string); ok {
				reqBody = []byte(str)
			} else {
				reqBody, _ = json.Marshal(tt.body)
			}

			req := httptest.NewRequest(http.MethodPost, "/posts/create-post", bytes.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)
			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestPostHandler_GetPostById_TableDriven(t *testing.T) {
	tests := []struct {
		name           string
		postIdParam    string
		mockGet        func(ctx context.Context, id, userID int) (*post.GetPostResponse, error)
		expectedStatus int
	}{
		{
			name:        "Success - 200 OK",
			postIdParam: "1",
			mockGet: func(ctx context.Context, id, userID int) (*post.GetPostResponse, error) {
				return &post.GetPostResponse{
					DetailPost: post.Data{ID: 1, PostTitle: "Test"},
				}, nil
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Failed - 400 Bad Request Invalid Param",
			postIdParam:    "invalid",
			mockGet:        nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "Failed - 404 Not Found",
			postIdParam: "99",
			mockGet: func(ctx context.Context, id, userID int) (*post.GetPostResponse, error) {
				return nil, constants.ErrPostNotFound
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			svc := &mockPostService{GetPostByIdFunc: tt.mockGet}
			h := post.NewPostHandler(svc, &configs.Config{})

			r.GET("/posts/get-post-by-id/:postId", h.GetPostById)

			req := httptest.NewRequest(http.MethodGet, "/posts/get-post-by-id/"+tt.postIdParam, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)
			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
