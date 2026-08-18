package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// LoginRateLimiter provides simple in-memory rate limiting for login attempts.
type LoginRateLimiter struct {
	mu       sync.Mutex
	attempts map[string][]time.Time
	limit    int
	window   time.Duration
}

var loginLimiter = &LoginRateLimiter{
	attempts: make(map[string][]time.Time),
	limit:    5,               // max attempts
	window:   5 * time.Minute, // per 5 minutes
}

// LoginRateLimit returns a gin middleware that limits login attempts per IP.
func LoginRateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !loginLimiter.Allow(ip) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many login attempts. Please try again later.",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

func (rl *LoginRateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-rl.window)

	// Clean old entries for this key.
	attempts := rl.attempts[key]
	valid := attempts[:0]
	for _, t := range attempts {
		if t.After(cutoff) {
			valid = append(valid, t)
		}
	}

	if len(valid) >= rl.limit {
		rl.attempts[key] = valid
		return false
	}

	rl.attempts[key] = append(valid, now)

	// Opportunistic sweep to prevent unbounded memory growth from idle keys.
	if len(rl.attempts) > 1024 {
		for k, v := range rl.attempts {
			if len(v) == 0 || !v[len(v)-1].After(cutoff) {
				delete(rl.attempts, k)
			}
		}
	}
	return true
}
