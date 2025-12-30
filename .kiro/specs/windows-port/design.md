# Design Document: GDU Windows Port

## Overview

Этот документ описывает архитектуру и дизайн портирования GDU (Go Disk Usage) для Windows. GDU - быстрый анализатор дискового пространства с TUI интерфейсом, изначально разработанный для Linux/Unix систем.

## Architecture

### Высокоуровневая архитектура

```
┌─────────────────────────────────────────────────────────────┐
│                        cmd/gdu/main.go                       │
│                    (CLI Entry Point)                         │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                      cmd/gdu/app/app.go                      │
│                    (Application Logic)                       │
└─────────────────────────────────────────────────────────────┘
                              │
              ┌───────────────┼───────────────┐
              ▼               ▼               ▼
┌─────────────────┐ ┌─────────────────┐ ┌─────────────────┐
│   tui/tui.go    │ │  stdout/stdout  │ │ report/export   │
│ (Interactive)   │ │ (Non-interactive)│ │   (JSON)        │
└─────────────────┘ └─────────────────┘ └─────────────────┘
              │               │               │
              └───────────────┼───────────────┘
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                      pkg/analyze/                            │
│              (Parallel/Sequential Analyzer)                  │
└─────────────────────────────────────────────────────────────┘
                              │
              ┌───────────────┼───────────────┐
              ▼               ▼               ▼
┌─────────────────┐ ┌─────────────────┐ ┌─────────────────┐
│  pkg/device/    │ │    pkg/fs/      │ │   pkg/remove/   │
│ (Device Info)   │ │ (File System)   │ │  (Deletion)     │
└─────────────────┘ └─────────────────┘ └─────────────────┘
```

### Platform-specific Components

```
┌─────────────────────────────────────────────────────────────┐
│                    Platform Abstraction                      │
├─────────────────┬─────────────────┬─────────────────────────┤
│    Windows      │     Linux       │      Plan9/Other        │
├─────────────────┼─────────────────┼─────────────────────────┤
│ dev_windows.go  │ dev_linux.go    │ dev_other.go            │
│ dir_windows.go  │ dir_linux*.go   │ dir_other.go            │
│ exec_windows.go │ exec_other.go   │ exec_other.go           │
└─────────────────┴─────────────────┴─────────────────────────┘
```

## Components and Interfaces

### 1. Device Detection (pkg/device/)

**Interface:**
```go
type DevicesInfoGetter interface {
    GetDevicesInfo() (Devices, error)
    GetMounts() (Devices, error)
}

type Device struct {
    Name       string
    MountPoint string
    Fstype     string
    Size       int64
    Free       int64
}
```

**Windows Implementation (`dev_windows.go`):**
- Uses `windows.GetLogicalDrives()` to enumerate drives
- Uses `windows.GetDiskFreeSpaceEx()` for disk space info
- Uses `windows.GetVolumeInformation()` for volume names
- Filters to DRIVE_FIXED and DRIVE_REMOVABLE types

### 2. File Analysis (pkg/analyze/)

**Interface:**
```go
type Analyzer interface {
    AnalyzeDir(path string, ignore ShouldDirBeIgnored, constGC bool) fs.Item
    GetProgressChan() chan CurrentProgress
    GetDone() SignalGroup
    ResetProgress()
}
```

**Windows-specific (`dir_windows.go`):**
```go
func setPlatformSpecificAttrs(file *File, f os.FileInfo) {
    stat := f.Sys().(*syscall.Win32FileAttributeData)
    file.Mtime = time.Unix(0, stat.LastWriteTime.Nanoseconds())
}
```

### 3. Path Handling (pkg/path/)

**Cross-platform path shortening:**
```go
func ShortenPath(path string, maxLen int) string {
    sep := string(filepath.Separator)
    parts := strings.SplitAfter(path, sep)
    // ... truncation logic using platform separator
}
```

### 4. Shell Execution (tui/)

**Windows (`exec_windows.go`):**
```go
func getShellBin() string {
    shellbin, ok := os.LookupEnv("COMSPEC")
    if !ok {
        shellbin = "C:\\WINDOWS\\System32\\cmd.exe"
    }
    return shellbin
}
```

## Data Models

### Device Model
```go
type Device struct {
    Name       string  // "Local Disk (C:)" or "C:"
    MountPoint string  // "C:\\"
    Fstype     string  // "NTFS", "FAT32", etc.
    Size       int64   // Total bytes
    Free       int64   // Free bytes
}
```

### File/Directory Model
```go
type File struct {
    Name   string
    Flag   rune      // ' ', 'e' (empty), '!' (error), '@' (symlink)
    Size   int64
    Usage  int64     // Disk usage (not available on Windows)
    Mtime  time.Time
    Mli    uint64    // Multi-link inode (not available on Windows)
    Parent fs.Item
}

type Dir struct {
    *File
    BasePath  string
    ItemCount int
    Files     fs.Files
}
```

## Correctness Properties

*A property is a characteristic or behavior that should hold true across all valid executions of a system—essentially, a formal statement about what the system should do. Properties serve as the bridge between human-readable specifications and machine-verifiable correctness guarantees.*

### Property 1: Platform Defaults Correctness
*For any* Windows system, `getDefaultLogFile()` should return "NUL" and `getDefaultStoragePath()` should return a path starting with the system temp directory.
**Validates: Requirements 1.4, 1.5**

### Property 2: Device Info Completeness
*For any* device returned by `GetDevicesInfo()` on Windows, the device should have non-empty Name, MountPoint containing drive letter, non-empty Fstype, and positive Size value.
**Validates: Requirements 2.2, 2.3**

### Property 3: Drive Type Filtering
*For any* drive enumerated by Windows API, only drives of type DRIVE_FIXED or DRIVE_REMOVABLE should be included in the result.
**Validates: Requirements 2.4**

### Property 4: Path Separator Handling
*For any* Windows path containing backslashes, `ShortenPath()` should correctly split and rejoin using backslash separator.
**Validates: Requirements 3.4, 3.5**

### Property 5: JSON Export/Import Round-trip
*For any* analyzed directory tree, exporting to JSON and re-importing should produce an equivalent structure with preserved Windows paths.
**Validates: Requirements 5.3**

## Error Handling

### Device Access Errors
- If a drive cannot be accessed (permissions, not ready), skip it silently
- Log error but continue with other drives

### Path Errors
- Invalid paths return appropriate error messages
- Permission denied errors are displayed to user

### File System Errors
- Read errors mark directory with '!' flag
- Continue analysis of other directories

## Testing Strategy

### Unit Tests
- Test `getDefaultLogFile()` returns "NUL" on Windows
- Test `getDefaultStoragePath()` uses temp directory
- Test `ShortenPath()` with Windows paths
- Test device filtering logic

### Property-Based Tests
Using Go's `testing/quick` or `gopter` library:

1. **Property 1**: Platform defaults test
   - Generate random calls, verify consistent return values

2. **Property 2**: Device info completeness
   - For each device, verify all required fields are populated

3. **Property 4**: Path separator handling
   - Generate random Windows paths, verify correct splitting

5. **Property 5**: JSON round-trip
   - Generate random directory structures, export/import, compare

### Integration Tests
- Full analysis of test directory
- Export and import JSON files
- TUI interaction tests (manual)

### Test Configuration
- Minimum 100 iterations per property test
- Tag format: **Feature: windows-port, Property N: description**
