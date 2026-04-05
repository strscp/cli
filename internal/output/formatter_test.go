package output

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFormatList_Table(t *testing.T) {
	f := NewFormatter(FormatTable, true) // no color
	buf := &bytes.Buffer{}
	f.SetWriter(buf)

	headers := []string{"ID", "NAME", "STATUS"}
	rows := [][]string{
		{"1", "Alpha", "active"},
		{"2", "Beta", "inactive"},
	}

	err := f.FormatList(headers, rows)
	require.NoError(t, err)

	out := buf.String()
	assert.Contains(t, out, "ID")
	assert.Contains(t, out, "Alpha")
	assert.Contains(t, out, "Beta")
}

func TestFormatList_Table_Empty(t *testing.T) {
	f := NewFormatter(FormatTable, true)
	buf := &bytes.Buffer{}
	f.SetWriter(buf)

	err := f.FormatList([]string{"ID"}, [][]string{})
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "No results found.")
}

func TestFormatList_JSON(t *testing.T) {
	f := NewFormatter(FormatJSON, false)
	buf := &bytes.Buffer{}
	f.SetWriter(buf)

	headers := []string{"ID", "NAME"}
	rows := [][]string{
		{"1", "Alpha"},
	}

	err := f.FormatList(headers, rows)
	require.NoError(t, err)

	out := buf.String()
	assert.Contains(t, out, `"ID"`)
	assert.Contains(t, out, `"Alpha"`)
}

func TestFormatListWithMeta_JSON(t *testing.T) {
	f := NewFormatter(FormatJSON, false)
	buf := &bytes.Buffer{}
	f.SetWriter(buf)

	headers := []string{"ID", "NAME"}
	rows := [][]string{{"1", "Alpha"}}
	meta := &ListMeta{CurrentPage: 1, LastPage: 3, PerPage: 15, Total: 42, From: 1, To: 1}

	err := f.FormatListWithMeta(headers, rows, meta)
	require.NoError(t, err)

	out := buf.String()
	assert.Contains(t, out, `"data"`)
	assert.Contains(t, out, `"meta"`)
	assert.Contains(t, out, `"total": 42`)
	assert.Contains(t, out, `"last_page": 3`)
}

func TestFormatList_CSV(t *testing.T) {
	f := NewFormatter(FormatCSV, false)
	buf := &bytes.Buffer{}
	f.SetWriter(buf)

	headers := []string{"ID", "NAME"}
	rows := [][]string{
		{"1", "Alpha"},
		{"2", "Beta"},
	}

	err := f.FormatList(headers, rows)
	require.NoError(t, err)

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	assert.Len(t, lines, 3) // header + 2 rows
	assert.Equal(t, "ID,NAME", lines[0])
	assert.Equal(t, "1,Alpha", lines[1])
}

func TestFormatDetail_Table(t *testing.T) {
	f := NewFormatter(FormatTable, true)
	buf := &bytes.Buffer{}
	f.SetWriter(buf)

	fields := []Field{
		{Label: "Name", Value: "Test"},
		{Label: "Status", Value: "active"},
	}

	err := f.FormatDetail(fields)
	require.NoError(t, err)

	out := buf.String()
	assert.Contains(t, out, "Name")
	assert.Contains(t, out, "Test")
	assert.Contains(t, out, "Status")
	assert.Contains(t, out, "active")
}

func TestFormatDetail_JSON(t *testing.T) {
	f := NewFormatter(FormatJSON, false)
	buf := &bytes.Buffer{}
	f.SetWriter(buf)

	fields := []Field{
		{Label: "Name", Value: "Test"},
	}

	err := f.FormatDetail(fields)
	require.NoError(t, err)
	assert.Contains(t, buf.String(), `"Name": "Test"`)
}

func TestParseFormat(t *testing.T) {
	tests := []struct {
		input    string
		expected Format
		wantErr  bool
	}{
		{"table", FormatTable, false},
		{"json", FormatJSON, false},
		{"csv", FormatCSV, false},
		{"markdown", FormatMarkdown, false},
		{"md", FormatMarkdown, false},
		{"", FormatTable, false},
		{"xml", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ParseFormat(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, got)
			}
		})
	}
}

