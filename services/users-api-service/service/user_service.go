package service

import (
	"context"

	"users-api-service/repository"
	"users-api-service/model"

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

func (s *UserService) GetUser(ctx context.Context, username string) (model.User, bool, error) {
	tr := otel.Tracer("users-api-service")
	ctx, span := tr.Start(ctx, "UserService.GetUser")
	span.SetAttributes(attribute.String("app.username", username))
	defer span.End()

	return s.repo.GetUserByUsername(ctx, username)
}