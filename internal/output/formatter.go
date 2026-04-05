package output

import (
	"fmt"
	"io"
	"os"
)

// Format represents an output format type.
type Format string

const (
	FormatTable    Format = "table"
	FormatJSON     Format = "json"
	FormatCSV      Format = "csv"
	FormatMarkdown Format = "markdown"
)

// Field represents a key-value pair for detail output.
type Field struct {
	Label string
	Value string
}

// ListMeta holds pagination metadata for list output.
type ListMeta struct {
	CurrentPage int `json:"current_page"`
	LastPage    int `json:"last_page"`
	PerPage     int `json:"per_page"`
	Total       int `json:"total"`
	From        int `json:"from"`
	To          int `json:"to"`
}

// ColumnStyleFunc colorizes a cell value for table output.
// It receives the cell string and noColor flag, and returns the display string.
type ColumnStyleFunc func(value string, noColor bool) string

// Formatter writes structured data to an output writer.
type Formatter struct {
	format       Format
	writer       io.Writer
	noColor      bool
	columnStyles map[int]ColumnStyleFunc
}

// NewFormatter creates a new Formatter for the given format.
func NewFormatter(format Format, noColor bool) *Formatter {
	return &Formatter{
		format:  format,
		writer:  os.Stdout,
		noColor: noColor,
	}
}

// SetWriter sets the output writer (useful for testing).
func (f *Formatter) SetWriter(w io.Writer) {
	f.writer = w
}

// WithColumnStyle registers a color function for a column index (table format only).
func (f *Formatter) WithColumnStyle(colIndex int, styleFn ColumnStyleFunc) *Formatter {
	if f.columnStyles == nil {
		f.columnStyles = make(map[int]ColumnStyleFunc)
	}
	f.columnStyles[colIndex] = styleFn
	return f
}

// Format returns the current output format.
func (f *Formatter) Format() Format {
	return f.format
}

// FormatList writes a list of rows with headers.
func (f *Formatter) FormatList(headers []string, rows [][]string) error {
	switch f.format {
	case FormatJSON:
		return writeJSONList(f.writer, headers, rows, nil)
	case FormatCSV:
		return writeCSV(f.writer, headers, rows)
	case FormatMarkdown:
		return writeMarkdownList(f.writer, headers, rows)
	default:
		return writeTable(f.writer, headers, rows, f.noColor, f.columnStyles)
	}
}

// FormatListWithMeta writes a list with pagination metadata (JSON only).
func (f *Formatter) FormatListWithMeta(headers []string, rows [][]string, meta *ListMeta) error {
	switch f.format {
	case FormatJSON:
		return writeJSONList(f.writer, headers, rows, meta)
	case FormatCSV:
		return writeCSV(f.writer, headers, rows)
	case FormatMarkdown:
		return writeMarkdownList(f.writer, headers, rows)
	default:
		return writeTable(f.writer, headers, rows, f.noColor, f.columnStyles)
	}
}

// FormatDetail writes a single resource as key-value pairs.
func (f *Formatter) FormatDetail(fields []Field) error {
	switch f.format {
	case FormatJSON:
		return writeJSONDetail(f.writer, fields)
	case FormatCSV:
		headers := make([]string, len(fields))
		values := make([]string, len(fields))
		for i, field := range fields {
			headers[i] = field.Label
			values[i] = field.Value
		}
		return writeCSV(f.writer, headers, [][]string{values})
	case FormatMarkdown:
		return writeMarkdownDetail(f.writer, fields)
	default:
		return writeDetailTable(f.writer, fields, f.noColor)
	}
}

// FormatRaw writes raw JSON to the output.
func (f *Formatter) FormatRaw(data []byte) error {
	_, err := fmt.Fprintln(f.writer, string(data))
	return err
}

// ParseFormat parses a string into a Format, returning an error for invalid values.
func ParseFormat(s string) (Format, error) {
	switch s {
	case "table", "":
		return FormatTable, nil
	case "json":
		return FormatJSON, nil
	case "csv":
		return FormatCSV, nil
	case "markdown", "md":
		return FormatMarkdown, nil
	default:
		return "", fmt.Errorf("invalid output format %q (valid: table, json, csv, markdown)", s)
	}
}
