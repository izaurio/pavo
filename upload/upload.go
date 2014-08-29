package upload

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
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

	// define saver
	if meta.MediaType == "multipart/form-data" {
		// fetch list OrirginalFile
		files, err := SaveFilesFromMultipart(body.Body, meta.Boundary)
		if err != nil {
			return nil, err
		}

		return files, nil
	}

	// load single file from binary body
	return nil, errors.New("Implement load single file")
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

			// identify base mime
		}
	}

	return files, nil
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
