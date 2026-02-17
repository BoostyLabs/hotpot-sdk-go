package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	// authHeader defines the header used for authentication in API requests.
	authHeader string = "X-Api-Key"

	// affiliatePrefix defines the string prefix used for identifying or tagging affiliate-related entities or operations.
	affiliatePrefix string = "affiliate"
	// versionPrefix defines the prefix used to indicate the version of the application or API.
	versionPrefix string = "v1"
)

// Client represents a client for interacting with the API.
type Client struct {
	baseURL string
	apiKey  string

	inner *http.Client
}

// New creates a new Client.
func New(baseURL string, apiKey string, client *http.Client) *Client {
	return &Client{
		baseURL: strings.TrimSuffix(baseURL, "/"),
		apiKey:  apiKey,
		inner:   client,
	}
}

// NewDefault creates a new Client with a default timeout.
func NewDefault(baseURL string, apiKey string) *Client {
	return New(baseURL, apiKey, &http.Client{Timeout: 10 * time.Second})
}

// doRequest performs an HTTP request with the specified context, method, endpoint, body, and response placeholders.
func (c *Client) doRequest(ctx context.Context, method, endpoint string, body any, response any) (err error) {
	var reqBody io.Reader = http.NoBody
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal body: %w", err)
		}

		reqBody = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, endpoint, reqBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(authHeader, c.apiKey)

	resp, err := c.inner.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer func() { err = errors.Join(err, resp.Body.Close()) }()

	if resp.StatusCode != http.StatusOK {
		var responseBody []byte
		responseBody, err = io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %w", err)
		}

		var errResp = &ApiError{}
		if err = json.Unmarshal(responseBody, errResp); err != nil {
			return fmt.Errorf("failed to decode error response: %w; raw response: %s", err, string(responseBody))
		}

		return errResp
	}

	if response != nil {
		if err = json.NewDecoder(resp.Body).Decode(response); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}

// buildURL builds a URL for the given pattern with substitutions.
//
// Examples:
//
//	buildURL("intents/%s?limit=%d", "abcd", 1) -> "<baseURL>/affiliate/v1/intents/abcd?limit=1"
//	buildURL("intents") -> "<baseURL>/affiliate/v1/intents"
func (c *Client) buildURL(pathPattern string, args ...any) string {
	path := fmt.Sprintf(pathPattern, args...)
	return fmt.Sprintf("%s/%s/%s/%s", c.baseURL, affiliatePrefix, versionPrefix, path)
}
