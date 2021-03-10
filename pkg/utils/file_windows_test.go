// +build windows

package utils

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Due to the way some functions work some tests need to be run
// with data for the running OS

// windows path tests for IsPathInDir (path separator is evaluated at runtime)
var windowsTestPaths = []struct {
	basePath string
	path     string
	isInDir  bool
}{
	{"c:\\media\\.previews", "c:\\media", false}, //1
	{"c:\\media", "c:\\media\\.previews", true},
	{"c:\\", "c:\\my\\media\\path", true},
	{"\\\\netshare\\stash\\media", "\\\\netshare\\opt\\media", false},
	{"\\\\netshare", "\\\\netshare\\stash\\media", true},
	{"c:\\user\\data\\2", "c:\\user\\data\\2", true}, // 6
}

func TestIsPathInDirWindows(t *testing.T) {
	assert := assert.New(t)

	for i, tp := range windowsTestPaths {
		assert.Equal(tp.isInDir, IsPathInDir(tp.basePath, tp.path), "Test "+strconv.Itoa(i+1)+":"+tp.basePath+"... failed")
	}
}
