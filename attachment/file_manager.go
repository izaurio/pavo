package main

type FileManager interface {
	Convert(*OriginalFile, string) error
}

type FileBaseManager struct {
	Version  string
	Filename string
}

// Return FileManager for given base mime and version.
func NewFileManager(mime_base, version string) FileManager {
	switch mime_base {
	case "image":
		return &FileImageManager{FileBaseManager: FileBaseManager{Version: version}}
	}

	return nil
}
