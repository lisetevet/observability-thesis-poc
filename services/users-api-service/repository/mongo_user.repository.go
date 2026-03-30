package repository

import (
	"context"

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
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", false, nil
		}
		return "", false, err
	}
	return doc.UUID, true, nil
}