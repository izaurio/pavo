package upload

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

func writeMPBody(fname string, mw *multipart.Writer) error {
	fw, _ := mw.CreateFormFile("files[]", filepath.Base(fname))
	f, err := os.Open(fname)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(fw, f)
	if err != nil {
		return err
	}

	return nil
}
