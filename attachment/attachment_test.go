package attachment

import (
	"testing"

	"github.com/kavkaz/pavo/upload"
	"github.com/stretchr/testify/assert"
)

func TestCreateAttachment(t *testing.T) {
	assert := assert.New(t)

	ofile := originalImageFile()
	storage := "../dummy/root_storage"
	converts := map[string]string{"original": "", "thumbnail": "120x80"}

	attachment, err := Create(storage, ofile, converts)
	assert.Nil(err)
	assert.Equal(len(attachment.Versions), 2)

	data := attachment.ToJson()
	assert.Equal(data["type"], "image")
}

func originalImageFile() *upload.OriginalFile {
	return &upload.OriginalFile{
		BaseMime: "image",
		Filepath: "../dummy/32509211_news_bigpic.jpg",
		Filename: "32509211_news_bigpic.jpg",
	}
}

func originalPdfFile() *upload.OriginalFile {
	return &upload.OriginalFile{
		BaseMime: "application",
		Filepath: "../dummy/Learning-Go-latest.pdf",
		Filename: "Learning-Go-latest.pdf",
	}
}
