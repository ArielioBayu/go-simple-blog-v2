package follow

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/constants"
	"github.com/ArielioBayu/go-simple-blog-v2/pkg/utils"
)

type FollowRepository interface {
	CheckUserExists(ctx context.Context, userID int) (bool, error)
	IsUserPrivate(ctx context.Context, userID int) (bool, error)
	GetFollowRecord(ctx context.Context, followerID, followingID int) (*FollowModel, error)
	FollowUser(ctx context.Context, followerID, followingID int, status string) error
	UnfollowUser(ctx context.Context, followerID, followingID int) error
	RemoveFollower(ctx context.Context, currentUserID, followerID int) error
	GetPendingFollowRequests(ctx context.Context, followingID, limit, offset int) ([]FollowRequestUser, int, error)
	AcceptFollowRequest(ctx context.Context, followingID, followerID int) error
	RejectFollowRequest(ctx context.Context, followingID, followerID int) error
	CanViewUserContent(ctx context.Context, viewerID, targetUserID int) (bool, error)
	GetFollowers(ctx context.Context, currentUserID, targetUserID, limit, offset int, search string) ([]FollowUserData, int, error)
	GetFollowing(ctx context.Context, currentUserID, targetUserID, limit, offset int, search string) ([]FollowUserData, int, error)
	GetRelationship(ctx context.Context, currentUserID, targetUserID int) (*RelationshipResponse, error)
}

type followRepository struct {
	db *sql.DB
}

func NewFollowRepository(db *sql.DB) FollowRepository {
	return &followRepository{db: db}
}

func (r *followRepository) CheckUserExists(ctx context.Context, userID int) (bool, error) {
	query := `SELECT 1 FROM users WHERE id = ? LIMIT 1`
	var exists int
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("repository CheckUserExists: %w", err)
	}
	return true, nil
}

func (r *followRepository) IsUserPrivate(ctx context.Context, userID int) (bool, error) {
	query := `SELECT COALESCE(is_private, false) FROM users WHERE id = ? LIMIT 1`
	var isPrivate bool
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&isPrivate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, constants.ErrUserNotFound
		}
		return false, fmt.Errorf("repository IsUserPrivate: %w", err)
	}
	return isPrivate, nil
}

func (r *followRepository) GetFollowRecord(ctx context.Context, followerID, followingID int) (*FollowModel, error) {
	query := `SELECT id, follower_id, following_id, status, created_at, updated_at 
				FROM user_follows 
				WHERE follower_id = ? AND following_id = ? LIMIT 1`
	var model FollowModel
	err := r.db.QueryRowContext(ctx, query, followerID, followingID).Scan(
		&model.ID,
		&model.FollowerID,
		&model.FollowingID,
		&model.Status,
		&model.CreatedAt,
		&model.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("repository GetFollowRecord: %w", err)
	}
	return &model, nil
}

func (r *followRepository) FollowUser(ctx context.Context, followerID, followingID int, status string) error {
	query := `INSERT INTO user_follows (follower_id, following_id, status)
				VALUES (?, ?, ?)
				ON DUPLICATE KEY UPDATE status = VALUES(status), updated_at = NOW()`
	_, err := r.db.ExecContext(ctx, query, followerID, followingID, status)
	if err != nil {
		return fmt.Errorf("repository FollowUser: %w", err)
	}
	return nil
}

func (r *followRepository) UnfollowUser(ctx context.Context, followerID, followingID int) error {
	query := `DELETE FROM user_follows WHERE follower_id = ? AND following_id = ?`
	res, err := r.db.ExecContext(ctx, query, followerID, followingID)
	if err != nil {
		return fmt.Errorf("repository UnfollowUser: %w", err)
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return constants.ErrFollowNotFound
	}
	return nil
}

func (r *followRepository) RemoveFollower(ctx context.Context, currentUserID, followerID int) error {
	// currentUserID is the following_id, followerID is the follower_id
	query := `DELETE FROM user_follows WHERE follower_id = ? AND following_id = ?`
	res, err := r.db.ExecContext(ctx, query, followerID, currentUserID)
	if err != nil {
		return fmt.Errorf("repository RemoveFollower: %w", err)
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return constants.ErrFollowNotFound
	}
	return nil
}

