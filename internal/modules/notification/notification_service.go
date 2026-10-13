package notification

import (
	"context"
	"errors"
	"fmt"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/configs"
)

type NotificationService interface {
	CreateNotification(ctx context.Context, input CreateNotificationInput) error
	DeleteNotification(ctx context.Context, actorID int64, notifType NotificationType, entityID int64) error
	GetNotifications(ctx context.Context, userID int64, page, limit int) (*GetNotificationsResponse, error)
	GetUnreadCount(ctx context.Context, userID int64) (*UnreadCountResponse, error)
	MarkAllAsRead(ctx context.Context, userID int64) error
	MarkAsReadByID(ctx context.Context, userID int64, notificationID int64) error
}

type notificationService struct {
	cfg  *configs.Config
	repo NotificationRepository
}

func NewNotificationService(cfg *configs.Config, repo NotificationRepository) NotificationService {
	return &notificationService{
		cfg:  cfg,
		repo: repo,
	}
}

func (s *notificationService) CreateNotification(ctx context.Context, input CreateNotificationInput) error {
	if input.ActorID == input.UserID {
		return nil
	}

	if input.UserID <= 0 || input.ActorID <= 0 || input.Type == "" || input.EntityID <= 0 {
		return errors.New("invalid notification input")
	}

	if err := s.repo.CreateNotification(ctx, input); err != nil {
		return fmt.Errorf("service CreateNotification: %w", err)
	}

	return nil
}

func (s *notificationService) DeleteNotification(ctx context.Context, actorID int64, notifType NotificationType, entityID int64) error {
	if actorID <= 0 || notifType == "" || entityID <= 0 {
		return nil
	}

	if err := s.repo.DeleteNotification(ctx, actorID, notifType, entityID); err != nil {
		return fmt.Errorf("service DeleteNotification: %w", err)
	}

	return nil
}

func (s *notificationService) GetNotifications(ctx context.Context, userID int64, page, limit int) (*GetNotificationsResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit
	items, totalData, err := s.repo.GetNotifications(ctx, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("service GetNotifications: %w", err)
	}

	totalPage := 0
	if totalData > 0 {
		totalPage = (totalData + limit - 1) / limit
	}

	return &GetNotificationsResponse{
		Data: items,
		Pagination: NotificationPagination{
			Page:      page,
			Limit:     limit,
			TotalPage: totalPage,
			TotalData: totalData,
		},
	}, nil
}

func (s *notificationService) GetUnreadCount(ctx context.Context, userID int64) (*UnreadCountResponse, error) {
	count, err := s.repo.CountUnread(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("service GetUnreadCount: %w", err)
	}

	return &UnreadCountResponse{
		UnreadCount: count,
	}, nil
}

func (s *notificationService) MarkAllAsRead(ctx context.Context, userID int64) error {
	if err := s.repo.MarkAllAsRead(ctx, userID); err != nil {
		return fmt.Errorf("service MarkAllAsRead: %w", err)
	}

	return nil
}

func (s *notificationService) MarkAsReadByID(ctx context.Context, userID int64, notificationID int64) error {
	if err := s.repo.MarkAsReadByID(ctx, userID, notificationID); err != nil {
		return fmt.Errorf("service MarkAsReadByID: %w", err)
	}

	return nil
}
