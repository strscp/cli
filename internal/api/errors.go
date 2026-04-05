package api

import (
	"fmt"
	"time"
)

// Exit codes for different error types.
const (
	ExitOK             = 0
	ExitAuth           = 1
	ExitForbidden      = 2
	ExitNotFound       = 3
	ExitValidation     = 4
	ExitRateLimit      = 5
	ExitAPIError       = 6
	ExitConfigError    = 7
)

// CLIError is implemented by errors that provide structured info for agents.
type CLIError interface {
	error
	ErrorCode() string
	ExitCode() int
}

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

func (e *APIError) ErrorCode() string { return "api_error" }
func (e *APIError) ExitCode() int     { return ExitAPIError }

// AuthenticationError indicates invalid or missing credentials.
type AuthenticationError struct{}

func (e *AuthenticationError) Error() string     { return "authentication failed: invalid or missing API token" }
func (e *AuthenticationError) ErrorCode() string { return "authentication_failed" }
func (e *AuthenticationError) ExitCode() int     { return ExitAuth }

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

func (e *ForbiddenError) ErrorCode() string { return "forbidden" }
func (e *ForbiddenError) ExitCode() int     { return ExitForbidden }

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

func (e *NotFoundError) ErrorCode() string { return "not_found" }
func (e *NotFoundError) ExitCode() int     { return ExitNotFound }

// RateLimitError indicates the API rate limit has been exceeded.
type RateLimitError struct {
	RetryAfter time.Duration
}

func (e *RateLimitError) Error() string {
	return fmt.Sprintf("rate limit exceeded, retry after %s", e.RetryAfter)
}

func (e *RateLimitError) ErrorCode() string { return "rate_limit_exceeded" }
func (e *RateLimitError) ExitCode() int     { return ExitRateLimit }

// ValidationError indicates invalid request parameters.
type ValidationError struct {
	Errors map[string][]string
}

func (e *ValidationError) Error() string {
	if len(e.Errors) == 0 {
		return "validation failed"
	}
	msg := "validation failed:"
	for field, errs := range e.Errors {
		for _, err := range errs {
			msg += fmt.Sprintf("\n  %s: %s", field, err)
		}
	}
	return msg
}

func (e *ValidationError) ErrorCode() string { return "validation_failed" }
func (e *ValidationError) ExitCode() int     { return ExitValidation }

// apiErrorResponse matches the JSON error shape from the Starscope API.
type apiErrorResponse struct {
	Message string              `json:"message"`
	Errors  map[string][]string `json:"errors,omitempty"`
}
