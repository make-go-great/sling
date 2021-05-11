package http

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type JSONResponseDecoder struct{}

var _ ResponseDecoder = (*JSONResponseDecoder)(nil)

func (rd *JSONResponseDecoder) Decode(rsp *http.Response, v interface{}) error {
	if err := json.NewDecoder(rsp.Body).Decode(v); err != nil {
		return fmt.Errorf("failed to decode json: %w", err)
	}

	return nil
}
