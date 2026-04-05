package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/strscp/cli/pkg/models"
)

// FetchAllPages iterates through all pages of a paginated endpoint,
// collecting all items. It paces requests at 1/sec to respect rate limits.
func FetchAllPages[T any](ctx context.Context, c *Client, path string, params url.Values) ([]T, error) {
	var all []T
	page := 1

	for {
		if ctx.Err() != nil {
			return all, ctx.Err()
		}

		p := cloneParams(params)
		p.Set("page", fmt.Sprintf("%d", page))

		var resp models.PaginatedResponse[T]
		if err := c.get(ctx, path, p, &resp); err != nil {
			return all, err
		}

		all = append(all, resp.Data...)

		if resp.Meta.CurrentPage >= resp.Meta.LastPage {
			break
		}

		page++

		// Pace requests to stay under rate limit
		select {
		case <-ctx.Done():
			return all, ctx.Err()
		case <-time.After(1 * time.Second):
		}
	}

	return all, nil
}

// FetchPage fetches a single page and returns the raw JSON for passthrough.
func FetchPage(ctx context.Context, c *Client, path string, params url.Values) (json.RawMessage, error) {
	raw, err := c.getRaw(ctx, path, params)
	if err != nil {
		return nil, err
	}
	return json.RawMessage(raw), nil
}

func cloneParams(params url.Values) url.Values {
	p := make(url.Values)
	for k, v := range params {
		p[k] = append([]string{}, v...)
	}
	return p
}
