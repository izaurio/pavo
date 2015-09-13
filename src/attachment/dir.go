package attachment

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// Directory mananger
type DirManager struct {
	Root string
	Path string
}

// Prepare DirManager given root, mime.
func CreateDir(root, mime string) (*DirManager, error) {
	dm := NewDirManager(root)

	dm.CalcPath(mime)
	if err := dm.Create(); err != nil {
		return nil, err
	}

	return dm, nil
}

// Check path and return DirManager.
func CheckDir(root, path string) (*DirManager, error) {
	dm := NewDirManager(root)

	if m, _ := filepath.Match("/[a-z]*/[0-9]*/[0-9a-z]*/[0-9a-z]*", path); m != true {
		return nil, errors.New("dir: path does not match the pattern")
	}
	dm.Path = path

	return dm, nil
}

// NewDirManager returns a new DirManager given a root.
func NewDirManager(root string) *DirManager {
	return &DirManager{Root: root}
}

// Return absolute path for directory
func (dm *DirManager) Abs() string {
	return filepath.Join(dm.Root, dm.Path)
}

// Create directory obtained by concatenating the root and path.
func (dm *DirManager) Create() error {
	return os.MkdirAll(dm.Root+dm.Path, 0755)
}

// Generate path given mime and date.
func (dm *DirManager) CalcPath(mime string) {
	date := time.Now()
	dm.Path = fmt.Sprintf("/%s/%d/%s/%s", mime, date.Year(), yearDay(date), containerName(date))
}

func yearDay(t time.Time) string {
	return strconv.FormatInt(int64(t.YearDay()), 36)
}

func containerName(t time.Time) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(1000)
	seconds := t.Hour()*3600 + t.Minute()*60 + t.Second()

	return strconv.FormatInt(int64(seconds*1000+r), 36)
}
