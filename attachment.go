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

type Attachment struct {
	BaseMime string
	Dir      *DirManager
	Files    []FileManager
}

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
		files = append(files, fm)

	}
	return files, nil
}

func saveVersion(ofile *OriginalFile, dm *DirManager, version string, convert string) (FileManager, error) {
	fm := GetFileManager(ofile.BaseMime, version)
	if err := fm.Convert(ofile, convert); err != nil {
		return nil, err
	}

	return fm, nil
}

// ---

// ---
