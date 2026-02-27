package posts

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/model/posts"
)

type PostsRepository interface {
	CreatePost(ctx context.Context, model posts.PostModel) error
	CreateComment(ctx context.Context, model posts.CommentModel) error
	CreateActivities(ctx context.Context, model posts.ActivityModel) error
	GetActivities(ctx context.Context, postId, userId int) (*posts.ActivityModel, error)
	GetAllPost(ctx context.Context, limit, offset int) (posts.GetAllPostResponse, error)
	UpdateActivities(ctx context.Context, model posts.ActivityModel) error
}

type postsRepository struct {
	DB *sql.DB
}

func NewPostsRepository(db *sql.DB) *postsRepository {
	return &postsRepository{
		DB: db,
	}
}

func (r *postsRepository) CreatePost(ctx context.Context, model posts.PostModel) error {
	query := `INSERT INTO posts (user_id, post_title, post_content, post_hashtags, created_at, updated_at,
				created_by, updated_by) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := r.DB.ExecContext(ctx, query, model.UserId, model.PostTitle, model.PostContent, model.PostHashtags, model.CreatedAt, model.UpdatedAt,
		model.CreatedBy, model.UpdatedBy)
	if err != nil {
		return err
	}

	return nil
}

func (r *postsRepository) GetAllPost(ctx context.Context, limit, offset int) (posts.GetAllPostResponse, error) {
	var response posts.GetAllPostResponse
	query := `SELECT p.id, p.user_id, u.username, p.post_title, p.post_content, p.post_hashtags, p.created_at, p.updated_at
				FROM posts as p
				JOIN users as u ON p.user_id = u.id LIMIT ? OFFSET ?`
	rows, err := r.DB.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return response, fmt.Errorf("Repository GetAllPost: %w", err)
	}
	defer rows.Close()

	data := make([]posts.Data, 0) //	Buat slice kosong
	for rows.Next() {
		var username string
		var model posts.PostModel

		err = rows.Scan(
			&model.ID,
			&model.UserId,
			&username,
			&model.PostTitle,
			&model.PostContent,
			&model.PostHashtags,
			&model.CreatedAt,
			&model.UpdatedAt,
		)
		if err != nil {
			return response, fmt.Errorf("Repository GetAllPost: %w", err)
		}

		data = append(data, posts.Data{
			ID:           model.ID,
			UserId:       model.UserId,
			Username:     username,
			PostTitle:    model.PostTitle,
			PostContent:  model.PostContent,
			PostHashtags: strings.Split(model.PostHashtags, ","),
			UpdatedAt:    model.UpdatedAt,
			CreatedAt:    model.CreatedAt,
		})
	}
	response.Data = data
	response.Pagination = posts.Pagination{
		Limit:  limit,
		Offset: offset,
	}

	return response, nil
}
