package middleware

import (
	"log"
	"net/http"
	"runtime/debug"

	"github.com/ArielioBayu/go-simple-blog-v2/pkg/response"
	"github.com/gin-gonic/gin"
)

// RecoveryMiddleware provides production-grade panic recovery for Gin.
// It logs the internal panic and full stack trace for developers with the Request ID,
// while safely shielding the public client with a sanitized HTTP 500 response.
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				reqID := GetRequestID(c)
				stack := string(debug.Stack())

				// Log internal panic details for developer debugging
				log.Printf("[PANIC RECOVERED] request_id=%s method=%s path=%s err=%v\nstack:\n%s",
					reqID, c.Request.Method, c.Request.URL.Path, r, stack)

				// Never expose stack trace or raw panic message to client
				c.AbortWithStatusJSON(http.StatusInternalServerError, response.BaseResponse{
					Status:    http.StatusInternalServerError,
					Message:   "internal server error",
					RequestID: reqID,
				})
			}
		}()
		c.Next()
	}
}
