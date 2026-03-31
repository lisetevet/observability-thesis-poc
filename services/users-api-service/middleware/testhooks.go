package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// TestHooks provides controlled failure/latency injection for experiments.
// Usage examples:
//   ?delayMs=250   -> sleeps 250ms before continuing
//   ?fail=true     -> returns 500 immediately
func TestHooks() gin.HandlerFunc {
	return func(c *gin.Context) {
		// fail injection
		if c.Query("fail") == "true" {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "injected failure",
			})
			c.Abort()
			return
		}

		// delay injection
		if v := c.Query("delayMs"); v != "" {
			ms, err := strconv.Atoi(v)
			if err == nil && ms > 0 && ms <= 60000 {
				time.Sleep(time.Duration(ms) * time.Millisecond)
			}
		}

		c.Next()
	}
}