package notification

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/constants"
	"github.com/ArielioBayu/go-simple-blog-v2/pkg/utils"
)

type NotificationRepository interface {
	CreateNotification(ctx context.Context, input CreateNotificationInput) error
	DeleteNotification(ctx context.Context, actorID int64, notifType NotificationType, entityID int64) error
	DeleteNotificationsByEntity(ctx context.Context, notifType NotificationType, entityID int64) error
	GetNotifications(ctx context.Context, userID int64, limit, offset int) ([]NotificationItemResponse, int, error)
	CountUnread(ctx context.Context, userID int64) (int, error)
	MarkAllAsRead(ctx context.Context, userID int64) error
	MarkAsReadByID(ctx context.Context, userID int64, notificationID int64) error
}

type notificationRepository struct {
	db *sql.DB
}

func NewNotificationRepository(db *sql.DB) NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) CreateNotification(ctx context.Context, input CreateNotificationInput) error {
	query := `INSERT INTO notifications (user_id, actor_id, type, entity_id, entity_type, message, is_read, created_at)
			  VALUES (?, ?, ?, ?, ?, ?, FALSE, NOW())`
	_, err := r.db.ExecContext(ctx, query,
		input.UserID,
		input.ActorID,
		string(input.Type),
		input.EntityID,
		string(input.EntityType),
		input.Message,
	)
	if err != nil {
		return fmt.Errorf("repository CreateNotification: %w", err)
	}

	return nil
}

func (r *notificationRepository) DeleteNotification(ctx context.Context, actorID int64, notifType NotificationType, entityID int64) error {
	query := `DELETE FROM notifications WHERE actor_id = ? AND type = ? AND entity_id = ? ORDER BY id DESC LIMIT 1`
	_, err := r.db.ExecContext(ctx, query, actorID, string(notifType), entityID)
	if err != nil {
		return fmt.Errorf("repository DeleteNotification: %w", err)
	}

	return nil
}

func (r *notificationRepository) DeleteNotificationsByEntity(ctx context.Context, notifType NotificationType, entityID int64) error {
	query := `DELETE FROM notifications WHERE type = ? AND entity_id = ?`
	_, err := r.db.ExecContext(ctx, query, string(notifType), entityID)
	if err != nil {
		return fmt.Errorf("repository DeleteNotificationsByEntity: %w", err)
	}

	return nil
}

func (r *notificationRepository) CountUnread(ctx context.Context, userID int64) (int, error) {
	query := `SELECT COUNT(1) FROM notifications WHERE user_id = ? AND is_read = FALSE`
	var count int
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("repository CountUnread: %w", err)
	}

	return count, nil
}

func (r *notificationRepository) MarkAllAsRead(ctx context.Context, userID int64) error {
	query := `UPDATE notifications SET is_read = TRUE, read_at = NOW() WHERE user_id = ? AND is_read = FALSE`
	_, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("repository MarkAllAsRead: %w", err)
	}

	return nil
}

func (r *notificationRepository) MarkAsReadByID(ctx context.Context, userID int64, notificationID int64) error {
	query := `UPDATE notifications SET is_read = TRUE, read_at = NOW() WHERE id = ? AND user_id = ?`
	res, err := r.db.ExecContext(ctx, query, notificationID, userID)
	if err != nil {
		return fmt.Errorf("repository MarkAsReadByID: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("repository MarkAsReadByID rowsAffected: %w", err)
	}

	if rowsAffected == 0 {
		return constants.ErrNotificationNotFound
	}

	return nil
}

func (r *notificationRepository) GetNotifications(ctx context.Context, userID int64, limit, offset int) ([]NotificationItemResponse, int, error) {
	countQuery := `SELECT COUNT(1) FROM notifications WHERE user_id = ?`
	var totalData int
	if err := r.db.QueryRowContext(ctx, countQuery, userID).Scan(&totalData); err != nil {
		return nil, 0, fmt.Errorf("repository GetNotifications count: %w", err)
	}

	if totalData == 0 {
		return []NotificationItemResponse{}, 0, nil
	}

	query := `SELECT n.id, n.actor_id, u.username, COALESCE(u.avatar_url, ''),
					 n.type, n.entity_id, n.entity_type, n.message, n.is_read,
					 n.created_at, n.read_at
			  FROM notifications n
			  JOIN users u ON n.actor_id = u.id
			  WHERE n.user_id = ?
			  ORDER BY n.created_at DESC
			  LIMIT ? OFFSET ?`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("repository GetNotifications query: %w", err)
	}
	defer rows.Close()

	items := make([]NotificationItemResponse, 0, limit)
	for rows.Next() {
		var (
			item      NotificationItemResponse
			rawType   string
			rawEntity string
			createdAt time.Time
			readAt    sql.NullTime
		)

		err := rows.Scan(
			&item.ID,
			&item.Actor.ID,
			&item.Actor.Username,
			&item.Actor.AvatarURL,
			&rawType,
			&item.EntityID,
			&rawEntity,
			&item.Message,
			&item.IsRead,
			&createdAt,
			&readAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("repository GetNotifications scan: %w", err)
		}

		item.Type = NotificationType(rawType)
		item.EntityType = EntityType(rawEntity)
		item.CreatedAt = utils.JsonTime(createdAt)
		if readAt.Valid {
			rAt := utils.JsonTime(readAt.Time)
			item.ReadAt = &rAt
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("repository GetNotifications rows: %w", err)
	}

	return items, totalData, nil
}
