package api

import (
	"context"
	"fmt"
	"net/url"

	"github.com/strscp/cli/pkg/models"
)

// ReviewsService handles review API operations.
type ReviewsService struct {
	client *Client
}

// ListReviewsParams defines filters for listing reviews.
type ListReviewsParams struct {
	ConnectionID int
	Platform     string
	RatingMin    int
	RatingMax    int
	DateFrom     string
	DateTo       string
	PerPage      int
	Page         int
}

func (p ListReviewsParams) toValues() url.Values {
	v := make(url.Values)
	if p.ConnectionID > 0 {
		v.Set("connection_id", fmt.Sprintf("%d", p.ConnectionID))
	}
	if p.Platform != "" {
		v.Set("platform", p.Platform)
	}
	if p.RatingMin > 0 {
		v.Set("rating_min", fmt.Sprintf("%d", p.RatingMin))
	}
	if p.RatingMax > 0 {
		v.Set("rating_max", fmt.Sprintf("%d", p.RatingMax))
	}
	if p.DateFrom != "" {
		v.Set("date_from", p.DateFrom)
	}
	if p.DateTo != "" {
		v.Set("date_to", p.DateTo)
	}
	if p.PerPage > 0 {
		v.Set("per_page", fmt.Sprintf("%d", p.PerPage))
	}
	if p.Page > 0 {
		v.Set("page", fmt.Sprintf("%d", p.Page))
	}
	return v
}

// List returns a paginated list of reviews.
func (s *ReviewsService) List(ctx context.Context, params ListReviewsParams) (*models.PaginatedResponse[models.Review], error) {
	var resp models.PaginatedResponse[models.Review]
	if err := s.client.get(ctx, "/reviews", params.toValues(), &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListAll fetches all pages of reviews.
func (s *ReviewsService) ListAll(ctx context.Context, params ListReviewsParams) ([]models.Review, error) {
	return FetchAllPages[models.Review](ctx, s.client, "/reviews", params.toValues())
}

// Get returns a single review by ID.
func (s *ReviewsService) Get(ctx context.Context, id int) (*models.Review, error) {
	var wrapper struct {
		Data models.Review `json:"data"`
	}
	if err := s.client.get(ctx, fmt.Sprintf("/reviews/%d", id), nil, &wrapper); err != nil {
		return nil, err
	}
	return &wrapper.Data, nil
}
