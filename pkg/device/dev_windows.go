//go:build windows

package device

import (
	"golang.org/x/sys/windows"
)

// WindowsDevicesInfoGetter returns info for Windows devices
type WindowsDevicesInfoGetter struct{}

// Getter is current instance of DevicesInfoGetter
var Getter DevicesInfoGetter = WindowsDevicesInfoGetter{}

// GetDevicesInfo returns result of GetMounts with usage info about mounted devices
func (t WindowsDevicesInfoGetter) GetDevicesInfo() (Devices, error) {
	drives, err := getLogicalDrives()
	if err != nil {
		return nil, err
	}

	devices := make(Devices, 0, len(drives))
	for _, drive := range drives {
		device, err := getDriveInfo(drive)
		if err != nil {
			continue // Skip drives that can't be accessed
		}
		devices = append(devices, device)
	}

	return devices, nil
}

// GetMounts returns all mounted filesystems (same as GetDevicesInfo on Windows)
func (t WindowsDevicesInfoGetter) GetMounts() (Devices, error) {
	return t.GetDevicesInfo()
}

// getLogicalDrives returns a list of available drive letters
func getLogicalDrives() ([]string, error) {
	mask, err := windows.GetLogicalDrives()
	if err != nil {
		return nil, err
	}

	drives := make([]string, 0)
	for i := 0; i < 26; i++ {
		if mask&(1<<uint(i)) != 0 {
			drive := string(rune('A'+i)) + ":\\"
			drives = append(drives, drive)
		}
	}

	return drives, nil
}

// getDriveInfo returns device info for a specific drive
func getDriveInfo(drive string) (*Device, error) {
	drivePtr, err := windows.UTF16PtrFromString(drive)
	if err != nil {
		return nil, err
	}

	driveType := windows.GetDriveType(drivePtr)

	// Skip network drives, CD-ROMs, and unknown types for now
	if driveType != windows.DRIVE_FIXED && driveType != windows.DRIVE_REMOVABLE {
		return nil, nil
	}

	var freeBytesAvailable, totalBytes, totalFreeBytes uint64
	err = windows.GetDiskFreeSpaceEx(
		drivePtr,
		&freeBytesAvailable,
		&totalBytes,
		&totalFreeBytes,
	)
	if err != nil {
		return nil, err
	}

	// Get volume name
	volumeName := make([]uint16, windows.MAX_PATH+1)
	fsName := make([]uint16, windows.MAX_PATH+1)
	var serialNumber, maxComponentLen, fsFlags uint32

	err = windows.GetVolumeInformation(
		drivePtr,
		&volumeName[0],
		uint32(len(volumeName)),
		&serialNumber,
		&maxComponentLen,
		&fsFlags,
		&fsName[0],
		uint32(len(fsName)),
	)

	name := drive[:2] // e.g., "C:"
	if err == nil && volumeName[0] != 0 {
		name = windows.UTF16ToString(volumeName) + " (" + drive[:2] + ")"
	}

	return &Device{
		Name:       name,
		MountPoint: drive[:2] + "\\",
		Fstype:     windows.UTF16ToString(fsName),
		Size:       int64(totalBytes),
		Free:       int64(totalFreeBytes),
	}, nil
}
