package http

import "io"

type Body interface {
	ContentType() string
	Output() (io.Reader, error)
}
