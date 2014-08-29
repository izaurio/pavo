package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSaveAttachment(t *testing.T) {
	assert := assert.New(t)

	ofile := originalFile()
	storage := "dummy/root_storage"
	converts := &Convert{"origin": "", "thumbnail": "80x80"}

	attachment, err := SaveAttachment(storage, ofile, converts)
	assert.Nil(err)
	assert.Equal(len(attachment.Files), 2)
}

func originalFile() *OriginalFile {
	return &OriginalFile{
		BaseMime: "image",
		Filepath: "dummy/32509211_news_bigpic.jpg",
		Filename: "32509211_news_bigpic.jpg",
	}
}
