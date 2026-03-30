package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"context"

	"mobile-api-service/config"
	"mobile-api-service/service"

	"github.com/gin-gonic/gin"
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
	orch := service.NewOrchestrator(client, cfg.UsersServiceURL, cfg.ProfileServiceURL)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(otelgin.Middleware("mobile-api-service"))

	r.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	r.GET(cfg.BasePath+"/profile/:username", func(c *gin.Context) {
		username := c.Param("username")

		status, contentType, body, err := orch.FetchProfileByUsername(username)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}

		c.Data(status, contentType, body)
	})

	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("starting mobile-api-service on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}