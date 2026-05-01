package usersclient

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

type Client struct {
	httpClient *http.Client
	baseURL    string
}

func New(httpClient *http.Client, baseURL string) *Client {
	return &Client{httpClient: httpClient, baseURL: baseURL}
}

func (c *Client) Get(
	ctx *gin.Context,
	path string,
	query url.Values,
) (statusCode int, contentType string, responseBody []byte, err error) {
	base := strings.TrimRight(c.baseURL, "/")
	cleanPath := "/" + strings.TrimLeft(path, "/")

	u, err := url.Parse(base + cleanPath)
	if err != nil {
		return 0, "", nil, fmt.Errorf("invalid users-service url: %w", err)
	}

	u.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx.Request.Context(), http.MethodGet, u.String(), nil)
	if err != nil {
		return 0, "", nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, "", nil, fmt.Errorf("users-service request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, "", nil, fmt.Errorf("failed to read users-service response body: %w", err)
	}

	contentType = resp.Header.Get("Content-Type")
	return resp.StatusCode, contentType, body, nil
}
