package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

const jsonContentType = "application/json"

type jsonBody struct {
	data interface{}
}

var _ Body = (*jsonBody)(nil)

func (b *jsonBody) ContentType() string {
	return jsonContentType
}

func (b *jsonBody) Output() (io.Reader, error) {
	buf := &bytes.Buffer{}

	if err := json.NewEncoder(buf).Encode(b.data); err != nil {
		return nil, fmt.Errorf("failed to new json encoder: %w", err)
	}

	return buf, nil
}
