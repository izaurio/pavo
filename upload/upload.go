package upload

import (
	"fmt"
	"os/exec"
	"strings"
)

// Original uploaded file.
type OriginalFile struct {
	BaseMime string
	Filepath string
	Filename string
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
