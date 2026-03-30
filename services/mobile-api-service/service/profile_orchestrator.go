package service

import (
	"net/http"

	"mobile-api-service/pkg/profileclient"
	"mobile-api-service/pkg/usersclient"
)

type Orchestrator struct {
	users   *usersclient.Client
	profile *profileclient.Client
}

func NewOrchestrator(httpClient *http.Client, usersURL, profileURL string) *Orchestrator {
	return &Orchestrator{
		users:   usersclient.New(httpClient, usersURL),
		profile: profileclient.New(httpClient, profileURL),
	}
}

func (o *Orchestrator) FetchProfileByUsername(username string) (int, string, []byte, error) {
	// 1) get uuid from users-service
	status, ct, body, uuid, err := o.users.GetUUIDByUsername(username)
	if err != nil {
		return 0, "", nil, err
	}
	if status != http.StatusOK {
		// pass-through users error (e.g., 404 user not found)
		return status, ct, body, nil
	}

	// 2) get profile from profile-service
	pStatus, pCT, pBody, err := o.profile.GetProfileByUUID(uuid)
	if err != nil {
		return 0, "", nil, err
	}
	return pStatus, pCT, pBody, nil
}