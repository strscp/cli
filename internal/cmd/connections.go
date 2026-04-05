package cmd

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/strscp/cli/internal/output"
	"github.com/strscp/cli/pkg/models"
)

var connectionsCmd = &cobra.Command{
	Use:     "connections",
	Aliases: []string{"c", "conn"},
	Short:   "Manage review platform connections",
}

var connectionsListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List connections",
	Example: `  starscope-cli connections list
  starscope-cli c ls --output json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient()
		if err != nil {
			return err
		}
		f, err := newFormatter()
		if err != nil {
			return err
		}

		resp, err := client.Connections.List(context.Background())
		if err != nil {
			return err
		}

		return formatConnectionList(f, resp.Data)
	},
}

var connectionsShowCmd = &cobra.Command{
	Use:     "show <id>",
	Aliases: []string{"s"},
	Short:   "Show a connection",
	Example: `  starscope-cli connections show 5`,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid connection ID: %s", args[0])
		}

		client, err := newAPIClient()
		if err != nil {
			return err
		}
		f, err := newFormatter()
		if err != nil {
			return err
		}

		conn, err := client.Connections.Get(context.Background(), id)
		if err != nil {
			return err
		}

		return formatConnectionDetail(f, conn)
	},
}

func formatConnectionList(f *output.Formatter, connections []models.Connection) error {
	headers := []string{"ID", "PLATFORM", "NAME", "ACTIVE", "REVIEWS", "RATING", "LAST SYNCED"}
	rows := make([][]string, len(connections))
	for i, c := range connections {
		rating := "-"
		if c.Rating != nil {
			rating = fmt.Sprintf("%.1f", *c.Rating)
		}
		lastSynced := "-"
		if c.LastSyncedAt != nil {
			lastSynced = *c.LastSyncedAt
		}
		rows[i] = []string{
			strconv.Itoa(c.ID),
			c.Platform,
			c.Name,
			fmt.Sprintf("%t", c.IsActive),
			strconv.Itoa(c.ReviewCount),
			rating,
			lastSynced,
		}
	}
	return f.FormatList(headers, rows)
}

func formatConnectionDetail(f *output.Formatter, c *models.Connection) error {
	rating := "-"
	if c.Rating != nil {
		rating = fmt.Sprintf("%.1f", *c.Rating)
	}

	fields := []output.Field{
		{Label: "ID", Value: strconv.Itoa(c.ID)},
		{Label: "Platform", Value: c.Platform},
		{Label: "Name", Value: c.Name},
		{Label: "External ID", Value: c.ExternalID},
		{Label: "External URL", Value: deref(c.ExternalURL)},
		{Label: "Active", Value: fmt.Sprintf("%t", c.IsActive)},
		{Label: "Review Count", Value: strconv.Itoa(c.ReviewCount)},
		{Label: "Rating", Value: rating},
		{Label: "Last Synced", Value: deref(c.LastSyncedAt)},
	}
	return f.FormatDetail(fields)
}

func init() {
	connectionsCmd.AddCommand(connectionsListCmd)
	connectionsCmd.AddCommand(connectionsShowCmd)
	defaultToList(connectionsCmd, connectionsListCmd)
}
