package saves

import (
	"net/http"
	"strconv"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/configs"
	"github.com/ArielioBayu/go-simple-blog-v2/pkg/apperror"
	"github.com/ArielioBayu/go-simple-blog-v2/pkg/response"
	"github.com/gin-gonic/gin"
)

type SavesHandler struct {
	savesService SavesService
	cfg          *configs.Config
}

func NewSavesHandler(savesService SavesService, cfg *configs.Config) *SavesHandler {
	return &SavesHandler{
		savesService: savesService,
		cfg:          cfg,
	}
}

func (h *SavesHandler) SavePost(c *gin.Context) {
	ctx := c.Request.Context()

	var request SaveRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		response.Error(c, err)
		return
	}

	postId, err := strconv.Atoi(c.Param("postId"))
	if err != nil {
		response.Error(c, apperror.NewBadRequest("Invalid Post Id", err))
		return
	}

	userId, exists := c.Get("id")
	if !exists || userId.(int) == 0 {
		response.Error(c, apperror.NewUnauthorized("Unauthorized", nil))
		return
	}

	err = h.savesService.SavePost(ctx, postId, userId.(int), request)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, "success")
}

func (h *SavesHandler) GetSavedPosts(c *gin.Context) {
	ctx := c.Request.Context()

	userID := c.GetInt("id")
	if userID == 0 {
		response.Error(c, apperror.NewUnauthorized("Unauthorized", nil))
		return
	}

	page := 1
	if pStr := c.Query("page"); pStr != "" {
		if p, err := strconv.Atoi(pStr); err == nil && p > 0 {
			page = p
		}
	} else if pStr := c.Query("pageIndex"); pStr != "" {
		if p, err := strconv.Atoi(pStr); err == nil && p > 0 {
			page = p
		}
	}

	limit := 10
	if lStr := c.Query("limit"); lStr != "" {
		if l, err := strconv.Atoi(lStr); err == nil && l > 0 {
			limit = l
		}
	} else if lStr := c.Query("pageSize"); lStr != "" {
		if l, err := strconv.Atoi(lStr); err == nil && l > 0 {
			limit = l
		}
	}

	res, err := h.savesService.GetSavedPosts(ctx, userID, page, limit)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Paginated(c, http.StatusOK, "success get saved posts", res.Data, res.Pagination)
}

func (h *SavesHandler) GetSavedPostIDs(c *gin.Context) {
	ctx := c.Request.Context()

	userID := c.GetInt("id")
	if userID == 0 {
		response.Error(c, apperror.NewUnauthorized("Unauthorized", nil))
		return
	}

	ids, err := h.savesService.GetSavedPostIDs(ctx, userID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Data(c, http.StatusOK, "success get saved post ids", ids)
}

