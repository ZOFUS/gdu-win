//go:build windows

package tui

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecute(t *testing.T) {
	// Use cmd.exe /c echo on Windows instead of Unix 'true' command
	err := Execute("cmd.exe", []string{"/c", "echo", "test"}, []string{})

	assert.Nil(t, err)
}
