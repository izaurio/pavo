package main

// Original uploaded file.
type OriginalFile struct {
	BaseMime string
	Filepath string
	Filename string
}

// Convert params
type Convert map[string]string

// For range operator return native map.
func (c *Convert) Map() map[string]string {
	return map[string]string(*c)
}

// Attachment contain info about directory, base mime type and all files saved.
type Attachment struct {
	BaseMime string
	Dir      *DirManager
	Files    []FileManager
}

// Function recieve root directory, original file, convertaion parametrs.
// Return Attachment saved.
func SaveAttachment(storage string, ofile *OriginalFile, converts *Convert) (*Attachment, error) {
	dm, err := PrepareDir(storage, ofile.BaseMime)
	if err != nil {
		return nil, err
	}

	attachment := &Attachment{BaseMime: ofile.BaseMime, Dir: dm, Files: make([]FileManager, 0)}

	for version, convert_arg := range converts.Map() {
		fm, err := saveVersion(ofile, dm, version, convert_arg)
		if err != nil {
			return nil, err
		}

		attachment.Files = append(attachment.Files, fm)
	}

	return attachment, nil
}

// Directly save single version and return FileManager.
func saveVersion(ofile *OriginalFile, dm *DirManager, version string, convert string) (FileManager, error) {
	fm := NewFileManager(ofile.BaseMime, version)

	if err := fm.Convert(ofile, convert); err != nil {
		return nil, err
	}

	return fm, nil
}
