package upload

import (
	"net/http"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/configs"
	"github.com/ArielioBayu/go-simple-blog-v2/pkg/apperror"
	"github.com/ArielioBayu/go-simple-blog-v2/pkg/response"
	"github.com/ArielioBayu/go-simple-blog-v2/pkg/utils"
	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	uploadService UploadService
	cfg           *configs.Config
}

func NewUploadHandler(uploadService UploadService, cfg *configs.Config) *UploadHandler {
	return &UploadHandler{
		uploadService: uploadService,
		cfg:           cfg,
	}
}

func (h *UploadHandler) UploadFile(c *gin.Context) {
	ctx := c.Request.Context()

	userID := c.GetInt("id")
	if userID == 0 {
		response.Error(c, apperror.NewUnauthorized("unauthorized", nil))
		return
	}

	// Batasi ukuran request body ke 5MB + 1MB multipart overhead untuk proteksi DoS
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, utils.MaxUploadSize+1024*1024)

	file, err := c.FormFile("file")
	if err != nil {
		response.Error(c, apperror.NewBadRequest("File is required (key: 'file') and must be under 5MB", err))
		return
	}

	scheme := "http"
	if c.Request.TLS != nil || c.Request.Header.Get("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}
	host := c.Request.Host
	if host == "" {
		host = "localhost" + h.cfg.Service.Port
	}

	res, err := h.uploadService.UploadImage(ctx, userID, file, scheme, host)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Data(c, http.StatusCreated, "file uploaded successfully", res)
}

func (h *UploadHandler) GetMyUploads(c *gin.Context) {
	ctx := c.Request.Context()

	userID := c.GetInt("id")
	if userID == 0 {
		response.Error(c, apperror.NewUnauthorized("unauthorized", nil))
		return
	}

	uploads, err := h.uploadService.GetUserUploads(ctx, userID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Data(c, http.StatusOK, "success get user uploads", uploads)
}
