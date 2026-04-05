package cmd

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/strscp/cli/internal/api"
	"github.com/strscp/cli/internal/output"
	"github.com/strscp/cli/pkg/models"
)

var analyticsCmd = &cobra.Command{
	Use:     "analytics",
	Aliases: []string{"a", "ana"},
	Short:   "View analytics and metrics",
}

var (
	analyticsPeriod       int
	analyticsConnectionID int
)

var analyticsOverviewCmd = &cobra.Command{
	Use:   "overview",
	Short: "Show analytics overview",
	Example: `  starscope-cli analytics overview
  starscope-cli a overview --connection-id 5`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}
		f, err := newFormatter()
		if err != nil {
			return err
		}

		overview, err := client.Analytics.Overview(context.Background(), api.OverviewParams{
			ConnectionID: analyticsConnectionID,
		})
		if err != nil {
			return err
		}

		return formatOverview(f, overview)
	},
}

var analyticsMetricsCmd = &cobra.Command{
	Use:   "metrics",
	Short: "Show analytics metrics",
	Example: `  starscope-cli analytics metrics
  starscope-cli a metrics --period 30 --connection-id 5`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}
		f, err := newFormatter()
		if err != nil {
			return err
		}

		metrics, err := client.Analytics.Metrics(context.Background(), api.MetricsParams{
			Period:       analyticsPeriod,
			ConnectionID: analyticsConnectionID,
		})
		if err != nil {
			return err
		}

		return formatMetrics(f, metrics)
	},
}

func formatOverview(f *output.Formatter, o *models.AnalyticsOverview) error {
	fields := []output.Field{
		{Label: "Total Reviews", Value: strconv.Itoa(o.TotalReviews)},
		{Label: "Average Rating", Value: fmt.Sprintf("%.2f", o.AverageRating)},
		{Label: "Response Rate", Value: fmt.Sprintf("%.1f%%", o.ResponseRate*100)},
	}

	if len(o.SentimentBreakdown) > 0 {
		for k, v := range o.SentimentBreakdown {
			fields = append(fields, output.Field{
				Label: fmt.Sprintf("Sentiment: %s", k),
				Value: strconv.Itoa(v),
			})
		}
	}

	if len(o.RatingDistribution) > 0 {
		for k, v := range o.RatingDistribution {
			fields = append(fields, output.Field{
				Label: fmt.Sprintf("Rating %s", k),
				Value: strconv.Itoa(v),
			})
		}
	}

	return f.FormatDetail(fields)
}

func formatMetrics(f *output.Formatter, m *models.AnalyticsMetrics) error {
	headers := []string{"DATE", "REVIEWS", "AVG RATING", "SENTIMENT"}
	f.WithColumnStyle(2, output.StyleRating)
	rows := make([][]string, len(m.DataPoints))
	for i, dp := range m.DataPoints {
		rows[i] = []string{
			dp.Date,
			strconv.Itoa(dp.ReviewCount),
			fmt.Sprintf("%.2f", dp.AverageRating),
			fmt.Sprintf("%.2f", dp.Sentiment),
		}
	}
	return f.FormatList(headers, rows)
}

func init() {
	analyticsOverviewCmd.Flags().IntVar(&analyticsConnectionID, "connection-id", 0, "filter by connection ID")

	analyticsMetricsCmd.Flags().IntVar(&analyticsPeriod, "period", 7, "number of days (1-90)")
	analyticsMetricsCmd.Flags().IntVar(&analyticsConnectionID, "connection-id", 0, "filter by connection ID")

	analyticsCmd.AddCommand(analyticsOverviewCmd)
	analyticsCmd.AddCommand(analyticsMetricsCmd)
}
