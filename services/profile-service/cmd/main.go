package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"profile-service/config"
	"profile-service/repository"
	"profile-service/service"
	"profile-service/controller"

	observability "users-observability"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
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
    ctrl := controller.NewProfileController(svc)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(otelgin.Middleware("profile-service"))

	r.GET("/health", ctrl.Health)
	r.GET(cfg.BasePath+"/profiles/:uuid", ctrl.GetProfile)

	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("starting profile-service on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}