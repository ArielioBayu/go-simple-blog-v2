package comment

import (
	"context"
	"database/sql"
	"fmt"
)

type CommentRepository interface {
	CreateComment(ctx context.Context, model CommentModel) error
	GetCommentById(ctx context.Context, postId int) ([]GetComment, error)
	CountCommentsByPostID(ctx context.Context, postId int) (int, error)
}

type commentRepository struct {
	db *sql.DB
}

func NewCommentRepository(db *sql.DB) CommentRepository {
	return &commentRepository{db: db}
}

func (r *commentRepository) GetCommentById(ctx context.Context, postId int) ([]GetComment, error) {
	query := `SELECT c.id, c.user_id, u.username, c.comment_content
				FROM comments as c
				JOIN users as u ON c.user_id = u.id
				WHERE c.post_id = ?`
	rows, err := r.db.QueryContext(ctx, query, postId)
	if err != nil {
		return nil, fmt.Errorf("repository GetCommentById: %w", err)
	}
	defer rows.Close()

	response := make([]GetComment, 0)
	for rows.Next() {
		var model GetComment
		var username string

		err := rows.Scan(
			&model.ID,
			&model.UserId,
			&username,
			&model.CommentContent,
		)
		if err != nil {
			return nil, fmt.Errorf("repository GetCommentById: %w", err)
		}

		response = append(response, GetComment{
			ID:             model.ID,
			UserId:         model.UserId,
			Username:       username,
			CommentContent: model.CommentContent,
		})
	}
	return response, nil
}

func (r *commentRepository) CreateComment(ctx context.Context, model CommentModel) error {
	query := `INSERT INTO comments (post_id, user_id, comment_content, created_at, updated_at, created_by, updated_by)
				VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, model.PostId, model.UserId, model.CommentContent, model.CreatedAt, model.UpdatedAt,
		model.CreatedBy, model.UpdatedBy)
	if err != nil {
		return fmt.Errorf("repository CreateComment: %w", err)
	}

	return nil
}

func (r *commentRepository) CountCommentsByPostID(ctx context.Context, postId int) (int, error) {
	query := `SELECT COUNT(id) FROM comments WHERE post_id = ?`
	row := r.db.QueryRowContext(ctx, query, postId)

	var count int
	err := row.Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("repository CountCommentsByPostID: %w", err)
	}

	return count, nil
}
