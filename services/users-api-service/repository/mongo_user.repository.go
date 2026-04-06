package repository

import (
	"context"
	"errors"
	"log"

	"users-api-service/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoUserRepository struct {
	coll *mongo.Collection
}

func NewMongoUserRepository(coll *mongo.Collection) *MongoUserRepository {
	return &MongoUserRepository{coll: coll}
}

func (r *MongoUserRepository) GetUUIDByUsername(ctx context.Context, username string) (string, bool, error) {
	var doc model.User
	err := r.coll.FindOne(ctx, bson.M{"username": username}).Decode(&doc)

	if err == nil {
		return doc.UUID, true, nil
	}

	if errors.Is(err, mongo.ErrNoDocuments) {
		log.Printf("No documents found.")
		return "", false, nil
	}

	log.Printf("mongo GetUUIDByUsername failed (username=%s): %v", username, err)
	return "", false, err
}