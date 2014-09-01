package upload

import (
	"crypto/md5"
	"encoding/hex"
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

// Error Incomplete returned by uploader when loaded non-last chunk.
var Incomplete = errors.New("Incomplete")

// Structure describes the state of the original file.
type OriginalFile struct {
	BaseMime string
	Filepath string
	Filename string
	Size     int64
}

func (ofile *OriginalFile) Ext() string {
	return strings.ToLower(filepath.Ext(ofile.Filename))
}

// Downloading files from the received request.
// The root directory of storage, storage,  used to temporarily store chunks.
// Returns an array of the original files and error.
// If you load a portion of the file, chunk, it will be stored in err error Incomplete,
// and in an array of a single file. File size will fit the current size.
func Process(req *http.Request, storage string) ([]*OriginalFile, error) {
	meta, err := ParseMeta(req)
	if err != nil {
		return nil, err
	}

	body, err := NewBody(req.Header.Get("X-File"), req.Body)
	if err != nil {
		return nil, err
	}
	up := &Uploader{Root: storage, Meta: meta, Body: body}

	files, err := up.SaveFiles()
	if err == Incomplete {
		return files, err
	}
	if err != nil {
		return nil, err
	}

	return files, nil
}

// Upload manager.
type Uploader struct {
	Root string
	Meta *Meta
	Body *Body
}

// Function SaveFiles sequentially loads the original files or chunk's.
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

// Function loads one or download the original file chunk.
// Asks for the starting position in the body of the request to read the next file.
// Asks for a temporary file.
// Writes data from the request body into a temporary file.
// Specifies the size of the resulting temporary file.
// If the query specified header Content-Range,
// and the size of the resulting file does not match, it returns an error Incomplete.
// Otherwise, defines the basic mime type, and returns the original file.
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

	if up.Meta.Range != nil && ofile.Size != up.Meta.Range.Size {
		return ofile, Incomplete
	}

	ofile.BaseMime, err = IdentifyMime(ofile.Filepath)
	if err != nil {
		return nil, err
	}

	return ofile, nil
}

// Returns the reader to read the file or chunk of request body and the original file name.
// If the request header Conent-Type is multipart/form-data, returns the next copy part.
// If all of part read the case of binary loading read the request body, an error is returned io.EOF.
func (up *Uploader) Reader() (io.Reader, string, error) {
	if up.Meta.MediaType == "multipart/form-data" {
		if up.Body.MR == nil {
			up.Body.MR = multipart.NewReader(up.Body.Body, up.Meta.Boundary)
		}
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

	if up.Body.Available == false {
		return nil, "", io.EOF
	}

	up.Body.Available = false

	return up.Body.Body, up.Meta.Filename, nil
}

// Returns a temporary file to download the file or resume chunk.
func (up *Uploader) TempFile() (*os.File, error) {
	if up.Meta.Range == nil {
		return TempFile()
	}
	return TempFileChunks(up.Meta.Range.Start, up.Root, up.Meta.UploadSid, up.Meta.Filename)
}

// Returns the newly created temporary file.
func TempFile() (*os.File, error) {
	return ioutil.TempFile(os.TempDir(), "pavo")
}

// Returns a temporary file to download chunk.
// To calculate a unique file name used cookie named pavo and the original file name.
// File located in the directory chunks storage root directory.
// Before returning the file pointer is shifted by the value of offset,
// in a situation where the pieces are loaded from the second to the last.
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

// The function writes a temporary file value from reader.
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
