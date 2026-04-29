package service

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"mobile-api-service/client/profileclient"
	"mobile-api-service/model"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type Orchestrator struct {
	profile *profileclient.Client
}

func NewOrchestrator(profile *profileclient.Client) *Orchestrator {
	return &Orchestrator{profile: profile}
}

func (o *Orchestrator) FetchProfileByUsername(
	ctx context.Context,
	username string,
	query model.ProfileLookupQuery,
) (int, string, []byte, error) {
	tr := otel.Tracer("mobile-api-service")
	ctx, span := tr.Start(ctx, "Orchestrator.FetchProfileByUsername")
	span.SetAttributes(
		attribute.String("app.username", username),
		attribute.String("test.usersDelayMs", query.UsersDelayMs),
		attribute.String("test.usersFail", query.UsersFail),
		attribute.String("test.profileDelayMs", query.ProfileDelayMs),
		attribute.String("test.profileFail", query.ProfileFail),
	)
	defer span.End()

	status, contentType, body, err := o.profile.Get(
		ctx,
		url.PathEscape(username),
		query.ToProfileServiceQuery(),
	)
	span.SetAttributes(attribute.Int("downstream.profile.status", status))

	if err != nil {
		log.Printf("profile-service lookup failed (username=%s): %v", username, err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return 0, "", nil, err
	}

	if status >= 400 {
		log.Printf("profile-service returned non-success status (username=%s status=%d)", username, status)
		span.SetStatus(codes.Error, fmt.Sprintf("profile-service returned %d", status))
	}

	return status, contentType, body, nil
}
