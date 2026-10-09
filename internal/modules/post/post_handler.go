package post

import (
	"net/http"
	"strconv"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/configs"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/constants"
	"github.com/ArielioBayu/go-simple-blog-v2/pkg/apperror"
	"github.com/ArielioBayu/go-simple-blog-v2/pkg/response"
	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	postService PostService
	cfg         *configs.Config
}

func NewPostHandler(postService PostService, cfg *configs.Config) *PostHandler {
	return &PostHandler{
		postService: postService,
		cfg:         cfg,
	}
}

func (h *PostHandler) CreatePost(c *gin.Context) {
	var request PostRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		response.Error(c, err)
		return
	}

	userId := c.GetInt("id")
	if userId == 0 {
		response.Error(c, constants.ErrUnauthorized)
		return
	}

	err = h.postService.CreatePost(c.Request.Context(), userId, request)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "Success Create Post")
}

func (h *PostHandler) GetAllPost(c *gin.Context) {
	ctx := c.Request.Context()

	pageIndex := 1
	if pageIndexStr := c.Query("pageIndex"); pageIndexStr != "" {
		p, err := strconv.Atoi(pageIndexStr)
		if err != nil || p <= 0 {
			response.Error(c, apperror.NewBadRequest("Invalid page index", nil))
			return
		}
		pageIndex = p
	}

	pageSize := 10
	if pageSizeStr := c.Query("pageSize"); pageSizeStr != "" {
		s, err := strconv.Atoi(pageSizeStr)
		if err != nil || s <= 0 {
			response.Error(c, apperror.NewBadRequest("Invalid page size", nil))
			return
		}
		pageSize = s
	}

	userID := c.GetInt("id")
	postResp, err := h.postService.GetAllPost(ctx, pageSize, pageIndex, userID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Paginated(c, http.StatusOK, "success get all post", postResp.Data, postResp.Pagination)
}

func (h *PostHandler) GetPersonalizedFeed(c *gin.Context) {
	ctx := c.Request.Context()
	userID := c.GetInt("id")
	if userID == 0 {
		response.Error(c, constants.ErrUnauthorized)
		return
	}

	pageIndex := 1
	if pageIndexStr := c.Query("pageIndex"); pageIndexStr != "" {
		p, err := strconv.Atoi(pageIndexStr)
		if err == nil && p > 0 {
			pageIndex = p
		}
	} else if pageStr := c.Query("page"); pageStr != "" {
		p, err := strconv.Atoi(pageStr)
		if err == nil && p > 0 {
			pageIndex = p
		}
	}

	pageSize := 10
	if pageSizeStr := c.Query("pageSize"); pageSizeStr != "" {
		s, err := strconv.Atoi(pageSizeStr)
		if err == nil && s > 0 {
			pageSize = s
		}
	} else if limitStr := c.Query("limit"); limitStr != "" {
		s, err := strconv.Atoi(limitStr)
		if err == nil && s > 0 {
			pageSize = s
		}
	}

	postResp, err := h.postService.GetPersonalizedFeed(ctx, pageSize, pageIndex, userID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Paginated(c, http.StatusOK, "success get personalized feed", postResp.Data, postResp.Pagination)
}

func (h *PostHandler) GetPostById(c *gin.Context) {
	ctx := c.Request.Context()
	postId := c.Param("postId")

	postIdInt, err := strconv.Atoi(postId)
	if err != nil || postIdInt <= 0 {
		response.Error(c, apperror.NewBadRequest("invalid post id", nil))
		return
	}

	userID := c.GetInt("id")
	data, err := h.postService.GetPostById(ctx, postIdInt, userID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Data(c, http.StatusOK, "success get post", data)
}

func (h *PostHandler) GetPostsByUserID(c *gin.Context) {
	ctx := c.Request.Context()

	userIdParam := c.Param("userId")
	if userIdParam == "" {
		userIdParam = c.Param("user_id")
	}

	targetUserID, err := strconv.Atoi(userIdParam)
	if err != nil || targetUserID <= 0 {
		response.Error(c, apperror.NewBadRequest("invalid user id", nil))
		return
	}

	pageIndex := 1
	if pageIndexStr := c.Query("pageIndex"); pageIndexStr != "" {
		p, err := strconv.Atoi(pageIndexStr)
		if err != nil || p <= 0 {
			response.Error(c, apperror.NewBadRequest("Invalid page index", nil))
			return
		}
		pageIndex = p
	} else if pageStr := c.Query("page"); pageStr != "" {
		p, err := strconv.Atoi(pageStr)
		if err == nil && p > 0 {
			pageIndex = p
		}
	}

	pageSize := 10
	if pageSizeStr := c.Query("pageSize"); pageSizeStr != "" {
		s, err := strconv.Atoi(pageSizeStr)
		if err != nil || s <= 0 {
			response.Error(c, apperror.NewBadRequest("Invalid page size", nil))
			return
		}
		pageSize = s
	} else if limitStr := c.Query("limit"); limitStr != "" {
		s, err := strconv.Atoi(limitStr)
		if err == nil && s > 0 {
			pageSize = s
		}
	}

	currentUserID := c.GetInt("id")
	postResp, err := h.postService.GetPostsByUserID(ctx, targetUserID, currentUserID, pageSize, pageIndex)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Paginated(c, http.StatusOK, "success get posts by user id", postResp.Data, postResp.Pagination)
}

func (h *PostHandler) DeletePost(c *gin.Context) {
	ctx := c.Request.Context()

	postId := c.Param("postId")
	postIdInt, err := strconv.Atoi(postId)
	if err != nil || postIdInt <= 0 {
		response.Error(c, apperror.NewBadRequest("invalid post id", nil))
		return
	}

	userId := c.GetInt("id")
	err = h.postService.DeletePost(ctx, postIdInt, userId)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, "success delete post")
}
