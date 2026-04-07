package client

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"tech-ip-sem2/shared/httpx"
	"tech-ip-sem2/shared/requestctx"

	"go.uber.org/zap"
)

type HttpAuthenticationClient struct {
	httpClient *httpx.HTTPClient
	log        *zap.Logger
}

func NewHttpAuthenticationClient(baseURL string, timeout time.Duration, log *zap.Logger) *HttpAuthenticationClient {
	return &HttpAuthenticationClient{
		httpClient: httpx.NewHTTPClient(baseURL, 3*time.Second),
		log:        log,
	}
}

func (c *HttpAuthenticationClient) Verify(ctx context.Context, token string) (*string, error) {
	headers := map[string]string{
		"Authorization": "Bearer " + token,
		"X-Request-ID":  requestctx.RequestID(ctx),
	}

	out, err := c.httpClient.Get(ctx, "/v1/auth/verify", headers)
	if err != nil {
		return nil, fmt.Errorf("Authentication service request failed: %w", err)
	}
	defer out.Body.Close()

	var verifyResponse verifyResponse
	if err := json.NewDecoder(out.Body).Decode(&verifyResponse); err != nil {
		c.log.Error("failed to decode authentication response",
			zap.String("component", "auth_client"),
			zap.String("request_id", requestctx.RequestID(ctx)),
			zap.Error(err),
		)
		return nil, fmt.Errorf("Authentication service returned invalid JSON: %w", err)
	}
	if !verifyResponse.Valid || verifyResponse.Subject == nil {
		return nil, nil
	}
	return verifyResponse.Subject, nil
}

type verifyResponse struct {
	Valid   bool    `json:"valid"`
	Subject *string `json:"subject,omitempty"`
	Error   *string `json:"error,omitempty"`
}
