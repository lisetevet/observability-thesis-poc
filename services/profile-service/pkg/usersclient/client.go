package usersclient

import (
	"context"
	"encoding/json"
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

type lookupResponse struct {
	Username string `json:"username"`
	UUID     string `json:"uuid"`
}

func (c *Client) GetUUIDByUsername(ctx context.Context, username, delayMs, fail string) (uuid string, ok bool, err error) {
	base := strings.TrimRight(c.baseURL, "/")
	rawURL := fmt.Sprintf("%s/%s", base, url.PathEscape(username))

	u, err := url.Parse(rawURL)
	if err != nil {
		return "", false, fmt.Errorf("invalid users-service url: %w", err)
	}

	q := u.Query()
	if delayMs != "" {
		q.Set("delayMs", delayMs)
	}
	if fail == "true" {
		q.Set("fail", "true")
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return "", false, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", false, fmt.Errorf("users-service request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", false, fmt.Errorf("failed to read users-service response body: %w", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return "", false, nil
	}

	if resp.StatusCode != http.StatusOK {
		return "", false, fmt.Errorf("users-service returned %d: %s", resp.StatusCode, string(body))
	}

	var lr lookupResponse
	if err := json.Unmarshal(body, &lr); err != nil {
		return "", false, fmt.Errorf("invalid users-service response: %s", string(body))
	}

	if lr.UUID == "" {
		return "", false, fmt.Errorf("users-service response did not contain uuid: %s", string(body))
	}

	return lr.UUID, true, nil
}