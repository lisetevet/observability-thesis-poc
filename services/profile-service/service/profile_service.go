package service

import (
	"context"

	"profile-service/repository"
	"profile-service/pkg/usersclient"
	"profile-service/model"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

type ProfileService struct {
	repo    repository.ProfileRepository
	usersCl *usersclient.Client
}

func NewProfileService(repo repository.ProfileRepository, usersCl *usersclient.Client) *ProfileService {
	return &ProfileService{repo: repo, usersCl: usersCl}
}

func (s *ProfileService) GetProfile(ctx context.Context, uuid string) (model.Profile, bool, error) {
	tr := otel.Tracer("profile-service")
	ctx, span := tr.Start(ctx, "ProfileService.GetProfile")
	span.SetAttributes(attribute.String("app.uuid", uuid))
	defer span.End()

	return s.repo.GetByUUID(ctx, uuid)
}

func (s *ProfileService) GetProfileByUsername(ctx context.Context, username string) (model.Profile, bool, error) {
	// 1) resolve uuid via users-service
	uuid, ok, err := s.usersCl.GetUUIDByUsername(ctx, username)
	if err != nil || !ok {
		return model.Profile{}, false, err
	}

	// 2) fetch profile by uuid from profile DB
	return s.repo.GetByUUID(ctx, uuid)
}