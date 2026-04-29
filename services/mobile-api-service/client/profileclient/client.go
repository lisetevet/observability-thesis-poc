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

func (c *Client) Get(ctx context.Context, path string, query url.Values) (int, string, []byte, error) {
	base := strings.TrimRight(c.baseURL, "/")
	cleanPath := "/" + strings.TrimLeft(path, "/")

	u, err := url.Parse(base + cleanPath)
	if err != nil {
		return 0, "", nil, fmt.Errorf("invalid profile-service url: %w", err)
	}

	u.RawQuery = query.Encode()

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
