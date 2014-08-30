package upload

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUploadChunked(t *testing.T) {
	assert := assert.New(t)
	storage := "../dummy/root_storage"
	fname := "../dummy/kino.jpg"
	f, _ := os.Open(fname)
	defer f.Close()

	cookie := &http.Cookie{Name: "pavo", Value: "abcdef"}

	req := createChunkRequest(f, 0, 24999)
	req.AddCookie(cookie)
	_, err := SaveChunk(req, storage)
	assert.Nil(err)

	req = createChunkRequest(f, 25000, 49999)
	req.AddCookie(cookie)
	_, err = SaveChunk(req, storage)
	assert.Nil(err)

	req = createChunkRequest(f, 50000, 52096)
	req.AddCookie(cookie)
	chunk, err := SaveChunk(req, storage)
	assert.Nil(err)

	assert.Equal(52097, chunk.Size)
}

func createChunkRequest(f *os.File, start int64, end int64) *http.Request {
	var body bytes.Buffer
	mw := multipart.NewWriter(&body)
	fi, _ := f.Stat()
	fw, _ := mw.CreateFormFile("files[]", fi.Name())

	io.CopyN(fw, f, end-start+1)
	mw.Close()

	req, _ := http.NewRequest("POST", "/files", &body)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	req.Header.Set("Content-Disposition", `attachment; filename="`+fi.Name()+`"`)
	req.Header.Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, fi.Size()))

	return req
}

func TestGetTempFileChunks(t *testing.T) {
	assert := assert.New(t)

	file, err := GetTempFileChunks("../dummy/root_storage", "abcdef", "kino.jpg")
	assert.Nil(err)
	assert.NotNil(file)
}
