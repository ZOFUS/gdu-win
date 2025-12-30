# GDU Windows Port - Requirements

## Introduction

Porting GDU (Go Disk Usage) for full Windows PowerShell support. GDU is a fast disk usage analyzer written in Go.

## Glossary

- **GDU**: Go Disk Usage - disk space analyzer
- **TUI**: Text User Interface
- **PowerShell**: Windows command shell

## Requirements

### Requirement 1: Basic Launch

**User Story:** As a Windows user, I want to run GDU in PowerShell, so that I can analyze disk usage.

#### Acceptance Criteria

1. WHEN a user runs `gdu.exe -h` THEN THE System SHALL display help information
2. WHEN a user runs `gdu.exe -v` THEN THE System SHALL display version information
3. WHEN a user runs `gdu.exe` without arguments THEN THE System SHALL start in interactive mode
4. THE System SHALL use `NUL` instead of `/dev/null` for log file on Windows
5. THE System SHALL use `os.TempDir()` for storage path on Windows

### Requirement 2: Windows Drive Display

**User Story:** As a Windows user, I want to see available drives (C:, D:, etc.), so that I can choose which drive to analyze.

#### Acceptance Criteria

1. WHEN a user runs `gdu.exe -d` THEN THE System SHALL display list of Windows drives
2. THE System SHALL show drive letter, volume name, and filesystem type
3. THE System SHALL show total size and free space for each drive
4. THE System SHALL skip network drives and CD-ROMs by default

### Requirement 3: Directory Analysis

**User Story:** As a Windows user, I want to analyze directories, so that I can find large files.

#### Acceptance Criteria

1. WHEN a user runs `gdu.exe -n -p .` THEN THE System SHALL analyze current directory
2. WHEN a user runs `gdu.exe -n -p C:\Users\...` THEN THE System SHALL analyze specified path
3. WHEN a user runs `gdu.exe -s -n .` THEN THE System SHALL show total size only
4. THE System SHALL handle Windows paths with backslashes correctly
5. THE System SHALL use `filepath.Separator` for path operations

### Requirement 4: Interactive Mode (TUI)

**User Story:** As a Windows user, I want to use interactive interface, so that I can navigate and manage files.

#### Acceptance Criteria

1. WHEN a user runs `gdu.exe` THEN THE System SHALL display TUI interface
2. THE System SHALL support keyboard navigation (arrows, Enter, q)
3. THE System SHALL support file deletion (d key)
4. THE System SHALL support shell spawning (b key) using cmd.exe
5. THE System SHALL support file viewing (v key)
6. THE System SHALL support export to JSON (e key)

### Requirement 5: Export/Import

**User Story:** As a Windows user, I want to export results to JSON, so that I can save or share analysis.

#### Acceptance Criteria

1. WHEN a user runs `gdu.exe -o report.json -n .` THEN THE System SHALL create JSON file
2. WHEN a user runs `gdu.exe -f report.json -n` THEN THE System SHALL read JSON file
3. THE System SHALL store Windows paths correctly in JSON

### Requirement 6: Filtering and Sorting

**User Story:** As a Windows user, I want to filter and sort results, so that I can find specific files.

#### Acceptance Criteria

1. THE System SHALL support filtering by name (/ key in TUI)
2. THE System SHALL support sorting by size, name, count, mtime
3. THE System SHALL support time-based filtering (--since, --until, --max-age, --min-age)
4. THE System SHALL support ignoring directories (-i flag)
5. THE System SHALL support ignoring hidden files (-H flag)

### Requirement 7: Additional Features

**User Story:** As a Windows user, I want all GDU features to work, so that I have full functionality.

#### Acceptance Criteria

1. THE System SHALL support archive browsing (--archive-browsing)
2. THE System SHALL support path collapsing (--collapse-path)
3. THE System SHALL support mouse input (--mouse)
4. THE System SHALL support configuration file (~/.gdu.yaml)
5. THE System SHALL support parallel deletion for better performance

## Technical Requirements

### TR-1: Platform-specific Build Tags

| File | Build Tag | Description |
|------|-----------|-------------|
| `pkg/device/dev_windows.go` | `windows` | Windows device detection |
| `pkg/device/dev_linux.go` | `linux` | Linux device detection |
| `pkg/device/dev_other.go` | `plan9` | Plan9 fallback |
| `pkg/analyze/dir_windows.go` | `windows` | Windows file attributes |
| `pkg/analyze/dir_linux-openbsd.go` | `linux \|\| openbsd` | Linux/OpenBSD file attributes |
| `pkg/analyze/dir_other.go` | `plan9` | Plan9 fallback |
| `tui/exec_windows.go` | `windows` | Windows shell spawning |
| `tui/exec_other.go` | `!windows` | Unix shell spawning |

### TR-2: Completed Fixes

| File | Fix | Status |
|------|-----|--------|
| `cmd/gdu/main.go` | NUL instead of /dev/null | ✅ Done |
| `cmd/gdu/main.go` | Empty ignore dirs for Windows | ✅ Done |
| `cmd/gdu/main.go` | os.TempDir() for storage path | ✅ Done |
| `pkg/device/dev_windows.go` | Windows API for drives | ✅ Done |
| `pkg/device/dev_other.go` | Build tag plan9 only | ✅ Done |
| `pkg/analyze/dir_windows.go` | Windows file attributes | ✅ Done |
| `pkg/analyze/dir_other.go` | Build tag plan9 only | ✅ Done |
| `pkg/path/path.go` | filepath.Separator | ✅ Done |

