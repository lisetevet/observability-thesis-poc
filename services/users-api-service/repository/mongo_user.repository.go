package repository

import (
	"context"
	"errors"
	"log"

	"users-api-service/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type MongoUserRepository struct {
	coll *mongo.Collection
}

func NewMongoUserRepository(coll *mongo.Collection) *MongoUserRepository {
	return &MongoUserRepository{coll: coll}
}

func (r *MongoUserRepository) GetUUIDByUsername(ctx context.Context, username string) (string, bool, error) {
	tr := otel.Tracer("users-api-service")
	ctx, span := tr.Start(ctx, "MongoUserRepository.GetUUIDByUsername")
	span.SetAttributes(attribute.String("app.username", username))
	defer span.End()

	var doc model.User
	err := r.coll.FindOne(ctx, bson.M{"username": username}).Decode(&doc)

	if err == nil {
		span.SetAttributes(attribute.Bool("db.found", true))
		return doc.UUID, true, nil
	}

	if errors.Is(err, mongo.ErrNoDocuments) {
		span.SetAttributes(attribute.Bool("db.found", false))
		return "", false, nil
	}

	span.RecordError(err)
	span.SetStatus(codes.Error, "mongo find failed")
	log.Printf("mongo GetUUIDByUsername failed (username=%s): %v", username, err)
	return "", false, err
}

func (r *MongoUserRepository) GetUserByUsername(ctx context.Context, username string) (model.User, bool, error) {
	tr := otel.Tracer("users-api-service")
	ctx, span := tr.Start(ctx, "MongoUserRepository.GetUserByUsername")
	span.SetAttributes(attribute.String("app.username", username))
	defer span.End()

	var doc model.User
	err := r.coll.FindOne(ctx, bson.M{"username": username}).Decode(&doc)

	if err == nil {
		span.SetAttributes(attribute.Bool("db.found", true))
		return doc, true, nil
	}

	if errors.Is(err, mongo.ErrNoDocuments) {
		span.SetAttributes(attribute.Bool("db.found", false))
		return model.User{}, false, nil
	}

	span.RecordError(err)
	span.SetStatus(codes.Error, "mongo find failed")
	log.Printf("mongo GetUserByUsername failed (username=%s): %v", username, err)
	return model.User{}, false, err
}