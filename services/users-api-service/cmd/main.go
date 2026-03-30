package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"users-api-service/config"
	"users-api-service/controller"
	"users-api-service/repository"
	"users-api-service/service"
	
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
	shutdown := observability.InitProvider(ctx, "users-api-service")
	defer func() { _ = shutdown(ctx) }()

	seedUsers := map[string]string{
		"chris": "11111111-1111-1111-1111-111111111111",
		"lissu": "22222222-2222-2222-2222-222222222222",
	}

	repo := repository.NewInMemoryUserRepository(seedUsers)
	svc := service.NewUserService(repo)
	ctrl := controller.NewUsersController(svc)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(otelgin.Middleware("users-api-service"))

	r.GET("/health", ctrl.Health)
	r.GET(cfg.BasePath+"/:username", ctrl.GetUserUUID)

	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("starting users-api-service on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}