package service

import (
	"context"

	"users-api-service/repository"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUUID(ctx context.Context, username string) (string, bool, error) {
	tr := otel.Tracer("users-api-service")
	ctx, span := tr.Start(ctx, "UserService.GetUUID")
	span.SetAttributes(attribute.String("app.username", username))
	defer span.End()

	return s.repo.GetUUIDByUsername(ctx, username)
}