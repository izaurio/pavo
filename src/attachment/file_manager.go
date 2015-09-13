package attachment

import (
	"path/filepath"
	"strconv"
	"time"
)

type FileManager interface {
	Convert(string, string) error
	SetFilename(string)
	ToJson() map[string]interface{}
}

type FileBaseManager struct {
	Dir      *DirManager
	Version  string
	Filename string
}

// Return FileManager for given base mime and version.
func NewFileManager(dm *DirManager, mime_base, version string) FileManager {
	fbm := &FileBaseManager{Dir: dm, Version: version}
	switch mime_base {
	case "image":
		return &FileImageManager{FileBaseManager: fbm}
	default:
		return &FileDefaultManager{FileBaseManager: fbm}
	}

	return nil
}

func (fbm *FileBaseManager) SetFilename(ext string) {
	salt := strconv.FormatInt(seconds(), 36)
	fbm.Filename = fbm.Version + "-" + salt + ext
}

func (fbm *FileBaseManager) Filepath() string {
	return filepath.Join(fbm.Dir.Abs(), fbm.Filename)
}

func (fbm *FileBaseManager) Url() string {
	return filepath.Join(fbm.Dir.Path, fbm.Filename)
}

func seconds() int64 {
	t := time.Now()
	return int64(t.Hour()*3600 + t.Minute()*60 + t.Second())
}
