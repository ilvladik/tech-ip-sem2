package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"tech-ip-sem2/shared/httpx"
)

type Client struct {
	httpClient *httpx.HTTPClient
}

type VerifyResponse struct {
	Valid   bool   `json:"valid"`
	Subject string `json:"subject,omitempty"`
	Error   string `json:"error,omitempty"`
}

func NewClient(baseURL string, timeout time.Duration) *Client {
	return &Client{
		httpClient: httpx.NewHTTPClient(baseURL, timeout),
	}
}

func (c *Client) VerifyToken(ctx context.Context, token string, requestId string) (bool, string, error) {
	headers := map[string]string{
		"Authorization": "Bearer " + token,
		"X-Request-ID":  requestId,
	}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	resp, err := c.httpClient.Get(ctx, "/v1/auth/verify", headers)
	if err != nil {
		return false, "", fmt.Errorf("auth service unavailable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, "", nil
	}

	var verifyResp VerifyResponse
	if err := json.NewDecoder(resp.Body).Decode(&verifyResp); err != nil {
		return false, "", fmt.Errorf("failed to decode auth response: %w", err)
	}

	return verifyResp.Valid, verifyResp.Subject, nil
}
