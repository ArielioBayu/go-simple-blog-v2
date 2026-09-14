package activity_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/configs"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/modules/activity"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockActivityService struct {
	activity.ActivityService
	InsertUpdateActivitiesFunc func(ctx context.Context, postId, userId int, request activity.ActivityRequest) error
}

func (m *mockActivityService) InsertUpdateActivities(ctx context.Context, postId, userId int, request activity.ActivityRequest) error {
	if m.InsertUpdateActivitiesFunc != nil {
		return m.InsertUpdateActivitiesFunc(ctx, postId, userId, request)
	}
	return nil
}

func TestActivityHandler_InsertUpdateActivities_TableDriven(t *testing.T) {
	tests := []struct {
		name           string
		postIdParam    string
		hasContextUser bool
		contextUserID  int
		body           any
		mockInsert     func(ctx context.Context, postId, userId int, request activity.ActivityRequest) error
		expectedStatus int
	}{
		{
			name:           "Success - 200 OK Like Post",
			postIdParam:    "10",
			hasContextUser: true,
			contextUserID:  2,
			body:           activity.ActivityRequest{IsLiked: true},
			mockInsert: func(ctx context.Context, postId, userId int, request activity.ActivityRequest) error {
				return nil
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Failed - 400 Bad Request Invalid JSON",
			postIdParam:    "10",
			hasContextUser: true,
			contextUserID:  2,
			body:           "invalid-json",
			mockInsert:     nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Failed - 400 Bad Request Invalid Post ID",
			postIdParam:    "invalid-post",
			hasContextUser: true,
			contextUserID:  2,
			body:           activity.ActivityRequest{IsLiked: true},
			mockInsert:     nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Failed - 401 Unauthorized Missing User ID in Context",
			postIdParam:    "10",
			hasContextUser: false,
			contextUserID:  0,
			body:           activity.ActivityRequest{IsLiked: true},
			mockInsert:     nil,
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			svc := &mockActivityService{InsertUpdateActivitiesFunc: tt.mockInsert}
			h := activity.NewActivityHandler(svc, &configs.Config{})

			r.POST("/posts/user-activity/:postId", func(c *gin.Context) {
				if tt.hasContextUser {
					c.Set("id", tt.contextUserID)
				}
				h.InsertUpdateActivities(c)
			})

			var reqBody []byte
			if str, ok := tt.body.(string); ok {
				reqBody = []byte(str)
			} else {
				reqBody, _ = json.Marshal(tt.body)
			}

			req := httptest.NewRequest(http.MethodPost, "/posts/user-activity/"+tt.postIdParam, bytes.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)
			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
