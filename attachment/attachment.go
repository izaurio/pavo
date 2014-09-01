package attachment

import "github.com/kavkaz/pavo/upload"

// Attachment contain info about directory, base mime type and all files saved.
type Attachment struct {
	OriginalFile *upload.OriginalFile
	Dir          *DirManager
	Versions     map[string]FileManager
}

// Function recieve root directory, original file, convertaion parametrs.
// Return Attachment saved.
func Create(storage string, ofile *upload.OriginalFile, converts map[string]string) (*Attachment, error) {
	dm, err := CreateDir(storage, ofile.BaseMime)
	if err != nil {
		return nil, err
	}

	attachment := &Attachment{
		OriginalFile: ofile,
		Dir:          dm,
		Versions:     make(map[string]FileManager),
	}

	if ofile.BaseMime == "image" {
		converts["thumbnail"] = "120x90"
	}
	for version, convert_opt := range converts {
		fm, err := attachment.CreateVersion(version, convert_opt)
		if err != nil {
			return nil, err
		}

		attachment.Versions[version] = fm
	}

	return attachment, nil
}

// Directly save single version and return FileManager.
func (attachment *Attachment) CreateVersion(version string, convert string) (FileManager, error) {
	fm := NewFileManager(attachment.Dir, attachment.OriginalFile.BaseMime, version)
	fm.SetFilename(attachment.OriginalFile.Ext())

	if err := fm.Convert(attachment.OriginalFile.Filepath, convert); err != nil {
		return nil, err
	}

	return fm, nil
}

func (attachment *Attachment) ToJson() map[string]interface{} {
	data := make(map[string]interface{})
	data["type"] = attachment.OriginalFile.BaseMime
	data["dir"] = attachment.Dir.Path
	data["name"] = attachment.OriginalFile.Filename
	versions := make(map[string]interface{})
	for version, fm := range attachment.Versions {
		versions[version] = fm.ToJson()
	}
	data["versions"] = versions

	return data
}