### TR-3: Verified Working

| Feature | Command | Status |
|---------|---------|--------|
| Help | `gdu.exe -h` | ✅ Works |
| Version | `gdu.exe -v` | ✅ Works |
| List drives | `gdu.exe -d -n` | ✅ Works |
| Analyze dir | `gdu.exe -n -p .` | ✅ Works |
| Summarize | `gdu.exe -s -n .` | ✅ Works |
| Export JSON | `gdu.exe -o file.json -n .` | ✅ Works |
| Import JSON | `gdu.exe -f file.json -n` | ✅ Works |
| Top files | `gdu.exe -n --top 5 .` | ✅ Works |
| Ignore dirs | `gdu.exe -n -i ".git,pkg" .` | ✅ Works |
| Ignore hidden | `gdu.exe -n -H .` | ✅ Works |
| Interactive TUI | `gdu.exe` | ✅ Works |
| Navigation | TUI arrows, Enter, q | ✅ Works |
| Delete files | TUI 'd' key | ✅ Works |
| Spawn shell | TUI 'b' key (cmd.exe) | ✅ Works |
| View file | TUI 'v' key | ✅ Works |
| Export TUI | TUI 'e' key | ✅ Works |

### TR-4: Bug Fixes Applied

| File | Fix | Description |
|------|-----|-------------|
| `tui/actions.go` | File handle leak | Added `defer file.Close()` after `os.Create()` in export |
| `pkg/analyze/zipdir.go` | Path separators | Fixed to use forward slashes in ZIP archives |
| `internal/common/ignore.go` | Regex escaping | Escape backslashes in Windows paths for regex |
| `tui/*_test.go` | File locking | Fixed tests to properly close file handles |
| `pkg/path/path_test.go` | Cross-platform | Rewrote tests for Windows path separators |
| `stdout/stdout_test.go` | Path separators | Updated to use `filepath.Join()` |
| `report/import_test.go` | Path separators | Updated to use `filepath.Join()` |
| `tui/exec_test.go` | Build tag | Added `!windows` tag, created Windows-specific test |
| `pkg/analyze/dir_test.go` | Symlinks/Hardlinks | Skip tests requiring elevated privileges |
| `pkg/analyze/sequential_test.go` | Symlinks/Hardlinks | Skip tests requiring elevated privileges |
| `pkg/analyze/stored_test.go` | File usage | Handle Windows file usage differences |
| `cmd/gdu/app/app_test.go` | Error messages | Handle Windows error message format |

## Files Modified/Created

| File | Status | Description |
|------|--------|-------------|
| `cmd/gdu/main.go` | Modified | Windows-specific defaults (NUL, TempDir) |
| `pkg/device/dev_windows.go` | Created | Windows device detection using Windows API |
| `pkg/device/dev_other.go` | Modified | Build tag plan9 only |
| `pkg/analyze/dir_windows.go` | Created | Windows-specific file attributes |
| `pkg/analyze/dir_other.go` | Modified | Build tag plan9 only |
| `pkg/analyze/zipdir.go` | Modified | Fixed path separators in ZIP archives |
| `pkg/path/path.go` | Modified | Use filepath.Separator |
| `tui/actions.go` | Modified | Fixed file handle leak in export |
| `tui/exec_windows_test.go` | Created | Windows-specific shell test |
| `tui/exec_test.go` | Modified | Added !windows build tag |
| `tui/*_test.go` | Modified | Fixed file locking issues |
| `pkg/path/path_test.go` | Modified | Cross-platform path tests |
| `stdout/stdout_test.go` | Modified | Cross-platform path tests |
| `report/import_test.go` | Modified | Cross-platform path tests |
| `pkg/remove/remove_test.go` | Modified | Windows error message handling |
| `pkg/analyze/symlink_test.go` | Modified | Skip symlink tests on Windows |
| `pkg/analyze/parallel_coverage_test.go` | Modified | Skip symlink tests on Windows |
| `pkg/analyze/dir_test.go` | Modified | Skip symlink/hardlink tests on Windows |
| `pkg/analyze/sequential_test.go` | Modified | Skip symlink/hardlink tests on Windows |
| `pkg/analyze/stored_test.go` | Modified | Handle Windows file usage differences |
| `cmd/gdu/app/app_test.go` | Modified | Handle Windows error messages |
| `internal/common/ignore.go` | Modified | Escape backslashes in regex patterns |

## Windows-specific Test Files

| File | Description |
|------|-------------|
| `cmd/gdu/main_windows_test.go` | Platform defaults tests |
| `pkg/device/dev_windows_test.go` | Device detection tests |
| `pkg/path/path_windows_test.go` | Path handling tests |
| `pkg/analyze/dir_windows_test.go` | Directory analysis tests |
| `tui/exec_windows_test.go` | Shell execution tests |
| `tui/tui_windows_test.go` | TUI integration tests |
| `stdout/stdout_windows_test.go` | Non-interactive mode tests |
| `report/roundtrip_windows_test.go` | JSON export/import tests |
| `internal/common/ignore_windows_test.go` | Ignore pattern tests |
