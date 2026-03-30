package service

import (
	"context"

	"users-api-service/repository"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUUID(ctx context.Context, username string) (string, bool, error) {
	return s.repo.GetUUIDByUsername(ctx, username)
}