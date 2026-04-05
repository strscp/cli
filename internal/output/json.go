package output

import (
	"encoding/json"
	"io"
)

func writeJSONList(w io.Writer, headers []string, rows [][]string, meta *ListMeta) error {
	items := make([]map[string]string, len(rows))
	for i, row := range rows {
		item := make(map[string]string)
		for j, header := range headers {
			if j < len(row) {
				item[header] = row[j]
			}
		}
		items[i] = item
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")

	if meta != nil {
		return enc.Encode(map[string]any{
			"data": items,
			"meta": meta,
		})
	}

	return enc.Encode(items)
}

func writeJSONDetail(w io.Writer, fields []Field) error {
	item := make(map[string]string)
	for _, f := range fields {
		item[f.Label] = f.Value
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(item)
}

// WriteJSONError writes a structured JSON error to the writer.
func WriteJSONError(w io.Writer, code string, message string, exitCode int) {
	enc := json.NewEncoder(w)
	_ = enc.Encode(map[string]any{
		"error":     code,
		"message":   message,
		"exit_code": exitCode,
	})
}
