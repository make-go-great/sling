package http

import "io"

type BodyProvider interface {
	ContentType() string
	Body() (io.Reader, error)
}
