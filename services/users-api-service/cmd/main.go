package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "os"

    "users-api-service/config"
    "users-api-service/repository"
    "users-api-service/service"

    observability "users-observability"

    "github.com/gin-gonic/gin"
    "go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

type UserResponse struct {
	Username string `json:"username"`
	UUID     string `json:"uuid"`
}

func main() {
	// Load configuration (allow override via CONFIG_PATH)
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config/config.yaml"
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	ctx := context.Background()
	shutdown := observability.InitProvider(ctx, "users-api-service")
	defer func() { _ = shutdown(ctx) }()

	seedUsers := map[string]string{
		"chris": "11111111-1111-1111-1111-111111111111",
		"lissu": "22222222-2222-2222-2222-222222222222",
	}

	repo := repository.NewInMemoryUserRepository(seedUsers)
	svc := service.NewUserService(repo)
	
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(otelgin.Middleware("users-api-service"))

	r.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	// GET {base_path}/:username -> { username, uuid } or 404
	r.GET(cfg.BasePath+"/:username", func(c *gin.Context) {
		username := c.Param("username")

		uuid, ok, err := svc.GetUUID(c.Request.Context(), username)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error":"failed to fetch user uuid"})
			return
		}
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{
				"error":    "user not found",
				"username": username,
			})
			return
		}

		c.JSON(http.StatusOK, UserResponse{
			Username: username,
			UUID:     uuid,
		})
	})

	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("starting users-api-service on %s", addr)

	if err := r.Run(addr); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}