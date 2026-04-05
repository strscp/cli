package api

import (
	"context"
	"fmt"
	"net/url"

	"github.com/strscp/cli/pkg/models"
)

// TopicsService handles topic API operations.
type TopicsService struct {
	client *Client
}

// ListTopicsParams defines filters for listing topics.
type ListTopicsParams struct {
	ConnectionID int
	PerPage      int
	Page         int
}

func (p ListTopicsParams) ToValues() url.Values {
	v := make(url.Values)
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

// List returns a paginated list of topics.
func (s *TopicsService) List(ctx context.Context, params ListTopicsParams) (*models.PaginatedResponse[models.Topic], error) {
	var resp models.PaginatedResponse[models.Topic]
	if err := s.client.get(ctx, "/topics", params.ToValues(), &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListAll fetches all pages of topics.
func (s *TopicsService) ListAll(ctx context.Context, params ListTopicsParams) ([]models.Topic, error) {
	return FetchAllPages[models.Topic](ctx, s.client, "/topics", params.ToValues())
}

// Get returns a single topic by ID.
func (s *TopicsService) Get(ctx context.Context, id int) (*models.Topic, error) {
	var wrapper struct {
		Data models.Topic `json:"data"`
	}
	if err := s.client.get(ctx, fmt.Sprintf("/topics/%d", id), nil, &wrapper); err != nil {
		return nil, err
	}
	return &wrapper.Data, nil
}

// Reviews returns a paginated list of reviews for a topic.
func (s *TopicsService) Reviews(ctx context.Context, id int, perPage, page int) (*models.PaginatedResponse[models.Review], error) {
	params := make(url.Values)
	if perPage > 0 {
		params.Set("per_page", fmt.Sprintf("%d", perPage))
	}
	if page > 0 {
		params.Set("page", fmt.Sprintf("%d", page))
	}

	var resp models.PaginatedResponse[models.Review]
	if err := s.client.get(ctx, fmt.Sprintf("/topics/%d/reviews", id), params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
