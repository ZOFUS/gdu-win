//go:build windows

package path

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Feature: windows-port, Property 4: Path Separator Handling
// Validates: Requirements 3.4, 3.5

func TestShortenPath_WindowsPaths(t *testing.T) {
	// Property 4: For any Windows path containing backslashes,
	// ShortenPath() should correctly split and rejoin using backslash separator

	// Test short path (no truncation needed)
	assert.Equal(t, "C:\\Users", ShortenPath("C:\\Users", 20))

	// Test path that needs truncation
	result := ShortenPath("C:\\Users\\dundee\\Documents\\file.txt", 20)
	assert.True(t, strings.Contains(result, "...\\"),
		"Truncated path should contain '...\\', got: %s", result)
	assert.True(t, strings.HasSuffix(result, "file.txt"),
		"Truncated path should end with filename, got: %s", result)

	// Test path that fits exactly
	shortPath := "C:\\foo"
	assert.Equal(t, shortPath, ShortenPath(shortPath, 50))
}

func TestShortenPath_WindowsBackslashSeparator(t *testing.T) {
	// Property 4: Verify backslash separator is used correctly
	path := "C:\\Users\\dundee\\Documents\\Projects\\file.txt"
	result := ShortenPath(path, 25)

	// Result should use backslash as separator
	assert.True(t, strings.Contains(result, "\\"),
		"Result should contain backslash separator, got: %s", result)

	// Result should not contain forward slashes
	assert.False(t, strings.Contains(result, "/"),
		"Result should not contain forward slashes on Windows, got: %s", result)
}

func TestShortenPath_WindowsDeepPath(t *testing.T) {
	// Test deeply nested Windows path
	deepPath := "C:\\Users\\dundee\\Documents\\Projects\\MyProject\\src\\main\\java\\com\\example\\App.java"
	result := ShortenPath(deepPath, 30)

	// Should be truncated
	assert.LessOrEqual(t, len(result), len(deepPath),
		"Result should be shorter or equal to original")

	// Should preserve the filename
	assert.True(t, strings.HasSuffix(result, "App.java"),
		"Should preserve filename, got: %s", result)

	// Should contain truncation indicator
	assert.True(t, strings.Contains(result, "..."),
		"Should contain truncation indicator, got: %s", result)
}

func TestShortenPath_WindowsRootPath(t *testing.T) {
	// Test root path
	assert.Equal(t, "C:\\", ShortenPath("C:\\", 10))
	assert.Equal(t, "D:\\", ShortenPath("D:\\", 10))
}

func TestFilepathSeparator_Windows(t *testing.T) {
	// Verify that filepath.Separator is backslash on Windows
	assert.Equal(t, '\\', filepath.Separator,
		"filepath.Separator should be backslash on Windows")
	assert.Equal(t, "\\", string(filepath.Separator),
		"string(filepath.Separator) should be backslash on Windows")
}

func TestShortenPath_WindowsUNCPath(t *testing.T) {
	// Test UNC path (network path)
	uncPath := "\\\\server\\share\\folder\\file.txt"
	result := ShortenPath(uncPath, 50)

	// Short enough, should not be truncated
	assert.Equal(t, uncPath, result)
}

func TestShortenPath_WindowsPreservesStructure(t *testing.T) {
	// Property 4: Verify path structure is preserved after shortening
	path := "C:\\Users\\dundee\\file.txt"

	// When path fits, it should be unchanged
	result := ShortenPath(path, 100)
	assert.Equal(t, path, result)

	// When truncated, should still be valid path structure
	result = ShortenPath(path, 15)
	parts := strings.Split(result, "\\")
	assert.GreaterOrEqual(t, len(parts), 2,
		"Truncated path should have at least 2 parts, got: %v", parts)
}
