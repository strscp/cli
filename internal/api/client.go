package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// Client is the Starscope API HTTP client.
type Client struct {
	httpClient  *http.Client
	baseURL     string
	token       string
	workspaceID int
	userAgent   string

	Reviews     *ReviewsService
	Insights    *InsightsService
	Topics      *TopicsService
	Connections *ConnectionsService
	Analytics   *AnalyticsService
	Workspace   *WorkspaceService
}

// ClientOption configures the Client.
type ClientOption func(*Client)

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(hc *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = hc
	}
}

// WithWorkspaceID sets the workspace ID header.
func WithWorkspaceID(id int) ClientOption {
	return func(c *Client) {
		c.workspaceID = id
	}
}

// WithUserAgent sets a custom User-Agent.
func WithUserAgent(ua string) ClientOption {
	return func(c *Client) {
		c.userAgent = ua
	}
}

// NewClient creates a new API client.
func NewClient(baseURL, token string, opts ...ClientOption) *Client {
	c := &Client{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		baseURL:    baseURL,
		token:      token,
		userAgent:  "starscope-cli",
	}

	for _, opt := range opts {
		opt(c)
	}

	c.Reviews = &ReviewsService{client: c}
	c.Insights = &InsightsService{client: c}
	c.Topics = &TopicsService{client: c}
	c.Connections = &ConnectionsService{client: c}
	c.Analytics = &AnalyticsService{client: c}
	c.Workspace = &WorkspaceService{client: c}

	return c
}

// newRequest creates an authenticated HTTP request.
func (c *Client) newRequest(ctx context.Context, method, path string, params url.Values) (*http.Request, error) {
	u := c.baseURL + path
	if params != nil && len(params) > 0 {
		u += "?" + params.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, method, u, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.userAgent)

	if c.workspaceID > 0 {
		req.Header.Set("X-Workspace-Id", strconv.Itoa(c.workspaceID))
	}

	return req, nil
}

// do executes a request and handles common error responses.
func (c *Client) do(req *http.Request) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return resp, nil

	case http.StatusUnauthorized:
		resp.Body.Close()
		return nil, &AuthenticationError{}

	case http.StatusForbidden:
		defer resp.Body.Close()
		var errResp apiErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err == nil {
			return nil, &ForbiddenError{Message: errResp.Message}
		}
		return nil, &ForbiddenError{}

	case http.StatusNotFound:
		resp.Body.Close()
		return nil, &NotFoundError{}

	case http.StatusUnprocessableEntity:
		defer resp.Body.Close()
		var errResp apiErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err == nil {
			return nil, &ValidationError{Errors: errResp.Errors}
		}
		return nil, &ValidationError{}

	case http.StatusTooManyRequests:
		resp.Body.Close()
		retryAfter := 60 * time.Second
		if ra := resp.Header.Get("Retry-After"); ra != "" {
			if secs, err := strconv.Atoi(ra); err == nil {
				retryAfter = time.Duration(secs) * time.Second
			}
		}
		return nil, &RateLimitError{RetryAfter: retryAfter}

	default:
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		return nil, &APIError{
			StatusCode: resp.StatusCode,
			Message:    string(body),
		}
	}
}

// get performs a GET request and decodes the JSON response.
func (c *Client) get(ctx context.Context, path string, params url.Values, v any) error {
	req, err := c.newRequest(ctx, http.MethodGet, path, params)
	if err != nil {
		return err
	}

	resp, err := c.do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if v != nil {
		if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
			return fmt.Errorf("decoding response: %w", err)
		}
	}

	return nil
}

// getRaw performs a GET request and returns the raw JSON bytes.
func (c *Client) getRaw(ctx context.Context, path string, params url.Values) ([]byte, error) {
	req, err := c.newRequest(ctx, http.MethodGet, path, params)
	if err != nil {
		return nil, err
	}

	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
