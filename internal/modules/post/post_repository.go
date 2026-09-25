package post

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/constants"
	"github.com/ArielioBayu/go-simple-blog-v2/pkg/utils"
)

type PostRepository interface {
	CreatePost(ctx context.Context, model PostModel, uploadIDs []int64) error
	GetAllPost(ctx context.Context, limit, offset, userID int) (GetAllPostResponse, error)
	GetPostsByUserID(ctx context.Context, targetUserID, currentUserID, limit, offset int) (GetAllPostResponse, error)
	GetPostById(ctx context.Context, id, userID int) (*Data, error)
	DeletePost(ctx context.Context, id int) error
	CheckPostExists(ctx context.Context, id int) (bool, error)
}

type postRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) PostRepository {
	return &postRepository{db: db}
}

func (r *postRepository) CreatePost(ctx context.Context, model PostModel, uploadIDs []int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("repository CreatePost begin tx: %w", err)
	}
	defer func() {
		if tx != nil {
			_ = tx.Rollback()
		}
	}()

	query := `INSERT INTO posts (user_id, post_title, post_content, post_hashtags, upload_id, created_at, updated_at,
				created_by, updated_by) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	res, err := tx.ExecContext(ctx, query, model.UserId, model.PostTitle, model.PostContent, model.PostHashtags, model.UploadID, model.CreatedAt, model.UpdatedAt,
		model.CreatedBy, model.UpdatedBy)
	if err != nil {
		return fmt.Errorf("repository CreatePost insert post: %w", err)
	}

	postID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("repository CreatePost last insert id: %w", err)
	}

	if len(uploadIDs) > 0 {
		mediaQuery := `INSERT INTO post_media (post_id, upload_id, sort_order) VALUES (?, ?, ?)`
		mediaStmt, err := tx.PrepareContext(ctx, mediaQuery)
		if err != nil {
			return fmt.Errorf("repository CreatePost prepare media: %w", err)
		}
		defer mediaStmt.Close()

		for idx, uploadID := range uploadIDs {
			_, err = mediaStmt.ExecContext(ctx, postID, uploadID, idx)
			if err != nil {
				return fmt.Errorf("repository CreatePost insert media: %w", err)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("repository CreatePost commit: %w", err)
	}
	tx = nil

	return nil
}

func (r *postRepository) getMediaByPostIDs(ctx context.Context, postIDs []int) (map[int][]PostMedia, error) {
	result := make(map[int][]PostMedia)
	if len(postIDs) == 0 {
		return result, nil
	}

	placeholders := make([]string, len(postIDs))
	args := make([]interface{}, len(postIDs))
	for i, id := range postIDs {
		placeholders[i] = "?"
		args[i] = id
	}

	query := fmt.Sprintf(`SELECT pm.id, pm.post_id, pm.upload_id, pm.sort_order,
				COALESCE(up.file_path, ''), COALESCE(up.file_type, ''), COALESCE(up.file_size, 0),
				pm.created_at
			FROM post_media pm
			JOIN uploads up ON pm.upload_id = up.id
			WHERE pm.post_id IN (%s)
			ORDER BY pm.sort_order ASC, pm.id ASC`, strings.Join(placeholders, ","))

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("getMediaByPostIDs: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			media     PostMedia
			postID    int
			createdAt time.Time
		)
		if err := rows.Scan(
			&media.ID,
			&postID,
			&media.UploadID,
			&media.SortOrder,
			&media.FilePath,
			&media.FileType,
			&media.FileSize,
			&createdAt,
		); err != nil {
			return nil, fmt.Errorf("getMediaByPostIDs scan: %w", err)
		}
		media.PostID = postID
		media.CreatedAt = utils.JsonTime(createdAt)

		result[postID] = append(result[postID], media)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("getMediaByPostIDs rows: %w", err)
	}

	return result, nil
}

func (r *postRepository) GetAllPost(ctx context.Context, limit, offset, userID int) (GetAllPostResponse, error) {
	var response GetAllPostResponse
	query := `SELECT p.id, p.user_id, u.username, COALESCE(u.avatar_url, ''), p.post_title, p.post_content, p.post_hashtags, 
				COALESCE(act.is_liked, false), (usp.id IS NOT NULL), p.created_at, p.updated_at,
				COALESCE(up.file_path, ''), COALESCE(up.file_type, ''), COALESCE(up.file_size, 0),
				COALESCE(p.upload_id, up.id, 0)
				FROM posts as p
				JOIN users as u ON p.user_id = u.id
				LEFT JOIN activities as act ON p.id = act.post_id AND act.user_id = ?
				LEFT JOIN user_saved_posts as usp ON p.id = usp.post_id AND usp.user_id = ?
				LEFT JOIN uploads as up ON (p.upload_id = up.id OR (p.upload_id IS NULL AND p.post_content LIKE CONCAT('%', up.system_filename, '%')))
				ORDER BY p.created_at DESC 
				LIMIT ? OFFSET ?`
	rows, err := r.db.QueryContext(ctx, query, userID, userID, limit, offset)
	if err != nil {
		return response, fmt.Errorf("repository GetAllPost: %w", err)
	}
	defer rows.Close()

	data := make([]Data, 0)
	postIDs := make([]int, 0)
	for rows.Next() {
		var username, avatarURL string
		var isliked, issaved bool
		var model PostModel
		var filePath, fileType string
		var fileSize int64
		var uploadID int64

		err = rows.Scan(
			&model.ID,
			&model.UserId,
			&username,
			&avatarURL,
			&model.PostTitle,
			&model.PostContent,
			&model.PostHashtags,
			&isliked,
			&issaved,
			&model.CreatedAt,
			&model.UpdatedAt,
			&filePath,
			&fileType,
			&fileSize,
			&uploadID,
		)
		if err != nil {
			return response, fmt.Errorf("repository GetAllPost: %w", err)
		}

		var uploadIDPtr *int64
		if uploadID > 0 {
			uploadIDPtr = &uploadID
		}

		postIDs = append(postIDs, model.ID)
		data = append(data, Data{
			ID:           model.ID,
			UserId:       model.UserId,
			Username:     username,
			AvatarURL:    avatarURL,
			PostTitle:    model.PostTitle,
			PostContent:  model.PostContent,
			PostHashtags: strings.Split(model.PostHashtags, ","),
			IsLiked:      isliked,
			IsSaved:      issaved,
			UploadID:     uploadIDPtr,
			FilePath:     filePath,
			Filepath:     filePath,
			FileType:     fileType,
			FileSize:     fileSize,
			Media:        []PostMedia{},
			UpdatedAt:    model.UpdatedAt,
			CreatedAt:    model.CreatedAt,
		})
	}

	if len(postIDs) > 0 {
		mediaMap, err := r.getMediaByPostIDs(ctx, postIDs)
		if err != nil {
			return response, fmt.Errorf("repository GetAllPost media: %w", err)
		}
		for i := range data {
			if medias, ok := mediaMap[data[i].ID]; ok && len(medias) > 0 {
				data[i].Media = medias
			} else if data[i].UploadID != nil && *data[i].UploadID > 0 {
				data[i].Media = []PostMedia{
					{
						PostID:    data[i].ID,
						UploadID:  *data[i].UploadID,
						FilePath:  data[i].FilePath,
						FileType:  data[i].FileType,
						FileSize:  data[i].FileSize,
						SortOrder: 0,
						CreatedAt: data[i].CreatedAt,
					},
				}
			} else if data[i].FilePath != "" {
				data[i].Media = []PostMedia{
					{
						PostID:    data[i].ID,
						FilePath:  data[i].FilePath,
						FileType:  data[i].FileType,
						FileSize:  data[i].FileSize,
						SortOrder: 0,
						CreatedAt: data[i].CreatedAt,
					},
				}
			}
		}
	}

	response.Data = data
	response.Pagination = Pagination{
		Limit:  limit,
		Offset: offset,
	}

	return response, nil
}

func (r *postRepository) GetPostsByUserID(ctx context.Context, targetUserID, currentUserID, limit, offset int) (GetAllPostResponse, error) {
	var response GetAllPostResponse
	query := `SELECT p.id, p.user_id, u.username, COALESCE(u.avatar_url, ''), p.post_title, p.post_content, p.post_hashtags, 
				COALESCE(act.is_liked, false), (usp.id IS NOT NULL), p.created_at, p.updated_at,
				COALESCE(up.file_path, ''), COALESCE(up.file_type, ''), COALESCE(up.file_size, 0),
				COALESCE(p.upload_id, up.id, 0)
				FROM posts as p
				JOIN users as u ON p.user_id = u.id
				LEFT JOIN activities as act ON p.id = act.post_id AND act.user_id = ?
				LEFT JOIN user_saved_posts as usp ON p.id = usp.post_id AND usp.user_id = ?
				LEFT JOIN uploads as up ON (p.upload_id = up.id OR (p.upload_id IS NULL AND p.post_content LIKE CONCAT('%', up.system_filename, '%')))
				WHERE p.user_id = ?
				ORDER BY p.created_at DESC 
				LIMIT ? OFFSET ?`
	rows, err := r.db.QueryContext(ctx, query, currentUserID, currentUserID, targetUserID, limit, offset)
	if err != nil {
		return response, fmt.Errorf("repository GetPostsByUserID: %w", err)
	}
	defer rows.Close()

	data := make([]Data, 0)
	postIDs := make([]int, 0)
	for rows.Next() {
		var username, avatarURL string
		var isliked, issaved bool
		var model PostModel
		var filePath, fileType string
		var fileSize int64
		var uploadID int64

		err = rows.Scan(
			&model.ID,
			&model.UserId,
			&username,
			&avatarURL,
			&model.PostTitle,
			&model.PostContent,
			&model.PostHashtags,
			&isliked,
			&issaved,
			&model.CreatedAt,
			&model.UpdatedAt,
			&filePath,
			&fileType,
			&fileSize,
			&uploadID,
		)
		if err != nil {
			return response, fmt.Errorf("repository GetPostsByUserID: %w", err)
		}

		var uploadIDPtr *int64
		if uploadID > 0 {
			uploadIDPtr = &uploadID
		}

		postIDs = append(postIDs, model.ID)
		data = append(data, Data{
			ID:           model.ID,
			UserId:       model.UserId,
			Username:     username,
			AvatarURL:    avatarURL,
			PostTitle:    model.PostTitle,
			PostContent:  model.PostContent,
			PostHashtags: strings.Split(model.PostHashtags, ","),
			IsLiked:      isliked,
			IsSaved:      issaved,
			UploadID:     uploadIDPtr,
			FilePath:     filePath,
			Filepath:     filePath,
			FileType:     fileType,
			FileSize:     fileSize,
			Media:        []PostMedia{},
			UpdatedAt:    model.UpdatedAt,
			CreatedAt:    model.CreatedAt,
		})
	}

	if len(postIDs) > 0 {
		mediaMap, err := r.getMediaByPostIDs(ctx, postIDs)
		if err != nil {
			return response, fmt.Errorf("repository GetPostsByUserID media: %w", err)
		}
		for i := range data {
			if medias, ok := mediaMap[data[i].ID]; ok && len(medias) > 0 {
				data[i].Media = medias
			} else if data[i].UploadID != nil && *data[i].UploadID > 0 {
				data[i].Media = []PostMedia{
					{
						PostID:    data[i].ID,
						UploadID:  *data[i].UploadID,
						FilePath:  data[i].FilePath,
						FileType:  data[i].FileType,
						FileSize:  data[i].FileSize,
						SortOrder: 0,
						CreatedAt: data[i].CreatedAt,
					},
				}
			} else if data[i].FilePath != "" {
				data[i].Media = []PostMedia{
					{
						PostID:    data[i].ID,
						FilePath:  data[i].FilePath,
						FileType:  data[i].FileType,
						FileSize:  data[i].FileSize,
						SortOrder: 0,
						CreatedAt: data[i].CreatedAt,
					},
				}
			}
		}
	}

	response.Data = data
	response.Pagination = Pagination{
		Limit:  limit,
		Offset: offset,
	}

	return response, nil
}

func (r *postRepository) GetPostById(ctx context.Context, id, userID int) (*Data, error) {
	query := `SELECT p.id, p.user_id, u.username, COALESCE(u.avatar_url, ''), p.post_title, p.post_content, p.post_hashtags, 
				COALESCE(act.is_liked, false), (usp.id IS NOT NULL), p.created_at, p.updated_at,
				COALESCE(up.file_path, ''), COALESCE(up.file_type, ''), COALESCE(up.file_size, 0),
				COALESCE(p.upload_id, up.id, 0)
				FROM posts as p
				JOIN users as u ON p.user_id = u.id
				LEFT JOIN activities as act ON p.id = act.post_id AND act.user_id = ?
				LEFT JOIN user_saved_posts as usp ON p.id = usp.post_id AND usp.user_id = ?
				LEFT JOIN uploads as up ON (p.upload_id = up.id OR (p.upload_id IS NULL AND p.post_content LIKE CONCAT('%', up.system_filename, '%')))
				WHERE p.id = ?
				LIMIT 1`

	row := r.db.QueryRowContext(ctx, query, userID, userID, id)
	var (
		model                        PostModel
		username, avatarURL          string
		isliked, issaved             bool
		filePath, fileType           string
		fileSize                     int64
		uploadID                     int64
	)
	err := row.Scan(
		&model.ID,
		&model.UserId,
		&username,
		&avatarURL,
		&model.PostTitle,
		&model.PostContent,
		&model.PostHashtags,
		&isliked,
		&issaved,
		&model.CreatedAt,
		&model.UpdatedAt,
		&filePath,
		&fileType,
		&fileSize,
		&uploadID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, constants.ErrPostNotFound
		}
		return nil, fmt.Errorf("repository GetPostById: %w", err)
	}

	var uploadIDPtr *int64
	if uploadID > 0 {
		uploadIDPtr = &uploadID
	}

	data := Data{
		ID:           model.ID,
		UserId:       model.UserId,
		Username:     username,
		AvatarURL:    avatarURL,
		PostTitle:    model.PostTitle,
		PostContent:  model.PostContent,
		PostHashtags: strings.Split(model.PostHashtags, ","),
		IsLiked:      isliked,
		IsSaved:      issaved,
		UploadID:     uploadIDPtr,
		FilePath:     filePath,
		Filepath:     filePath,
		FileType:     fileType,
		FileSize:     fileSize,
		Media:        []PostMedia{},
		CreatedAt:    model.CreatedAt,
		UpdatedAt:    model.UpdatedAt,
	}

	mediaMap, err := r.getMediaByPostIDs(ctx, []int{id})
	if err != nil {
		return nil, fmt.Errorf("repository GetPostById media: %w", err)
	}

	if medias, ok := mediaMap[id]; ok && len(medias) > 0 {
		data.Media = medias
	} else if data.UploadID != nil && *data.UploadID > 0 {
		data.Media = []PostMedia{
			{
				PostID:    data.ID,
				UploadID:  *data.UploadID,
				FilePath:  data.FilePath,
				FileType:  data.FileType,
				FileSize:  data.FileSize,
				SortOrder: 0,
				CreatedAt: data.CreatedAt,
			},
		}
	} else if data.FilePath != "" {
		data.Media = []PostMedia{
			{
				PostID:    data.ID,
				FilePath:  data.FilePath,
				FileType:  data.FileType,
				FileSize:  data.FileSize,
				SortOrder: 0,
				CreatedAt: data.CreatedAt,
			},
		}
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

func (r *postRepository) DeletePost(ctx context.Context, id int) error {
	query := `DELETE FROM posts WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("repository DeletePost: %w", err)
	}
	return nil
}
