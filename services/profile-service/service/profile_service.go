package service

import (
	"context"

	"profile-service/repository"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

type ProfileService struct {
	repo repository.ProfileRepository
}

func NewProfileService(repo repository.ProfileRepository) *ProfileService {
	return &ProfileService{repo: repo}
}

func (s *ProfileService) GetProfile(ctx context.Context, uuid string) (repository.Profile, bool, error) {
	tr := otel.Tracer("profile-service")
	ctx, span := tr.Start(ctx, "ProfileService.GetProfile")
	span.SetAttributes(attribute.String("app.uuid", uuid))
	defer span.End()

	return s.repo.GetByUUID(ctx, uuid)
}