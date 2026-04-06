package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"users-api-service/config"
	"users-api-service/controller"
	"users-api-service/repository"
	"users-api-service/service"
	"users-api-service/router"
	"users-api-service/middleware"
	
	observability "users-observability"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"
	"go.mongodb.org/mongo-driver/bson"
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

	type seedUser struct {
		UUID         string
		Name         string
		Surname      string
		Email        string
		PersonalCode string
	}

	seedUsers := map[string]seedUser{
		"chris": {
			UUID:         "11111111-1111-1111-1111-111111111111",
			Name:         "Chris",
			Surname:      "Example",
			Email:        "chris@example.com",
			PersonalCode: "12345678901",
		},
		"lissu": {
			UUID:         "22222222-2222-2222-2222-222222222222",
			Name:         "Lissu",
			Surname:      "Example",
			Email:        "lissu@example.com",
			PersonalCode: "10987654321",
		},
	}

	for username, u := range seedUsers {
		_, err := coll.UpdateOne(
			ctx,
			bson.M{"username": username},
			bson.M{"$set": bson.M{
				"username":       username,
				"uuid":           u.UUID,
				"name":           u.Name,
				"surname":        u.Surname,
				"email":          u.Email,
				"personal_code":  u.PersonalCode,
			}},
			options.Update().SetUpsert(true),
		)
		if err != nil {
			log.Fatalf("failed to seed user %s: %v", username, err)
		}
	}

	repo := repository.NewMongoUserRepository(coll)
	svc := service.NewUserService(repo)
	ctrl := controller.NewUsersController(svc)

	rt := router.New()
	rt.Engine().Use(otelgin.Middleware("users-api-service"))
	rt.Engine().Use(middleware.TestHooks())
	rt.Setup(ctrl, cfg.BasePath)

	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("starting users-api-service on %s", addr)
	if err := rt.Engine().Run(addr); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}