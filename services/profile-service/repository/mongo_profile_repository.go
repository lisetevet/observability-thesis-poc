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

func (r *MongoProfileRepository) GetByUUID(ctx context.Context, uuid string) (model.Profile, bool, error) {
    tr := otel.Tracer("profile-service")
    ctx, span := tr.Start(ctx, "MongoProfileRepository.GetByUUID")
    span.SetAttributes(attribute.String("app.uuid", uuid))
    defer span.End()

    var doc model.Profile
    err := r.coll.FindOne(ctx, bson.M{"uuid": uuid}).Decode(&doc)

    if err == nil {
        return doc, true, nil
    }

    if errors.Is(err, mongo.ErrNoDocuments) {
        log.Printf("no profile found (uuid=%s)", uuid)
        return model.Profile{}, false, nil
    }

    log.Printf("mongo GetByUUID failed (uuid=%s): %v", uuid, err)
    return model.Profile{}, false, err
}

func (r *MongoProfileRepository) UpsertProfile(ctx context.Context, p model.Profile) error {
	tr := otel.Tracer("profile-service")
	ctx, span := tr.Start(ctx, "MongoProfileRepository.UpsertProfile")
	span.SetAttributes(attribute.String("db.collection", "profiles"))
	span.SetAttributes(attribute.String("app.uuid", p.UUID))
	defer span.End()

	_, err := r.coll.UpdateOne(
		ctx,
		bson.M{"uuid": p.UUID},
		bson.M{"$set": bson.M{
			"username":      p.Username,
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

func (r *MongoProfileRepository) GetByUsername(ctx context.Context, username string) (model.Profile, bool, error) {
	tr := otel.Tracer("profile-service")
	ctx, span := tr.Start(ctx, "MongoProfileRepository.GetByUsername")
	span.SetAttributes(attribute.String("app.username", username))
	defer span.End()

	var doc model.Profile
	err := r.coll.FindOne(ctx, bson.M{"username": username}).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Printf("profile not found (username=%s)", username)
			return model.Profile{}, false, nil
		}
		log.Printf("mongo GetByUsername failed (username=%s): %v", username, err)
		return model.Profile{}, false, err
	}
	return doc, true, nil
}