package api

import (
	"context"

	"github.com/strscp/cli/pkg/models"
)

// WorkspaceService handles workspace API operations.
type WorkspaceService struct {
	client *Client
}

// Show returns the current workspace info.
func (s *WorkspaceService) Show(ctx context.Context) (*models.Workspace, error) {
	var wrapper struct {
		Data models.Workspace `json:"data"`
	}
	if err := s.client.get(ctx, "/workspace", nil, &wrapper); err != nil {
		return nil, err
	}
	return &wrapper.Data, nil
}