func (r *followRepository) GetPendingFollowRequests(ctx context.Context, followingID, limit, offset int) ([]FollowRequestUser, int, error) {
	countQuery := `SELECT COUNT(id) FROM user_follows WHERE following_id = ? AND status = 'pending'`
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, followingID).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("repository GetPendingFollowRequests count: %w", err)
	}

	if total == 0 {
		return []FollowRequestUser{}, 0, nil
	}

	query := `SELECT uf.id, u.id, u.username, COALESCE(u.avatar_url, ''), COALESCE(u.bio, ''), uf.created_at
				FROM user_follows uf
				JOIN users u ON uf.follower_id = u.id
				WHERE uf.following_id = ? AND uf.status = 'pending'
				ORDER BY uf.created_at DESC
				LIMIT ? OFFSET ?`

	rows, err := r.db.QueryContext(ctx, query, followingID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("repository GetPendingFollowRequests rows: %w", err)
	}
	defer rows.Close()

	requests := make([]FollowRequestUser, 0)
	for rows.Next() {
		var req FollowRequestUser
		var createdAt time.Time
		if err := rows.Scan(
			&req.ID,
			&req.UserID,
			&req.Username,
			&req.AvatarURL,
			&req.Bio,
			&createdAt,
		); err != nil {
			return nil, 0, fmt.Errorf("repository GetPendingFollowRequests scan: %w", err)
		}
		req.RequestedAt = utils.JsonTime(createdAt)
		requests = append(requests, req)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("repository GetPendingFollowRequests: %w", err)
	}

	return requests, total, nil
}

func (r *followRepository) AcceptFollowRequest(ctx context.Context, followingID, followerID int) error {
	query := `UPDATE user_follows SET status = 'accepted', updated_at = NOW() 
				WHERE follower_id = ? AND following_id = ? AND status = 'pending'`
	res, err := r.db.ExecContext(ctx, query, followerID, followingID)
	if err != nil {
		return fmt.Errorf("repository AcceptFollowRequest: %w", err)
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return constants.ErrFollowRequestNotFound
	}
	return nil
}

func (r *followRepository) RejectFollowRequest(ctx context.Context, followingID, followerID int) error {
	query := `DELETE FROM user_follows 
				WHERE follower_id = ? AND following_id = ? AND status = 'pending'`
	res, err := r.db.ExecContext(ctx, query, followerID, followingID)
	if err != nil {
		return fmt.Errorf("repository RejectFollowRequest: %w", err)
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return constants.ErrFollowRequestNotFound
	}
	return nil
}

func (r *followRepository) CanViewUserContent(ctx context.Context, viewerID, targetUserID int) (bool, error) {
	if viewerID == targetUserID {
		return true, nil
	}

	query := `
		SELECT 
			COALESCE(u.is_private, false),
			EXISTS(
				SELECT 1 FROM user_follows 
				WHERE follower_id = ? AND following_id = ? AND status = 'accepted'
			) AS is_following
		FROM users u 
		WHERE u.id = ? LIMIT 1`

	var isPrivate, isFollowing bool
	err := r.db.QueryRowContext(ctx, query, viewerID, targetUserID, targetUserID).Scan(&isPrivate, &isFollowing)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, constants.ErrUserNotFound
		}
		return false, fmt.Errorf("repository CanViewUserContent: %w", err)
	}

	if !isPrivate || isFollowing {
		return true, nil
	}

	return false, nil
}

func (r *followRepository) GetFollowers(ctx context.Context, currentUserID, targetUserID, limit, offset int, search string) ([]FollowUserData, int, error) {
	countQuery := `SELECT COUNT(uf.id) 
					FROM user_follows uf 
					JOIN users u ON uf.follower_id = u.id 
					WHERE uf.following_id = ? AND uf.status = 'accepted'`
	countArgs := []interface{}{targetUserID}

	if search != "" {
		countQuery += ` AND u.username LIKE ?`
		countArgs = append(countArgs, "%"+search+"%")
	}

	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("repository GetFollowers count: %w", err)
	}

	if total == 0 {
		return []FollowUserData{}, 0, nil
	}

	query := `SELECT 
				u.id, 
				u.username, 
				COALESCE(u.avatar_url, ''), 
				COALESCE(u.bio, ''),
				EXISTS(SELECT 1 FROM user_follows WHERE follower_id = ? AND following_id = u.id AND status = 'accepted') AS is_following,
				EXISTS(SELECT 1 FROM user_follows WHERE follower_id = u.id AND following_id = ? AND status = 'accepted') AS is_followed_by
			  FROM user_follows uf
			  JOIN users u ON uf.follower_id = u.id
			  WHERE uf.following_id = ? AND uf.status = 'accepted'`
	args := []interface{}{currentUserID, currentUserID, targetUserID}

	if search != "" {
		query += ` AND u.username LIKE ?`
		args = append(args, "%"+search+"%")
	}

	query += ` ORDER BY uf.created_at DESC LIMIT ? OFFSET ?`
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("repository GetFollowers rows: %w", err)
	}
	defer rows.Close()

	followers := make([]FollowUserData, 0)
	for rows.Next() {
		var user FollowUserData
		if err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.AvatarURL,
			&user.Bio,
			&user.IsFollowing,
			&user.IsFollowedBy,
		); err != nil {
			return nil, 0, fmt.Errorf("repository GetFollowers scan: %w", err)
		}
		followers = append(followers, user)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("repository GetFollowers: %w", err)
	}

	return followers, total, nil
}

