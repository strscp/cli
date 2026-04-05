package api

import (
	"fmt"
	"time"
)

// APIError represents a generic API error response.
type APIError struct {
	StatusCode int
	Message    string
	Errors     map[string][]string
}

func (e *APIError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("API error (%d): %s", e.StatusCode, e.Message)
	}
	return fmt.Sprintf("API error (%d)", e.StatusCode)
}

// AuthenticationError indicates invalid or missing credentials.
type AuthenticationError struct{}

func (e *AuthenticationError) Error() string {
	return "authentication failed: invalid or missing API token"
}

// ForbiddenError indicates the request is not allowed.
type ForbiddenError struct {
	Message string
}

func (e *ForbiddenError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("forbidden: %s", e.Message)
	}
	return "forbidden: you don't have access to this resource"
}

// NotFoundError indicates the requested resource was not found.
type NotFoundError struct {
	Resource string
	ID       int
}

func (e *NotFoundError) Error() string {
	if e.Resource != "" {
		return fmt.Sprintf("%s with ID %d not found", e.Resource, e.ID)
	}
	return "resource not found"
}

// RateLimitError indicates the API rate limit has been exceeded.
type RateLimitError struct {
	RetryAfter time.Duration
}

func (e *RateLimitError) Error() string {
	return fmt.Sprintf("rate limit exceeded, retry after %s", e.RetryAfter)
}

// ValidationError indicates invalid request parameters.
type ValidationError struct {
	Errors map[string][]string
}

func (e *ValidationError) Error() string {
	return "validation failed"
}

// apiErrorResponse matches the JSON error shape from the Starscope API.
type apiErrorResponse struct {
	Message string              `json:"message"`
	Errors  map[string][]string `json:"errors,omitempty"`
}
