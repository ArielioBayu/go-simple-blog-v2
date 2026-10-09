package follow_test

import (
	"context"
	"errors"
	"testing"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/configs"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/constants"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/modules/follow"
	"github.com/ArielioBayu/go-simple-blog-v2/pkg/utils"
)

type mockFollowRepo struct {
	checkUserExistsFunc         func(ctx context.Context, userID int) (bool, error)
	isUserPrivateFunc           func(ctx context.Context, userID int) (bool, error)
	getFollowRecordFunc         func(ctx context.Context, followerID, followingID int) (*follow.FollowModel, error)
	followUserFunc              func(ctx context.Context, followerID, followingID int, status string) error
	unfollowUserFunc            func(ctx context.Context, followerID, followingID int) error
	removeFollowerFunc          func(ctx context.Context, currentUserID, followerID int) error
	getPendingFollowRequestsFunc func(ctx context.Context, followingID, limit, offset int) ([]follow.FollowRequestUser, int, error)
	acceptFollowRequestFunc     func(ctx context.Context, followingID, followerID int) error
	rejectFollowRequestFunc     func(ctx context.Context, followingID, followerID int) error
	canViewUserContentFunc      func(ctx context.Context, viewerID, targetUserID int) (bool, error)
	getFollowersFunc            func(ctx context.Context, currentUserID, targetUserID, limit, offset int, search string) ([]follow.FollowUserData, int, error)
	getFollowingFunc            func(ctx context.Context, currentUserID, targetUserID, limit, offset int, search string) ([]follow.FollowUserData, int, error)
	getRelationshipFunc         func(ctx context.Context, currentUserID, targetUserID int) (*follow.RelationshipResponse, error)
}

func (m *mockFollowRepo) CheckUserExists(ctx context.Context, userID int) (bool, error) {
	if m.checkUserExistsFunc != nil {
		return m.checkUserExistsFunc(ctx, userID)
	}
	return true, nil
}

func (m *mockFollowRepo) IsUserPrivate(ctx context.Context, userID int) (bool, error) {
	if m.isUserPrivateFunc != nil {
		return m.isUserPrivateFunc(ctx, userID)
	}
	return false, nil
}

func (m *mockFollowRepo) GetFollowRecord(ctx context.Context, followerID, followingID int) (*follow.FollowModel, error) {
	if m.getFollowRecordFunc != nil {
		return m.getFollowRecordFunc(ctx, followerID, followingID)
	}
	return nil, nil
}

func (m *mockFollowRepo) FollowUser(ctx context.Context, followerID, followingID int, status string) error {
	if m.followUserFunc != nil {
		return m.followUserFunc(ctx, followerID, followingID, status)
	}
	return nil
}

func (m *mockFollowRepo) UnfollowUser(ctx context.Context, followerID, followingID int) error {
	if m.unfollowUserFunc != nil {
		return m.unfollowUserFunc(ctx, followerID, followingID)
	}
	return nil
}

func (m *mockFollowRepo) RemoveFollower(ctx context.Context, currentUserID, followerID int) error {
	if m.removeFollowerFunc != nil {
		return m.removeFollowerFunc(ctx, currentUserID, followerID)
	}
	return nil
}

func (m *mockFollowRepo) GetPendingFollowRequests(ctx context.Context, followingID, limit, offset int) ([]follow.FollowRequestUser, int, error) {
	if m.getPendingFollowRequestsFunc != nil {
		return m.getPendingFollowRequestsFunc(ctx, followingID, limit, offset)
	}
	return nil, 0, nil
}

func (m *mockFollowRepo) AcceptFollowRequest(ctx context.Context, followingID, followerID int) error {
	if m.acceptFollowRequestFunc != nil {
		return m.acceptFollowRequestFunc(ctx, followingID, followerID)
	}
	return nil
}

func (m *mockFollowRepo) RejectFollowRequest(ctx context.Context, followingID, followerID int) error {
	if m.rejectFollowRequestFunc != nil {
		return m.rejectFollowRequestFunc(ctx, followingID, followerID)
	}
	return nil
}

func (m *mockFollowRepo) CanViewUserContent(ctx context.Context, viewerID, targetUserID int) (bool, error) {
	if m.canViewUserContentFunc != nil {
		return m.canViewUserContentFunc(ctx, viewerID, targetUserID)
	}
	return true, nil
}

