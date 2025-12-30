//go:build windows

package main

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Feature: windows-port, Property 1: Platform Defaults Correctness
// Validates: Requirements 1.4, 1.5

func TestGetDefaultLogFile_Windows(t *testing.T) {
	// Property 1: For any Windows system, getDefaultLogFile() should return "NUL"
	result := getDefaultLogFile()
	assert.Equal(t, "NUL", result, "getDefaultLogFile() should return 'NUL' on Windows")
}

func TestGetDefaultStoragePath_Windows(t *testing.T) {
	// Property 1: For any Windows system, getDefaultStoragePath() should return
	// a path starting with the system temp directory
	result := getDefaultStoragePath()
	tempDir := os.TempDir()

	assert.True(t, strings.HasPrefix(result, tempDir),
		"getDefaultStoragePath() should start with temp directory '%s', got '%s'", tempDir, result)
	assert.True(t, strings.HasSuffix(result, "gdu-badger"),
		"getDefaultStoragePath() should end with 'gdu-badger', got '%s'", result)
}

func TestGetDefaultIgnoreDirs_Windows(t *testing.T) {
	// On Windows, default ignore dirs should be empty (no /proc, /dev, etc.)
	result := getDefaultIgnoreDirs()
	assert.Empty(t, result, "getDefaultIgnoreDirs() should return empty slice on Windows")
}
