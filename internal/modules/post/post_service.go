package post

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/configs"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/constants"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/modules/activity"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/modules/comment"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/modules/upload"
	"github.com/ArielioBayu/go-simple-blog-v2/pkg/utils"
)

type PostService interface {
	CreatePost(ctx context.Context, userId int, request PostRequest) error
	GetAllPost(ctx context.Context, pageSize, pageIndex, userID int) (GetAllPostResponse, error)
	GetPostsByUserID(ctx context.Context, targetUserID, currentUserID, pageSize, pageIndex int) (GetAllPostResponse, error)
	GetPostById(ctx context.Context, id, userID int) (*GetPostResponse, error)
	DeletePost(ctx context.Context, postId, userId int) error
}

type postService struct {
	cfg          *configs.Config
	postRepo     PostRepository
	commentRepo  comment.CommentRepository
	activityRepo activity.ActivityRepository
	uploadRepo   upload.UploadRepository
}

func NewPostService(
	cfg *configs.Config,
	postRepo PostRepository,
	commentRepo comment.CommentRepository,
	activityRepo activity.ActivityRepository,
	uploadRepo upload.UploadRepository,
) PostService {
	return &postService{
		cfg:          cfg,
		postRepo:     postRepo,
		commentRepo:  commentRepo,
		activityRepo: activityRepo,
		uploadRepo:   uploadRepo,
	}
}

func (s *postService) CreatePost(ctx context.Context, userId int, request PostRequest) error {
	postHashtags := strings.Join(request.PostHashtags, ",")
	now := time.Now()

	model := PostModel{
		UserId:       userId,
		PostTitle:    request.PostTitle,
		PostContent:  request.PostContent,
		PostHashtags: postHashtags,
		UploadID:     request.UploadID,
		CreatedAt:    utils.JsonTime(now),
		UpdatedAt:    utils.JsonTime(now),
		CreatedBy:    strconv.Itoa(userId),
		UpdatedBy:    strconv.Itoa(userId),
	}

	err := s.postRepo.CreatePost(ctx, model)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (s *postService) GetAllPost(ctx context.Context, pageSize, pageIndex, userID int) (GetAllPostResponse, error) {
	limit := pageSize
	offset := pageSize * (pageIndex - 1)

	response, err := s.postRepo.GetAllPost(ctx, limit, offset, userID)
	if err != nil {
		return response, fmt.Errorf("service GetAllPost: %w", err)
	}

	return response, nil
}

func (s *postService) GetPostsByUserID(ctx context.Context, targetUserID, currentUserID, pageSize, pageIndex int) (GetAllPostResponse, error) {
	limit := pageSize
	offset := pageSize * (pageIndex - 1)

	response, err := s.postRepo.GetPostsByUserID(ctx, targetUserID, currentUserID, limit, offset)
	if err != nil {
		return response, fmt.Errorf("service GetPostsByUserID: %w", err)
	}

	return response, nil
}

func (s *postService) GetPostById(ctx context.Context, id, userID int) (*GetPostResponse, error) {
	data, err := s.postRepo.GetPostById(ctx, id, userID)
	if err != nil {
		return nil, fmt.Errorf("service GetPostById: %w", err)
	}

	var counts int
	if s.activityRepo != nil {
		counts, err = s.activityRepo.CountLikedByPostID(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("service CountLikedByPostID: %w", err)
		}
	}

	var comments []comment.GetComment
	if s.commentRepo != nil {
		comments, err = s.commentRepo.GetCommentById(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("service GetCommentById: %w", err)
		}
	}

	result := &GetPostResponse{
		DetailPost: *data,
		LikedCount: counts,
		Comments:   comments,
	}

	return result, nil
}

func (s *postService) DeletePost(ctx context.Context, postId, userId int) error {
	post, err := s.postRepo.GetPostById(ctx, postId, userId)
	if err != nil {
		if errors.Is(err, constants.ErrPostNotFound) {
			return constants.ErrPostNotFound
		}
		return fmt.Errorf("service DeletePost: %w", err)
	}

	if post.UserId != userId {
		return constants.ErrForbidden
	}

	var uploadItem *upload.UploadModel
	if s.uploadRepo != nil {
		if post.UploadID != nil && *post.UploadID > 0 {
			uploadItem, err = s.uploadRepo.GetUploadByID(ctx, int(*post.UploadID))
			if err != nil {
				log.Printf("warning: failed to fetch upload record %d: %v", *post.UploadID, err)
			}
		} else if post.FilePath != "" {
			uploadItem, err = s.uploadRepo.GetUploadByPathOrFilename(ctx, post.FilePath)
			if err != nil {
				log.Printf("warning: failed to fetch upload record by path %s: %v", post.FilePath, err)
			}
		}
	}

	err = s.postRepo.DeletePost(ctx, postId)
	if err != nil {
		return fmt.Errorf("service DeletePost: %w", err)
	}

	if uploadItem != nil {
		if err := s.uploadRepo.DeleteUpload(ctx, int(uploadItem.ID)); err != nil {
			log.Printf("warning: failed to delete upload record %d: %v", uploadItem.ID, err)
		}
		fileTarget := uploadItem.SystemFilename
		if fileTarget == "" {
			fileTarget = uploadItem.FilePath
		}
		if err := utils.RemoveUploadFile(fileTarget); err != nil {
			log.Printf("warning: failed to remove physical upload file (%s): %v", fileTarget, err)
		}
	} else if post.FilePath != "" {
		if err := utils.RemoveUploadFile(post.FilePath); err != nil {
			log.Printf("warning: failed to remove physical upload file (%s): %v", post.FilePath, err)
		}
	} else if strings.Contains(post.PostContent, "uploads") && s.uploadRepo != nil {
		if contentUpload, err := s.uploadRepo.GetUploadByPathOrFilename(ctx, post.PostContent); err == nil && contentUpload != nil {
			_ = s.uploadRepo.DeleteUpload(ctx, int(contentUpload.ID))
			_ = utils.RemoveUploadFile(contentUpload.SystemFilename)
		}
	}

	return nil
}
