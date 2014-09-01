package attachment

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrepareDir(t *testing.T) {
	assert := assert.New(t)
	root := "dummy/root_storage"

	dm, err := CreateDir(root, "image")
	assert.Nil(err)
	assert.Equal(root, dm.Root)

	dm, err = CheckDir(root, "/image/2014/2a/q1b12")
	assert.Nil(err)
}
