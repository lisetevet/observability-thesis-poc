package usersclient

import (
	"encoding/json"
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

type LookupResponse struct {
	Username string `json:"username"`
	UUID     string `json:"uuid"`
}

// GetUUIDByUsername returns: statusCode, contentType, bodyBytes, uuid(if status 200)
func (c *Client) GetUUIDByUsername(username string) (int, string, []byte, string, error) {
	url := fmt.Sprintf("%s/%s", c.baseURL, username)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return 0, "", nil, "", fmt.Errorf("users-service request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	ct := resp.Header.Get("Content-Type")

	if resp.StatusCode != http.StatusOK {
		return resp.StatusCode, ct, body, "", nil
	}

	var lr LookupResponse
	if err := json.Unmarshal(body, &lr); err != nil || lr.UUID == "" {
		return http.StatusBadGateway, "application/json", []byte(`{"error":"invalid users-service response"}`), "", nil
	}

	return resp.StatusCode, ct, body, lr.UUID, nil
}