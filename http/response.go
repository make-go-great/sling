package http

import (
	"net/http"
)

type ResponseDecoder interface {
	Decode(rsp *http.Response, v interface{}) error
}
