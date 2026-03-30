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
	"mobile-api-service/controller"
	"mobile-api-service/router"

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
	ctrl := controller.NewMobileController(orch)
	
	rt := router.New()
	rt.Engine().Use(otelgin.Middleware("mobile-api-service"))
	rt.Setup(ctrl, cfg.BasePath)

	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("starting mobile-api-service on %s", addr)
	if err := rt.Engine().Run(addr); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}