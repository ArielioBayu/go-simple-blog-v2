package posts

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/constants"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/model/posts"
)

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

func (r *postsRepository) GetAllPost(ctx context.Context, limit, offset, userID int) (posts.GetAllPostResponse, error) {
	var response posts.GetAllPostResponse
	query := `SELECT p.id, p.user_id, u.username, p.post_title, p.post_content, p.post_hashtags, COALESCE(act.is_liked, false), p.created_at, p.updated_at
				FROM posts as p
				JOIN users as u ON p.user_id = u.id
				LEFT JOIN activities as act ON p.id = act.post_id AND act.user_id = ?
				ORDER BY p.created_at DESC 
				LIMIT ? OFFSET ?`
	rows, err := r.DB.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return response, fmt.Errorf("Repository GetAllPost: %w", err)
	}
	defer rows.Close()

	data := make([]posts.Data, 0) //	Buat slice kosong
	for rows.Next() {
		var username string
		var isliked bool
		var model posts.PostModel

		err = rows.Scan(
			&model.ID,
			&model.UserId,
			&username,
			&model.PostTitle,
			&model.PostContent,
			&model.PostHashtags,
			&isliked,
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
			IsLiked:      isliked,
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

func (r *postsRepository) GetPostById(ctx context.Context, id, userID int) (*posts.Data, error) {
	query := `SELECT p.id, p.user_id, u.username, p.post_title, p.post_content, p.post_hashtags, COALESCE(act.is_liked, false), p.created_at, p.updated_at 
				FROM posts as p
				JOIN users as u ON p.user_id = u.id
				LEFT JOIN activities as act ON p.id = act.post_id AND act.user_id = ?
				WHERE p.id = ?
				LIMIT 1`

	row := r.DB.QueryRowContext(ctx, query, userID, id)
	var (
		model    posts.PostModel
		username string
		isliked  bool
	)
	err := row.Scan(
		&model.ID,
		&model.UserId,
		&username,
		&model.PostTitle,
		&model.PostContent,
		&model.PostHashtags,
		&isliked,
		&model.CreatedAt,
		&model.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, constants.ErrPostNotFound
		}
		return nil, fmt.Errorf("Repository GetPostById: %w", err)
	}

	data := posts.Data{
		ID:           model.ID,
		UserId:       model.UserId,
		Username:     username,
		PostTitle:    model.PostTitle,
		PostContent:  model.PostContent,
		PostHashtags: strings.Split(model.PostHashtags, ","),
		IsLiked:      isliked,
		CreatedAt:    model.CreatedAt,
		UpdatedAt:    model.UpdatedAt,
	}

	return &data, nil
}