func (m *mockFollowRepo) GetFollowers(ctx context.Context, currentUserID, targetUserID, limit, offset int, search string) ([]follow.FollowUserData, int, error) {
	if m.getFollowersFunc != nil {
		return m.getFollowersFunc(ctx, currentUserID, targetUserID, limit, offset, search)
	}
	return nil, 0, nil
}

func (m *mockFollowRepo) GetFollowing(ctx context.Context, currentUserID, targetUserID, limit, offset int, search string) ([]follow.FollowUserData, int, error) {
	if m.getFollowingFunc != nil {
		return m.getFollowingFunc(ctx, currentUserID, targetUserID, limit, offset, search)
	}
	return nil, 0, nil
}

func (m *mockFollowRepo) GetRelationship(ctx context.Context, currentUserID, targetUserID int) (*follow.RelationshipResponse, error) {
	if m.getRelationshipFunc != nil {
		return m.getRelationshipFunc(ctx, currentUserID, targetUserID)
	}
	return nil, nil
}

func TestFollowUser(t *testing.T) {
	ctx := context.Background()
	cfg := &configs.Config{}

	t.Run("Cannot follow self", func(t *testing.T) {
		repo := &mockFollowRepo{}
		svc := follow.NewFollowService(cfg, repo)

		_, err := svc.FollowUser(ctx, 1, 1)
		if !errors.Is(err, constants.ErrCannotFollowSelf) {
			t.Fatalf("expected ErrCannotFollowSelf, got %v", err)
		}
	})

	t.Run("Target user not found", func(t *testing.T) {
		repo := &mockFollowRepo{
			checkUserExistsFunc: func(ctx context.Context, userID int) (bool, error) {
				return false, nil
			},
		}
		svc := follow.NewFollowService(cfg, repo)

		_, err := svc.FollowUser(ctx, 1, 2)
		if !errors.Is(err, constants.ErrUserNotFound) {
			t.Fatalf("expected ErrUserNotFound, got %v", err)
		}
	})

	t.Run("Already following", func(t *testing.T) {
		repo := &mockFollowRepo{
			getFollowRecordFunc: func(ctx context.Context, followerID, followingID int) (*follow.FollowModel, error) {
				return &follow.FollowModel{Status: utils.FollowStatusAccepted}, nil
			},
		}
		svc := follow.NewFollowService(cfg, repo)

		_, err := svc.FollowUser(ctx, 1, 2)
		if !errors.Is(err, constants.ErrAlreadyFollowing) {
			t.Fatalf("expected ErrAlreadyFollowing, got %v", err)
		}
	})

	t.Run("Follow request already pending", func(t *testing.T) {
		repo := &mockFollowRepo{
			getFollowRecordFunc: func(ctx context.Context, followerID, followingID int) (*follow.FollowModel, error) {
				return &follow.FollowModel{Status: utils.FollowStatusPending}, nil
			},
		}
		svc := follow.NewFollowService(cfg, repo)

		_, err := svc.FollowUser(ctx, 1, 2)
		if !errors.Is(err, constants.ErrFollowRequestPending) {
			t.Fatalf("expected ErrFollowRequestPending, got %v", err)
		}
	})

	t.Run("Follow public user -> accepted status", func(t *testing.T) {
		repo := &mockFollowRepo{
			isUserPrivateFunc: func(ctx context.Context, userID int) (bool, error) {
				return false, nil
			},
			followUserFunc: func(ctx context.Context, followerID, followingID int, status string) error {
				if status != utils.FollowStatusAccepted {
					t.Fatalf("expected status accepted, got %s", status)
				}
				return nil
			},
		}
		svc := follow.NewFollowService(cfg, repo)

		res, err := svc.FollowUser(ctx, 1, 2)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res.Status != utils.FollowStatusAccepted || res.IsPending {
			t.Fatalf("unexpected response: %+v", res)
		}
	})

	t.Run("Follow private user -> pending status", func(t *testing.T) {
		repo := &mockFollowRepo{
			isUserPrivateFunc: func(ctx context.Context, userID int) (bool, error) {
				return true, nil
			},
			followUserFunc: func(ctx context.Context, followerID, followingID int, status string) error {
				if status != utils.FollowStatusPending {
					t.Fatalf("expected status pending, got %s", status)
				}
				return nil
			},
		}
		svc := follow.NewFollowService(cfg, repo)

		res, err := svc.FollowUser(ctx, 1, 2)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res.Status != utils.FollowStatusPending || !res.IsPending {
			t.Fatalf("unexpected response: %+v", res)
		}
	})
}

