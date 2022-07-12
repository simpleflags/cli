package main

import (
	"bytes"
	"encoding/json"
)

func jsonFormatter(prefix, indent string, value any) string {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent(prefix, indent)
	err := encoder.Encode(value)
	if err != nil {
		return ""
	}
	return buffer.String()
}
