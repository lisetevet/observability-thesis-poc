package repository

import "context"

type UserRepository interface {
	GetUUIDByUsername(ctx context.Context, username string) (string, bool, error)
}
