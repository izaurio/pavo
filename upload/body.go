package upload

import (
	"bufio"
	"io"
	"mime/multipart"
	"os"
)

// Upload body info.
type Body struct {
	XFile     *os.File
	Body      io.Reader
	MR        *multipart.Reader
	Available bool
}

// Check exists body in xfile and return Body.
func NewBody(xfile string, req_body io.Reader) (*Body, error) {
	if xfile == "" {
		return &Body{Body: req_body, Available: true}, nil
	}

	fh, err := os.Open(xfile)
	if err != nil {
		return nil, err
	}

	return &Body{XFile: fh, Body: bufio.NewReader(fh), Available: true}, nil
}

// Close filehandler of body if XFile exists.
func (body *Body) Close() error {
	if body.XFile != nil {
		return body.XFile.Close()
	}

	return nil
}
