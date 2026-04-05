package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"time"

	"github.com/strscp/cli/pkg/models"
)

// FetchAllPages iterates through all pages of a paginated endpoint,
// collecting all items. It paces requests at 1/sec to respect rate limits.
// Progress is printed to stderr.
func FetchAllPages[T any](ctx context.Context, c *Client, path string, params url.Values) ([]T, error) {
	return FetchAllPagesWithProgress[T](ctx, c, path, params, os.Stderr)
}

// FetchAllPagesWithProgress iterates through all pages with progress output.
func FetchAllPagesWithProgress[T any](ctx context.Context, c *Client, path string, params url.Values, progress io.Writer) ([]T, error) {
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

		if progress != nil && resp.Meta.LastPage > 1 {
			fmt.Fprintf(progress, "\rFetching page %d/%d (%d items)...", resp.Meta.CurrentPage, resp.Meta.LastPage, len(all))
		}

		if resp.Meta.CurrentPage >= resp.Meta.LastPage {
			if progress != nil && resp.Meta.LastPage > 1 {
				fmt.Fprintf(progress, "\rFetched %d items across %d pages.    \n", len(all), resp.Meta.LastPage)
			}
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
