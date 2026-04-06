package repository

import (
	"context"
	"errors"
	"log"

	"profile-service/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

type MongoProfileRepository struct {
	coll *mongo.Collection
}

func NewMongoProfileRepository(coll *mongo.Collection) *MongoProfileRepository {
	return &MongoProfileRepository{coll: coll}
}

func (r *MongoProfileRepository) GetByUUID(ctx context.Context, uuid string) (Profile, bool, error) {
    tr := otel.Tracer("profile-service")
    ctx, span := tr.Start(ctx, "MongoProfileRepository.GetByUUID")
    span.SetAttributes(attribute.String("app.uuid", uuid))
    defer span.End()

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

func (r *MongoProfileRepository) UpsertProfile(ctx context.Context, p Profile) error {
	tr := otel.Tracer("profile-service")
	ctx, span := tr.Start(ctx, "MongoProfileRepository.UpsertProfile")
	span.SetAttributes(attribute.String("db.collection", "profiles"))
	span.SetAttributes(attribute.String("app.uuid", p.UUID))
	defer span.End()

	_, err := r.coll.UpdateOne(
		ctx,
		bson.M{"uuid": p.UUID},
		bson.M{"$set": bson.M{
			"uuid":          p.UUID,
			"name":          p.Name,
			"surname":       p.Surname,
			"email":         p.Email,
			"personal_code": p.PersonalCode,
		}},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		log.Printf("mongo UpsertProfile failed (uuid=%s): %v", p.UUID, err)
		return err
	}
	return nil
}