package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

func TestHooks() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		tr := otel.Tracer("profile-service/testhooks")
		ctx, span := tr.Start(ctx, "TestHooks")
		c.Request = c.Request.WithContext(ctx)
		defer span.End()

		// query params: ?delayMs=... & ?fail=true
		delayMsStr := c.Query("delayMs")
		fail := c.Query("fail")

		if delayMsStr != "" {
			ms, err := strconv.Atoi(delayMsStr)
			if err != nil || ms < 0 {
				span.RecordError(fmt.Errorf("invalid delayMs=%q", delayMsStr))
				span.SetStatus(codes.Error, "invalid delayMs")
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid delayMs"})
				return
			}
			span.SetAttributes(attribute.Int("test.delay_ms", ms))
			time.Sleep(time.Duration(ms) * time.Millisecond)
		}

		if fail == "true" {
			err := fmt.Errorf("injected failure")
			span.RecordError(err)
			span.SetStatus(codes.Error, "injected failure")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "injected failure"})
			return
		}

		c.Next()
	}
}
