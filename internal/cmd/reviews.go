package cmd

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/strscp/cli/internal/api"
	"github.com/strscp/cli/internal/output"
	"github.com/strscp/cli/pkg/models"
)

var reviewsCmd = &cobra.Command{
	Use:     "reviews",
	Aliases: []string{"r", "rev"},
	Short:   "Manage reviews",
}

var (
	reviewConnectionID int
	reviewPlatform     string
	reviewRatingMin    int
	reviewRatingMax    int
	reviewDateFrom     string
	reviewDateTo       string
	reviewPerPage      int
	reviewPage         int
	reviewAll          bool
)

var reviewsListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List reviews",
	Example: `  starscope-cli reviews list
  starscope-cli reviews list --platform trustpilot --rating-min 4
  starscope-cli reviews list --date-from 2026-01-01 --output json
  starscope-cli r ls --all`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}
		f, err := newFormatter()
		if err != nil {
			return err
		}

		params := api.ListReviewsParams{
			ConnectionID: reviewConnectionID,
			Platform:     reviewPlatform,
			RatingMin:    reviewRatingMin,
			RatingMax:    reviewRatingMax,
			DateFrom:     reviewDateFrom,
			DateTo:       reviewDateTo,
			PerPage:      reviewPerPage,
			Page:         reviewPage,
		}

		ctx := context.Background()

		if reviewAll {
			reviews, err := client.Reviews.ListAll(ctx, params)
			if err != nil {
				return err
			}
			return formatReviewList(f, reviews)
		}

		resp, err := client.Reviews.List(ctx, params)
		if err != nil {
			return err
		}

		if err := formatReviewList(f, resp.Data); err != nil {
			return err
		}

		paginationFooter("reviews", resp.Meta)
		return nil
	},
}

var reviewsShowCmd = &cobra.Command{
	Use:     "show <id>",
	Aliases: []string{"s"},
	Short:   "Show a review",
	Example: `  starscope-cli reviews show 42`,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid review ID: %s", args[0])
		}

		client, err := newAPIClient()
		if err != nil {
			return err
		}
		f, err := newFormatter()
		if err != nil {
			return err
		}

		review, err := client.Reviews.Get(context.Background(), id)
		if err != nil {
			return err
		}

		return formatReviewDetail(f, review)
	},
}

func formatReviewList(f *output.Formatter, reviews []models.Review) error {
	headers := []string{"ID", "PLATFORM", "AUTHOR", "RATING", "PUBLISHED", "TITLE"}
	rows := make([][]string, len(reviews))
	for i, r := range reviews {
		title := ""
		if r.Title != nil {
			title = *r.Title
		}
		rows[i] = []string{
			strconv.Itoa(r.ID),
			r.Platform,
			r.AuthorName,
			strconv.Itoa(r.Rating),
			r.PublishedAt,
			title,
		}
	}
	return f.FormatList(headers, rows)
}

func formatReviewDetail(f *output.Formatter, r *models.Review) error {
	fields := []output.Field{
		{Label: "ID", Value: strconv.Itoa(r.ID)},
		{Label: "Platform", Value: r.Platform},
		{Label: "Author", Value: r.AuthorName},
		{Label: "Rating", Value: strconv.Itoa(r.Rating)},
		{Label: "Published", Value: r.PublishedAt},
		{Label: "Title", Value: deref(r.Title)},
		{Label: "Body", Value: deref(r.Body)},
		{Label: "Language", Value: deref(r.Language)},
		{Label: "Sentiment", Value: deref(r.Sentiment)},
		{Label: "Topics", Value: strings.Join(r.Topics, ", ")},
		{Label: "Keywords", Value: strings.Join(r.Keywords, ", ")},
		{Label: "Summary", Value: deref(r.Summary)},
	}

	if r.ResponseText != nil {
		fields = append(fields, output.Field{Label: "Response", Value: *r.ResponseText})
	}
	if r.RespondedAt != nil {
		fields = append(fields, output.Field{Label: "Responded At", Value: *r.RespondedAt})
	}

	return f.FormatDetail(fields)
}

func deref(s *string) string {
	if s == nil {
		return "-"
	}
	return *s
}

func init() {
	reviewsListCmd.Flags().IntVar(&reviewConnectionID, "connection-id", 0, "filter by connection ID")
	reviewsListCmd.Flags().StringVar(&reviewPlatform, "platform", "", "filter by platform (trustpilot, google, feedback_company)")
	reviewsListCmd.Flags().IntVar(&reviewRatingMin, "rating-min", 0, "minimum rating (1-5)")
	reviewsListCmd.Flags().IntVar(&reviewRatingMax, "rating-max", 0, "maximum rating (1-5)")
	reviewsListCmd.Flags().StringVar(&reviewDateFrom, "date-from", "", "start date (YYYY-MM-DD)")
	reviewsListCmd.Flags().StringVar(&reviewDateTo, "date-to", "", "end date (YYYY-MM-DD)")
	reviewsListCmd.Flags().IntVar(&reviewPerPage, "per-page", 15, "results per page (1-100)")
	reviewsListCmd.Flags().IntVar(&reviewPage, "page", 1, "page number")
	reviewsListCmd.Flags().BoolVar(&reviewAll, "all", false, "fetch all pages")

	reviewsCmd.AddCommand(reviewsListCmd)
	reviewsCmd.AddCommand(reviewsShowCmd)
	defaultToList(reviewsCmd, reviewsListCmd)
}
