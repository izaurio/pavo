package main

type FileImageManager struct {
	FileBaseManager
	Width  int
	Height int
	Size   int64
}

// Save version from original with convert command-line tool.
func (fim *FileImageManager) Convert(ofile *OriginalFile, convert string) error {
	// calc filename
	fim.Filename = "pic-abc.jpg"

	// convert image

	// identify sizes
	fim.Width = 120
	fim.Height = 90
	fim.Size = 20 * 1024

	return nil
}