func TestFormatList_Markdown(t *testing.T) {
	f := NewFormatter(FormatMarkdown, false)
	buf := &bytes.Buffer{}
	f.SetWriter(buf)

	headers := []string{"ID", "NAME"}
	rows := [][]string{
		{"1", "Alpha"},
		{"2", "Beta"},
	}

	err := f.FormatList(headers, rows)
	require.NoError(t, err)

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	assert.Len(t, lines, 4) // header + separator + 2 rows
	assert.Equal(t, "| ID | NAME |", lines[0])
	assert.Equal(t, "| --- | --- |", lines[1])
	assert.Equal(t, "| 1 | Alpha |", lines[2])
	assert.Equal(t, "| 2 | Beta |", lines[3])
}

func TestFormatList_Markdown_Empty(t *testing.T) {
	f := NewFormatter(FormatMarkdown, false)
	buf := &bytes.Buffer{}
	f.SetWriter(buf)

	err := f.FormatList([]string{"ID"}, [][]string{})
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "*No results found.*")
}

func TestFormatList_Markdown_EscapePipes(t *testing.T) {
	f := NewFormatter(FormatMarkdown, false)
	buf := &bytes.Buffer{}
	f.SetWriter(buf)

	headers := []string{"NAME"}
	rows := [][]string{{"foo|bar"}}

	err := f.FormatList(headers, rows)
	require.NoError(t, err)
	assert.Contains(t, buf.String(), `foo\|bar`)
}

func TestFormatDetail_Markdown(t *testing.T) {
	f := NewFormatter(FormatMarkdown, false)
	buf := &bytes.Buffer{}
	f.SetWriter(buf)

	fields := []Field{
		{Label: "Name", Value: "Test"},
		{Label: "Status", Value: "active"},
	}

	err := f.FormatDetail(fields)
	require.NoError(t, err)

	out := buf.String()
	assert.Contains(t, out, "**Name:** Test")
	assert.Contains(t, out, "**Status:** active")
}

func TestFormatList_Table_WithColumnStyles(t *testing.T) {
	f := NewFormatter(FormatTable, false) // colors enabled
	buf := &bytes.Buffer{}
	f.SetWriter(buf)
	f.WithColumnStyle(1, StyleRating)

	headers := []string{"NAME", "RATING"}
	rows := [][]string{
		{"Good", "5"},
		{"Bad", "1"},
	}

	err := f.FormatList(headers, rows)
	require.NoError(t, err)

	out := buf.String()
	assert.Contains(t, out, colorGreen+"5"+colorReset)
	assert.Contains(t, out, colorRed+"1"+colorReset)
}

func TestFormatList_Table_WithColumnStyles_NoColor(t *testing.T) {
	f := NewFormatter(FormatTable, true) // no color
	buf := &bytes.Buffer{}
	f.SetWriter(buf)
	f.WithColumnStyle(1, StyleRating)

	headers := []string{"NAME", "RATING"}
	rows := [][]string{{"Good", "5"}}

	err := f.FormatList(headers, rows)
	require.NoError(t, err)

	out := buf.String()
	assert.NotContains(t, out, colorGreen)
	assert.NotContains(t, out, colorReset)
}

func TestStyleRating(t *testing.T) {
	assert.Contains(t, StyleRating("5", false), colorGreen)
	assert.Contains(t, StyleRating("4", false), colorGreen)
	assert.Contains(t, StyleRating("4.2", false), colorGreen)
	assert.Contains(t, StyleRating("3", false), colorYellow)
	assert.Contains(t, StyleRating("3.5", false), colorYellow)
	assert.Contains(t, StyleRating("2", false), colorRed)
	assert.Contains(t, StyleRating("1", false), colorRed)
	assert.Equal(t, "-", StyleRating("-", false))
	assert.Equal(t, "5", StyleRating("5", true))
}

func TestStyleSeverity(t *testing.T) {
	assert.Contains(t, StyleSeverity("critical", false), colorRed)
	assert.Contains(t, StyleSeverity("warning", false), colorYellow)
	assert.Contains(t, StyleSeverity("info", false), colorGray)
	assert.Equal(t, "critical", StyleSeverity("critical", true))
}

func TestStyleBool(t *testing.T) {
	assert.Contains(t, StyleBool("true", false), colorGreen)
	assert.Contains(t, StyleBool("false", false), colorRed)
	assert.Equal(t, "true", StyleBool("true", true))
}
