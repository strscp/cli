//go:build integration

package integration

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/strscp/cli/internal/api"
)

func newTestClient(t *testing.T) *api.Client {
	t.Helper()

	token := os.Getenv("STARSCOPE_TEST_TOKEN")
	if token == "" {
		t.Skip("STARSCOPE_TEST_TOKEN not set, skipping integration tests")
	}

	apiURL := os.Getenv("STARSCOPE_TEST_API_URL")
	if apiURL == "" {
		apiURL = "https://starscope.app/api/v1"
	}

	return api.NewClient(apiURL, token, api.WithUserAgent("starscope-cli/integration-test"))
}

func TestIntegration_WorkspaceShow(t *testing.T) {
	client := newTestClient(t)

	ws, err := client.Workspace.Show(context.Background())
	require.NoError(t, err)

	assert.NotZero(t, ws.ID)
	assert.NotEmpty(t, ws.Name)
	assert.NotEmpty(t, ws.Slug)
	assert.NotEmpty(t, ws.PlanTier)
}

func TestIntegration_ConnectionsList(t *testing.T) {
	client := newTestClient(t)

	resp, err := client.Connections.List(context.Background())
	require.NoError(t, err)

	assert.NotNil(t, resp)
	assert.GreaterOrEqual(t, resp.Meta.Total, 0)

	for _, conn := range resp.Data {
		assert.NotZero(t, conn.ID)
		assert.NotEmpty(t, conn.Platform)
		assert.NotEmpty(t, conn.Name)
	}
}

func TestIntegration_ReviewsList(t *testing.T) {
	client := newTestClient(t)

	resp, err := client.Reviews.List(context.Background(), api.ListReviewsParams{
		PerPage: 5,
		Page:    1,
	})
	require.NoError(t, err)

	assert.NotNil(t, resp)
	assert.LessOrEqual(t, len(resp.Data), 5)

	for _, review := range resp.Data {
		assert.NotZero(t, review.ID)
		assert.NotEmpty(t, review.Platform)
		assert.NotEmpty(t, review.AuthorName)
		assert.GreaterOrEqual(t, review.Rating, 1)
		assert.LessOrEqual(t, review.Rating, 5)
		assert.NotEmpty(t, review.PublishedAt)
	}
}

func TestIntegration_ReviewsListWithFilters(t *testing.T) {
	client := newTestClient(t)

	resp, err := client.Reviews.List(context.Background(), api.ListReviewsParams{
		RatingMin: 4,
		PerPage:   5,
	})
	require.NoError(t, err)

	for _, review := range resp.Data {
		assert.GreaterOrEqual(t, review.Rating, 4)
	}
}

func TestIntegration_ReviewShow(t *testing.T) {
	client := newTestClient(t)

	// First get a review ID from the list
	resp, err := client.Reviews.List(context.Background(), api.ListReviewsParams{PerPage: 1})
	require.NoError(t, err)

	if len(resp.Data) == 0 {
		t.Skip("No reviews available")
	}

	review, err := client.Reviews.Get(context.Background(), resp.Data[0].ID)
	require.NoError(t, err)

	assert.Equal(t, resp.Data[0].ID, review.ID)
	assert.NotEmpty(t, review.Platform)
}

func TestIntegration_InsightsList(t *testing.T) {
	client := newTestClient(t)

	resp, err := client.Insights.List(context.Background(), api.ListInsightsParams{
		PerPage: 5,
	})
	require.NoError(t, err)

	assert.NotNil(t, resp)

	for _, insight := range resp.Data {
		assert.NotZero(t, insight.ID)
		assert.NotEmpty(t, insight.Type)
		assert.NotEmpty(t, insight.Title)
		assert.NotEmpty(t, insight.Severity)
	}
}

func TestIntegration_TopicsList(t *testing.T) {
	client := newTestClient(t)

	resp, err := client.Topics.List(context.Background(), api.ListTopicsParams{
		PerPage: 5,
	})
	require.NoError(t, err)

	assert.NotNil(t, resp)

	for _, topic := range resp.Data {
		assert.NotZero(t, topic.ID)
		assert.NotEmpty(t, topic.Label)
		assert.GreaterOrEqual(t, topic.ReviewCount, 0)
	}
}

func TestIntegration_AnalyticsOverview(t *testing.T) {
	client := newTestClient(t)

	overview, err := client.Analytics.Overview(context.Background(), api.OverviewParams{})
	require.NoError(t, err)

	assert.GreaterOrEqual(t, overview.TotalReviews, 0)
	assert.GreaterOrEqual(t, overview.AverageRating, 0.0)
}

func TestIntegration_AnalyticsMetrics(t *testing.T) {
	client := newTestClient(t)

	metrics, err := client.Analytics.Metrics(context.Background(), api.MetricsParams{
		Period: 7,
	})
	require.NoError(t, err)

	assert.NotNil(t, metrics)
}

func TestIntegration_InvalidToken(t *testing.T) {
	apiURL := os.Getenv("STARSCOPE_TEST_API_URL")
	if apiURL == "" {
		apiURL = "https://starscope.app/api/v1"
	}

	client := api.NewClient(apiURL, "invalid-token")
	_, err := client.Workspace.Show(context.Background())

	assert.Error(t, err)
	assert.IsType(t, &api.AuthenticationError{}, err)
}
