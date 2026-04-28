package profileclient

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	httpClient *http.Client
	baseURL    string
}

func New(httpClient *http.Client, baseURL string) *Client {
	return &Client{httpClient: httpClient, baseURL: baseURL}
}

func (c *Client) GetProfileByUsername(ctx context.Context, username, usersDelayMs, usersFail, profileDelayMs, profileFail string) (int, string, []byte, error) {
	base := strings.TrimRight(c.baseURL, "/")
	rawURL := fmt.Sprintf("%s/%s", base, url.PathEscape(username))

	u, err := url.Parse(rawURL)
	if err != nil {
		return 0, "", nil, fmt.Errorf("invalid profile-service url: %w", err)
	}

	q := u.Query()
	if usersDelayMs != "" {
		q.Set("usersDelayMs", usersDelayMs)
	}
	if usersFail == "true" {
		q.Set("usersFail", "true")
	}
	if profileDelayMs != "" {
		q.Set("delayMs", profileDelayMs)
	}
	if profileFail == "true" {
		q.Set("fail", "true")
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return 0, "", nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, "", nil, fmt.Errorf("profile-service request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, "", nil, fmt.Errorf("failed to read profile-service response body: %w", err)
	}

	contentType := resp.Header.Get("Content-Type")

	return resp.StatusCode, contentType, body, nil
}