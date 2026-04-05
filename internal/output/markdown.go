package output

import (
	"fmt"
	"io"
	"strings"
)

func writeMarkdownList(w io.Writer, headers []string, rows [][]string) error {
	if len(rows) == 0 {
		fmt.Fprintln(w, "*No results found.*")
		return nil
	}

	// Header row
	fmt.Fprintf(w, "| %s |\n", strings.Join(headers, " | "))

	// Alignment row
	seps := make([]string, len(headers))
	for i := range headers {
		seps[i] = "---"
	}
	fmt.Fprintf(w, "| %s |\n", strings.Join(seps, " | "))

	// Data rows
	for _, row := range rows {
		cells := make([]string, len(headers))
		for i := range headers {
			if i < len(row) {
				cells[i] = strings.ReplaceAll(row[i], "|", "\\|")
			}
		}
		fmt.Fprintf(w, "| %s |\n", strings.Join(cells, " | "))
	}

	return nil
}

func writeMarkdownDetail(w io.Writer, fields []Field) error {
	for _, f := range fields {
		fmt.Fprintf(w, "**%s:** %s\n", f.Label, f.Value)
	}
	return nil
}