func TestUnfollowUser(t *testing.T) {
	ctx := context.Background()
	cfg := &configs.Config{}

	t.Run("Follow record not found", func(t *testing.T) {
		repo := &mockFollowRepo{
			unfollowUserFunc: func(ctx context.Context, followerID, followingID int) error {
				return constants.ErrFollowNotFound
			},
		}
		svc := follow.NewFollowService(cfg, repo)

		err := svc.UnfollowUser(ctx, 1, 2)
		if !errors.Is(err, constants.ErrFollowNotFound) {
			t.Fatalf("expected ErrFollowNotFound, got %v", err)
		}
	})

	t.Run("Success unfollow", func(t *testing.T) {
		repo := &mockFollowRepo{
			unfollowUserFunc: func(ctx context.Context, followerID, followingID int) error {
				return nil
			},
		}
		svc := follow.NewFollowService(cfg, repo)

		err := svc.UnfollowUser(ctx, 1, 2)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestAcceptAndRejectFollowRequest(t *testing.T) {
	ctx := context.Background()
	cfg := &configs.Config{}

	t.Run("Accept request not found", func(t *testing.T) {
		repo := &mockFollowRepo{
			acceptFollowRequestFunc: func(ctx context.Context, followingID, followerID int) error {
				return constants.ErrFollowRequestNotFound
			},
		}
		svc := follow.NewFollowService(cfg, repo)

		err := svc.AcceptFollowRequest(ctx, 2, 1)
		if !errors.Is(err, constants.ErrFollowRequestNotFound) {
			t.Fatalf("expected ErrFollowRequestNotFound, got %v", err)
		}
	})

	t.Run("Success accept request", func(t *testing.T) {
		repo := &mockFollowRepo{
			acceptFollowRequestFunc: func(ctx context.Context, followingID, followerID int) error {
				return nil
			},
		}
		svc := follow.NewFollowService(cfg, repo)

		err := svc.AcceptFollowRequest(ctx, 2, 1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("Reject request not found", func(t *testing.T) {
		repo := &mockFollowRepo{
			rejectFollowRequestFunc: func(ctx context.Context, followingID, followerID int) error {
				return constants.ErrFollowRequestNotFound
			},
		}
		svc := follow.NewFollowService(cfg, repo)

		err := svc.RejectFollowRequest(ctx, 2, 1)
		if !errors.Is(err, constants.ErrFollowRequestNotFound) {
			t.Fatalf("expected ErrFollowRequestNotFound, got %v", err)
		}
	})

	t.Run("Success reject request", func(t *testing.T) {
		repo := &mockFollowRepo{
			rejectFollowRequestFunc: func(ctx context.Context, followingID, followerID int) error {
				return nil
			},
		}
		svc := follow.NewFollowService(cfg, repo)

		err := svc.RejectFollowRequest(ctx, 2, 1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestCanViewUserContent(t *testing.T) {
	ctx := context.Background()
	cfg := &configs.Config{}

	t.Run("Target user is private and viewer does not follow", func(t *testing.T) {
		repo := &mockFollowRepo{
			canViewUserContentFunc: func(ctx context.Context, viewerID, targetUserID int) (bool, error) {
				return false, nil
			},
		}
		svc := follow.NewFollowService(cfg, repo)

		canView, err := svc.CanViewUserContent(ctx, 1, 2)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if canView {
			t.Fatalf("expected false, got true")
		}
	})

	t.Run("Target user is public or viewer is follower", func(t *testing.T) {
		repo := &mockFollowRepo{
			canViewUserContentFunc: func(ctx context.Context, viewerID, targetUserID int) (bool, error) {
				return true, nil
			},
		}
		svc := follow.NewFollowService(cfg, repo)

		canView, err := svc.CanViewUserContent(ctx, 1, 2)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !canView {
			t.Fatalf("expected true, got false")
		}
	})
}
