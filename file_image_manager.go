package main

type FileImageManager struct {
	FileBaseManager
	Width  int
	Height int
	Size   int64
}

// Save version from original with convert command-line tool.
func (fim *FileImageManager) Convert(ofile *OriginalFile, convert string) error {
	return nil
}
