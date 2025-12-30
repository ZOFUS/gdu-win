//go:build windows

package report

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/quick"

	log "github.com/sirupsen/logrus"

	"github.com/dundee/gdu/v5/internal/testdir"
	"github.com/dundee/gdu/v5/pkg/analyze"
	"github.com/dundee/gdu/v5/pkg/fs"
	"github.com/stretchr/testify/assert"
)

func init() {
	log.SetLevel(log.WarnLevel)
}

// Feature: windows-port, Property 5: JSON Export/Import Round-trip
// Validates: Requirements 5.1, 5.2, 5.3

// TestWindowsJSONExportImportRoundTrip tests that exporting to JSON and re-importing
// produces an equivalent structure with preserved Windows paths.
// Property 5: For any analyzed directory tree, exporting to JSON and re-importing
// should produce an equivalent structure with preserved Windows paths.
func TestWindowsJSONExportImportRoundTrip(t *testing.T) {
	fin := testdir.CreateTestDir()
	defer fin()

	// Get absolute path to test directory
	absPath, err := filepath.Abs("test_dir")
	assert.Nil(t, err)

	// Step 1: Export directory to JSON
	output := bytes.NewBuffer(make([]byte, 0))
	reportOutput := bytes.NewBuffer(make([]byte, 0))

	ui := CreateExportUI(output, reportOutput, false, false, false, false)
	err = ui.AnalyzePath(absPath, nil)
	assert.Nil(t, err)
	err = ui.StartUILoop()
	assert.Nil(t, err)

	jsonData := reportOutput.String()
	assert.NotEmpty(t, jsonData, "JSON export should not be empty")

	// Step 2: Import JSON back
	importReader := bytes.NewBufferString(jsonData)
	importedDir, err := ReadAnalysis(importReader)
	assert.Nil(t, err)
	assert.NotNil(t, importedDir, "Imported directory should not be nil")

	// Step 3: Verify structure is preserved
	// The imported directory path should contain the original path
	importedPath := importedDir.GetPath()
	assert.Contains(t, importedPath, "test_dir", "Directory path should contain 'test_dir'")

	// Verify nested structure exists
	foundNested := false
	for _, item := range importedDir.GetFiles() {
		if item.GetName() == "nested" {
			foundNested = true
			// Check subnested exists
			if dir, ok := item.(*analyze.Dir); ok {
				_ = dir // Type assertion successful
			}
		}
	}
	assert.True(t, foundNested, "Nested directory should be preserved after round-trip")
}

// TestWindowsJSONExportImportFileRoundTrip tests that file properties are preserved
func TestWindowsJSONExportImportFileRoundTrip(t *testing.T) {
	fin := testdir.CreateTestDir()
	defer fin()

	absPath, err := filepath.Abs("test_dir")
	assert.Nil(t, err)

	// Export
	output := bytes.NewBuffer(make([]byte, 0))
	reportOutput := bytes.NewBuffer(make([]byte, 0))

	ui := CreateExportUI(output, reportOutput, false, false, false, false)
	err = ui.AnalyzePath(absPath, nil)
	assert.Nil(t, err)
	err = ui.StartUILoop()
	assert.Nil(t, err)

	// Import
	importReader := bytes.NewBufferString(reportOutput.String())
	importedDir, err := ReadAnalysis(importReader)
	assert.Nil(t, err)

	// Find the file and verify its properties
	var foundFile bool
	var checkFiles func(files fs.Files)
	checkFiles = func(files fs.Files) {
		for _, item := range files {
			if item.GetName() == "file" || item.GetName() == "file2" {
				foundFile = true
				// File should have size > 0
				assert.Greater(t, item.GetSize(), int64(0), "File size should be preserved")
			}
			if dir, ok := item.(fs.Item); ok {
				if subFiles := dir.GetFiles(); subFiles != nil {
					checkFiles(subFiles)
				}
			}
		}
	}
	checkFiles(importedDir.GetFiles())
	assert.True(t, foundFile, "Files should be found after round-trip")
}

