package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"context"
	"net"

	"mobile-api-service/config"
	"mobile-api-service/service"
	"mobile-api-service/controller"
	"mobile-api-service/router"
	"mobile-api-service/middleware"

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

	baseTransport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   3 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   10,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	
	client := &http.Client{
		Timeout: 5 * time.Second,
		Transport: otelhttp.NewTransport(
			baseTransport,
			otelhttp.WithSpanNameFormatter(func(_ string, r *http.Request) string {
				return fmt.Sprintf("HTTP %s %s", r.Method, r.URL.Path)
			}),
		),
	}
	orch := service.NewOrchestrator(client, cfg.UsersServiceURL, cfg.ProfileServiceURL)
	ctrl := controller.NewMobileController(orch)
	
	rt := router.New()
	rt.Engine().Use(otelgin.Middleware("mobile-api-service"))
	rt.Engine().Use(middleware.TestHooks())
	rt.Setup(ctrl, cfg.BasePath)

	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("starting mobile-api-service on %s", addr)
	if err := rt.Engine().Run(addr); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}