//go:build windows

package app

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	log "github.com/sirupsen/logrus"

	"github.com/dundee/gdu/v5/internal/testapp"
	"github.com/dundee/gdu/v5/internal/testdir"
	"github.com/dundee/gdu/v5/pkg/device"
	"github.com/stretchr/testify/assert"
)

func init() {
	log.SetLevel(log.WarnLevel)
}

// Feature: windows-port, Integration Tests
// Validates: Requirements 2.1, 3.1, 3.2, 3.3

// TestWindowsNonInteractiveAnalyzeCurrentDir tests non-interactive mode analyzing current directory
// Validates: Requirements 3.1
func TestWindowsNonInteractiveAnalyzeCurrentDir(t *testing.T) {
	fin := testdir.CreateTestDir()
	defer fin()

	out, err := runWindowsApp(
		&Flags{LogFile: "NUL"},
		[]string{"test_dir"},
		false,
		device.Getter,
	)

	assert.Nil(t, err)
	assert.Contains(t, out, "nested", "Output should contain 'nested' directory")
}

// TestWindowsNonInteractiveAnalyzeWithPath tests non-interactive mode with explicit path
// Validates: Requirements 3.2
func TestWindowsNonInteractiveAnalyzeWithPath(t *testing.T) {
	fin := testdir.CreateTestDir()
	defer fin()

	// Get absolute path with Windows backslashes
	absPath, err := filepath.Abs("test_dir")
	assert.Nil(t, err)

	out, err := runWindowsApp(
		&Flags{LogFile: "NUL"},
		[]string{absPath},
		false,
		device.Getter,
	)

	assert.Nil(t, err)
	assert.Contains(t, out, "nested", "Output should contain 'nested' directory")
}

// TestWindowsNonInteractiveSummary tests summary mode (-s flag equivalent)
// Validates: Requirements 3.3
func TestWindowsNonInteractiveSummary(t *testing.T) {
	fin := testdir.CreateTestDir()
	defer fin()

	out, err := runWindowsApp(
		&Flags{LogFile: "NUL", Summarize: true},
		[]string{"test_dir"},
		false,
		device.Getter,
	)

	assert.Nil(t, err)
	// Summary mode should show total size
	assert.NotEmpty(t, out, "Summary output should not be empty")
}

// TestWindowsListDrives tests listing Windows drives (-d flag)
// Validates: Requirements 2.1
func TestWindowsListDrives(t *testing.T) {
	out, err := runWindowsApp(
		&Flags{LogFile: "NUL", ShowDisks: true},
		[]string{},
		false,
		device.Getter,
	)

	assert.Nil(t, err)
	// Should show device information
	assert.Contains(t, out, "Device", "Output should contain 'Device' header")
}

// TestWindowsPathWithBackslashes tests that Windows paths with backslashes work correctly
// Validates: Requirements 3.4
func TestWindowsPathWithBackslashes(t *testing.T) {
	fin := testdir.CreateTestDir()
	defer fin()

	// Create a path with explicit backslashes
	path := "test_dir\\nested"

	out, err := runWindowsApp(
		&Flags{LogFile: "NUL"},
		[]string{path},
		false,
		device.Getter,
	)

	assert.Nil(t, err)
	assert.Contains(t, out, "subnested", "Output should contain 'subnested' directory")
}

func runWindowsApp(flags *Flags, args []string, istty bool, getter device.DevicesInfoGetter) (output string, err error) {
	buff := bytes.NewBufferString("")

	app := App{
		Flags:       flags,
		Args:        args,
		Istty:       istty,
		Writer:      buff,
		TermApp:     testapp.CreateMockedApp(false),
		Getter:      getter,
		PathChecker: os.Stat,
	}
	err = app.Run()

	return strings.TrimSpace(buff.String()), err
}
