package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"profile-service/config"
	"profile-service/repository"
	"profile-service/service"

	observability "users-observability"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

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
	shutdown := observability.InitProvider(ctx, "profile-service")
	defer func() { _ = shutdown(ctx) }()

	seedProfiles := map[string]repository.Profile{
		"11111111-1111-1111-1111-111111111111": {
			UUID:         "11111111-1111-1111-1111-111111111111",
			Name:         "Chris",
			Surname:      "Example",
			Email:        "chris@example.com",
			PersonalCode: "12345678901",
		},
	}

	repo := repository.NewInMemoryProfileRepository(seedProfiles)
	svc := service.NewProfileService(repo)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(otelgin.Middleware("profile-service"))

	r.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	r.GET(cfg.BasePath+"/profiles/:uuid", func(c *gin.Context) {
		uuid := c.Param("uuid")
		p, ok, err := svc.GetProfile(c.Request.Context(), uuid)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "repository error"})
			return
		}
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "no profile found for user",
				"uuid":  uuid,
			})
			return
		}
		c.JSON(http.StatusOK, ProfileResponse{
			UUID:         p.UUID,
			Name:         p.Name,
			Surname:      p.Surname,
			Email:        p.Email,
			PersonalCode: p.PersonalCode,
		})
	})

	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("starting profile-service on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}