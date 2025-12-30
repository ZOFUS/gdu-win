//go:build windows

package common_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/dundee/gdu/v5/internal/common"
	"github.com/stretchr/testify/assert"
)

// Feature: windows-port, Property: Windows Path Ignore Patterns
// Validates: Requirements 6.4, 6.5

func TestIgnoreByWindowsAbsPath(t *testing.T) {
	ui := &common.UI{}
	// Use Windows-style absolute path
	absPath, _ := filepath.Abs("test_dir")
	ui.SetIgnoreDirPaths([]string{absPath})
	shouldBeIgnored := ui.CreateIgnoreFunc()

	assert.True(t, shouldBeIgnored("test_dir", absPath))
}

func TestIgnoreByWindowsRelativePath(t *testing.T) {
	ui := &common.UI{}
	ui.SetIgnoreDirPaths([]string{"test_dir\\abc"})
	shouldBeIgnored := ui.CreateIgnoreFunc()

	assert.True(t, shouldBeIgnored("abc", "test_dir\\abc"))
}

func TestIgnorePatternWithWindowsPath(t *testing.T) {
	ui := &common.UI{}
	// Pattern should work with Windows paths
	err := ui.SetIgnoreDirPatterns([]string{"[abc]+"})
	assert.Nil(t, err)
	shouldBeIgnored := ui.CreateIgnoreFunc()

	assert.True(t, shouldBeIgnored("aaa", "aaa"))
	assert.True(t, shouldBeIgnored("abc", "abc"))
	assert.False(t, shouldBeIgnored("xyz", "xyz"))
}

func TestIgnoreFromFileWithWindowsPaths(t *testing.T) {
	// Create temp ignore file
	tempFile := filepath.Join(os.TempDir(), "gdu_test_ignore")
	file, err := os.OpenFile(tempFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		t.Fatal(err)
	}

	// Write patterns
	file.WriteString("test_dir\\aaa\n")
	file.WriteString("[abc]+\n")
	file.Close()
	defer os.Remove(tempFile)

	ui := &common.UI{}
	err = ui.SetIgnoreFromFile(tempFile)
	assert.Nil(t, err)
	shouldBeIgnored := ui.CreateIgnoreFunc()

	assert.True(t, shouldBeIgnored("aaa", "test_dir\\aaa"))
	assert.True(t, shouldBeIgnored("abc", "abc"))
}

func TestIgnoreHiddenOnWindows(t *testing.T) {
	ui := &common.UI{}
	ui.SetIgnoreHidden(true)
	shouldBeIgnored := ui.CreateIgnoreFunc()

	// Files starting with dot should be ignored
	assert.True(t, shouldBeIgnored(".git", "C:\\Users\\test\\.git"))
	assert.True(t, shouldBeIgnored(".hidden", "C:\\Users\\test\\.hidden"))
	assert.False(t, shouldBeIgnored("normal", "C:\\Users\\test\\normal"))
}

func TestIgnoreCombinedOnWindows(t *testing.T) {
	ui := &common.UI{}
	ui.SetIgnoreDirPaths([]string{"C:\\Users\\test\\ignore_me"})
	err := ui.SetIgnoreDirPatterns([]string{"[abc]+"})
	assert.Nil(t, err)
	ui.SetIgnoreHidden(true)
	shouldBeIgnored := ui.CreateIgnoreFunc()

	// Test all ignore methods
	assert.True(t, shouldBeIgnored("ignore_me", "C:\\Users\\test\\ignore_me"))
	assert.True(t, shouldBeIgnored("abc", "abc"))
	assert.True(t, shouldBeIgnored(".git", "C:\\Users\\test\\.git"))
	assert.False(t, shouldBeIgnored("normal", "C:\\Users\\test\\normal"))
}
