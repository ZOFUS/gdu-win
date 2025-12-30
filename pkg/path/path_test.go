package path

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShortenPath(t *testing.T) {
	sep := string(filepath.Separator)

	// Test with platform-native paths
	root := sep + "root"
	assert.Equal(t, root, ShortenPath(root, 10))

	homeDundeeFoo := sep + "home" + sep + "dundee" + sep + "foo"
	assert.Equal(t, sep+"home"+sep+"..."+sep+"foo", ShortenPath(homeDundeeFoo, 10))
	assert.Equal(t, homeDundeeFoo, ShortenPath(homeDundeeFoo, 50))

	homeDundeeFooBar := sep + "home" + sep + "dundee" + sep + "foo" + sep + "bar.txt"
	assert.Equal(t, sep+"home"+sep+"dundee"+sep+"..."+sep+"bar.txt", ShortenPath(homeDundeeFooBar, 20))
	assert.Equal(t, sep+"home"+sep+"..."+sep+"bar.txt", ShortenPath(homeDundeeFooBar, 15))
}
