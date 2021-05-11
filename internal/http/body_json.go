package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

const jsonContentType = "application/json"

type jsonBodyProvider struct {
	data interface{}
}

var _ BodyProvider = (*jsonBodyProvider)(nil)

func (b *jsonBodyProvider) ContentType() string {
	return jsonContentType
}

func (b *jsonBodyProvider) Body() (io.Reader, error) {
	buf := &bytes.Buffer{}

	if err := json.NewEncoder(buf).Encode(b.data); err != nil {
		return nil, fmt.Errorf("failed to new json encoder: %w", err)
	}

	return buf, nil
}
