package cmd

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"
)

const dashboardURL = "https://starscope.app"

var openTargets = map[string]string{
	"":            dashboardURL + "/dashboard",
	"dashboard":   dashboardURL + "/dashboard",
	"tokens":      dashboardURL + "/settings/api-tokens",
	"settings":    dashboardURL + "/settings",
	"connections": dashboardURL + "/connections",
}

var openCmd = &cobra.Command{
	Use:   "open [target]",
	Short: "Open Starscope in the browser",
	Long: `Open a Starscope page in your default browser.

Available targets: dashboard, tokens, settings, connections.
Defaults to dashboard if no target is specified.`,
	Example: `  starscope-cli open
  starscope-cli open tokens
  starscope-cli open connections`,
	Args:      cobra.MaximumNArgs(1),
	ValidArgs: []string{"dashboard", "tokens", "settings", "connections"},
	RunE: func(cmd *cobra.Command, args []string) error {
		target := ""
		if len(args) > 0 {
			target = args[0]
		}

		url, ok := openTargets[target]
		if !ok {
			return fmt.Errorf("unknown target %q (available: dashboard, tokens, settings, connections)", target)
		}

		fmt.Printf("Opening %s...\n", url)
		return openBrowser(url)
	},
}

func openBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "darwin":
		cmd = "open"
		args = []string{url}
	case "linux":
		cmd = "xdg-open"
		args = []string{url}
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start", url}
	default:
		return fmt.Errorf("unsupported platform, open manually: %s", url)
	}

	return exec.Command(cmd, args...).Start()
}

func init() {
	// registered in root.go
}
