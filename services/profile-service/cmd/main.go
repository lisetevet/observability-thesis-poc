package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
	"net/http"
	"net"

	"profile-service/config"
	"profile-service/repository"
	"profile-service/service"
	"profile-service/controller"
	"profile-service/router"
	"profile-service/middleware"
	"profile-service/pkg/usersclient"
	"profile-service/model"

	observability "users-observability"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"
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
	shutdown := observability.InitProvider(ctx, "profile-service")
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

	httpClient := &http.Client{
		Timeout: 5 * time.Second,
		Transport: otelhttp.NewTransport(
			baseTransport,
			otelhttp.WithSpanNameFormatter(func(_ string, r *http.Request) string {
				return fmt.Sprintf("HTTP %s %s", r.Method, r.URL.Path)
			}),
		),
	}
	usersCl := usersclient.New(httpClient, cfg.UsersServiceURL)

	// MongoDB client with OTEL monitoring
	mongoCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	clientOpts := options.Client().ApplyURI(cfg.MongoDB.URI)
	clientOpts.Monitor = otelmongo.NewMonitor()

	client, err := mongo.Connect(mongoCtx, clientOpts)
	if err != nil {
		log.Fatalf("failed to connect to MongoDB: %v", err)
	}
	defer func() { _ = client.Disconnect(context.Background()) }()

	db := client.Database(cfg.MongoDB.Database)
	coll := db.Collection(cfg.MongoDB.Collection)

	// Seed profiles (upsert)
	seed := model.Profile{
		Username:     "chris",
		UUID:         "11111111-1111-1111-1111-111111111111",
		Name:         "Chris",
		Surname:      "Example",
		Email:        "chris@example.com",
		PersonalCode: "12345678901",
	}

	_, err = coll.UpdateOne(
		ctx,
		bson.M{"uuid": seed.UUID},
		bson.M{"$set": bson.M{
			"username":      seed.Username,
			"uuid":          seed.UUID,
			"name":          seed.Name,
			"surname":       seed.Surname,
			"email":         seed.Email,
			"personal_code": seed.PersonalCode,
		}},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		log.Fatalf("failed to seed profile: %v", err)
	}

	repo := repository.NewMongoProfileRepository(coll)
	svc := service.NewProfileService(repo, usersCl)
	ctrl := controller.NewProfileController(svc)

	rt := router.New()
	rt.Engine().Use(otelgin.Middleware("profile-service"))
	rt.Engine().Use(middleware.TestHooks())
	rt.Setup(ctrl, cfg.BasePath)

	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("starting profile-service on %s", addr)
	if err := rt.Engine().Run(addr); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}