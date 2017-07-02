package models

import (
	"compress/gzip"
	"io"
)

type GZIPreadCloser struct {
	*gzip.Reader
	io.Closer
}

func (gz GZIPreadCloser) Close() error {
	return gz.Closer.Close()
}
