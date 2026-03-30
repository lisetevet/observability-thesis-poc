package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type UsersLookupResponse struct {
	Username string `json:"username"`
	UUID     string `json:"uuid"`
}

type Orchestrator struct {
	Client            *http.Client
	UsersServiceURL   string
	ProfileServiceURL string
}

func NewOrchestrator(client *http.Client, usersURL, profileURL string) *Orchestrator {
	return &Orchestrator{
		Client:            client,
		UsersServiceURL:   usersURL,
		ProfileServiceURL: profileURL,
	}
}

// FetchProfileByUsername returns: statusCode, contentType, bodyBytes, error
func (o *Orchestrator) FetchProfileByUsername(username string) (int, string, []byte, error) {
	// 1) users lookup
	usersURL := fmt.Sprintf("%s/%s", o.UsersServiceURL, username)

	usersResp, err := o.Client.Get(usersURL)
	if err != nil {
		return 0, "", nil, fmt.Errorf("users-service request failed: %w", err)
	}
	defer usersResp.Body.Close()

	usersBody, _ := io.ReadAll(usersResp.Body)
	if usersResp.StatusCode != http.StatusOK {
		return usersResp.StatusCode, usersResp.Header.Get("Content-Type"), usersBody, nil
	}

	var lookup UsersLookupResponse
	if err := json.Unmarshal(usersBody, &lookup); err != nil || lookup.UUID == "" {
		return http.StatusBadGateway, "application/json", []byte(`{"error":"invalid users-service response"}`), nil
	}

	// 2) profile lookup
	profileURL := fmt.Sprintf("%s/%s", o.ProfileServiceURL, lookup.UUID)

	profResp, err := o.Client.Get(profileURL)
	if err != nil {
		return 0, "", nil, fmt.Errorf("profile-service request failed: %w", err)
	}
	defer profResp.Body.Close()

	profBody, _ := io.ReadAll(profResp.Body)
	return profResp.StatusCode, profResp.Header.Get("Content-Type"), profBody, nil
}