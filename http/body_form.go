package http

import (
	"fmt"
	"io"
	"strings"

	"github.com/google/go-querystring/query"
)

const formContentType = "application/x-www-form-urlencoded"

type FormBodyProvider struct {
	Data interface{}
}

var _ BodyProvider = (*FormBodyProvider)(nil)

func (bp *FormBodyProvider) ContentType() string {
	return formContentType
}

func (bp *FormBodyProvider) Body() (io.Reader, error) {
	urlValues, err := query.Values(bp.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to query values: %w", err)
	}

	return strings.NewReader(urlValues.Encode()), nil
}
