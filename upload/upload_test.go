package upload

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUploadMultipart(t *testing.T) {
	assert := assert.New(t)

	var body bytes.Buffer
	mw := multipart.NewWriter(&body)

	if err := writeMPBody("../dummy/32509211_news_bigpic.jpg", mw); err != nil {
		assert.Error(err)
	}
	if err := writeMPBody("../dummy/kino.jpg", mw); err != nil {
		assert.Error(err)
	}

	mw.Close()

	req, _ := http.NewRequest("POST", "/files", &body)
	req.Header.Set("Content-Type", mw.FormDataContentType())

	t.Logf("%#v", req.Header)

}

func TestUploadBinary(t *testing.T) {
	//assert := assert.New(t)

	req, _ := http.NewRequest("POST", "/files", nil)
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("X-File", "../dummy/bin-data")
	req.Header.Set("Content-Disposition", `attachment; filename="basta.png"`)

	t.Logf("%#v", req.Header)

}

func writeMPBody(fname string, mw *multipart.Writer) error {
	fw, _ := mw.CreateFormFile("files[]", filepath.Base(fname))
	f, err := os.Open(fname)
	if err != nil {
		return err
	}

	_, err = io.Copy(fw, f)
	if err != nil {
		return err
	}

	return nil
}
