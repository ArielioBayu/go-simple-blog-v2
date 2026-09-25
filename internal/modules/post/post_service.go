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
	if len(request.UploadIDs) > 10 {
		return errors.New("maximum 10 images allowed per post")
	}

	var primaryUploadID *int64
	if len(request.UploadIDs) > 0 {
		primaryID := request.UploadIDs[0]
		primaryUploadID = &primaryID
	}

	postHashtags := strings.Join(request.PostHashtags, ",")
	now := time.Now()

	model := PostModel{
		UserId:       userId,
		PostTitle:    request.PostTitle,
		PostContent:  request.PostContent,
		PostHashtags: postHashtags,
		UploadID:     primaryUploadID,
		CreatedAt:    utils.JsonTime(now),
		UpdatedAt:    utils.JsonTime(now),
		CreatedBy:    strconv.Itoa(userId),
		UpdatedBy:    strconv.Itoa(userId),
	}

	err := s.postRepo.CreatePost(ctx, model, request.UploadIDs)
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

	type mediaToDelete struct {
		uploadID int
		filePath string
	}
	var filesToDelete []mediaToDelete

	if len(post.Media) > 0 {
		for _, m := range post.Media {
			filesToDelete = append(filesToDelete, mediaToDelete{
				uploadID: int(m.UploadID),
				filePath: m.FilePath,
			})
		}
	} else if post.UploadID != nil && *post.UploadID > 0 {
		filesToDelete = append(filesToDelete, mediaToDelete{
			uploadID: int(*post.UploadID),
			filePath: post.FilePath,
		})
	} else if post.FilePath != "" {
		filesToDelete = append(filesToDelete, mediaToDelete{
			uploadID: 0,
			filePath: post.FilePath,
		})
	}

	err = s.postRepo.DeletePost(ctx, postId)
	if err != nil {
		return fmt.Errorf("service DeletePost: %w", err)
	}

	for _, item := range filesToDelete {
		var uploadRecord *upload.UploadModel
		if s.uploadRepo != nil && item.uploadID > 0 {
			uploadRecord, err = s.uploadRepo.GetUploadByID(ctx, item.uploadID)
			if err != nil {
				log.Printf("warning: failed to fetch upload record %d: %v", item.uploadID, err)
			}
		} else if s.uploadRepo != nil && item.filePath != "" {
			uploadRecord, err = s.uploadRepo.GetUploadByPathOrFilename(ctx, item.filePath)
			if err != nil {
				log.Printf("warning: failed to fetch upload record by path %s: %v", item.filePath, err)
			}
		}

		if uploadRecord != nil {
			if s.uploadRepo != nil {
				if err := s.uploadRepo.DeleteUpload(ctx, int(uploadRecord.ID)); err != nil {
					log.Printf("warning: failed to delete upload record %d: %v", uploadRecord.ID, err)
				}
			}
			fileTarget := uploadRecord.SystemFilename
			if fileTarget == "" {
				fileTarget = uploadRecord.FilePath
			}
			if err := utils.RemoveUploadFile(fileTarget); err != nil {
				log.Printf("warning: failed to remove physical upload file (%s): %v", fileTarget, err)
			}
		} else if item.filePath != "" {
			if err := utils.RemoveUploadFile(item.filePath); err != nil {
				log.Printf("warning: failed to remove physical upload file (%s): %v", item.filePath, err)
			}
		}
	}

	if strings.Contains(post.PostContent, "uploads") && s.uploadRepo != nil {
		if contentUpload, err := s.uploadRepo.GetUploadByPathOrFilename(ctx, post.PostContent); err == nil && contentUpload != nil {
			_ = s.uploadRepo.DeleteUpload(ctx, int(contentUpload.ID))
			_ = utils.RemoveUploadFile(contentUpload.SystemFilename)
		}
	}

	return nil
}
