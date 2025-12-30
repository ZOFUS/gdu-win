//go:build plan9

package analyze

import (
	"os"
)

func setPlatformSpecificAttrs(file *File, f os.FileInfo) {
	file.Mtime = f.ModTime()
}

func setDirPlatformSpecificAttrs(dir *Dir, path string) {
	stat, err := os.Stat(path)
	if err != nil {
		return
	}
	dir.Mtime = stat.ModTime()
}
