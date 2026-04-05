package api

import (
	"context"
	"fmt"
	"net/url"

	"github.com/strscp/cli/pkg/models"
)

// AnalyticsService handles analytics API operations.
type AnalyticsService struct {
	client *Client
}

// OverviewParams defines filters for the analytics overview.
type OverviewParams struct {
	ConnectionID int
}

// MetricsParams defines filters for analytics metrics.
type MetricsParams struct {
	Period       int
	ConnectionID int
}

// Overview returns the analytics overview for the workspace.
func (s *AnalyticsService) Overview(ctx context.Context, params OverviewParams) (*models.AnalyticsOverview, error) {
	v := make(url.Values)
	if params.ConnectionID > 0 {
		v.Set("connection_id", fmt.Sprintf("%d", params.ConnectionID))
	}

	var wrapper struct {
		Data models.AnalyticsOverview `json:"data"`
	}
	if err := s.client.get(ctx, "/analytics/overview", v, &wrapper); err != nil {
		return nil, err
	}
	return &wrapper.Data, nil
}

// Metrics returns analytics metrics for the workspace.
func (s *AnalyticsService) Metrics(ctx context.Context, params MetricsParams) (*models.AnalyticsMetrics, error) {
	v := make(url.Values)
	if params.Period > 0 {
		v.Set("period", fmt.Sprintf("%d", params.Period))
	}
	if params.ConnectionID > 0 {
		v.Set("connection_id", fmt.Sprintf("%d", params.ConnectionID))
	}

	var wrapper struct {
		Data models.AnalyticsMetrics `json:"data"`
	}
	if err := s.client.get(ctx, "/analytics/metrics", v, &wrapper); err != nil {
		return nil, err
	}
	return &wrapper.Data, nil
}
