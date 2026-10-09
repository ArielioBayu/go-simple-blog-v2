package middleware_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/middleware"
	"github.com/ArielioBayu/go-simple-blog-v2/pkg/response"
	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestRequestIDMiddleware(t *testing.T) {
	t.Run("Generates Request ID when not provided", func(t *testing.T) {
		r := gin.New()
		r.Use(middleware.RequestIDMiddleware())
		r.GET("/test", func(c *gin.Context) {
			reqID := middleware.GetRequestID(c)
			if reqID == "" {
				t.Fatalf("expected non-empty request id in context")
			}
			c.String(http.StatusOK, "ok")
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		r.ServeHTTP(w, req)

		headerID := w.Header().Get("X-Request-ID")
		if headerID == "" {
			t.Fatalf("expected X-Request-ID header to be present")
		}
	})

	t.Run("Preserves client provided X-Request-ID", func(t *testing.T) {
		r := gin.New()
		r.Use(middleware.RequestIDMiddleware())
		r.GET("/test", func(c *gin.Context) {
			reqID := middleware.GetRequestID(c)
			if reqID != "client-id-123" {
				t.Fatalf("expected client-id-123, got %s", reqID)
			}
			c.String(http.StatusOK, "ok")
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("X-Request-ID", "client-id-123")
		r.ServeHTTP(w, req)

		headerID := w.Header().Get("X-Request-ID")
		if headerID != "client-id-123" {
			t.Fatalf("expected header client-id-123, got %s", headerID)
		}
	})
}

func TestRecoveryMiddleware(t *testing.T) {
	t.Run("Catches panic and returns safe 500 JSON with request_id", func(t *testing.T) {
		r := gin.New()
		r.Use(middleware.RequestIDMiddleware())
		r.Use(middleware.RecoveryMiddleware())
		r.GET("/panic", func(c *gin.Context) {
			panic("unexpected null pointer dereference")
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/panic", nil)
		req.Header.Set("X-Request-ID", "req-panic-test")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", w.Code)
		}

		var res response.ErrorResponse
		if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if res.Status != http.StatusInternalServerError {
			t.Fatalf("expected status 500, got %d", res.Status)
		}
		if res.Message != "internal server error" {
			t.Fatalf("expected message 'internal server error', got '%s'", res.Message)
		}
		if res.RequestID != "req-panic-test" {
			t.Fatalf("expected request_id 'req-panic-test', got '%s'", res.RequestID)
		}
	})
}
