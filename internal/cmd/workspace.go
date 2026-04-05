package cmd

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/strscp/cli/internal/output"
	"github.com/strscp/cli/pkg/models"
)

var workspaceCmd = &cobra.Command{
	Use:   "workspace",
	Short: "Show workspace information",
	RunE:  workspaceShowRun,
}

var workspaceShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show workspace information",
	RunE:  workspaceShowRun,
}

func workspaceShowRun(cmd *cobra.Command, args []string) error {
	client, err := newAPIClient()
	if err != nil {
		return err
	}
	f, err := newFormatter()
	if err != nil {
		return err
	}

	ws, err := client.Workspace.Show(context.Background())
	if err != nil {
		return err
	}

	return formatWorkspace(f, ws)
}

func formatWorkspace(f *output.Formatter, ws *models.Workspace) error {
	trialEnds := "-"
	if ws.TrialEndsAt != nil {
		trialEnds = *ws.TrialEndsAt
	}

	fields := []output.Field{
		{Label: "ID", Value: strconv.Itoa(ws.ID)},
		{Label: "Name", Value: ws.Name},
		{Label: "Slug", Value: ws.Slug},
		{Label: "Plan", Value: ws.PlanTier},
		{Label: "Connections", Value: fmt.Sprintf("%d / %d", ws.ConnectionCount, ws.ConnectionLimit)},
		{Label: "Trial Ends", Value: trialEnds},
	}
	return f.FormatDetail(fields)
}

func init() {
	workspaceCmd.AddCommand(workspaceShowCmd)
}
