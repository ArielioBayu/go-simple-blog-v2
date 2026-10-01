package comment

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ArielioBayu/go-simple-blog-v2/pkg/utils"
)

type CommentRepository interface {
	CreateComment(ctx context.Context, model CommentModel) error
	GetCommentById(ctx context.Context, postId int) ([]GetComment, error)
	GetCommentsByPostID(ctx context.Context, postId, limit, offset int) ([]GetComment, int, error)
	CheckCommentExists(ctx context.Context, commentId int) (bool, error)
	CountCommentsByPostID(ctx context.Context, postId int) (int, error)

	CreateReply(ctx context.Context, model CommentReplyModel) error
	GetRepliesByCommentIDs(ctx context.Context, commentIDs []int) (map[int][]GetCommentReply, error)
	GetRepliesByCommentID(ctx context.Context, commentId, limit, offset int) ([]GetCommentReply, int, error)

	GetCommentAuthInfo(ctx context.Context, commentId int) (*CommentAuthInfo, error)
	UpdateComment(ctx context.Context, commentId int, content string) error
	DeleteComment(ctx context.Context, commentId int) error

	GetReplyAuthInfo(ctx context.Context, replyId int) (*ReplyAuthInfo, error)
	UpdateReply(ctx context.Context, replyId int, content string) error
	DeleteReply(ctx context.Context, replyId int) error
}

type commentRepository struct {
	db *sql.DB
}

func NewCommentRepository(db *sql.DB) CommentRepository {
	return &commentRepository{db: db}
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

func (r *commentRepository) CheckCommentExists(ctx context.Context, commentId int) (bool, error) {
	query := `SELECT 1 FROM comments WHERE id = ? LIMIT 1`
	var exists int
	err := r.db.QueryRowContext(ctx, query, commentId).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("repository CheckCommentExists: %w", err)
	}
	return true, nil
}

