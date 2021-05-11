package http

import (
	"fmt"
	"io"
	"strings"

	"github.com/google/go-querystring/query"
)

const formContentType = "application/x-www-form-urlencoded"

type formBodyProvider struct {
	data interface{}
}

var _ BodyProvider = (*formBodyProvider)(nil)

func (b *formBodyProvider) ContentType() string {
	return formContentType
}

func (b *formBodyProvider) Body() (io.Reader, error) {
	urlValues, err := query.Values(b.data)
	if err != nil {
		return nil, fmt.Errorf("failed to query values: %w", err)
	}

	return strings.NewReader(urlValues.Encode()), nil
}
