package upload

import (
	"errors"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
)

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
