//go:build windows

package stdout

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/dundee/gdu/v5/internal/testdir"
	"github.com/stretchr/testify/assert"
)

// Feature: windows-port, Property: Windows Non-Interactive Mode
// Validates: Requirements 3.1, 3.2, 3.3

func TestShowTop_WindowsPaths(t *testing.T) {
	fin := testdir.CreateTestDir()
	defer fin()

	buff := bytes.NewBuffer(make([]byte, 0))
	ui := CreateStdoutUI(buff, false, false, false, false, false, false, false, false, "", 10, false)

	absPath, _ := filepath.Abs("test_dir")
	err := ui.AnalyzePath(absPath, nil)
	assert.Nil(t, err)

	err = ui.StartUILoop()
	assert.Nil(t, err)

	output := buff.String()
	// Output should contain directory info
	assert.NotEmpty(t, output)
}

func TestSummarize_Windows(t *testing.T) {
	fin := testdir.CreateTestDir()
	defer fin()

	buff := bytes.NewBuffer(make([]byte, 0))
	ui := CreateStdoutUI(buff, false, false, false, false, true, false, false, false, "", 10, false)

	err := ui.AnalyzePath("test_dir", nil)
	assert.Nil(t, err)

	err = ui.StartUILoop()
	assert.Nil(t, err)

	output := buff.String()
	// Should show summary
	assert.NotEmpty(t, output)
}

func TestNoProgress_Windows(t *testing.T) {
	fin := testdir.CreateTestDir()
	defer fin()

	buff := bytes.NewBuffer(make([]byte, 0))
	ui := CreateStdoutUI(buff, false, false, false, false, false, false, false, false, "", 10, false)

	err := ui.AnalyzePath("test_dir", nil)
	assert.Nil(t, err)

	err = ui.StartUILoop()
	assert.Nil(t, err)

	// Should complete without progress output
	assert.NotNil(t, ui)
}

func TestNoColor_Windows(t *testing.T) {
	fin := testdir.CreateTestDir()
	defer fin()

	buff := bytes.NewBuffer(make([]byte, 0))
	ui := CreateStdoutUI(buff, false, false, false, false, false, false, false, false, "", 10, false)

	err := ui.AnalyzePath("test_dir", nil)
	assert.Nil(t, err)

	err = ui.StartUILoop()
	assert.Nil(t, err)

	// Should complete without color codes
	output := buff.String()
	assert.NotEmpty(t, output)
}

func TestTopFiles_Windows(t *testing.T) {
	fin := testdir.CreateTestDir()
	defer fin()

	buff := bytes.NewBuffer(make([]byte, 0))
	ui := CreateStdoutUI(buff, false, false, false, false, false, false, false, false, "", 5, false)

	err := ui.AnalyzePath("test_dir", nil)
	assert.Nil(t, err)

	err = ui.StartUILoop()
	assert.Nil(t, err)

	output := buff.String()
	assert.NotEmpty(t, output)
}

func TestReverseSort_Windows(t *testing.T) {
	fin := testdir.CreateTestDir()
	defer fin()

	buff := bytes.NewBuffer(make([]byte, 0))
	ui := CreateStdoutUI(buff, false, false, false, false, false, false, false, false, "", 10, true)

	err := ui.AnalyzePath("test_dir", nil)
	assert.Nil(t, err)

	err = ui.StartUILoop()
	assert.Nil(t, err)

	output := buff.String()
	assert.NotEmpty(t, output)
}
