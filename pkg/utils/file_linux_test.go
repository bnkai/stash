// +build linux

package utils

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Due to the way some functions work some tests need to be run
// with data for the running OS

// linux path tests for IsPathInDir (path separator is evaluated at runtime)
var linuxTestPaths = []struct {
	basePath string
	path     string
	isInDir  bool
}{
	{"/media/.previews", "/media", false}, //1
	{"/media", "/media/.previews", true},
	{"/", "/my/media/path", true},
	{"/opt/stash/media", "/opt/media", false},
	{"/opt/", "/opt/stash/media", true},
	{"/opt/stash/", "/opt/stash", true}, // 6
}

func TestIsPathInDirLinux(t *testing.T) {
	assert := assert.New(t)

	for i, tp := range linuxTestPaths {
		assert.Equal(tp.isInDir, IsPathInDir(tp.basePath, tp.path), "Test "+strconv.Itoa(i+1)+":"+tp.basePath+"... failed")
	}
}
