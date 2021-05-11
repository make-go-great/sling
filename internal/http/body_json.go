package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

const jsonContentType = "application/json"

type JSONBodyProvider struct {
	Data interface{}
}

var _ BodyProvider = (*JSONBodyProvider)(nil)

func (b *JSONBodyProvider) ContentType() string {
	return jsonContentType
}

func (b *JSONBodyProvider) Body() (io.Reader, error) {
	buf := &bytes.Buffer{}

	if err := json.NewEncoder(buf).Encode(b.Data); err != nil {
		return nil, fmt.Errorf("failed to new json encoder: %w", err)
	}

	return buf, nil
}
