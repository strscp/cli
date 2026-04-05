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
