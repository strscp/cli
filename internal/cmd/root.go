package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/strscp/cli/internal/api"
	"github.com/strscp/cli/internal/config"
	"github.com/strscp/cli/internal/output"
	"github.com/strscp/cli/pkg/models"
	"golang.org/x/term"
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
	flagQuiet       bool
	flagRaw         bool
	flagWatch       string

	outputExplicit bool // true when user passed --output explicitly
)

var rootCmd = &cobra.Command{
	Use:   "starscope-cli",
	Short: "CLI for the Starscope review analytics API",
	Long:  "starscope-cli lets you interact with the Starscope API to manage reviews, insights, topics, connections, and analytics from the command line.",
	SilenceUsage:  true,
	SilenceErrors: true,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Auto-switch to JSON when stdout is not a TTY (piped/redirected),
		// unless user explicitly set --output.
		if !outputExplicit && !term.IsTerminal(int(os.Stdout.Fd())) {
			flagOutput = "json"
		}
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&flagAPIURL, "api-url", "", "API base URL (default: https://starscope.app/api/v1)")
	rootCmd.PersistentFlags().StringVarP(&flagToken, "token", "t", "", "API token (overrides stored token)")
	rootCmd.PersistentFlags().IntVarP(&flagWorkspaceID, "workspace-id", "w", 0, "workspace ID")
	rootCmd.PersistentFlags().StringVarP(&flagProfile, "profile", "p", "", "config profile name")
	rootCmd.PersistentFlags().StringVarP(&flagOutput, "output", "o", "table", "output format: table, json, csv, markdown")
	rootCmd.PersistentFlags().BoolVar(&flagNoColor, "no-color", false, "disable colored output")
	rootCmd.PersistentFlags().BoolVarP(&flagVerbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolVarP(&flagQuiet, "quiet", "q", false, "suppress non-data output")
	rootCmd.PersistentFlags().BoolVar(&flagRaw, "raw", false, "output raw API JSON response")
	rootCmd.PersistentFlags().StringVar(&flagWatch, "watch", "", "re-run command at interval (e.g. 10s, 1m, 5m)")

	// Track whether --output was explicitly set + handle --watch
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		outputExplicit = cmd.Flags().Changed("output")
		// Auto-switch to JSON when piped, unless user set --output
		if !outputExplicit && !term.IsTerminal(int(os.Stdout.Fd())) {
			flagOutput = "json"
		}

		// Handle --watch: wrap the leaf command's RunE with a watch loop
		if flagWatch != "" {
			interval, err := time.ParseDuration(flagWatch)
			if err != nil {
				return fmt.Errorf("invalid --watch interval %q: %w", flagWatch, err)
			}
			if interval < 5*time.Second {
				return fmt.Errorf("--watch interval must be at least 5s")
			}
			if hasAllFlag(cmd) {
				return fmt.Errorf("--watch cannot be used with --all")
			}
			if cmd.RunE != nil {
				originalRunE := cmd.RunE
				cmd.RunE = func(c *cobra.Command, a []string) error {
					return runWithWatch(c, a, interval, originalRunE)
				}
			}
		}

		return nil
	}

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(authCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(reviewsCmd)
	rootCmd.AddCommand(insightsCmd)
	rootCmd.AddCommand(topicsCmd)
	rootCmd.AddCommand(connectionsCmd)
	rootCmd.AddCommand(analyticsCmd)
	rootCmd.AddCommand(workspaceCmd)
	rootCmd.AddCommand(openCmd)
	rootCmd.AddCommand(completionCmd)
}

// Execute runs the root command and returns the appropriate exit code.
func Execute() error {
	err := rootCmd.Execute()
	if err != nil {
		handleError(err)
	}
	return err
}

// ExitCodeForError returns the exit code for an error.
func ExitCodeForError(err error) int {
	if err == nil {
		return 0
	}
	if cliErr, ok := err.(api.CLIError); ok {
		return cliErr.ExitCode()
	}
	return api.ExitConfigError
}

// handleError writes the error to stderr, using JSON when output is JSON.
func handleError(err error) {
	f, _ := resolveFormat()
	if f == output.FormatJSON {
		code := "error"
		exitCode := api.ExitConfigError
		if cliErr, ok := err.(api.CLIError); ok {
			code = cliErr.ErrorCode()
			exitCode = cliErr.ExitCode()
		}
		output.WriteJSONError(os.Stderr, code, err.Error(), exitCode)
	} else {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
	}
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

// paginationFooter prints pagination info and a next-page hint for table output.
func paginationFooter(resource string, meta models.PaginationMeta) {
	if flagQuiet || flagOutput != "table" {
		return
	}
	fmt.Printf("\nShowing %d-%d of %d %s (page %d/%d)\n",
		meta.From, meta.To, meta.Total, resource,
		meta.CurrentPage, meta.LastPage)
	if meta.CurrentPage < meta.LastPage {
		fmt.Printf("Use --page %d for the next page, or --all to fetch everything.\n", meta.CurrentPage+1)
	}
}

// listMeta converts pagination meta to output ListMeta.
func listMeta(meta models.PaginationMeta) *output.ListMeta {
	return &output.ListMeta{
		CurrentPage: meta.CurrentPage,
		LastPage:    meta.LastPage,
		PerPage:     meta.PerPage,
		Total:       meta.Total,
		From:        meta.From,
		To:          meta.To,
	}
}

// defaultToList makes the parent command run its "list" subcommand when no subcommand is given.
func defaultToList(parent *cobra.Command, listCmd *cobra.Command) {
	parent.RunE = func(cmd *cobra.Command, args []string) error {
		return listCmd.RunE(cmd, args)
	}
}

// newFormatter creates an output formatter from flags.
func newFormatter() (*output.Formatter, error) {
	format, err := resolveFormat()
	if err != nil {
		return nil, err
	}

	noColor := flagNoColor || os.Getenv("NO_COLOR") != ""
	return output.NewFormatter(format, noColor), nil
}

func resolveFormat() (output.Format, error) {
	return output.ParseFormat(flagOutput)
}
