package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"mobile-api-service/config"

	"github.com/gin-gonic/gin"
	"context"
	observability "users-observability"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func main() {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config/config.yaml"
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	ctx := context.Background()
	shutdown := observability.InitProvider(ctx, "mobile-api-service")
	defer func() { _ = shutdown(ctx) }()

	client := &http.Client{
		Timeout:   5 * time.Second,
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(otelgin.Middleware("mobile-api-service"))

	r.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	// GET /api/v1/profile/:username -> calls users-service and returns its response
	r.GET(cfg.BasePath+"/profile/:username", func(c *gin.Context) {
		username := c.Param("username")
		url := fmt.Sprintf("%s/%s", cfg.UsersServiceURL, username)

		resp, err := client.Get(url)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "users-service request failed"})
			return
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), body)
	})

	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("starting mobile-api-service on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}