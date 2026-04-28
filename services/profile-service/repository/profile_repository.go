package repository

import (
	"context"

	"profile-service/model"
)

type ProfileRepository interface {
	GetByUUID(ctx context.Context, uuid string) (model.Profile, bool, error)
	GetByUsername(ctx context.Context, username string) (model.Profile, bool, error)
	UpsertProfile(ctx context.Context, p model.Profile) error
}
