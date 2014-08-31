package upload

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTempFileChunks(t *testing.T) {
	assert := assert.New(t)

	file, err := TempFileChunks(0, "../dummy/root_storage", "abcdef", "kino.jpg")
	assert.Nil(err)
	assert.NotNil(file)
}

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
