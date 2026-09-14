package comment_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/configs"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/constants"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/modules/comment"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockCommentService struct {
	comment.CommentService
	CreateCommentFunc func(ctx context.Context, postId, userId int, request comment.CommentRequest) error
}

func (m *mockCommentService) CreateComment(ctx context.Context, postId, userId int, request comment.CommentRequest) error {
	if m.CreateCommentFunc != nil {
		return m.CreateCommentFunc(ctx, postId, userId, request)
	}
	return nil
}

func TestCommentHandler_CreateComment_TableDriven(t *testing.T) {
	tests := []struct {
		name           string
		postIdParam    string
		userIdContext  int
		body           any
		mockCreate     func(ctx context.Context, postId, userId int, request comment.CommentRequest) error
		expectedStatus int
	}{
		{
			name:          "Success - 201 Created",
			postIdParam:   "10",
			userIdContext: 1,
			body:          comment.CommentRequest{CommentContent: "Nice blog!"},
			mockCreate: func(ctx context.Context, postId, userId int, request comment.CommentRequest) error {
				return nil
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "Failed - 400 Bad Request Invalid JSON",
			postIdParam:    "10",
			userIdContext:  1,
			body:           "invalid-json",
			mockCreate:     nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Failed - 400 Bad Request Invalid Post ID Param",
			postIdParam:    "invalid-id",
			userIdContext:  1,
			body:           comment.CommentRequest{CommentContent: "Nice blog!"},
			mockCreate:     nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:          "Failed - 404 Not Found Post",
			postIdParam:   "99",
			userIdContext: 1,
			body:          comment.CommentRequest{CommentContent: "Nice blog!"},
			mockCreate: func(ctx context.Context, postId, userId int, request comment.CommentRequest) error {
				return constants.ErrPostNotFound
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			svc := &mockCommentService{CreateCommentFunc: tt.mockCreate}
			h := comment.NewCommentHandler(svc, &configs.Config{})

			r.POST("/posts/create-comment/:postId", func(c *gin.Context) {
				c.Set("id", tt.userIdContext)
				h.CreateComment(c)
			})

			var reqBody []byte
			if str, ok := tt.body.(string); ok {
				reqBody = []byte(str)
			} else {
				reqBody, _ = json.Marshal(tt.body)
			}

			req := httptest.NewRequest(http.MethodPost, "/posts/create-comment/"+tt.postIdParam, bytes.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)
			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
