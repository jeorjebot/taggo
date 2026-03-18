# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [v1.2.3] - 2026-03-18

### Changed
- Changelog updates are no longer auto-committed; taggo only writes `CHANGELOG.md` to disk.

## [v1.2.2] - 2026-03-18

## Fixed
- `--version` now shows the correct version when installed via `go install ...@<version>`.

## [v1.2.1] - 2026-03-18

### Added
- Re-added `taggo init` as a positional subcommand for backward compatibility.

## [v1.2.0] - 2026-03-18

### Added
- Added `-v`/`--version` flag to show taggo version.
- Added `-l`/`--list` command to list tags in the current branch with dates.
- Added automatic `CHANGELOG.md` management in [Keep a Changelog](https://keepachangelog.com/en/1.0.0/) format on tag creation.
- Added `--no-changelog` flag to skip automatic changelog update.
- Handled pre-release tags with the `-n` or `--pre-release` flag.
- Added Homebrew Cask for macOS installation (ARM and Intel). Contributed by @thetombrider.

### Changed
- Tags without the `v` prefix can now be created.

### Fixed
- Fixed `CurrentTag` to work cross-platform without `tail` dependency.
- Restored gopher image in the README.

## [v1.1.0] - 2023-05-02

### Added
- Added commands for creating a specific tag.

### Fixed
- Fixed windows get current tag command, without using `tail`.

## [v1.0.0] - 2023-04-30

### Added
- Added commands for creating major, minor, and patch tags.
- Added command for deleting the latest tag.

[unreleased]: https://github.com/jeorjebot/taggo/compare/v1.2.3...HEAD
[v1.2.3]: https://github.com/jeorjebot/taggo/compare/v1.2.2...v1.2.3
[v1.2.2]: https://github.com/jeorjebot/taggo/compare/v1.2.1...v1.2.2
[v1.2.1]: https://github.com/jeorjebot/taggo/compare/v1.2.0...v1.2.1
[v1.2.0]: https://github.com/jeorjebot/taggo/compare/v1.1.0...v1.2.0
[v1.1.0]: https://github.com/jeorjebot/taggo/compare/v1.0.0...v1.1.0
[v1.0.0]: https://github.com/jeorjebot/taggo/compare/v0.0.0...v1.0.0
