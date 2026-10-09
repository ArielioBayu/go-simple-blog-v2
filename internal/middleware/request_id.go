package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	RequestIDHeader = "X-Request-ID"
	RequestIDKey    = "request_id"
)

// GenerateRequestID generates a clean, URL-safe random request ID (e.g. req-4a92e10c7bf3d8a5).
func GenerateRequestID() string {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return "req-default"
	}
	return "req-" + hex.EncodeToString(b)
}

// RequestIDMiddleware ensures every incoming HTTP request has a unique Request ID.
// If the client sends X-Request-ID, it uses it; otherwise, it generates a new one.
// The Request ID is set in the gin.Context and sent back in the HTTP response header.
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqID := strings.TrimSpace(c.Request.Header.Get(RequestIDHeader))
		if reqID == "" {
			reqID = GenerateRequestID()
		}

		c.Set(RequestIDKey, reqID)
		c.Header(RequestIDHeader, reqID)
		c.Next()
	}
}

// GetRequestID retrieves the request ID from gin.Context.
func GetRequestID(c *gin.Context) string {
	if c == nil {
		return ""
	}
	if id, exists := c.Get(RequestIDKey); exists {
		if idStr, ok := id.(string); ok {
			return idStr
		}
	}
	return ""
}