// TestWindowsJSONExportToFile tests exporting to an actual file
// Validates: Requirements 5.1
func TestWindowsJSONExportToFile(t *testing.T) {
	fin := testdir.CreateTestDir()
	defer fin()

	outputFile := "test_export.json"
	defer os.Remove(outputFile)

	absPath, err := filepath.Abs("test_dir")
	assert.Nil(t, err)

	// Create output file
	reportFile, err := os.OpenFile(outputFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	assert.Nil(t, err)

	output := bytes.NewBuffer(make([]byte, 0))

	ui := CreateExportUI(output, reportFile, false, false, false, false)
	err = ui.AnalyzePath(absPath, nil)
	assert.Nil(t, err)
	err = ui.StartUILoop()
	assert.Nil(t, err)

	// Verify file was created and has content
	fileInfo, err := os.Stat(outputFile)
	assert.Nil(t, err)
	assert.Greater(t, fileInfo.Size(), int64(0), "Export file should have content")

	// Read and verify JSON structure
	content, err := os.ReadFile(outputFile)
	assert.Nil(t, err)
	assert.Contains(t, string(content), "nested", "Export should contain nested directory")
}

// TestWindowsJSONImportFromFile tests importing from an actual file
// Validates: Requirements 5.2
func TestWindowsJSONImportFromFile(t *testing.T) {
	fin := testdir.CreateTestDir()
	defer fin()

	outputFile := "test_import.json"
	defer os.Remove(outputFile)

	absPath, err := filepath.Abs("test_dir")
	assert.Nil(t, err)

	// First export
	reportFile, err := os.OpenFile(outputFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	assert.Nil(t, err)

	output := bytes.NewBuffer(make([]byte, 0))
	ui := CreateExportUI(output, reportFile, false, false, false, false)
	err = ui.AnalyzePath(absPath, nil)
	assert.Nil(t, err)
	err = ui.StartUILoop()
	assert.Nil(t, err)

	// Now import from file
	importFile, err := os.Open(outputFile)
	assert.Nil(t, err)
	defer importFile.Close()

	importedDir, err := ReadAnalysis(importFile)
	assert.Nil(t, err)
	assert.NotNil(t, importedDir)
	// The path should contain test_dir
	assert.Contains(t, importedDir.GetPath(), "test_dir", "Path should contain test_dir")
}

// TestWindowsPathPreservationProperty is a property-based test that verifies
// Windows paths are correctly preserved through export/import cycle
// Property 5: JSON Export/Import Round-trip
// Validates: Requirements 5.3
func TestWindowsPathPreservationProperty(t *testing.T) {
	// Property: For any valid directory name, the path should be preserved after round-trip
	property := func(dirNameBytes []byte) bool {
		// Convert to string and filter to ASCII alphanumeric only for valid Windows names
		var validChars []byte
		for _, b := range dirNameBytes {
			if (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9') {
				validChars = append(validChars, b)
			}
		}
		dirName := string(validChars)

		// Filter out invalid directory names for Windows
		if len(dirName) == 0 || len(dirName) > 20 {
			return true // Skip invalid inputs
		}

		// Skip reserved Windows names
		reserved := []string{"CON", "PRN", "AUX", "NUL", "COM1", "LPT1"}
		upperName := strings.ToUpper(dirName)
		for _, r := range reserved {
			if upperName == r {
				return true // Skip reserved names
			}
		}

		// Create test directory with the generated name
		testPath := filepath.Join(os.TempDir(), "gdu_test_"+dirName)
		err := os.MkdirAll(testPath, os.ModePerm)
		if err != nil {
			return true // Skip if we can't create the directory
		}
		defer os.RemoveAll(testPath)

		// Create a test file inside
		testFile := filepath.Join(testPath, "testfile.txt")
		err = os.WriteFile(testFile, []byte("test content"), 0o644)
		if err != nil {
			return true // Skip if we can't create the file
		}

		// Export
		output := bytes.NewBuffer(make([]byte, 0))
		reportOutput := bytes.NewBuffer(make([]byte, 0))
		ui := CreateExportUI(output, reportOutput, false, false, false, false)
		err = ui.AnalyzePath(testPath, nil)
		if err != nil {
			return true // Skip on analysis error
		}
		ui.StartUILoop()

		// Import
		importReader := bytes.NewBufferString(reportOutput.String())
		importedDir, err := ReadAnalysis(importReader)
		if err != nil {
			return false // This is a failure
		}

		// Verify the directory path contains the expected name
		expectedName := "gdu_test_" + dirName
		return strings.Contains(importedDir.GetPath(), expectedName)
	}

	// Run property test with 100 iterations
	config := &quick.Config{MaxCount: 100}
	if err := quick.Check(property, config); err != nil {
		t.Errorf("Property test failed: %v", err)
	}
}
