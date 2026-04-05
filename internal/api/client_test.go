package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/strscp/cli/pkg/models"
)

func testServer(t *testing.T, handler http.HandlerFunc) (*httptest.Server, *Client) {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	client := NewClient(srv.URL, "test-token", WithWorkspaceID(1))
	return srv, client
}

func TestClient_AuthHeader(t *testing.T) {
	var gotAuth string
	_, client := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{"id": 1, "name": "Test"}})
	})

	_, err := client.Workspace.Show(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "Bearer test-token", gotAuth)
}

func TestClient_WorkspaceHeader(t *testing.T) {
	var gotHeader string
	_, client := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotHeader = r.Header.Get("X-Workspace-Id")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{"id": 1, "name": "Test"}})
	})

	_, err := client.Workspace.Show(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "1", gotHeader)
}

func TestClient_NoWorkspaceHeader(t *testing.T) {
	var gotHeader string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotHeader = r.Header.Get("X-Workspace-Id")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{"id": 1, "name": "Test"}})
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "test-token") // no workspace ID
	_, err := client.Workspace.Show(context.Background())
	require.NoError(t, err)
	assert.Empty(t, gotHeader)
}

func TestClient_401(t *testing.T) {
	_, client := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	})

	_, err := client.Workspace.Show(context.Background())
	assert.IsType(t, &AuthenticationError{}, err)
}

func TestClient_403(t *testing.T) {
	_, client := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]any{"message": "API access requires a Pro plan"})
	})

	_, err := client.Workspace.Show(context.Background())
	assert.IsType(t, &ForbiddenError{}, err)
	assert.Contains(t, err.Error(), "Pro plan")
}

func TestClient_404(t *testing.T) {
	_, client := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := client.Reviews.Get(context.Background(), 999)
	assert.IsType(t, &NotFoundError{}, err)
}

func TestClient_422(t *testing.T) {
	_, client := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]any{
			"message": "Validation failed",
			"errors":  map[string]any{"rating_min": []string{"must be between 1 and 5"}},
		})
	})

	_, err := client.Reviews.List(context.Background(), ListReviewsParams{RatingMin: 10})
	assert.IsType(t, &ValidationError{}, err)
}

func TestClient_429(t *testing.T) {
	_, client := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Retry-After", "30")
		w.WriteHeader(http.StatusTooManyRequests)
	})

	_, err := client.Reviews.List(context.Background(), ListReviewsParams{})
	assert.IsType(t, &RateLimitError{}, err)
	rle := err.(*RateLimitError)
	assert.Equal(t, 30, int(rle.RetryAfter.Seconds()))
}

func TestReviewsList(t *testing.T) {
	_, client := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/reviews", r.URL.Path)
		assert.Equal(t, "trustpilot", r.URL.Query().Get("platform"))
		assert.Equal(t, "4", r.URL.Query().Get("rating_min"))

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(models.PaginatedResponse[models.Review]{
			Data: []models.Review{
				{ID: 1, Platform: "trustpilot", AuthorName: "John", Rating: 5},
			},
			Meta: models.PaginationMeta{CurrentPage: 1, LastPage: 1, PerPage: 15, Total: 1},
		})
	})

	resp, err := client.Reviews.List(context.Background(), ListReviewsParams{
		Platform:  "trustpilot",
		RatingMin: 4,
	})
	require.NoError(t, err)
	assert.Len(t, resp.Data, 1)
	assert.Equal(t, "trustpilot", resp.Data[0].Platform)
	assert.Equal(t, 5, resp.Data[0].Rating)
}

func TestReviewsGet(t *testing.T) {
	_, client := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/reviews/42", r.URL.Path)

		title := "Great service"
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"data": models.Review{ID: 42, Platform: "google", AuthorName: "Jane", Rating: 4, Title: &title},
		})
	})

	review, err := client.Reviews.Get(context.Background(), 42)
	require.NoError(t, err)
	assert.Equal(t, 42, review.ID)
	assert.Equal(t, "google", review.Platform)
	assert.Equal(t, "Great service", *review.Title)
}

func TestConnectionsList(t *testing.T) {
	_, client := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/connections", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(models.PaginatedResponse[models.Connection]{
			Data: []models.Connection{
				{ID: 1, Platform: "trustpilot", Name: "My Store", IsActive: true, ReviewCount: 150},
			},
			Meta: models.PaginationMeta{CurrentPage: 1, LastPage: 1, PerPage: 15, Total: 1},
		})
	})

	resp, err := client.Connections.List(context.Background())
	require.NoError(t, err)
	assert.Len(t, resp.Data, 1)
	assert.Equal(t, "My Store", resp.Data[0].Name)
}

func TestWorkspaceShow(t *testing.T) {
	_, client := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/workspace", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"data": models.Workspace{ID: 1, Name: "Test Workspace", Slug: "test", PlanTier: "pro"},
		})
	})

	ws, err := client.Workspace.Show(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "Test Workspace", ws.Name)
	assert.Equal(t, "pro", ws.PlanTier)
}
