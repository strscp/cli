package main

import (
	"os"

	"github.com/strscp/cli/internal/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(cmd.ExitCodeForError(err))
	}
}
