package upload

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

type ChunkFile struct {
	Filepath string
	Size     int64
}

func SaveChunk(req *http.Request, storage string) (*ChunkFile, error) {
	meta, err := ParseMeta(req)
	if err != nil {
		return nil, err
	}

	body, err := NewBody(req.Header.Get("X-File"), req.Body)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	cookie_pavo, err := req.Cookie("pavo")
	if cookie_pavo == nil || err != nil {
		return nil, fmt.Errorf("Error fetch cookie pavo; err: %s", err)
	}

	temp_file, err := GetTempFileChunks(storage, cookie_pavo.Value, meta.Filename)
	if err != nil {
		return nil, err
	}
	defer temp_file.Close()

	if _, err = temp_file.Seek(meta.Range.Start, 0); err != nil {
		return nil, err
	}

	chunk_size := meta.Range.End - meta.Range.Start + 1
	if meta.MediaType == "multipart/form-data" {
		if err = SaveChunkFromMultipart(temp_file, chunk_size, body.Body, meta.Boundary); err != nil {
			return nil, err
		}
	} else {
		if err = SaveChunkFromOctetStream(temp_file, chunk_size, body.Body); err != nil {
			return nil, err
		}
	}

	fi, err := temp_file.Stat()
	if err != nil {
		return nil, err
	}

	return &ChunkFile{Filepath: temp_file.Name(), Size: fi.Size()}, nil
}

func SaveChunkFromMultipart(temp_file *os.File, size int64, body io.Reader, boundary string) error {
	mr := multipart.NewReader(body, boundary)

	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if part.FormName() == "files[]" {
			_, err = io.CopyN(temp_file, part, size)
			if err != nil {
				return err
			}

			return nil
		}
	}

	return errors.New("chunk: file not present in body")

}

func SaveChunkFromOctetStream(temp_file *os.File, size int64, body io.Reader) error {
	_, err := io.CopyN(temp_file, body, size)
	if err != nil {
		return err
	}

	return nil
}

func GetTempFileChunks(storage, upload_sid, user_filename string) (*os.File, error) {
	hasher := md5.New()
	hasher.Write([]byte(upload_sid + user_filename))
	filename := hex.EncodeToString(hasher.Sum(nil))

	path := filepath.Join(storage, "chunks")

	err := os.MkdirAll(path, 0755)
	if err != nil {
		return nil, err
	}

	file, err := os.OpenFile(filepath.Join(path, filename), os.O_CREATE|os.O_WRONLY, 0664)
	if err != nil {
		return nil, err
	}

	return file, nil
}
