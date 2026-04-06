package repository

import (
	"context"
	"errors"
    "log"

	"profile-service/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoProfileRepository struct {
	coll *mongo.Collection
}

func NewMongoProfileRepository(coll *mongo.Collection) *MongoProfileRepository {
	return &MongoProfileRepository{coll: coll}
}

func (r *MongoProfileRepository) GetByUUID(ctx context.Context, uuid string) (Profile, bool, error) {
    var doc model.Profile
    err := r.coll.FindOne(ctx, bson.M{"uuid": uuid}).Decode(&doc)

    if err == nil {
        return Profile{
            UUID:         doc.UUID,
            Name:         doc.Name,
            Surname:      doc.Surname,
            Email:        doc.Email,
            PersonalCode: doc.PersonalCode,
        }, true, nil
    }

    if errors.Is(err, mongo.ErrNoDocuments) {
		log.Printf("No documents found.")
        return Profile{}, false, nil
    }

    log.Printf("mongo GetByUUID failed (uuid=%s): %v", uuid, err)
    return Profile{}, false, err
}