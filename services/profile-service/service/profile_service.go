package service

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"fmt"

	"profile-service/repository"
	"profile-service/pkg/usersclient"
	"profile-service/model"

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

	// 1) cache lookup (profile DB)
	p, ok, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return model.Profile{}, false, err
	}
	if ok {
		span.SetAttributes(attribute.Bool("cache.hit", true))
		return p, true, nil
	}
	span.SetAttributes(attribute.Bool("cache.hit", false))

	// 2) fallback: fetch profile seed from users-service
	status, ct, body, err := s.usersCl.GetProfileByUsername(ctx, username, usersDelayMs, usersFail)
	span.SetAttributes(attribute.Int("users.status_code", status))
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return model.Profile{}, false, err
	}

	if status == http.StatusNotFound {
		return model.Profile{}, false, nil
	}
	if status != http.StatusOK {
		e := fmt.Errorf("users-service returned %d (%s): %s", status, ct, string(body))
		span.RecordError(e)
		span.SetStatus(codes.Error, e.Error())
		return model.Profile{}, false, e
	}

	var fromUsers model.Profile
	if err := json.Unmarshal(body, &fromUsers); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return model.Profile{}, false, err
	}
	// kindluse mõttes
	if fromUsers.Username == "" {
		fromUsers.Username = username
	}

	// 3) async cache write-back (goroutine)
	ctx2 := context.WithoutCancel(ctx)
	go func(p model.Profile) {
		tr2 := otel.Tracer("profile-service")
		_, span2 := tr2.Start(ctx2, "ProfileService.AsyncUpsert")
		span2.SetAttributes(
			attribute.String("app.username", p.Username),
			attribute.String("app.uuid", p.UUID),
		)
		defer span2.End()

		if err := s.repo.UpsertProfile(ctx2, p); err != nil {
			log.Printf("async UpsertProfile failed (username=%s uuid=%s): %v", p.Username, p.UUID, err)
			span2.RecordError(err)
			span2.SetStatus(codes.Error, err.Error())
		}
	}(fromUsers)

	// return fetched profile immediately
	return fromUsers, true, nil
}