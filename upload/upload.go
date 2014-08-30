package upload

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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

// Original uploaded file.
type OriginalFile struct {
	BaseMime string
	Filepath string
	Filename string
}

// For input Request define algoritm
func SaveFiles(req *http.Request) ([]*OriginalFile, error) {
	meta, err := ParseMeta(req)
	if err != nil {
		return nil, err
	}

	body, err := NewBody(req.Header.Get("X-File"), req.Body)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	if meta.MediaType == "multipart/form-data" {
		files, err := SaveFilesFromMultipart(body.Body, meta.Boundary)
		if err != nil {
			return nil, err
		}

		return files, nil
	}

	file, err := SaveFileFromOctetStream(body.Body, meta.Filename)
	if err != nil {
		return nil, err
	}

	return []*OriginalFile{file}, nil
}

func SaveFilesFromMultipart(body io.Reader, boundary string) ([]*OriginalFile, error) {
	mr := multipart.NewReader(body, boundary)
	files := make([]*OriginalFile, 0)

	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if part.FormName() == "files[]" {
			original_file, err := saveTempFile(part)
			if err != nil {
				return nil, err
			}
			original_file.Filename = part.FileName()
			files = append(files, original_file)

			original_file.BaseMime, err = IdentifyMime(original_file.Filepath)
			if err != nil {
				return nil, err
			}
		}
	}

	return files, nil
}

func SaveFileFromOctetStream(body io.Reader, filename string) (*OriginalFile, error) {
	original_file, err := saveTempFile(body)
	if err != nil {
		return nil, err
	}

	if filename == "" {
		return nil, errors.New("upload: undefined filename")
	}
	original_file.Filename = filename

	original_file.BaseMime, err = IdentifyMime(original_file.Filepath)
	if err != nil {
		return nil, err
	}

	return original_file, nil
}

func saveTempFile(src io.Reader) (*OriginalFile, error) {
	temp_file, err := ioutil.TempFile(os.TempDir(), "pavo")
	if err != nil {
		return nil, err
	}
	defer temp_file.Close()

	_, err = io.Copy(temp_file, src)
	if err != nil {
		return nil, err
	}

	return &OriginalFile{Filepath: temp_file.Name()}, nil
}

// TODO: move function
// Get parameters for convert from Request query string
func GetConvertParams(req *http.Request) (map[string]string, error) {
	raw_converts := req.URL.Query().Get("converts")

	if raw_converts == "" {
		raw_converts = "{}"
	}

	convert := make(map[string]string)

	err := json.Unmarshal([]byte(raw_converts), &convert)
	if err != nil {
		return nil, err
	}

	return convert, nil
}

// Get base mime type
//		$ file --mime-type pic.jpg
//		pic.jpg: image/jpeg
func IdentifyMime(file string) (string, error) {
	out, err := exec.Command("file", "--mime-type", file).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("Identify Mime: file --mime-type %s; err: %s; detail: %s", file, err, string(out))
	}

	mime := strings.Split(strings.Split(string(out), ": ")[1], "/")[0]

	return mime, nil
}
