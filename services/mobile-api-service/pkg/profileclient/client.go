package profileclient

import (
	"fmt"
	"io"
	"net/http"
	"context"
)

type Client struct {
	httpClient *http.Client
	baseURL    string
}

func New(httpClient *http.Client, baseURL string) *Client {
	return &Client{httpClient: httpClient, baseURL: baseURL}
}

// GetProfileByUUID returns: statusCode, contentType, bodyBytes
func (c *Client) GetProfileByUUID(ctx context.Context, uuid, delayMs, fail string) (int, string, []byte, error) {
	url := fmt.Sprintf("%s/%s", c.baseURL, uuid)
	q := ""
	if delayMs != "" {
		q += "delayMs=" + delayMs
	}
	if fail == "true" {
		if q != "" { q += "&" }
		q += "fail=true"
	}
	if q != "" {
		url += "?" + q
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, "", nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, "", nil, fmt.Errorf("profile-service request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	ct := resp.Header.Get("Content-Type")

	return resp.StatusCode, ct, body, nil
}