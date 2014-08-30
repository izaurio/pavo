package upload

import (
	"bytes"
	"mime/multipart"
	"net/http"
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

	files, err := SaveFiles(req)
	assert.Nil(err)
	assert.Equal("kino.jpg", files[1].Filename)
	assert.Equal("image", files[1].BaseMime)

}

func TestUploadBinary(t *testing.T) {
	assert := assert.New(t)

	req, _ := http.NewRequest("POST", "/files", nil)
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("X-File", "../dummy/bin-data")
	req.Header.Set("Content-Disposition", `attachment; filename="basta.png"`)

	files, err := SaveFiles(req)
	assert.Nil(err)
	assert.Equal("basta.png", files[0].Filename)
	assert.Equal("image", files[0].BaseMime)

}
