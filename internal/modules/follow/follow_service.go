package follow

import (
	"context"
	"fmt"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/configs"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/constants"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/modules/notification"
)

type FollowService interface {
	FollowUser(ctx context.Context, currentUserID, targetUserID int) (*FollowActionResponse, error)
	UnfollowUser(ctx context.Context, currentUserID, targetUserID int) error
	RemoveFollower(ctx context.Context, currentUserID, followerUserID int) error
	GetPendingRequests(ctx context.Context, currentUserID, page, limit int) (*FollowRequestsResponse, error)
	AcceptFollowRequest(ctx context.Context, currentUserID, followerUserID int) error
	RejectFollowRequest(ctx context.Context, currentUserID, followerUserID int) error
	GetFollowers(ctx context.Context, currentUserID, targetUserID, page, limit int, search string) (*FollowListResponse, error)
	GetFollowing(ctx context.Context, currentUserID, targetUserID, page, limit int, search string) (*FollowListResponse, error)
	GetRelationshipStatus(ctx context.Context, currentUserID, targetUserID int) (*RelationshipResponse, error)
	CanViewUserContent(ctx context.Context, viewerID, targetUserID int) (bool, error)
}

type followService struct {
	cfg        *configs.Config
	followRepo FollowRepository
	notifSvc   notification.NotificationService
}

func NewFollowService(cfg *configs.Config, followRepo FollowRepository, notifSvc ...notification.NotificationService) FollowService {
	var svc notification.NotificationService
	if len(notifSvc) > 0 {
		svc = notifSvc[0]
	}
	return &followService{
		cfg:        cfg,
		followRepo: followRepo,
		notifSvc:   svc,
	}
}

func (s *followService) FollowUser(ctx context.Context, currentUserID, targetUserID int) (*FollowActionResponse, error) {
	if currentUserID == targetUserID {
		return nil, constants.ErrCannotFollowSelf
	}

	exists, err := s.followRepo.CheckUserExists(ctx, targetUserID)
	if err != nil {
		return nil, fmt.Errorf("service FollowUser check exists: %w", err)
	}
	if !exists {
		return nil, constants.ErrUserNotFound
	}

	// Check existing relationship
	existing, err := s.followRepo.GetFollowRecord(ctx, currentUserID, targetUserID)
	if err != nil {
		return nil, fmt.Errorf("service FollowUser check record: %w", err)
	}
	if existing != nil {
		if existing.Status == constants.FollowStatusAccepted {
			return nil, constants.ErrAlreadyFollowing
		}
		if existing.Status == constants.FollowStatusPending {
			return nil, constants.ErrFollowRequestPending
		}
	}

	isPrivate, err := s.followRepo.IsUserPrivate(ctx, targetUserID)
	if err != nil {
		return nil, fmt.Errorf("service FollowUser check privacy: %w", err)
	}

	status := constants.FollowStatusAccepted
	if isPrivate {
		status = constants.FollowStatusPending
	}

	if err := s.followRepo.FollowUser(ctx, currentUserID, targetUserID, status); err != nil {
		return nil, fmt.Errorf("service FollowUser: %w", err)
	}

	if s.notifSvc != nil {
		if status == constants.FollowStatusAccepted {
			_ = s.notifSvc.CreateNotification(ctx, notification.CreateNotificationInput{
				UserID:     int64(targetUserID),
				ActorID:    int64(currentUserID),
				Type:       notification.NotificationTypeFollow,
				EntityID:   int64(currentUserID),
				EntityType: notification.EntityTypeUser,
				Message:    constants.MsgFollow,
			})
		} else if status == constants.FollowStatusPending {
			_ = s.notifSvc.CreateNotification(ctx, notification.CreateNotificationInput{
				UserID:     int64(targetUserID),
				ActorID:    int64(currentUserID),
				Type:       notification.NotificationTypeFollowRequest,
				EntityID:   int64(currentUserID),
				EntityType: notification.EntityTypeUser,
				Message:    constants.MsgFollowRequest,
			})
		}
	}

	return &FollowActionResponse{
		TargetUserID: targetUserID,
		Status:       status,
		IsFollowing:  status == constants.FollowStatusAccepted,
		IsPending:    status == constants.FollowStatusPending,
	}, nil
}

func (s *followService) UnfollowUser(ctx context.Context, currentUserID, targetUserID int) error {
	if currentUserID == targetUserID {
		return constants.ErrCannotFollowSelf
	}
	if err := s.followRepo.UnfollowUser(ctx, currentUserID, targetUserID); err != nil {
		return err
	}
	if s.notifSvc != nil {
		_ = s.notifSvc.DeleteNotification(ctx, int64(currentUserID), notification.NotificationTypeFollow, int64(currentUserID))
		_ = s.notifSvc.DeleteNotification(ctx, int64(currentUserID), notification.NotificationTypeFollowRequest, int64(currentUserID))
	}
	return nil
}

func (s *followService) RemoveFollower(ctx context.Context, currentUserID, followerUserID int) error {
	if currentUserID == followerUserID {
		return constants.ErrCannotFollowSelf
	}
	return s.followRepo.RemoveFollower(ctx, currentUserID, followerUserID)
}

