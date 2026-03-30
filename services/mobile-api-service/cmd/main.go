package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
	"context"
	"encoding/json"

	"mobile-api-service/config"

	"github.com/gin-gonic/gin"
	observability "users-observability"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type UsersLookupResponse struct {
	Username string `json:"username"`
	UUID     string `json:"uuid"`
}

type ProfileResponse struct {
	UUID         string `json:"uuid"`
	Name         string `json:"name"`
	Surname      string `json:"surname"`
	Email        string `json:"email"`
	PersonalCode string `json:"personal_code"`
}

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

	r.GET(cfg.BasePath+"/profile/:username", func(c *gin.Context) {
	username := c.Param("username")

	// 1) Call users-service to get UUID
	usersURL := fmt.Sprintf("%s/%s", cfg.UsersServiceURL, username)

	usersResp, err := client.Get(usersURL)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "users-service request failed"})
		return
	}
	defer usersResp.Body.Close()

	usersBody, _ := io.ReadAll(usersResp.Body)
	if usersResp.StatusCode != http.StatusOK {
		// Pass through errors (e.g., 404 user not found)
		c.Data(usersResp.StatusCode, usersResp.Header.Get("Content-Type"), usersBody)
		return
	}

	var lookup UsersLookupResponse
	if err := json.Unmarshal(usersBody, &lookup); err != nil || lookup.UUID == "" {
		c.JSON(http.StatusBadGateway, gin.H{"error": "invalid users-service response"})
		return
	}

	// 2) Call profile-service to get profile by UUID
	profileURL := fmt.Sprintf("%s/%s", cfg.ProfileServiceURL, lookup.UUID)

	profResp, err := client.Get(profileURL)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "profile-service request failed"})
		return
	}
	defer profResp.Body.Close()

	profBody, _ := io.ReadAll(profResp.Body)

	// Pass through profile-service response (200 or 404 no profile found)
	c.Data(profResp.StatusCode, profResp.Header.Get("Content-Type"), profBody)
})

	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("starting mobile-api-service on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}