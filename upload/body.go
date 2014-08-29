package upload

import (
	"bufio"
	"io"
	"os"
)

// Upload body info.
type Body struct {
	xfile *os.File
	Body  io.Reader
}

// Check exists body in xfile and return Body.
func NewBody(xfile string, req_body io.Reader) (*Body, error) {
	if xfile == "" {
		return &Body{Body: req_body}, nil
	}

	fh, err := os.Open(xfile)
	if err != nil {
		return nil, err
	}

	return &Body{xfile: fh, Body: bufio.NewReader(fh)}, nil
}

// Close filehandler of body if XFile exists.
func (body *Body) Close() error {
	if body.xfile != nil {
		return body.xfile.Close()
	}

	return nil
}
