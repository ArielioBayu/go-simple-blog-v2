package post

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/constants"
)

type PostRepository interface {
	CreatePost(ctx context.Context, model PostModel) error
	GetAllPost(ctx context.Context, limit, offset, userID int) (GetAllPostResponse, error)
	GetPostById(ctx context.Context, id, userID int) (*Data, error)
	CheckPostExists(ctx context.Context, id int) (bool, error)
}

type postRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) PostRepository {
	return &postRepository{db: db}
}

func (r *postRepository) CreatePost(ctx context.Context, model PostModel) error {
	query := `INSERT INTO posts (user_id, post_title, post_content, post_hashtags, upload_id, created_at, updated_at,
				created_by, updated_by) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, model.UserId, model.PostTitle, model.PostContent, model.PostHashtags, model.UploadID, model.CreatedAt, model.UpdatedAt,
		model.CreatedBy, model.UpdatedBy)
	if err != nil {
		return fmt.Errorf("repository CreatePost: %w", err)
	}

	return nil
}

func (r *postRepository) GetAllPost(ctx context.Context, limit, offset, userID int) (GetAllPostResponse, error) {
	var response GetAllPostResponse
	query := `SELECT p.id, p.user_id, u.username, p.post_title, p.post_content, p.post_hashtags, 
				COALESCE(act.is_liked, false), p.created_at, p.updated_at,
				COALESCE(up.file_path, ''), COALESCE(up.file_type, ''), COALESCE(up.file_size, 0)
				FROM posts as p
				JOIN users as u ON p.user_id = u.id
				LEFT JOIN activities as act ON p.id = act.post_id AND act.user_id = ?
				LEFT JOIN uploads as up ON (p.upload_id = up.id OR (p.upload_id IS NULL AND p.post_content LIKE CONCAT('%', up.system_filename, '%')))
				ORDER BY p.created_at DESC 
				LIMIT ? OFFSET ?`
	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return response, fmt.Errorf("repository GetAllPost: %w", err)
	}
	defer rows.Close()

	data := make([]Data, 0)
	for rows.Next() {
		var username string
		var isliked bool
		var model PostModel
		var filePath, fileType string
		var fileSize int64

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
			&filePath,
			&fileType,
			&fileSize,
		)
		if err != nil {
			return response, fmt.Errorf("repository GetAllPost: %w", err)
		}

		data = append(data, Data{
			ID:           model.ID,
			UserId:       model.UserId,
			Username:     username,
			PostTitle:    model.PostTitle,
			PostContent:  model.PostContent,
			PostHashtags: strings.Split(model.PostHashtags, ","),
			IsLiked:      isliked,
			FilePath:     filePath,
			Filepath:     filePath,
			FileType:     fileType,
			FileSize:     fileSize,
			UpdatedAt:    model.UpdatedAt,
			CreatedAt:    model.CreatedAt,
		})
	}
	response.Data = data
	response.Pagination = Pagination{
		Limit:  limit,
		Offset: offset,
	}

	return response, nil
}

func (r *postRepository) GetPostById(ctx context.Context, id, userID int) (*Data, error) {
	query := `SELECT p.id, p.user_id, u.username, p.post_title, p.post_content, p.post_hashtags, 
				COALESCE(act.is_liked, false), p.created_at, p.updated_at,
				COALESCE(up.file_path, ''), COALESCE(up.file_type, ''), COALESCE(up.file_size, 0)
				FROM posts as p
				JOIN users as u ON p.user_id = u.id
				LEFT JOIN activities as act ON p.id = act.post_id AND act.user_id = ?
				LEFT JOIN uploads as up ON (p.upload_id = up.id OR (p.upload_id IS NULL AND p.post_content LIKE CONCAT('%', up.system_filename, '%')))
				WHERE p.id = ?
				LIMIT 1`

	row := r.db.QueryRowContext(ctx, query, userID, id)
	var (
		model              PostModel
		username           string
		isliked            bool
		filePath, fileType string
		fileSize           int64
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
		&filePath,
		&fileType,
		&fileSize,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, constants.ErrPostNotFound
		}
		return nil, fmt.Errorf("repository GetPostById: %w", err)
	}

	data := Data{
		ID:           model.ID,
		UserId:       model.UserId,
		Username:     username,
		PostTitle:    model.PostTitle,
		PostContent:  model.PostContent,
		PostHashtags: strings.Split(model.PostHashtags, ","),
		IsLiked:      isliked,
		FilePath:     filePath,
		Filepath:     filePath,
		FileType:     fileType,
		FileSize:     fileSize,
		CreatedAt:    model.CreatedAt,
		UpdatedAt:    model.UpdatedAt,
	}

	return &data, nil
}

func (r *postRepository) CheckPostExists(ctx context.Context, id int) (bool, error) {
	query := `SELECT 1 FROM posts WHERE id = ? LIMIT 1`
	var exists int
	err := r.db.QueryRowContext(ctx, query, id).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("repository CheckPostExists: %w", err)
	}
	return true, nil
}
