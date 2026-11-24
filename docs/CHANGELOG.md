# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## Unreleased

## 1.4.0 - 2026-11-23

### Added

- In-cluster snapshot removal with _cluster snaprm [snapshot_name] cluster_name_
- In-cluster snapshot creation with _cluster snapset snapshot_name [parent_snapshot_name] cluster_name_

### Removed

- Database support, which hasn't yet made it yet in the final release

### Changed
- More consistent coloured output in various the snapshot commands


## 1.3.0 - 2026-11-22

### Added

- Shortened command aliases
- `-k` (`--kill`) switch to the VM shutdown commands

## 1.2.5 - 2026-11-17

### Changed

- Another round of TUI standardizations
- Refactored the `snapshot` package
- Refactored `vm setvmem`

## 1.2.4 - 2026-11-17

### Changed

- TUI enhancements
- Refactored (renamed) files across packages

## 1.2.3 - 2026-11-17

### Removed

- Short flag `-d`, as it conflicted with other flags; `--db` is kept

## 1.2.2 - 2026-11-17

### Changed

- Minor TUI enhancements and fixes

### Removed

- Support for archlinux, macOS and Windows

## 1.2.1 - 2026-09-17

### Changed

- Parameter ordering in the `cluster` command

## 1.2.0 - 2026-09-17

### Changed

- Stop/Start VM(s) now pause 3 seconds between operations

## 1.1.0 - 2026-09-14

### Added

- New `cluster` command group: `ls`, `define`, `remove`, `start`, `stop`, `reset`
- `cluster CLUSTER_NAME snaprev [SNAPSHOT_NAME]` to revert every node in a cluster to a snapshot
- `vol attach VM POOL NAME` command: attaches an existing volume to a VM as a new virtio disk, shutting the VM down first if it is running
- Missing `--help`/`Example` text across `src/cmd/` (notably the `vm` and `conn` command groups, which had none)

### Fixed

- Stale/copy-pasted flag and command descriptions in the `conn` group

## 1.0.0 - 2026-09-14

First stable release.

### Added

- New `vol`/`volume` command group: `create`, `rm`, `list`/`ls`
- `pool` command group extended with `create`, `start`, `stop`, `rm` (previously `list` only)

### Changed

- Installed binary renamed from `vmman4` to `vmman` (the package name stays `vmman4` on every format)

## 0.30.0 - 2026-09-09

### Added

- Windows build support

### Changed

- Version numbering is now SemVer-aligned
- Go version bumped to 1.27.1