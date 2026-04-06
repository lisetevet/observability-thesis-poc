package service

import (
	"net/http"
	"context"

	"mobile-api-service/pkg/profileclient"
	"mobile-api-service/pkg/usersclient"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type Orchestrator struct {
	users   *usersclient.Client
	profile *profileclient.Client
}

func NewOrchestrator(users *usersclient.Client, profile *profileclient.Client) *Orchestrator {
	return &Orchestrator{users: users, profile: profile}
}

func (o *Orchestrator) FetchProfileByUsername(ctx context.Context, username, usersDelayMs, usersFail, profileDelayMs, profileFail string) (int, string, []byte, error) {
	tr := otel.Tracer("mobile-api-service")
	ctx, span := tr.Start(ctx, "Orchestrator.FetchProfileByUsername")
	span.SetAttributes(
		attribute.String("app.username", username),
		attribute.String("test.usersDelayMs", usersDelayMs),
		attribute.String("test.usersFail", usersFail),
		attribute.String("test.profileDelayMs", profileDelayMs),
		attribute.String("test.profileFail", profileFail),
	)
	defer span.End()
	
	// 1) users lookup
	status, ct, body, uuid, err := o.users.GetUUIDByUsername(ctx, username, usersDelayMs, usersFail)
	span.SetAttributes(attribute.Int("downstream.users.status", status))

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "users lookup failed")
		return 0, "", nil, err
	}
	if status != http.StatusOK {
		// downstream tagastab error body -> pass-through
		span.SetStatus(codes.Error, "users returned non-200")
		return status, ct, body, nil
	}

	// 2) profile lookup
	pStatus, pCT, pBody, err := o.profile.GetProfileByUUID(ctx, uuid, profileDelayMs, profileFail)
	span.SetAttributes(attribute.Int("downstream.profile.status", pStatus))

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "profile lookup failed")
		return 0, "", nil, err
	}
	if pStatus != http.StatusOK {
		span.SetStatus(codes.Error, "profile returned non-200")
	}

	return pStatus, pCT, pBody, nil
}