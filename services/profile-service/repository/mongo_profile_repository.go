package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoProfileRepository struct {
	coll *mongo.Collection
}

func NewMongoProfileRepository(coll *mongo.Collection) *MongoProfileRepository {
	return &MongoProfileRepository{coll: coll}
}

type profileDoc struct {
	UUID         string `bson:"uuid"`
	Name         string `bson:"name"`
	Surname      string `bson:"surname"`
	Email        string `bson:"email"`
	PersonalCode string `bson:"personal_code"`
}

func (r *MongoProfileRepository) GetByUUID(ctx context.Context, uuid string) (Profile, bool, error) {
	var doc profileDoc
	err := r.coll.FindOne(ctx, bson.M{"uuid": uuid}).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return Profile{}, false, nil
		}
		return Profile{}, false, err
	}

	return Profile{
		UUID:         doc.UUID,
		Name:         doc.Name,
		Surname:      doc.Surname,
		Email:        doc.Email,
		PersonalCode: doc.PersonalCode,
	}, true, nil
}