package upload

import (
	"errors"
	"fmt"
	"mime"
	"net/http"
)

// Info about request headers
type Meta struct {
	MediaType string
	Boundary  string
	Range     *Range
	Filename  string
	UploadSid string
}

type Range struct {
	Start int64
	End   int64
	Size  int64
}

// Parse request headers and make Meta.
func ParseMeta(req *http.Request) (*Meta, error) {
	meta := &Meta{}

	if err := meta.parseContentType(req.Header.Get("Content-Type")); err != nil {
		return nil, err
	}

	if err := meta.parseContentRange(req.Header.Get("Content-Range")); err != nil {
		return nil, err
	}

	if err := meta.parseContentDisposition(req.Header.Get("Content-Disposition")); err != nil {
		return nil, err
	}

	cookie_pavo, err := req.Cookie("pavo")
	if err != nil {
		return nil, err
	}
	if cookie_pavo != nil {
		meta.UploadSid = cookie_pavo.Value
	}

	return meta, nil
}

func (meta *Meta) parseContentType(ct string) error {
	if ct == "" {
		meta.MediaType = "application/octet-stream"
		return nil
	}

	mediatype, params, err := mime.ParseMediaType(ct)
	if err != nil {
		return err
	}

	if mediatype == "multipart/form-data" {
		boundary, ok := params["boundary"]
		if !ok {
			return errors.New("meta: boundary not defined")
		}

		meta.MediaType = mediatype
		meta.Boundary = boundary
	} else {
		meta.MediaType = "application/octet-stream"
	}

	return nil
}

func (meta *Meta) parseContentRange(cr string) error {
	if cr == "" {
		return nil
	}

	var start, end, size int64

	_, err := fmt.Sscanf(cr, "bytes %d-%d/%d", &start, &end, &size)
	if err != nil {
		return err
	}

	meta.Range = &Range{Start: start, End: end, Size: size}

	return nil
}

func (meta *Meta) parseContentDisposition(cd string) error {
	if cd == "" {
		return nil
	}

	_, params, err := mime.ParseMediaType(cd)
	if err != nil {
		return err
	}

	filename, ok := params["filename"]
	if !ok {
		return errors.New("meta: filename in Content-Disposition not defined")
	}

	meta.Filename = filename

	return nil
}
