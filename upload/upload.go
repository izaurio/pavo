package upload

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

// Original uploaded file.
type OriginalFile struct {
	BaseMime string
	Filepath string
	Filename string
}

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
		switch part.FormName() {
		case "files[]", "files", "file":
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
