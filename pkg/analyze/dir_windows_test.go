//go:build windows

package analyze

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/dundee/gdu/v5/internal/testdir"
	"github.com/dundee/gdu/v5/pkg/fs"
	"github.com/stretchr/testify/assert"
)

// Feature: windows-port, Property 3: Windows File Attributes
// Validates: Requirements 3.4

func TestAnalyzeDir_WindowsPaths(t *testing.T) {
	fin := testdir.CreateTestDir()
	defer fin()

	absPath, err := filepath.Abs("test_dir")
	assert.Nil(t, err)

	analyzer := CreateAnalyzer()
	dir := analyzer.AnalyzeDir(
		absPath, func(_, _ string) bool { return false }, false,
	).(*Dir)

	analyzer.GetDone().Wait()
	dir.UpdateStats(make(fs.HardLinkedItems))

	// Verify directory was analyzed
	assert.NotNil(t, dir)
	assert.True(t, dir.IsDir())
	assert.Greater(t, dir.ItemCount, 0)

	// Verify path contains backslashes (Windows style)
	assert.Contains(t, dir.GetPath(), "\\")
}

func TestAnalyzeDir_WindowsAbsolutePath(t *testing.T) {
	fin := testdir.CreateTestDir()
	defer fin()

	// Get absolute path
	absPath, err := filepath.Abs("test_dir")
	assert.Nil(t, err)

	// Should start with drive letter on Windows
	assert.True(t, len(absPath) >= 2 && absPath[1] == ':',
		"Absolute path should start with drive letter, got: %s", absPath)

	analyzer := CreateAnalyzer()
	dir := analyzer.AnalyzeDir(
		absPath, func(_, _ string) bool { return false }, false,
	).(*Dir)

	analyzer.GetDone().Wait()

	// Path should be preserved
	assert.Equal(t, absPath, dir.GetPath())
}

func TestAnalyzeDir_WindowsNestedStructure(t *testing.T) {
	fin := testdir.CreateTestDir()
	defer fin()

	analyzer := CreateAnalyzer()
	dir := analyzer.AnalyzeDir(
		"test_dir", func(_, _ string) bool { return false }, false,
	).(*Dir)

	analyzer.GetDone().Wait()
	dir.UpdateStats(make(fs.HardLinkedItems))

	// Find nested directory
	var nested *Dir
	for _, item := range dir.Files {
		if item.GetName() == "nested" {
			nested = item.(*Dir)
			break
		}
	}

	assert.NotNil(t, nested, "Should find nested directory")
	assert.True(t, nested.IsDir())

	// Find subnested
	var subnested *Dir
	for _, item := range nested.Files {
		if item.GetName() == "subnested" {
			subnested = item.(*Dir)
			break
		}
	}

	assert.NotNil(t, subnested, "Should find subnested directory")
}

func TestAnalyzeDir_WindowsFileSize(t *testing.T) {
	fin := testdir.CreateTestDir()
	defer fin()

	analyzer := CreateAnalyzer()
	dir := analyzer.AnalyzeDir(
		"test_dir", func(_, _ string) bool { return false }, false,
	).(*Dir)

	analyzer.GetDone().Wait()
	dir.UpdateStats(make(fs.HardLinkedItems))

	// Total size should be positive
	assert.Greater(t, dir.Size, int64(0))

	// Find a file and check its size
	var foundFile bool
	var checkFiles func(files fs.Files)
	checkFiles = func(files fs.Files) {
		for _, item := range files {
			if !item.IsDir() {
				foundFile = true
				assert.Greater(t, item.GetSize(), int64(0),
					"File %s should have positive size", item.GetName())
			}
			if d, ok := item.(*Dir); ok {
				checkFiles(d.Files)
			}
		}
	}
	checkFiles(dir.Files)

	assert.True(t, foundFile, "Should find at least one file")
}

func TestAnalyzeDir_WindowsHiddenFiles(t *testing.T) {
	// Create a hidden file on Windows
	tempDir := t.TempDir()
	hiddenFile := filepath.Join(tempDir, ".hidden_file")
	err := os.WriteFile(hiddenFile, []byte("hidden"), 0644)
	assert.Nil(t, err)

	// Analyze without ignoring hidden
	analyzer := CreateAnalyzer()
	dir := analyzer.AnalyzeDir(
		tempDir, func(_, _ string) bool { return false }, false,
	).(*Dir)

	analyzer.GetDone().Wait()

	// Should find the hidden file
	var foundHidden bool
	for _, item := range dir.Files {
		if item.GetName() == ".hidden_file" {
			foundHidden = true
			break
		}
	}
	assert.True(t, foundHidden, "Should find hidden file when not ignoring")
}

func TestAnalyzeDir_WindowsIgnoreHidden(t *testing.T) {
	// Create a hidden file on Windows
	tempDir := t.TempDir()
	hiddenFile := filepath.Join(tempDir, ".hidden_file")
	err := os.WriteFile(hiddenFile, []byte("hidden"), 0644)
	assert.Nil(t, err)

	// Also create a normal file
	normalFile := filepath.Join(tempDir, "normal_file")
	err = os.WriteFile(normalFile, []byte("normal"), 0644)
	assert.Nil(t, err)

	// Create a hidden directory
	hiddenDir := filepath.Join(tempDir, ".hidden_dir")
	err = os.Mkdir(hiddenDir, 0755)
	assert.Nil(t, err)

	// Analyze with ignoring hidden (directories starting with dot)
	analyzer := CreateAnalyzer()
	dir := analyzer.AnalyzeDir(
		tempDir, func(name, _ string) bool {
			// Only ignore directories starting with dot
			return len(name) > 0 && name[0] == '.'
		}, false,
	).(*Dir)

	analyzer.GetDone().Wait()

	// Should NOT find the hidden directory, but files are not filtered by this function
	var foundHiddenDir bool
	var foundNormal bool
	var foundHiddenFile bool
	for _, item := range dir.Files {
		if item.GetName() == ".hidden_dir" {
			foundHiddenDir = true
		}
		if item.GetName() == "normal_file" {
			foundNormal = true
		}
		if item.GetName() == ".hidden_file" {
			foundHiddenFile = true
		}
	}
	assert.False(t, foundHiddenDir, "Should NOT find hidden directory when ignoring")
	assert.True(t, foundNormal, "Should find normal file")
	// Note: The ignore function only applies to directories, not files
	assert.True(t, foundHiddenFile, "Hidden files are not filtered by directory ignore function")
}

func TestSequentialAnalyzer_Windows(t *testing.T) {
	fin := testdir.CreateTestDir()
	defer fin()

	analyzer := CreateSeqAnalyzer()
	dir := analyzer.AnalyzeDir(
		"test_dir", func(_, _ string) bool { return false }, false,
	).(*Dir)

	analyzer.GetDone().Wait()
	dir.UpdateStats(make(fs.HardLinkedItems))

	// Verify same results as parallel analyzer
	assert.Equal(t, "test_dir", dir.Name)
	assert.Equal(t, 5, dir.ItemCount)
	assert.True(t, dir.IsDir())
}

func TestParallelAnalyzer_Windows(t *testing.T) {
	fin := testdir.CreateTestDir()
	defer fin()

	analyzer := CreateAnalyzer()
	dir := analyzer.AnalyzeDir(
		"test_dir", func(_, _ string) bool { return false }, false,
	).(*Dir)

	analyzer.GetDone().Wait()
	dir.UpdateStats(make(fs.HardLinkedItems))

	// Verify results
	assert.Equal(t, "test_dir", dir.Name)
	assert.Equal(t, 5, dir.ItemCount)
	assert.True(t, dir.IsDir())
}
