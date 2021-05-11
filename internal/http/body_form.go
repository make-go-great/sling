package http

import (
	"fmt"
	"io"
	"strings"

	"github.com/google/go-querystring/query"
)

const formContentType = "application/x-www-form-urlencoded"

type formBody struct {
	data interface{}
}

var _ Body = (*formBody)(nil)

func (b *formBody) ContentType() string {
	return formContentType
}

func (b *formBody) Output() (io.Reader, error) {
	urlValues, err := query.Values(b.data)
	if err != nil {
		return nil, fmt.Errorf("failed to query values: %w", err)
	}

	return strings.NewReader(urlValues.Encode()), nil
}
