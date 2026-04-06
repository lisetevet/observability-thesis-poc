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
	span.SetAttributes(attribute.String("app.username", username))
	defer span.End()
	
	// 1) users-service lookup (with optional injection)
	status, ct, body, uuid, err := o.users.GetUUIDByUsername(ctx, username, usersDelayMs, usersFail)
	span.SetAttributes(attribute.Int("users.status_code", status))
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return 0, "", nil, err
	}
	if status != http.StatusOK {
		// pass-through error response from users-service
		return status, ct, body, nil
	}

	// 2) profile-service lookup (with optional injection)
	pStatus, pCT, pBody, err := o.profile.GetProfileByUUID(ctx, uuid, profileDelayMs, profileFail)
	span.SetAttributes(attribute.Int("profile.status_code", pStatus))
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return 0, "", nil, err
	}

	return pStatus, pCT, pBody, nil
}