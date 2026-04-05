package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/strscp/cli/internal/api"
	"github.com/strscp/cli/internal/config"
	"github.com/strscp/cli/internal/output"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var (
	flagAPIURL      string
	flagToken       string
	flagWorkspaceID int
	flagProfile     string
	flagOutput      string
	flagNoColor     bool
	flagVerbose     bool
)

var rootCmd = &cobra.Command{
	Use:   "starscope-cli",
	Short: "CLI for the Starscope review analytics API",
	Long:  "starscope-cli lets you interact with the Starscope API to manage reviews, insights, topics, connections, and analytics from the command line.",
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	rootCmd.PersistentFlags().StringVar(&flagAPIURL, "api-url", "", "API base URL (default: https://starscope.app/api/v1)")
	rootCmd.PersistentFlags().StringVarP(&flagToken, "token", "t", "", "API token (overrides stored token)")
	rootCmd.PersistentFlags().IntVarP(&flagWorkspaceID, "workspace-id", "w", 0, "workspace ID")
	rootCmd.PersistentFlags().StringVarP(&flagProfile, "profile", "p", "", "config profile name")
	rootCmd.PersistentFlags().StringVarP(&flagOutput, "output", "o", "table", "output format: table, json, csv")
	rootCmd.PersistentFlags().BoolVar(&flagNoColor, "no-color", false, "disable colored output")
	rootCmd.PersistentFlags().BoolVarP(&flagVerbose, "verbose", "v", false, "verbose output")

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(authCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(reviewsCmd)
	rootCmd.AddCommand(insightsCmd)
	rootCmd.AddCommand(topicsCmd)
	rootCmd.AddCommand(connectionsCmd)
	rootCmd.AddCommand(analyticsCmd)
	rootCmd.AddCommand(workspaceCmd)
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}

// loadConfig loads the application config.
func loadConfig() (*config.Config, error) {
	cfg, err := config.New()
	if err != nil {
		return nil, fmt.Errorf("loading config: %w", err)
	}

	if flagProfile != "" {
		cfg.SetCurrentProfile(flagProfile)
	}

	return cfg, nil
}

// newAPIClient creates an authenticated API client from config.
func newAPIClient() (*api.Client, error) {
	cfg, err := loadConfig()
	if err != nil {
		return nil, err
	}

	token := cfg.ResolveToken(flagToken)
	if token == "" {
		return nil, fmt.Errorf("no API token configured. Run 'starscope-cli auth login' to authenticate")
	}

	apiURL := cfg.ResolveAPIURL(flagAPIURL)
	workspaceID := cfg.ResolveWorkspaceID(flagWorkspaceID)

	opts := []api.ClientOption{
		api.WithUserAgent(fmt.Sprintf("starscope-cli/%s", version)),
	}
	if workspaceID > 0 {
		opts = append(opts, api.WithWorkspaceID(workspaceID))
	}

	return api.NewClient(apiURL, token, opts...), nil
}

// newFormatter creates an output formatter from flags.
func newFormatter() (*output.Formatter, error) {
	format, err := output.ParseFormat(flagOutput)
	if err != nil {
		return nil, err
	}

	noColor := flagNoColor || os.Getenv("NO_COLOR") != ""
	return output.NewFormatter(format, noColor), nil
}
