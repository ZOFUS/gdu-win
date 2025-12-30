//go:build windows

package device

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Feature: windows-port, Property 2: Device Info Completeness
// Validates: Requirements 2.2, 2.3

func TestGetDevicesInfo_Windows(t *testing.T) {
	getter := WindowsDevicesInfoGetter{}
	devices, err := getter.GetDevicesInfo()

	assert.Nil(t, err, "GetDevicesInfo() should not return an error")
	assert.IsType(t, Devices{}, devices, "GetDevicesInfo() should return Devices type")
	// At least one drive should exist on any Windows system (typically C:)
	assert.NotEmpty(t, devices, "GetDevicesInfo() should return at least one device")
}

func TestGetMounts_Windows(t *testing.T) {
	getter := WindowsDevicesInfoGetter{}
	devices, err := getter.GetMounts()

	assert.Nil(t, err, "GetMounts() should not return an error")
	assert.IsType(t, Devices{}, devices, "GetMounts() should return Devices type")
}

func TestDeviceFieldsPopulated_Windows(t *testing.T) {
	// Property 2: For any device returned by GetDevicesInfo() on Windows,
	// the device should have non-empty Name, MountPoint containing drive letter,
	// non-empty Fstype, and positive Size value.
	getter := WindowsDevicesInfoGetter{}
	devices, err := getter.GetDevicesInfo()

	assert.Nil(t, err)

	for _, device := range devices {
		// Name should not be empty
		assert.NotEmpty(t, device.Name,
			"Device Name should not be empty")

		// MountPoint should contain drive letter (e.g., "C:\")
		assert.NotEmpty(t, device.MountPoint,
			"Device MountPoint should not be empty")
		assert.True(t, len(device.MountPoint) >= 2,
			"MountPoint should be at least 2 characters (e.g., 'C:\\')")
		assert.True(t, device.MountPoint[1] == ':',
			"MountPoint should contain drive letter format (e.g., 'C:\\')")

		// Fstype should not be empty (e.g., "NTFS", "FAT32")
		assert.NotEmpty(t, device.Fstype,
			"Device Fstype should not be empty, got device: %s", device.Name)

		// Size should be positive
		assert.Greater(t, device.Size, int64(0),
			"Device Size should be positive for device: %s", device.Name)

		// Free space should be non-negative and less than or equal to Size
		assert.GreaterOrEqual(t, device.Free, int64(0),
			"Device Free should be non-negative for device: %s", device.Name)
		assert.LessOrEqual(t, device.Free, device.Size,
			"Device Free should not exceed Size for device: %s", device.Name)
	}
}

func TestDriveLetterFormat_Windows(t *testing.T) {
	// Property 2: MountPoint should contain valid Windows drive letter
	getter := WindowsDevicesInfoGetter{}
	devices, err := getter.GetDevicesInfo()

	assert.Nil(t, err)

	for _, device := range devices {
		// Drive letter should be A-Z
		driveLetter := device.MountPoint[0]
		assert.True(t, driveLetter >= 'A' && driveLetter <= 'Z',
			"Drive letter should be A-Z, got: %c", driveLetter)

		// MountPoint should end with backslash
		assert.True(t, strings.HasSuffix(device.MountPoint, "\\"),
			"MountPoint should end with backslash, got: %s", device.MountPoint)
	}
}

func TestGetLogicalDrives_Windows(t *testing.T) {
	drives, err := getLogicalDrives()

	assert.Nil(t, err, "getLogicalDrives() should not return an error")
	assert.NotEmpty(t, drives, "getLogicalDrives() should return at least one drive")

	for _, drive := range drives {
		// Each drive should be in format "X:\"
		assert.Len(t, drive, 3, "Drive should be 3 characters (e.g., 'C:\\')")
		assert.True(t, drive[0] >= 'A' && drive[0] <= 'Z',
			"Drive letter should be A-Z, got: %s", drive)
		assert.Equal(t, ":\\", drive[1:], "Drive should end with ':\\'")
	}
}
