package repository

import (
	"context"

	"users-api-service/model"
)

type UserRepository interface {
	GetUUIDByUsername(ctx context.Context, username string) (string, bool, error)
	GetUserByUsername(ctx context.Context, username string) (model.User, bool, error)
}