func (r *followRepository) GetFollowing(ctx context.Context, currentUserID, targetUserID, limit, offset int, search string) ([]FollowUserData, int, error) {
	countQuery := `SELECT COUNT(uf.id) 
					FROM user_follows uf 
					JOIN users u ON uf.following_id = u.id 
					WHERE uf.follower_id = ? AND uf.status = 'accepted'`
	countArgs := []interface{}{targetUserID}

	if search != "" {
		countQuery += ` AND u.username LIKE ?`
		countArgs = append(countArgs, "%"+search+"%")
	}

	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("repository GetFollowing count: %w", err)
	}

	if total == 0 {
		return []FollowUserData{}, 0, nil
	}

	query := `SELECT 
				u.id, 
				u.username, 
				COALESCE(u.avatar_url, ''), 
				COALESCE(u.bio, ''),
				EXISTS(SELECT 1 FROM user_follows WHERE follower_id = ? AND following_id = u.id AND status = 'accepted') AS is_following,
				EXISTS(SELECT 1 FROM user_follows WHERE follower_id = u.id AND following_id = ? AND status = 'accepted') AS is_followed_by
			  FROM user_follows uf
			  JOIN users u ON uf.following_id = u.id
			  WHERE uf.follower_id = ? AND uf.status = 'accepted'`
	args := []interface{}{currentUserID, currentUserID, targetUserID}

	if search != "" {
		query += ` AND u.username LIKE ?`
		args = append(args, "%"+search+"%")
	}

	query += ` ORDER BY uf.created_at DESC LIMIT ? OFFSET ?`
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("repository GetFollowing rows: %w", err)
	}
	defer rows.Close()

	following := make([]FollowUserData, 0)
	for rows.Next() {
		var user FollowUserData
		if err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.AvatarURL,
			&user.Bio,
			&user.IsFollowing,
			&user.IsFollowedBy,
		); err != nil {
			return nil, 0, fmt.Errorf("repository GetFollowing scan: %w", err)
		}
		following = append(following, user)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("repository GetFollowing: %w", err)
	}

	return following, total, nil
}

func (r *followRepository) GetRelationship(ctx context.Context, currentUserID, targetUserID int) (*RelationshipResponse, error) {
	query := `SELECT 
				COALESCE(u.is_private, false),
				(SELECT status FROM user_follows WHERE follower_id = ? AND following_id = ? LIMIT 1) AS my_follow_status,
				(SELECT status FROM user_follows WHERE follower_id = ? AND following_id = ? LIMIT 1) AS their_follow_status
			  FROM users u
			  WHERE u.id = ? LIMIT 1`

	var isPrivate bool
	var myFollowStatus, theirFollowStatus sql.NullString
	err := r.db.QueryRowContext(ctx, query, currentUserID, targetUserID, targetUserID, currentUserID, targetUserID).Scan(
		&isPrivate,
		&myFollowStatus,
		&theirFollowStatus,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, constants.ErrUserNotFound
		}
		return nil, fmt.Errorf("repository GetRelationship: %w", err)
	}

	isSelf := currentUserID == targetUserID
	if isSelf {
		return &RelationshipResponse{
			TargetUserID:   targetUserID,
			IsSelf:         true,
			IsPrivate:      isPrivate,
			IsFollowing:    false,
			IsPending:      false,
			IsFollowedBy:   false,
			CanViewContent: true,
		}, nil
	}

	isFollowing := myFollowStatus.Valid && myFollowStatus.String == utils.FollowStatusAccepted
	isPending := myFollowStatus.Valid && myFollowStatus.String == utils.FollowStatusPending
	isFollowedBy := theirFollowStatus.Valid && theirFollowStatus.String == utils.FollowStatusAccepted
	canViewContent := !isPrivate || isFollowing

	return &RelationshipResponse{
		TargetUserID:   targetUserID,
		IsSelf:         false,
		IsPrivate:      isPrivate,
		IsFollowing:    isFollowing,
		IsPending:      isPending,
		IsFollowedBy:   isFollowedBy,
		CanViewContent: canViewContent,
	}, nil
}

