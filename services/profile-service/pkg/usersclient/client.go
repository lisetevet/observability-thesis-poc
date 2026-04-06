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

// GetUUIDByUsername calls users-api-service and returns uuid (or ok=false if 404).
func (c *Client) GetUUIDByUsername(ctx context.Context, username string) (uuid string, ok bool, err error) {
	url := fmt.Sprintf("%s/%s", c.baseURL, username)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", false, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", false, fmt.Errorf("users-service request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusNotFound {
		return "", false, nil
	}
	if resp.StatusCode != http.StatusOK {
		// pass through error details (useful in logs)
		return "", false, fmt.Errorf("users-service returned %d: %s", resp.StatusCode, string(body))
	}

	var lr lookupResponse
	if err := json.Unmarshal(body, &lr); err != nil || lr.UUID == "" {
		return "", false, fmt.Errorf("invalid users-service response: %s", string(body))
	}

	return lr.UUID, true, nil
}

func (c *Client) GetProfileByUsername(ctx context.Context, username, delayMs, fail string) (int, string, []byte, error) {
	base := strings.TrimRight(c.baseURL, "/")
	rawURL := fmt.Sprintf("%s/%s/profile", base, username)

	u, err := url.Parse(rawURL)
	if err != nil {
		return 0, "", nil, fmt.Errorf("invalid url: %w", err)
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
		return 0, "", nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, "", nil, fmt.Errorf("users-service profile request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	ct := resp.Header.Get("Content-Type")

	return resp.StatusCode, ct, body, nil
}