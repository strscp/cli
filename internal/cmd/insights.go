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

var insightsCmd = &cobra.Command{
	Use:   "insights",
	Short: "Manage AI insights",
}

var (
	insightType         string
	insightSeverity     string
	insightConnectionID int
	insightPerPage      int
	insightPage         int
	insightAll          bool
)

var insightsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List insights",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}
		f, err := newFormatter()
		if err != nil {
			return err
		}

		params := api.ListInsightsParams{
			Type:         insightType,
			Severity:     insightSeverity,
			ConnectionID: insightConnectionID,
			PerPage:      insightPerPage,
			Page:         insightPage,
		}

		ctx := context.Background()

		if insightAll {
			insights, err := client.Insights.ListAll(ctx, params)
			if err != nil {
				return err
			}
			return formatInsightList(f, insights)
		}

		resp, err := client.Insights.List(ctx, params)
		if err != nil {
			return err
		}

		if err := formatInsightList(f, resp.Data); err != nil {
			return err
		}

		if flagOutput == "table" {
			fmt.Printf("\nShowing %d-%d of %d insights (page %d/%d)\n",
				resp.Meta.From, resp.Meta.To, resp.Meta.Total,
				resp.Meta.CurrentPage, resp.Meta.LastPage)
		}

		return nil
	},
}

var insightsShowCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Show an insight",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid insight ID: %s", args[0])
		}

		client, err := newAPIClient()
		if err != nil {
			return err
		}
		f, err := newFormatter()
		if err != nil {
			return err
		}

		insight, err := client.Insights.Get(context.Background(), id)
		if err != nil {
			return err
		}

		return formatInsightDetail(f, insight)
	},
}

func formatInsightList(f *output.Formatter, insights []models.Insight) error {
	headers := []string{"ID", "TYPE", "SEVERITY", "TITLE", "GENERATED"}
	rows := make([][]string, len(insights))
	for i, ins := range insights {
		rows[i] = []string{
			strconv.Itoa(ins.ID),
			ins.Type,
			ins.Severity,
			ins.Title,
			ins.GeneratedAt,
		}
	}
	return f.FormatList(headers, rows)
}

func formatInsightDetail(f *output.Formatter, ins *models.Insight) error {
	fields := []output.Field{
		{Label: "ID", Value: strconv.Itoa(ins.ID)},
		{Label: "Type", Value: ins.Type},
		{Label: "Severity", Value: ins.Severity},
		{Label: "Title", Value: ins.Title},
		{Label: "Body", Value: ins.Body},
		{Label: "Confidence", Value: fmt.Sprintf("%.0f%%", ins.Confidence*100)},
		{Label: "Actions", Value: strings.Join(ins.SuggestedActions, "; ")},
		{Label: "Period", Value: ins.PeriodStart + " to " + ins.PeriodEnd},
		{Label: "Generated", Value: ins.GeneratedAt},
		{Label: "Read", Value: fmt.Sprintf("%t", ins.IsRead)},
	}
	return f.FormatDetail(fields)
}

func init() {
	insightsListCmd.Flags().StringVar(&insightType, "type", "", "filter by type (trend, anomaly, opportunity, summary)")
	insightsListCmd.Flags().StringVar(&insightSeverity, "severity", "", "filter by severity (info, warning, critical)")
	insightsListCmd.Flags().IntVar(&insightConnectionID, "connection-id", 0, "filter by connection ID")
	insightsListCmd.Flags().IntVar(&insightPerPage, "per-page", 15, "results per page (1-100)")
	insightsListCmd.Flags().IntVar(&insightPage, "page", 1, "page number")
	insightsListCmd.Flags().BoolVar(&insightAll, "all", false, "fetch all pages")

	insightsCmd.AddCommand(insightsListCmd)
	insightsCmd.AddCommand(insightsShowCmd)
}
