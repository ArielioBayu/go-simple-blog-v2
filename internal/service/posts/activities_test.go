package posts_test

import (
	"context"
	"testing"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/model/posts"
	service "github.com/ArielioBayu/go-simple-blog-v2/internal/service/posts"
)

type mockPostsRepo struct {
	upsertCalled bool
	lastModel    posts.ActivityModel
}

func (m *mockPostsRepo) CountLikedByPostID(ctx context.Context, postId int) (int, error) {
	return 0, nil
}
func (m *mockPostsRepo) CreatePost(ctx context.Context, model posts.PostModel) error {
	return nil
}
func (m *mockPostsRepo) CreateComment(ctx context.Context, model posts.CommentModel) error {
	return nil
}
func (m *mockPostsRepo) CreateActivities(ctx context.Context, model posts.ActivityModel) error {
	return nil
}
func (m *mockPostsRepo) GetActivities(ctx context.Context, postId, userId int) (*posts.ActivityModel, error) {
	return nil, nil
}
func (m *mockPostsRepo) GetAllPost(ctx context.Context, limit, offset, userID int) (posts.GetAllPostResponse, error) {
	return posts.GetAllPostResponse{}, nil
}
func (m *mockPostsRepo) GetPostById(ctx context.Context, id, userID int) (*posts.Data, error) {
	return nil, nil
}
func (m *mockPostsRepo) GetCommentById(ctx context.Context, postId int) ([]posts.GetComment, error) {
	return nil, nil
}
func (m *mockPostsRepo) UpdateActivities(ctx context.Context, model posts.ActivityModel) error {
	return nil
}
func (m *mockPostsRepo) UpsertActivities(ctx context.Context, model posts.ActivityModel) error {
	m.upsertCalled = true
	m.lastModel = model
	return nil
}

func TestInsertUpdateActivities_AtomicUpsert(t *testing.T) {
	mockRepo := &mockPostsRepo{}
	srv := service.NewPostsService(nil, mockRepo)

	err := srv.InsertUpdateActivities(context.Background(), 10, 50, posts.ActivityRequest{
		IsLiked: true,
	})

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if !mockRepo.upsertCalled {
		t.Errorf("expected UpsertActivities to be called directly")
	}

	if mockRepo.lastModel.PostId != 10 || mockRepo.lastModel.UserId != 50 || !mockRepo.lastModel.IsLiked {
		t.Errorf("unexpected model content: %+v", mockRepo.lastModel)
	}
}