func (r *commentRepository) GetCommentById(ctx context.Context, postId int) ([]GetComment, error) {
	query := `SELECT c.id, c.post_id, c.user_id, u.username, COALESCE(u.avatar_url, ''), c.comment_content, c.created_at
				FROM comments as c
				JOIN users as u ON c.user_id = u.id
				WHERE c.post_id = ?
				ORDER BY c.created_at ASC`
	rows, err := r.db.QueryContext(ctx, query, postId)
	if err != nil {
		return nil, fmt.Errorf("repository GetCommentById: %w", err)
	}
	defer rows.Close()

	comments := make([]GetComment, 0)
	commentIDs := make([]int, 0)

	for rows.Next() {
		var (
			id, postID, userID int
			username, avatar   string
			content            string
			createdAt          time.Time
		)

		err := rows.Scan(
			&id,
			&postID,
			&userID,
			&username,
			&avatar,
			&content,
			&createdAt,
		)
		if err != nil {
			return nil, fmt.Errorf("repository GetCommentById scan: %w", err)
		}

		commentIDs = append(commentIDs, id)
		comments = append(comments, GetComment{
			ID:             id,
			PostId:         postID,
			UserId:         userID,
			Username:       username,
			AvatarURL:      avatar,
			CommentContent: content,
			CreatedAt:      utils.JsonTime(createdAt),
			RepliesCount:   0,
			Replies:        []GetCommentReply{},
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("repository GetCommentById rows: %w", err)
	}

	if len(commentIDs) > 0 {
		repliesMap, err := r.GetRepliesByCommentIDs(ctx, commentIDs)
		if err != nil {
			return nil, fmt.Errorf("repository GetCommentById replies: %w", err)
		}

		for i := range comments {
			if replies, ok := repliesMap[comments[i].ID]; ok {
				comments[i].Replies = replies
				comments[i].RepliesCount = len(replies)
			}
		}
	}

	return comments, nil
}

func (r *commentRepository) GetCommentsByPostID(ctx context.Context, postId, limit, offset int) ([]GetComment, int, error) {
	countQuery := `SELECT COUNT(id) FROM comments WHERE post_id = ?`
	var totalData int
	if err := r.db.QueryRowContext(ctx, countQuery, postId).Scan(&totalData); err != nil {
		return nil, 0, fmt.Errorf("GetCommentsByPostID count: %w", err)
	}

	if totalData == 0 {
		return []GetComment{}, 0, nil
	}

	query := `SELECT c.id, c.post_id, c.user_id, u.username, COALESCE(u.avatar_url, ''), c.comment_content, c.created_at
				FROM comments as c
				JOIN users as u ON c.user_id = u.id
				WHERE c.post_id = ?
				ORDER BY c.created_at ASC
				LIMIT ? OFFSET ?`
	rows, err := r.db.QueryContext(ctx, query, postId, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("repository GetCommentsByPostID: %w", err)
	}
	defer rows.Close()

	comments := make([]GetComment, 0)
	commentIDs := make([]int, 0)

	for rows.Next() {
		var (
			id, postID, userID int
			username, avatar   string
			content            string
			createdAt          time.Time
		)

		err := rows.Scan(
			&id,
			&postID,
			&userID,
			&username,
			&avatar,
			&content,
			&createdAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("repository GetCommentsByPostID scan: %w", err)
		}

		commentIDs = append(commentIDs, id)
		comments = append(comments, GetComment{
			ID:             id,
			PostId:         postID,
			UserId:         userID,
			Username:       username,
			AvatarURL:      avatar,
			CommentContent: content,
			CreatedAt:      utils.JsonTime(createdAt),
			RepliesCount:   0,
			Replies:        []GetCommentReply{},
		})
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("repository GetCommentsByPostID rows: %w", err)
	}

	if len(commentIDs) > 0 {
		repliesMap, err := r.GetRepliesByCommentIDs(ctx, commentIDs)
		if err != nil {
			return nil, 0, fmt.Errorf("repository GetCommentsByPostID replies: %w", err)
		}

		for i := range comments {
			if replies, ok := repliesMap[comments[i].ID]; ok {
				comments[i].Replies = replies
				comments[i].RepliesCount = len(replies)
			}
		}
	}

	return comments, totalData, nil
}

func (r *commentRepository) CreateReply(ctx context.Context, model CommentReplyModel) error {
	query := `INSERT INTO comment_replies (comment_id, user_id, reply_to_user_id, reply_content, created_at, updated_at, created_by, updated_by)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, model.CommentID, model.UserID, model.ReplyToUserID, model.ReplyContent, model.CreatedAt, model.UpdatedAt,
		model.CreatedBy, model.UpdatedBy)
	if err != nil {
		return fmt.Errorf("repository CreateReply: %w", err)
	}

	return nil
}

func (r *commentRepository) GetRepliesByCommentIDs(ctx context.Context, commentIDs []int) (map[int][]GetCommentReply, error) {
	result := make(map[int][]GetCommentReply)
	if len(commentIDs) == 0 {
		return result, nil
	}

	placeholders := make([]string, len(commentIDs))
	args := make([]interface{}, len(commentIDs))
	for i, id := range commentIDs {
		placeholders[i] = "?"
		args[i] = id
	}

	query := fmt.Sprintf(`SELECT cr.id, cr.comment_id, cr.user_id, u.username, COALESCE(u.avatar_url, ''),
				cr.reply_to_user_id, COALESCE(ru.username, ''), cr.reply_content, cr.created_at
			FROM comment_replies cr
			JOIN users u ON cr.user_id = u.id
			LEFT JOIN users ru ON cr.reply_to_user_id = ru.id
			WHERE cr.comment_id IN (%s)
			ORDER BY cr.created_at ASC`, strings.Join(placeholders, ","))

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("GetRepliesByCommentIDs: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			reply            GetCommentReply
			createdAt        time.Time
			rawReplyToUserID sql.NullInt64
		)

		err := rows.Scan(
			&reply.ID,
			&reply.CommentID,
			&reply.UserID,
			&reply.Username,
			&reply.AvatarURL,
			&rawReplyToUserID,
			&reply.ReplyToUsername,
			&reply.ReplyContent,
			&createdAt,
		)
		if err != nil {
			return nil, fmt.Errorf("GetRepliesByCommentIDs scan: %w", err)
		}

		if rawReplyToUserID.Valid {
			reply.ReplyToUserID = &rawReplyToUserID.Int64
		}
		reply.CreatedAt = utils.JsonTime(createdAt)

		result[reply.CommentID] = append(result[reply.CommentID], reply)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetRepliesByCommentIDs rows: %w", err)
	}

	return result, nil
}

func (r *commentRepository) GetRepliesByCommentID(ctx context.Context, commentId, limit, offset int) ([]GetCommentReply, int, error) {
	countQuery := `SELECT COUNT(id) FROM comment_replies WHERE comment_id = ?`
	var totalData int
	if err := r.db.QueryRowContext(ctx, countQuery, commentId).Scan(&totalData); err != nil {
		return nil, 0, fmt.Errorf("GetRepliesByCommentID count: %w", err)
	}

	if totalData == 0 {
		return []GetCommentReply{}, 0, nil
	}

	query := `SELECT cr.id, cr.comment_id, cr.user_id, u.username, COALESCE(u.avatar_url, ''),
				cr.reply_to_user_id, COALESCE(ru.username, ''), cr.reply_content, cr.created_at
			FROM comment_replies cr
			JOIN users u ON cr.user_id = u.id
			LEFT JOIN users ru ON cr.reply_to_user_id = ru.id
			WHERE cr.comment_id = ?
			ORDER BY cr.created_at ASC
			LIMIT ? OFFSET ?`

	rows, err := r.db.QueryContext(ctx, query, commentId, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("GetRepliesByCommentID: %w", err)
	}
	defer rows.Close()

	replies := make([]GetCommentReply, 0)
	for rows.Next() {
		var (
			reply            GetCommentReply
			createdAt        time.Time
			rawReplyToUserID sql.NullInt64
		)

		err := rows.Scan(
			&reply.ID,
			&reply.CommentID,
			&reply.UserID,
			&reply.Username,
			&reply.AvatarURL,
			&rawReplyToUserID,
			&reply.ReplyToUsername,
			&reply.ReplyContent,
			&createdAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("GetRepliesByCommentID scan: %w", err)
		}

		if rawReplyToUserID.Valid {
			reply.ReplyToUserID = &rawReplyToUserID.Int64
		}
		reply.CreatedAt = utils.JsonTime(createdAt)

		replies = append(replies, reply)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("GetRepliesByCommentID rows: %w", err)
	}

	return replies, totalData, nil
}

func (r *commentRepository) CountCommentsByPostID(ctx context.Context, postId int) (int, error) {
	query := `SELECT 
				(SELECT COUNT(id) FROM comments WHERE post_id = ?) +
				(SELECT COUNT(cr.id) FROM comment_replies cr JOIN comments c ON cr.comment_id = c.id WHERE c.post_id = ?)`
	row := r.db.QueryRowContext(ctx, query, postId, postId)

	var count int
	err := row.Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("repository CountCommentsByPostID: %w", err)
	}

	return count, nil
}

func (r *commentRepository) GetCommentAuthInfo(ctx context.Context, commentId int) (*CommentAuthInfo, error) {
	query := `SELECT c.id, c.user_id, p.user_id 
				FROM comments c 
				JOIN posts p ON c.post_id = p.id 
				WHERE c.id = ?`
	var info CommentAuthInfo
	err := r.db.QueryRowContext(ctx, query, commentId).Scan(&info.CommentID, &info.CommentAuthorID, &info.PostAuthorID)
	if err != nil {
		return nil, err
	}
	return &info, nil
}

func (r *commentRepository) UpdateComment(ctx context.Context, commentId int, content string) error {
	query := `UPDATE comments SET comment_content = ?, updated_at = NOW() WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, content, commentId)
	if err != nil {
		return fmt.Errorf("repository UpdateComment: %w", err)
	}
	return nil
}

func (r *commentRepository) DeleteComment(ctx context.Context, commentId int) error {
	query := `DELETE FROM comments WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, commentId)
	if err != nil {
		return fmt.Errorf("repository DeleteComment: %w", err)
	}
	return nil
}

func (r *commentRepository) GetReplyAuthInfo(ctx context.Context, replyId int) (*ReplyAuthInfo, error) {
	query := `SELECT cr.id, cr.user_id, p.user_id 
				FROM comment_replies cr 
				JOIN comments c ON cr.comment_id = c.id 
				JOIN posts p ON c.post_id = p.id 
				WHERE cr.id = ?`
	var info ReplyAuthInfo
	err := r.db.QueryRowContext(ctx, query, replyId).Scan(&info.ReplyID, &info.ReplyAuthorID, &info.PostAuthorID)
	if err != nil {
		return nil, err
	}
	return &info, nil
}

func (r *commentRepository) UpdateReply(ctx context.Context, replyId int, content string) error {
	query := `UPDATE comment_replies SET reply_content = ?, updated_at = NOW() WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, content, replyId)
	if err != nil {
		return fmt.Errorf("repository UpdateReply: %w", err)
	}
	return nil
}

func (r *commentRepository) DeleteReply(ctx context.Context, replyId int) error {
	query := `DELETE FROM comment_replies WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, replyId)
	if err != nil {
		return fmt.Errorf("repository DeleteReply: %w", err)
	}
	return nil
}
