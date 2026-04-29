package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func TestHooks() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Query("fail") == "true" {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "injected failure",
			})
			c.Abort()
			return
		}

		if v := c.Query("delayMs"); v != "" {
			ms, err := strconv.Atoi(v)
			if err == nil && ms > 0 && ms <= 60000 {
				time.Sleep(time.Duration(ms) * time.Millisecond)
			}
		}

		c.Next()
	}
}
