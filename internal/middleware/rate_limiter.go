package middleware

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/constants"
	"github.com/ArielioBayu/go-simple-blog-v2/pkg/response"
	"github.com/gin-gonic/gin"
)

type clientRecord struct {
	count     int
	resetTime time.Time
}

type ipRateLimiter struct {
	mu          sync.Mutex
	records     map[string]*clientRecord
	maxRequests int
	window      time.Duration
}

func newIPRateLimiter(maxRequests int, window time.Duration) *ipRateLimiter {
	limiter := &ipRateLimiter{
		records:     make(map[string]*clientRecord),
		maxRequests: maxRequests,
		window:      window,
	}

	// Rutin background untuk membersihkan IP yang sudah tidak aktif agar mencegah kebocoran memori (memory leak)
	go func() {
		ticker := time.NewTicker(window * 2)
		for range ticker.C {
			limiter.mu.Lock()
			now := time.Now()
			for key, rec := range limiter.records {
				if now.After(rec.resetTime) {
					delete(limiter.records, key)
				}
			}
			limiter.mu.Unlock()
		}
	}()

	return limiter
}

func (l *ipRateLimiter) allow(key string) (bool, time.Duration) {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	rec, exists := l.records[key]

	if !exists || now.After(rec.resetTime) {
		l.records[key] = &clientRecord{
			count:     1,
			resetTime: now.Add(l.window),
		}
		return true, 0
	}

	if rec.count >= l.maxRequests {
		retryAfter := time.Until(rec.resetTime)
		if retryAfter < 0 {
			retryAfter = 0
		}
		return false, retryAfter
	}

	rec.count++
	return true, 0
}

func RateLimiter(maxRequests int, window time.Duration) gin.HandlerFunc {
	limiter := newIPRateLimiter(maxRequests, window)

	return func(c *gin.Context) {
		key := c.ClientIP() + ":" + c.FullPath()
		allowed, retryAfter := limiter.allow(key)

		if !allowed {
			seconds := int(retryAfter.Seconds()) + 1
			c.Header("Retry-After", fmt.Sprintf("%d", seconds))
			c.AbortWithStatusJSON(http.StatusTooManyRequests, response.MessageResponse{
				Status:  http.StatusTooManyRequests,
				Message: constants.ErrTooManyRequests.Error(),
			})
			return
		}

		c.Next()
	}
}

func OTPRateLimiter() gin.HandlerFunc {
	return RateLimiter(3, 1*time.Minute)
}
