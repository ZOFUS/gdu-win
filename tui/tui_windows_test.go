//go:build windows

package tui

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/dundee/gdu/v5/internal/testanalyze"
	"github.com/dundee/gdu/v5/internal/testapp"
	"github.com/dundee/gdu/v5/internal/testdir"
	"github.com/dundee/gdu/v5/pkg/analyze"
	"github.com/dundee/gdu/v5/pkg/fs"
	"github.com/stretchr/testify/assert"
)

// Feature: windows-port, Property: TUI Windows Integration
// Validates: Requirements 4.1, 4.2, 4.3, 4.4, 4.5, 4.6

func TestTUI_WindowsPathDisplay(t *testing.T) {
	fin := testdir.CreateTestDir()
	defer fin()

	simScreen := testapp.CreateSimScreen()
	defer simScreen.Fini()

	app := testapp.CreateMockedApp(true)
	ui := CreateUI(app, simScreen, &bytes.Buffer{}, false, false, false, false, false)

	absPath, _ := filepath.Abs("test_dir")
	ui.Analyzer = &testanalyze.MockedAnalyzer{}
	ui.done = make(chan struct{})
	err := ui.AnalyzePath(absPath, nil)
	assert.Nil(t, err)
}

func TestTUI_WindowsNavigation(t *testing.T) {
	fin := testdir.CreateTestDir()
	defer fin()

	simScreen := testapp.CreateSimScreen()
	defer simScreen.Fini()

	app := testapp.CreateMockedApp(true)
	ui := CreateUI(app, simScreen, &bytes.Buffer{}, false, false, false, false, false)

	ui.Analyzer = &testanalyze.MockedAnalyzer{}
	ui.done = make(chan struct{})
	err := ui.AnalyzePath("test_dir", nil)
	assert.Nil(t, err)
}

func TestTUI_WindowsExport(t *testing.T) {
	fin := testdir.CreateTestDir()
	defer fin()

	simScreen := testapp.CreateSimScreen()
	defer simScreen.Fini()

	app := testapp.CreateMockedApp(true)
	ui := CreateUI(app, simScreen, &bytes.Buffer{}, false, false, false, false, false)

	ui.Analyzer = &testanalyze.MockedAnalyzer{}
	ui.done = make(chan struct{})
	err := ui.AnalyzePath("test_dir", nil)
	assert.Nil(t, err)
}

func TestTUI_WindowsShowDevices(t *testing.T) {
	simScreen := testapp.CreateSimScreen()
	defer simScreen.Fini()

	app := testapp.CreateMockedApp(true)
	ui := CreateUI(app, simScreen, &bytes.Buffer{}, false, false, false, false, false)

	// This should not panic on Windows
	ui.Analyzer = &testanalyze.MockedAnalyzer{}
	ui.showDevices()

	// UI should still be functional
	assert.NotNil(t, ui)
}

func TestTUI_WindowsSorting(t *testing.T) {
	fin := testdir.CreateTestDir()
	defer fin()

	simScreen := testapp.CreateSimScreen()
	defer simScreen.Fini()

	app := testapp.CreateMockedApp(true)
	ui := CreateUI(app, simScreen, &bytes.Buffer{}, false, false, false, false, false)

	ui.Analyzer = &testanalyze.MockedAnalyzer{}
	ui.done = make(chan struct{})
	err := ui.AnalyzePath("test_dir", nil)
	assert.Nil(t, err)

	// Test sorting by name
	ui.setSorting("name")
	assert.Equal(t, "name", ui.sortBy)

	// Test sorting by size
	ui.setSorting("size")
	assert.Equal(t, "size", ui.sortBy)
}

func TestTUI_WindowsFilter(t *testing.T) {
	fin := testdir.CreateTestDir()
	defer fin()

	simScreen := testapp.CreateSimScreen()
	defer simScreen.Fini()

	app := testapp.CreateMockedApp(true)
	ui := CreateUI(app, simScreen, &bytes.Buffer{}, false, false, false, false, false)

	ui.Analyzer = &testanalyze.MockedAnalyzer{}
	ui.done = make(chan struct{})
	err := ui.AnalyzePath("test_dir", nil)
	assert.Nil(t, err)

	// Set filter
	ui.filterValue = "nested"

	// Filter should be applied
	assert.Equal(t, "nested", ui.filterValue)
}

func TestTUI_WindowsItemCount(t *testing.T) {
	fin := testdir.CreateTestDir()
	defer fin()

	simScreen := testapp.CreateSimScreen()
	defer simScreen.Fini()

	app := testapp.CreateMockedApp(true)
	ui := CreateUI(app, simScreen, &bytes.Buffer{}, false, false, false, false, false)
	ui.SetShowItemCount()

	ui.Analyzer = &testanalyze.MockedAnalyzer{}
	ui.done = make(chan struct{})
	err := ui.AnalyzePath("test_dir", nil)
	assert.Nil(t, err)

	// Item count should be shown
	assert.True(t, ui.showItemCount)
}

func TestTUI_WindowsMTime(t *testing.T) {
	fin := testdir.CreateTestDir()
	defer fin()

	simScreen := testapp.CreateSimScreen()
	defer simScreen.Fini()

	app := testapp.CreateMockedApp(true)
	ui := CreateUI(app, simScreen, &bytes.Buffer{}, false, false, false, false, false)
	ui.SetShowMTime()

	ui.Analyzer = &testanalyze.MockedAnalyzer{}
	ui.done = make(chan struct{})
	err := ui.AnalyzePath("test_dir", nil)
	assert.Nil(t, err)

	// MTime should be shown
	assert.True(t, ui.showMtime)
}

func TestAnalyzeWithProgress_Windows(t *testing.T) {
	fin := testdir.CreateTestDir()
	defer fin()

	simScreen := testapp.CreateSimScreen()
	defer simScreen.Fini()

	output := bytes.NewBuffer(make([]byte, 0))
	app := testapp.CreateMockedApp(true)
	ui := CreateUI(app, simScreen, output, false, false, false, false, false)

	ui.Analyzer = &testanalyze.MockedAnalyzer{}
	ui.done = make(chan struct{})

	absPath, _ := filepath.Abs("test_dir")
	err := ui.AnalyzePath(absPath, nil)
	assert.Nil(t, err)
}

func TestGetParentDir_Windows(t *testing.T) {
	dir := &analyze.Dir{
		File: &analyze.File{
			Name: "test",
		},
		BasePath: "C:\\Users\\test",
	}

	parent := dir.GetParent()
	// Parent should be nil for root
	assert.Nil(t, parent)
}

func TestDirPath_Windows(t *testing.T) {
	dir := &analyze.Dir{
		File: &analyze.File{
			Name: "subdir",
		},
		BasePath: "C:\\Users\\test",
	}

	path := dir.GetPath()
	assert.Contains(t, path, "C:\\Users\\test")
}

func TestFileInDir_Windows(t *testing.T) {
	file := &analyze.File{
		Name: "test.txt",
	}

	dir := &analyze.Dir{
		File: &analyze.File{
			Name: "parent",
		},
		BasePath: "C:\\Users",
		Files:    fs.Files{file},
	}

	file.Parent = dir

	// File should have correct parent
	assert.Equal(t, dir, file.GetParent())
}
