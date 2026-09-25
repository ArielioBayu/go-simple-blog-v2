package saves

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/modules/post"
	"github.com/ArielioBayu/go-simple-blog-v2/pkg/utils"
)

type SavesRepository interface {
	SavePost(ctx context.Context, userID, postID int) error
	UnsavePost(ctx context.Context, userID, postID int) error
	IsPostSaved(ctx context.Context, userID, postID int) (bool, error)
	CheckPostExists(ctx context.Context, postID int) (bool, error)
	GetSavedPosts(ctx context.Context, userID, limit, offset int) ([]SavedPostItem, int, error)
	GetSavedPostIDs(ctx context.Context, userID int) ([]int, error)
}

type savesRepository struct {
	db *sql.DB
}

func NewSavesRepository(db *sql.DB) SavesRepository {
	return &savesRepository{db: db}
}

func (r *savesRepository) CheckPostExists(ctx context.Context, postID int) (bool, error) {
	query := `SELECT 1 FROM posts WHERE id = ? LIMIT 1`
	var exists int
	err := r.db.QueryRowContext(ctx, query, postID).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("repository CheckPostExists: %w", err)
	}
	return true, nil
}

func (r *savesRepository) SavePost(ctx context.Context, userID, postID int) error {
	query := `INSERT INTO user_saved_posts (user_id, post_id, created_at)
	          VALUES (?, ?, ?)
	          ON DUPLICATE KEY UPDATE created_at = VALUES(created_at)`
	now := time.Now()
	_, err := r.db.ExecContext(ctx, query, userID, postID, now)
	if err != nil {
		return fmt.Errorf("repository SavePost: %w", err)
	}
	return nil
}

func (r *savesRepository) UnsavePost(ctx context.Context, userID, postID int) error {
	query := `DELETE FROM user_saved_posts WHERE user_id = ? AND post_id = ?`
	_, err := r.db.ExecContext(ctx, query, userID, postID)
	if err != nil {
		return fmt.Errorf("repository UnsavePost: %w", err)
	}
	return nil
}

func (r *savesRepository) IsPostSaved(ctx context.Context, userID, postID int) (bool, error) {
	query := `SELECT 1 FROM user_saved_posts WHERE user_id = ? AND post_id = ? LIMIT 1`
	var exists int
	err := r.db.QueryRowContext(ctx, query, userID, postID).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("repository IsPostSaved: %w", err)
	}
	return true, nil
}

