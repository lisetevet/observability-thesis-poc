package profileclient

import (
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	httpClient *http.Client
	baseURL    string
}

func New(httpClient *http.Client, baseURL string) *Client {
	return &Client{httpClient: httpClient, baseURL: baseURL}
}

// GetProfileByUUID returns: statusCode, contentType, bodyBytes
func (c *Client) GetProfileByUUID(uuid string) (int, string, []byte, error) {
	url := fmt.Sprintf("%s/%s", c.baseURL, uuid)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return 0, "", nil, fmt.Errorf("profile-service request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	ct := resp.Header.Get("Content-Type")

	return resp.StatusCode, ct, body, nil
}