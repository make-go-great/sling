package http

import (
	"fmt"
	"io"
	"net/http"
)

type RawResponseDecoder struct{}

var _ ResponseDecoder = (*RawResponseDecoder)(nil)

func (rd *RawResponseDecoder) Decode(rsp *http.Response, v interface{}) error {
	bytes, err := io.ReadAll(rsp.Body)
	if err != nil {
		return fmt.Errorf("failed to read all response body: %w", err)
	}

	vStr, ok := v.(*string)
	if !ok {
		return fmt.Errorf("v must be string pointer")
	}

	*vStr = string(bytes)

	return nil
}
