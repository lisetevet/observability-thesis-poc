package main

import (
	"fmt"
	"log"
	"os"

	"users-api-service/config"

	"github.com/gin-gonic/gin"
)

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

	r := gin.New()
	r.Use(gin.Recovery())

	r.GET("/health", func(c *gin.Context) {
		c.String(200, "ok")
	})

	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("starting users-api-service on %s", addr)

	if err := r.Run(addr); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}