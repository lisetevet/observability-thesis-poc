package service

import (
	"context"
	"log"

	"profile-service/model"
	"profile-service/pkg/usersclient"
	"profile-service/repository"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
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

func (s *ProfileService) GetProfileByUsername(ctx context.Context, username, usersDelayMs, usersFail string) (model.Profile, bool, error) {
	tr := otel.Tracer("profile-service")
	ctx, span := tr.Start(ctx, "ProfileService.GetProfileByUsername")
	span.SetAttributes(attribute.String("app.username", username))
	defer span.End()

	uuid, ok, err := s.usersCl.GetUUIDByUsername(ctx, username, usersDelayMs, usersFail)
	if err != nil {
		log.Printf("failed to resolve uuid from users-service (username=%s): %v", username, err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return model.Profile{}, false, err
	}

	if !ok {
		log.Printf("user not found in users-service (username=%s)", username)
		span.SetAttributes(attribute.Bool("users.found", false))
		return model.Profile{}, false, nil
	}

	span.SetAttributes(
		attribute.Bool("users.found", true),
		attribute.String("app.uuid", uuid),
	)

	profile, found, err := s.repo.GetByUUID(ctx, uuid)
	if err != nil {
		log.Printf("failed to fetch profile by uuid (username=%s uuid=%s): %v", username, uuid, err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return model.Profile{}, false, err
	}

	if !found {
		log.Printf("profile not found (username=%s uuid=%s)", username, uuid)
		span.SetAttributes(attribute.Bool("profile.found", false))
		return model.Profile{}, false, nil
	}

	span.SetAttributes(attribute.Bool("profile.found", true))
	return profile, true, nil
}
