//go:build !windows
// +build !windows

package tui

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecute(t *testing.T) {
	err := Execute("true", []string{}, []string{})

	assert.Nil(t, err)
}
