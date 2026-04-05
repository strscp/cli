package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/spf13/cobra"
)

// runWithWatch re-executes a command's RunE at a fixed interval, clearing the screen each time.
func runWithWatch(cmd *cobra.Command, args []string, interval time.Duration, runE func(*cobra.Command, []string) error) error {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	defer signal.Stop(sigCh)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		// Clear screen
		fmt.Fprint(os.Stdout, "\033[2J\033[H")

		if err := runE(cmd, args); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		}

		if !flagQuiet {
			fmt.Fprintf(os.Stderr, "\nRefreshing every %s... (Ctrl+C to stop)\n", interval)
		}

		select {
		case <-sigCh:
			return nil
		case <-ticker.C:
			// continue loop
		}
	}
}

// hasAllFlag checks if the --all flag was explicitly set on the command.
func hasAllFlag(cmd *cobra.Command) bool {
	f := cmd.Flags().Lookup("all")
	if f == nil {
		return false
	}
	return f.Changed
}
