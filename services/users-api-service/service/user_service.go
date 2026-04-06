package service

import (
	"context"

	"users-api-service/model"
	"users-api-service/repository"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
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

	uuid, ok, err := s.repo.GetUUIDByUsername(ctx, username)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "repository error")
		return "", false, err
	}
	return uuid, ok, nil
}

func (s *UserService) GetUser(ctx context.Context, username string) (model.User, bool, error) {
	tr := otel.Tracer("users-api-service")
	ctx, span := tr.Start(ctx, "UserService.GetUser")
	span.SetAttributes(attribute.String("app.username", username))
	defer span.End()

	u, ok, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "repository error")
		return model.User{}, false, err
	}
	return u, ok, nil
}