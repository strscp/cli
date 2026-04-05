package api

import (
	"context"
	"fmt"
	"net/url"

	"github.com/strscp/cli/pkg/models"
)

// InsightsService handles insight API operations.
type InsightsService struct {
	client *Client
}

// ListInsightsParams defines filters for listing insights.
type ListInsightsParams struct {
	Type         string
	Severity     string
	ConnectionID int
	PerPage      int
	Page         int
}

func (p ListInsightsParams) ToValues() url.Values {
	v := make(url.Values)
	if p.Type != "" {
		v.Set("type", p.Type)
	}
	if p.Severity != "" {
		v.Set("severity", p.Severity)
	}
	if p.ConnectionID > 0 {
		v.Set("connection_id", fmt.Sprintf("%d", p.ConnectionID))
	}
	if p.PerPage > 0 {
		v.Set("per_page", fmt.Sprintf("%d", p.PerPage))
	}
	if p.Page > 0 {
		v.Set("page", fmt.Sprintf("%d", p.Page))
	}
	return v
}

// List returns a paginated list of insights.
func (s *InsightsService) List(ctx context.Context, params ListInsightsParams) (*models.PaginatedResponse[models.Insight], error) {
	var resp models.PaginatedResponse[models.Insight]
	if err := s.client.get(ctx, "/insights", params.ToValues(), &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListAll fetches all pages of insights.
func (s *InsightsService) ListAll(ctx context.Context, params ListInsightsParams) ([]models.Insight, error) {
	return FetchAllPages[models.Insight](ctx, s.client, "/insights", params.ToValues())
}

// Get returns a single insight by ID.
func (s *InsightsService) Get(ctx context.Context, id int) (*models.Insight, error) {
	var wrapper struct {
		Data models.Insight `json:"data"`
	}
	if err := s.client.get(ctx, fmt.Sprintf("/insights/%d", id), nil, &wrapper); err != nil {
		return nil, err
	}
	return &wrapper.Data, nil
}