func (s *followService) GetPendingRequests(ctx context.Context, currentUserID, page, limit int) (*FollowRequestsResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 10
	}
	offset := (page - 1) * limit

	requests, total, err := s.followRepo.GetPendingFollowRequests(ctx, currentUserID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("service GetPendingRequests: %w", err)
	}

	totalPage := 0
	if total > 0 {
		totalPage = (total + limit - 1) / limit
	}

	return &FollowRequestsResponse{
		Data: requests,
		Pagination: FollowPagination{
			Page:      page,
			Limit:     limit,
			TotalPage: totalPage,
			TotalData: total,
		},
	}, nil
}

func (s *followService) AcceptFollowRequest(ctx context.Context, currentUserID, followerUserID int) error {
	if currentUserID == followerUserID {
		return constants.ErrCannotFollowSelf
	}
	if err := s.followRepo.AcceptFollowRequest(ctx, currentUserID, followerUserID); err != nil {
		return err
	}
	if s.notifSvc != nil {
		_ = s.notifSvc.DeleteNotification(ctx, int64(followerUserID), notification.NotificationTypeFollowRequest, int64(followerUserID))
		_ = s.notifSvc.CreateNotification(ctx, notification.CreateNotificationInput{
			UserID:     int64(currentUserID),
			ActorID:    int64(followerUserID),
			Type:       notification.NotificationTypeFollow,
			EntityID:   int64(followerUserID),
			EntityType: notification.EntityTypeUser,
			Message:    constants.MsgFollow,
		})
	}
	return nil
}

func (s *followService) RejectFollowRequest(ctx context.Context, currentUserID, followerUserID int) error {
	if currentUserID == followerUserID {
		return constants.ErrCannotFollowSelf
	}
	if err := s.followRepo.RejectFollowRequest(ctx, currentUserID, followerUserID); err != nil {
		return err
	}
	if s.notifSvc != nil {
		_ = s.notifSvc.DeleteNotification(ctx, int64(followerUserID), notification.NotificationTypeFollowRequest, int64(followerUserID))
	}
	return nil
}

func (s *followService) GetFollowers(ctx context.Context, currentUserID, targetUserID, page, limit int, search string) (*FollowListResponse, error) {
	exists, err := s.followRepo.CheckUserExists(ctx, targetUserID)
	if err != nil {
		return nil, fmt.Errorf("service GetFollowers check exists: %w", err)
	}
	if !exists {
		return nil, constants.ErrUserNotFound
	}

	canView, err := s.followRepo.CanViewUserContent(ctx, currentUserID, targetUserID)
	if err != nil {
		return nil, fmt.Errorf("service GetFollowers check privacy: %w", err)
	}
	if !canView {
		return nil, constants.ErrPrivateAccount
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 10
	}
	offset := (page - 1) * limit

	followers, total, err := s.followRepo.GetFollowers(ctx, currentUserID, targetUserID, limit, offset, search)
	if err != nil {
		return nil, fmt.Errorf("service GetFollowers: %w", err)
	}

	totalPage := 0
	if total > 0 {
		totalPage = (total + limit - 1) / limit
	}

	return &FollowListResponse{
		Data: followers,
		Pagination: FollowPagination{
			Page:      page,
			Limit:     limit,
			TotalPage: totalPage,
			TotalData: total,
		},
	}, nil
}

func (s *followService) GetFollowing(ctx context.Context, currentUserID, targetUserID, page, limit int, search string) (*FollowListResponse, error) {
	exists, err := s.followRepo.CheckUserExists(ctx, targetUserID)
	if err != nil {
		return nil, fmt.Errorf("service GetFollowing check exists: %w", err)
	}
	if !exists {
		return nil, constants.ErrUserNotFound
	}

	canView, err := s.followRepo.CanViewUserContent(ctx, currentUserID, targetUserID)
	if err != nil {
		return nil, fmt.Errorf("service GetFollowing check privacy: %w", err)
	}
	if !canView {
		return nil, constants.ErrPrivateAccount
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 10
	}
	offset := (page - 1) * limit

	following, total, err := s.followRepo.GetFollowing(ctx, currentUserID, targetUserID, limit, offset, search)
	if err != nil {
		return nil, fmt.Errorf("service GetFollowing: %w", err)
	}

	totalPage := 0
	if total > 0 {
		totalPage = (total + limit - 1) / limit
	}

	return &FollowListResponse{
		Data: following,
		Pagination: FollowPagination{
			Page:      page,
			Limit:     limit,
			TotalPage: totalPage,
			TotalData: total,
		},
	}, nil
}

func (s *followService) GetRelationshipStatus(ctx context.Context, currentUserID, targetUserID int) (*RelationshipResponse, error) {
	exists, err := s.followRepo.CheckUserExists(ctx, targetUserID)
	if err != nil {
		return nil, fmt.Errorf("service GetRelationshipStatus check exists: %w", err)
	}
	if !exists {
		return nil, constants.ErrUserNotFound
	}

	return s.followRepo.GetRelationship(ctx, currentUserID, targetUserID)
}

func (s *followService) CanViewUserContent(ctx context.Context, viewerID, targetUserID int) (bool, error) {
	return s.followRepo.CanViewUserContent(ctx, viewerID, targetUserID)
}


