package response

import (
	"log"

	"github.com/ArielioBayu/go-simple-blog-v2/pkg/apperror"
	"github.com/gin-gonic/gin"
)

// Pagination struct represents pagination metadata.
type Pagination struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

// BaseResponse struct defines the foundational envelope for all API responses.
type BaseResponse struct {
	Status    int    `json:"status"`
	Message   string `json:"message"`
	RequestID string `json:"request_id,omitempty"`
}

// DataResponse struct defines response containing data payload.
type DataResponse struct {
	Status    int    `json:"status"`
	Message   string `json:"message,omitempty"`
	Data      any    `json:"data,omitempty"`
	RequestID string `json:"request_id,omitempty"`
}

// MessageResponse is an alias to BaseResponse to ensure backward compatibility.
type MessageResponse = BaseResponse

// PaginationResponse struct defines response containing paginated data.
type PaginationResponse struct {
	Status     int    `json:"status"`
	Message    string `json:"message"`
	Pagination any    `json:"pagination,omitempty"`
	Data       any    `json:"data"`
	RequestID  string `json:"request_id,omitempty"`
}

// ErrorResponse represents a standardized public error response.
type ErrorResponse struct {
	Status    int               `json:"status"`
	Message   string            `json:"message"`
	Code      string            `json:"code,omitempty"`
	Errors    map[string]string `json:"errors,omitempty"`
	RequestID string            `json:"request_id,omitempty"`
}

// getRequestID retrieves request_id from gin context if present.
func getRequestID(c *gin.Context) string {
	if c == nil {
		return ""
	}
	if id, exists := c.Get("request_id"); exists {
		if s, ok := id.(string); ok {
			return s
		}
	}
	return ""
}

// Success sends a BaseResponse with the given HTTP status code and message.
func Success(c *gin.Context, statusCode int, message string) {
	reqID := getRequestID(c)
	c.JSON(statusCode, BaseResponse{
		Status:    statusCode,
		Message:   message,
		RequestID: reqID,
	})
}

// Data sends a DataResponse with the given HTTP status code, message, and data payload.
func Data(c *gin.Context, statusCode int, message string, data any) {
	reqID := getRequestID(c)
	c.JSON(statusCode, DataResponse{
		Status:    statusCode,
		Message:   message,
		Data:      data,
		RequestID: reqID,
	})
}

// Paginated sends a PaginationResponse with the given HTTP status code, message, data, and pagination metadata.
func Paginated(c *gin.Context, statusCode int, message string, data any, pagination any) {
	reqID := getRequestID(c)
	c.JSON(statusCode, PaginationResponse{
		Status:     statusCode,
		Message:    message,
		Pagination: pagination,
		Data:       data,
		RequestID:  reqID,
	})
}

// Error is the centralized error responder for all HTTP handlers.
// It maps internal errors into safe, consistent public responses,
// logs 5xx internal errors with Request ID, and never leaks raw technical details to clients.
func Error(c *gin.Context, err error) {
	if c == nil || err == nil {
		return
	}

	appErr := apperror.FromError(err)
	reqID := getRequestID(c)

	// Internal/unexpected error: log full context for developers
	if appErr.StatusCode >= 500 {
		log.Printf("[ERROR] request_id=%s method=%s path=%s err=%v",
			reqID, c.Request.Method, c.Request.URL.Path, appErr.Error())
	}

	c.JSON(appErr.StatusCode, ErrorResponse{
		Status:    appErr.StatusCode,
		Message:   appErr.Message,
		Code:      appErr.Code,
		Errors:    appErr.Errors,
		RequestID: reqID,
	})
}

