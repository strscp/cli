package api

import (
	"context"
	"fmt"

	"github.com/strscp/cli/pkg/models"
)

// ConnectionsService handles connection API operations.
type ConnectionsService struct {
	client *Client
}

// List returns a paginated list of connections.
func (s *ConnectionsService) List(ctx context.Context) (*models.PaginatedResponse[models.Connection], error) {
	var resp models.PaginatedResponse[models.Connection]
	if err := s.client.get(ctx, "/connections", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Get returns a single connection by ID.
func (s *ConnectionsService) Get(ctx context.Context, id int) (*models.Connection, error) {
	var wrapper struct {
		Data models.Connection `json:"data"`
	}
	if err := s.client.get(ctx, fmt.Sprintf("/connections/%d", id), nil, &wrapper); err != nil {
		return nil, err
	}
	return &wrapper.Data, nil
}
