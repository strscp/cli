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

func writeTable(w io.Writer, headers []string, rows [][]string, noColor bool, styles map[int]ColumnStyleFunc) error {
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
		if !noColor && len(styles) > 0 {
			fmt.Fprintln(w, formatRowWithStyles(row, widths, styles))
		} else {
			fmt.Fprintln(w, formatRow(row, widths))
		}
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

func formatRowWithStyles(cells []string, widths []int, styles map[int]ColumnStyleFunc) string {
	parts := make([]string, len(widths))
	for i, width := range widths {
		cell := ""
		if i < len(cells) {
			cell = cells[i]
		}
		if len(cell) > width {
			cell = cell[:width-1] + "…"
		}
		if styleFn, ok := styles[i]; ok {
			colored := styleFn(cell, false)
			padding := width - len(cell)
			if padding < 0 {
				padding = 0
			}
			parts[i] = colored + strings.Repeat(" ", padding)
		} else {
			parts[i] = fmt.Sprintf("%-*s", width, cell)
		}
	}
	return strings.Join(parts, "  ")
}

// StyleRating colors a rating value (1-5 scale, supports int and float strings).
func StyleRating(value string, noColor bool) string {
	if noColor || value == "-" {
		return value
	}
	var rating float64
	if _, err := fmt.Sscanf(value, "%f", &rating); err != nil {
		return value
	}
	switch {
	case rating >= 4.0:
		return colorGreen + value + colorReset
	case rating >= 3.0:
		return colorYellow + value + colorReset
	default:
		return colorRed + value + colorReset
	}
}

// StyleSeverity colors a severity string (critical/warning/info).
func StyleSeverity(value string, noColor bool) string {
	if noColor {
		return value
	}
	switch value {
	case "critical":
		return colorRed + value + colorReset
	case "warning":
		return colorYellow + value + colorReset
	default:
		return colorGray + value + colorReset
	}
}

// StyleBool colors a boolean string (true=green, false=red).
func StyleBool(value string, noColor bool) string {
	if noColor {
		return value
	}
	if value == "true" {
		return colorGreen + value + colorReset
	}
	return colorRed + value + colorReset
}
