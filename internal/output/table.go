package output

import (
	"fmt"
	"io"
	"strings"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
	colorGray   = "\033[90m"
)

func writeTable(w io.Writer, headers []string, rows [][]string, noColor bool) error {
	if len(rows) == 0 {
		fmt.Fprintln(w, "No results found.")
		return nil
	}

	// Calculate column widths
	widths := make([]int, len(headers))
	for i, h := range headers {
		widths[i] = len(h)
	}
	for _, row := range rows {
		for i, cell := range row {
			if i < len(widths) && len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}

	// Cap column widths
	for i := range widths {
		if widths[i] > 60 {
			widths[i] = 60
		}
	}

	// Print header
	headerLine := formatRow(headers, widths)
	if !noColor {
		fmt.Fprintln(w, colorCyan+headerLine+colorReset)
	} else {
		fmt.Fprintln(w, headerLine)
	}

	// Print separator
	parts := make([]string, len(widths))
	for i, width := range widths {
		parts[i] = strings.Repeat("-", width)
	}
	if !noColor {
		fmt.Fprintln(w, colorGray+strings.Join(parts, "  ")+colorReset)
	} else {
		fmt.Fprintln(w, strings.Join(parts, "  "))
	}

	// Print rows
	for _, row := range rows {
		fmt.Fprintln(w, formatRow(row, widths))
	}

	return nil
}

func writeDetailTable(w io.Writer, fields []Field, noColor bool) error {
	maxLabel := 0
	for _, f := range fields {
		if len(f.Label) > maxLabel {
			maxLabel = len(f.Label)
		}
	}

	for _, f := range fields {
		label := fmt.Sprintf("%-*s", maxLabel, f.Label)
		if !noColor {
			fmt.Fprintf(w, "%s%s%s  %s\n", colorCyan, label, colorReset, f.Value)
		} else {
			fmt.Fprintf(w, "%s  %s\n", label, f.Value)
		}
	}

	return nil
}

func formatRow(cells []string, widths []int) string {
	parts := make([]string, len(widths))
	for i, width := range widths {
		cell := ""
		if i < len(cells) {
			cell = cells[i]
		}
		if len(cell) > width {
			cell = cell[:width-1] + "…"
		}
		parts[i] = fmt.Sprintf("%-*s", width, cell)
	}
	return strings.Join(parts, "  ")
}

// ColorRating returns a color-coded rating string.
func ColorRating(rating int, noColor bool) string {
	s := fmt.Sprintf("%d", rating)
	if noColor {
		return s
	}
	switch {
	case rating >= 4:
		return colorGreen + s + colorReset
	case rating == 3:
		return colorYellow + s + colorReset
	default:
		return colorRed + s + colorReset
	}
}
