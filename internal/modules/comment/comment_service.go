package comment

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/configs"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/constants"
	"github.com/ArielioBayu/go-simple-blog-v2/pkg/utils"
)

type PostChecker interface {
	CheckPostExists(ctx context.Context, postId int) (bool, error)
}

type CommentService interface {
	CreateComment(ctx context.Context, postId, userId int, request CommentRequest) error
	GetCommentsByPostID(ctx context.Context, postId, page, limit int) (*GetCommentsResponse, error)
	CountComments(ctx context.Context, postId int) (int, error)

	CreateReply(ctx context.Context, commentId, userId int, request ReplyRequest) error
	GetRepliesByCommentID(ctx context.Context, commentId, page, limit int) (*GetRepliesResponse, error)

	UpdateComment(ctx context.Context, commentId, currentUserId int, request UpdateCommentRequest) error
	DeleteComment(ctx context.Context, commentId, currentUserId int) error
	UpdateReply(ctx context.Context, replyId, currentUserId int, request UpdateReplyRequest) error
	DeleteReply(ctx context.Context, replyId, currentUserId int) error
}

type commentService struct {
	cfg         *configs.Config
	commentRepo CommentRepository
	postChecker PostChecker
}

func NewCommentService(cfg *configs.Config, commentRepo CommentRepository, postChecker PostChecker) CommentService {
	return &commentService{
		cfg:         cfg,
		commentRepo: commentRepo,
		postChecker: postChecker,
	}
}

func (s *commentService) CreateComment(ctx context.Context, postId, userId int, request CommentRequest) error {
	if s.postChecker != nil {
		exists, err := s.postChecker.CheckPostExists(ctx, postId)
		if err != nil {
			return err
		}
		if !exists {
			return constants.ErrPostNotFound
		}
	}

	now := time.Now()
	model := CommentModel{
		PostId:         postId,
		UserId:         userId,
		CommentContent: request.CommentContent,
		CreatedAt:      utils.JsonTime(now),
		UpdatedAt:      utils.JsonTime(now),
		CreatedBy:      strconv.Itoa(userId),
		UpdatedBy:      strconv.Itoa(userId),
	}

	err := s.commentRepo.CreateComment(ctx, model)
	if err != nil {
		return fmt.Errorf("service CreateComment: %w", err)
	}

	return nil
}

func (s *commentService) GetCommentsByPostID(ctx context.Context, postId, page, limit int) (*GetCommentsResponse, error) {
	if s.postChecker != nil {
		exists, err := s.postChecker.CheckPostExists(ctx, postId)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, constants.ErrPostNotFound
		}
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 10
	}

	offset := (page - 1) * limit
	comments, totalData, err := s.commentRepo.GetCommentsByPostID(ctx, postId, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("service GetCommentsByPostID: %w", err)
	}

	totalPage := 0
	if totalData > 0 {
		totalPage = (totalData + limit - 1) / limit
	}

	return &GetCommentsResponse{
		Data: comments,
		Pagination: CommentPagination{
			Page:      page,
			Limit:     limit,
			TotalPage: totalPage,
			TotalData: totalData,
		},
	}, nil
}

func (s *commentService) CountComments(ctx context.Context, postId int) (int, error) {
	if s.postChecker != nil {
		exists, err := s.postChecker.CheckPostExists(ctx, postId)
		if err != nil {
			return 0, err
		}
		if !exists {
			return 0, constants.ErrPostNotFound
		}
	}

	return s.commentRepo.CountCommentsByPostID(ctx, postId)
}

func (s *commentService) CreateReply(ctx context.Context, commentId, userId int, request ReplyRequest) error {
	exists, err := s.commentRepo.CheckCommentExists(ctx, commentId)
	if err != nil {
		return fmt.Errorf("service CreateReply check comment: %w", err)
	}
	if !exists {
		return constants.ErrCommentNotFound
	}

	now := time.Now()
	model := CommentReplyModel{
		CommentID:     commentId,
		UserID:        userId,
		ReplyToUserID: request.ReplyToUserID,
		ReplyContent:  request.ReplyContent,
		CreatedAt:     utils.JsonTime(now),
		UpdatedAt:     utils.JsonTime(now),
		CreatedBy:     strconv.Itoa(userId),
		UpdatedBy:     strconv.Itoa(userId),
	}

	err = s.commentRepo.CreateReply(ctx, model)
	if err != nil {
		return fmt.Errorf("service CreateReply: %w", err)
	}

	return nil
}

func (s *commentService) GetRepliesByCommentID(ctx context.Context, commentId, page, limit int) (*GetRepliesResponse, error) {
	exists, err := s.commentRepo.CheckCommentExists(ctx, commentId)
	if err != nil {
		return nil, fmt.Errorf("service GetReplies check comment: %w", err)
	}
	if !exists {
		return nil, constants.ErrCommentNotFound
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 10
	}

	offset := (page - 1) * limit
	replies, totalData, err := s.commentRepo.GetRepliesByCommentID(ctx, commentId, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("service GetReplies: %w", err)
	}

	totalPage := 0
	if totalData > 0 {
		totalPage = (totalData + limit - 1) / limit
	}

	return &GetRepliesResponse{
		Data: replies,
		Pagination: ReplyPagination{
			Page:      page,
			Limit:     limit,
			TotalPage: totalPage,
			TotalData: totalData,
		},
	}, nil
}

func (s *commentService) UpdateComment(ctx context.Context, commentId, currentUserId int, request UpdateCommentRequest) error {
	if strings.TrimSpace(request.CommentContent) == "" {
		return fmt.Errorf("comment content cannot be empty")
	}

	auth, err := s.commentRepo.GetCommentAuthInfo(ctx, commentId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return constants.ErrCommentNotFound
		}
		return fmt.Errorf("service UpdateComment check auth: %w", err)
	}

	if auth.CommentAuthorID != currentUserId {
		return constants.ErrForbidden
	}

	return s.commentRepo.UpdateComment(ctx, commentId, request.CommentContent)
}

func (s *commentService) DeleteComment(ctx context.Context, commentId, currentUserId int) error {
	auth, err := s.commentRepo.GetCommentAuthInfo(ctx, commentId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return constants.ErrCommentNotFound
		}
		return fmt.Errorf("service DeleteComment check auth: %w", err)
	}

	if auth.CommentAuthorID != currentUserId && auth.PostAuthorID != currentUserId {
		return constants.ErrForbidden
	}

	return s.commentRepo.DeleteComment(ctx, commentId)
}

func (s *commentService) UpdateReply(ctx context.Context, replyId, currentUserId int, request UpdateReplyRequest) error {
	if strings.TrimSpace(request.ReplyContent) == "" {
		return fmt.Errorf("reply content cannot be empty")
	}

	auth, err := s.commentRepo.GetReplyAuthInfo(ctx, replyId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return constants.ErrReplyNotFound
		}
		return fmt.Errorf("service UpdateReply check auth: %w", err)
	}

	if auth.ReplyAuthorID != currentUserId {
		return constants.ErrForbidden
	}

	return s.commentRepo.UpdateReply(ctx, replyId, request.ReplyContent)
}

func (s *commentService) DeleteReply(ctx context.Context, replyId, currentUserId int) error {
	auth, err := s.commentRepo.GetReplyAuthInfo(ctx, replyId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return constants.ErrReplyNotFound
		}
		return fmt.Errorf("service DeleteReply check auth: %w", err)
	}

	if auth.ReplyAuthorID != currentUserId && auth.PostAuthorID != currentUserId {
		return constants.ErrForbidden
	}

	return s.commentRepo.DeleteReply(ctx, replyId)
}

