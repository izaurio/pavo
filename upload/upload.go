package upload

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Error Incomplete returned by uploader when loaded non-last chunk.
var Incomplete = errors.New("Incomplete")

// Original uploaded file.
type OriginalFile struct {
	BaseMime string
	Filepath string
	Filename string
	Size     int64
}

type Uploader struct {
	Root string
	Meta *Meta
	Body *Body
}

func (up *Uploader) Reader() (io.Reader, string, error) {
	if up.Meta.MediaType == "multipart/form-data" {
		up.Body.MR = multipart.NewReader(up.Body.Body, up.Meta.Boundary)
		for {
			part, err := up.Body.MR.NextPart()
			if err != nil {
				return nil, "", err
			}
			if part.FormName() == "files[]" {
				return part, part.FileName(), nil
			}
		}
	}

	up.Body.Available = false

	return up.Body.Body, up.Meta.Filename, nil
}

func (up *Uploader) TempFile() (*os.File, error) {
	if up.Meta.Range == nil {
		return TempFile()
	}
	return TempFileChunks(up.Meta.Range.Start, up.Root, up.Meta.UploadSid, up.Meta.Filename)
}

func (up *Uploader) Write(temp_file *os.File, body io.Reader) error {
	var err error
	if up.Meta.Range == nil {
		_, err = io.Copy(temp_file, body)
	} else {
		chunk_size := up.Meta.Range.End - up.Meta.Range.Start + 1
		_, err = io.CopyN(temp_file, body, chunk_size)
	}
	return err
}

func (up *Uploader) SaveFiles() ([]*OriginalFile, error) {
	files := make([]*OriginalFile, 0)
	for {
		ofile, err := up.SaveFile()
		if err == io.EOF {
			break
		}

		if err == Incomplete {
			files = append(files, ofile)
			return files, err
		}

		if err != nil {
			return nil, err
		}

		files = append(files, ofile)
	}

	return files, nil
}

func (up *Uploader) SaveFile() (*OriginalFile, error) {
	body, filename, err := up.Reader()
	if err != nil {
		return nil, err
	}

	temp_file, err := up.TempFile()
	if err != nil {
		return nil, err
	}
	defer temp_file.Close()

	if err = up.Write(temp_file, body); err != nil {
		return nil, err
	}

	fi, err := temp_file.Stat()
	if err != nil {
		return nil, err
	}

	ofile := &OriginalFile{Filename: filename, Filepath: temp_file.Name(), Size: fi.Size()}

	if ofile.Size != up.Meta.Range.Size {
		return ofile, Incomplete
	}

	return ofile, nil
}

func TempFile() (*os.File, err) {
	return ioutil.TempFile(os.TempDir(), "pavo")
}

func TempFileChunks(offset int64, storage, upload_sid, user_filename string) (*os.File, error) {
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

	if _, err = file.Seek(offset, 0); err != nil {
		return nil, err
	}

	return file, nil
}

// Get base mime type
//		$ file --mime-type pic.jpg
//		pic.jpg: image/jpeg
func IdentifyMime(file string) (string, error) {
	out, err := exec.Command("file", "--mime-type", file).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("identify: err: %s; detail: %s", err, string(out))
	}

	mime := strings.Split(strings.Split(string(out), ": ")[1], "/")[0]

	return mime, nil
}
