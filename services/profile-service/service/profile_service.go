package service

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

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

func (s *ProfileService) GetProfileByUsernameDBFirst(ctx context.Context, username, usersDelayMs, usersFail string) (model.Profile, bool, error) {
	// 1) try profile DB by username
	p, ok, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		return model.Profile{}, false, err
	}
	if ok {
		return p, true, nil
	}

	// 2) fallback: fetch profile seed from users-service
	status, _, body, err := s.usersCl.GetProfileByUsername(ctx, username, usersDelayMs, usersFail)
	if err != nil {
		return model.Profile{}, false, err
	}
	if status != http.StatusOK {
		// users returns 404 etc -> treat as "not found"
		return model.Profile{}, false, nil
	}

	var fromUsers model.Profile
	if err := json.Unmarshal(body, &fromUsers); err != nil {
		return model.Profile{}, false, err
	}
	// ensure username is set (users response might not include it)
	fromUsers.Username = username

	// 3) async cache write-back
	go func(p model.Profile) {
		bg := context.Background()
		if err := s.repo.UpsertProfile(bg, p); err != nil {
			log.Printf("async UpsertProfile failed (username=%s uuid=%s): %v", p.Username, p.UUID, err)
		}
	}(fromUsers)

	// return fetched profile immediately
	return fromUsers, true, nil
}