package upload

import (
	"errors"
	"log"
	"net/http"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/configs"
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
		c.JSON(http.StatusUnauthorized, response.MessageResponse{
			Status:  http.StatusUnauthorized,
			Message: "unauthorized",
		})
		return
	}

	// Batasi ukuran request body ke 5MB + 1MB multipart overhead untuk proteksi DoS
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, utils.MaxUploadSize+1024*1024)

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, response.MessageResponse{
			Status:  http.StatusBadRequest,
			Message: "File is required (key: 'file') and must be under 5MB",
		})
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
		if errors.Is(err, utils.ErrFileRequired) || errors.Is(err, utils.ErrFileTooLarge) || errors.Is(err, utils.ErrInvalidFormat) {
			c.JSON(http.StatusBadRequest, response.MessageResponse{
				Status:  http.StatusBadRequest,
				Message: err.Error(),
			})
			return
		}

		log.Printf("UploadFile error: %v", err)
		c.JSON(http.StatusInternalServerError, response.MessageResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to upload file",
		})
		return
	}

	c.JSON(http.StatusCreated, response.DataResponse{
		Status:  http.StatusCreated,
		Message: "file uploaded successfully",
		Data:    res,
	})
}

func (h *UploadHandler) GetMyUploads(c *gin.Context) {
	ctx := c.Request.Context()

	userID := c.GetInt("id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, response.MessageResponse{
			Status:  http.StatusUnauthorized,
			Message: "unauthorized",
		})
		return
	}

	uploads, err := h.uploadService.GetUserUploads(ctx, userID)
	if err != nil {
		log.Printf("GetMyUploads error: %v", err)
		c.JSON(http.StatusInternalServerError, response.MessageResponse{
			Status:  http.StatusInternalServerError,
			Message: "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, response.DataResponse{
		Status:  http.StatusOK,
		Message: "success get user uploads",
		Data:    uploads,
	})
}
