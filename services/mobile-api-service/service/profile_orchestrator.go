package service

import (
	"net/http"
	"context"

	"mobile-api-service/pkg/profileclient"
	"mobile-api-service/pkg/usersclient"
)

type Orchestrator struct {
	users   *usersclient.Client
	profile *profileclient.Client
}

func NewOrchestrator(users *usersclient.Client, profile *profileclient.Client) *Orchestrator {
	return &Orchestrator{users: users, profile: profile}
}

func (o *Orchestrator) FetchProfileByUsername(ctx context.Context, username, usersDelayMs, usersFail, profileDelayMs, profileFail string) (int, string, []byte, error) {
	// 1) users-service lookup (with optional injection)
	status, ct, body, uuid, err := o.users.GetUUIDByUsername(ctx, username, usersDelayMs, usersFail)
	if err != nil {
		return 0, "", nil, err
	}
	if status != http.StatusOK {
		return status, ct, body, nil
	}

	// 2) profile-service lookup (with optional injection)
	pStatus, pCT, pBody, err := o.profile.GetProfileByUUID(ctx, uuid, profileDelayMs, profileFail)
	if err != nil {
		return 0, "", nil, err
	}
	return pStatus, pCT, pBody, nil
}