func (r *savesRepository) GetSavedPostIDs(ctx context.Context, userID int) ([]int, error) {
	query := `SELECT post_id FROM user_saved_posts WHERE user_id = ? ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("repository GetSavedPostIDs: %w", err)
	}
	defer rows.Close()

	ids := make([]int, 0)
	for rows.Next() {
		var postID int
		if err := rows.Scan(&postID); err != nil {
			return nil, fmt.Errorf("repository GetSavedPostIDs scan: %w", err)
		}
		ids = append(ids, postID)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("repository GetSavedPostIDs rows: %w", err)
	}

	return ids, nil
}

func (r *savesRepository) GetSavedPosts(ctx context.Context, userID, limit, offset int) ([]SavedPostItem, int, error) {
	countQuery := `SELECT COUNT(id) FROM user_saved_posts WHERE user_id = ?`
	var totalData int
	if err := r.db.QueryRowContext(ctx, countQuery, userID).Scan(&totalData); err != nil {
		return nil, 0, fmt.Errorf("repository GetSavedPosts count: %w", err)
	}

	if totalData == 0 {
		return []SavedPostItem{}, 0, nil
	}

	query := `SELECT 
				p.id, p.user_id, u.username, COALESCE(u.avatar_url, ''),
				p.post_title, p.post_content, p.post_hashtags,
				COALESCE(act.is_liked, false), p.created_at, p.updated_at,
				COALESCE(up.file_path, ''), COALESCE(up.file_type, ''), COALESCE(up.file_size, 0),
				COALESCE(p.upload_id, up.id, 0),
				usp.created_at
			FROM user_saved_posts usp
			JOIN posts p ON usp.post_id = p.id
			JOIN users u ON p.user_id = u.id
			LEFT JOIN activities act ON p.id = act.post_id AND act.user_id = ?
			LEFT JOIN uploads up ON (p.upload_id = up.id OR (p.upload_id IS NULL AND p.post_content LIKE CONCAT('%', up.system_filename, '%')))
			WHERE usp.user_id = ?
			ORDER BY usp.created_at DESC
			LIMIT ? OFFSET ?`

	rows, err := r.db.QueryContext(ctx, query, userID, userID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("repository GetSavedPosts: %w", err)
	}
	defer rows.Close()

	items := make([]SavedPostItem, 0)
	for rows.Next() {
		var (
			postID             int
			postUserID         int
			username           string
			avatarURL          string
			postTitle          string
			postContent        string
			rawHashtags        string
			isLiked            bool
			createdAt          time.Time
			updatedAt          time.Time
			filePath, fileType string
			fileSize           int64
			uploadID           int64
			savedAt            time.Time
		)

		err := rows.Scan(
			&postID,
			&postUserID,
			&username,
			&avatarURL,
			&postTitle,
			&postContent,
			&rawHashtags,
			&isLiked,
			&createdAt,
			&updatedAt,
			&filePath,
			&fileType,
			&fileSize,
			&uploadID,
			&savedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("repository GetSavedPosts scan: %w", err)
		}

		var uploadIDPtr *int64
		if uploadID > 0 {
			uploadIDPtr = &uploadID
		}

		var hashtags []string
		if strings.TrimSpace(rawHashtags) != "" {
			hashtags = strings.Split(rawHashtags, ",")
		} else {
			hashtags = []string{}
		}

		items = append(items, SavedPostItem{
			ID:           postID,
			UserID:       postUserID,
			Username:     username,
			AvatarURL:    avatarURL,
			PostTitle:    postTitle,
			PostContent:  postContent,
			PostHashtags: hashtags,
			IsLiked:      isLiked,
			IsSaved:      true,
			UploadID:     uploadIDPtr,
			FilePath:     filePath,
			Filepath:     filePath,
			FileType:     fileType,
			FileSize:     fileSize,
			Media:        []post.PostMedia{},
			CreatedAt:    utils.JsonTime(createdAt),
			UpdatedAt:    utils.JsonTime(updatedAt),
			SavedAt:      utils.JsonTime(savedAt),
		})
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("repository GetSavedPosts rows: %w", err)
	}

	if len(items) > 0 {
		postIDs := make([]int, len(items))
		for i, it := range items {
			postIDs[i] = it.ID
		}

		mediaMap, err := r.getMediaByPostIDs(ctx, postIDs)
		if err != nil {
			return nil, 0, fmt.Errorf("repository GetSavedPosts media: %w", err)
		}

		for i := range items {
			if medias, ok := mediaMap[items[i].ID]; ok && len(medias) > 0 {
				items[i].Media = medias
			} else if items[i].UploadID != nil && *items[i].UploadID > 0 {
				items[i].Media = []post.PostMedia{
					{
						PostID:    items[i].ID,
						UploadID:  *items[i].UploadID,
						FilePath:  items[i].FilePath,
						FileType:  items[i].FileType,
						FileSize:  items[i].FileSize,
						SortOrder: 0,
						CreatedAt: items[i].CreatedAt,
					},
				}
			} else if items[i].FilePath != "" {
				items[i].Media = []post.PostMedia{
					{
						PostID:    items[i].ID,
						FilePath:  items[i].FilePath,
						FileType:  items[i].FileType,
						FileSize:  items[i].FileSize,
						SortOrder: 0,
						CreatedAt: items[i].CreatedAt,
					},
				}
			}
		}
	}

	return items, totalData, nil
}

func (r *savesRepository) getMediaByPostIDs(ctx context.Context, postIDs []int) (map[int][]post.PostMedia, error) {
	result := make(map[int][]post.PostMedia)
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
			media     post.PostMedia
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
